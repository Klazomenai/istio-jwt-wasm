package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/klazomenai/istio-jwt-wasm/pkg/jwt"
	"github.com/proxy-wasm/proxy-wasm-go-sdk/proxywasm"
	"github.com/proxy-wasm/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {}

func init() {
	proxywasm.SetVMContext(&vmContext{})
}

// Configuration loaded from WasmPlugin pluginConfig
type pluginConfig struct {
	Cluster      string `json:"cluster"`       // Envoy cluster name for jwt-token-service
	AuthorizeURL string `json:"authorizeUrl"`  // Path to authorize endpoint (e.g., "/authorize")
	Timeout      uint32 `json:"timeout"`       // HTTP call timeout in milliseconds
}

type vmContext struct {
	types.DefaultVMContext
}

func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{}
}

type pluginContext struct {
	types.DefaultPluginContext
	config *pluginConfig
}

func (ctx *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	proxywasm.LogInfo("Autonity JWT Revocation WasmPlugin starting...")

	// Parse plugin configuration
	configData, err := proxywasm.GetPluginConfiguration()
	if err != nil {
		proxywasm.LogCriticalf("Failed to get plugin configuration: %v", err)
		return types.OnPluginStartStatusFailed
	}

	// Default configuration (cluster MUST be provided via WasmPlugin CRD)
	config := &pluginConfig{
		Cluster:      "",
		AuthorizeURL: "/authorize",
		Timeout:      2000, // 2 seconds
	}

	// Parse JSON config if provided
	if len(configData) > 0 {
		if err := json.Unmarshal(configData, config); err != nil {
			proxywasm.LogCriticalf("Failed to parse plugin configuration: %v", err)
			return types.OnPluginStartStatusFailed
		}
	}

	// Validate required configuration
	if config.Cluster == "" {
		proxywasm.LogCritical("cluster is required in pluginConfig")
		return types.OnPluginStartStatusFailed
	}

	ctx.config = config

	proxywasm.LogInfof("✅ Plugin started with config: cluster=%s, authorizeUrl=%s, timeout=%dms",
		config.Cluster, config.AuthorizeURL, config.Timeout)

	return types.OnPluginStartStatusOK
}

func (ctx *pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &httpContext{
		contextID: contextID,
		config:    ctx.config,
	}
}

type httpContext struct {
	types.DefaultHttpContext
	contextID uint32
	config    *pluginConfig
	jti       string // Store JTI for logging
	token     string // Store full JWT token for Authorization header
}

func (ctx *httpContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	// 1. Extract Authorization header
	authHeader, err := proxywasm.GetHttpRequestHeader("Authorization")
	if err != nil {
		proxywasm.LogError("Missing Authorization header")
		return sendError(401, "Unauthorized: Missing Authorization header")
	}

	// 2. Extract Bearer token
	if !strings.HasPrefix(authHeader, "Bearer ") {
		proxywasm.LogError("Invalid Authorization header format")
		return sendError(401, "Unauthorized: Invalid Authorization header format")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	// 3. Extract JTI from JWT (for logging only)
	jti, err := jwt.ExtractJTI(token)
	if err != nil {
		proxywasm.LogErrorf("Failed to extract JTI: %v", err)
		return sendError(401, "Unauthorized: Invalid JWT token")
	}

	// Store both JTI and full token for callback
	ctx.jti = jti
	ctx.token = token

	proxywasm.LogInfof("🔍 Checking revocation status for JTI: %s", jti)

	// 4. Dispatch async HTTP call to authorization service
	// Send full JWT token in Authorization header (jwt-token-service expects this)
	// Extract hostname from cluster name (format: "outbound|PORT||HOSTNAME")
	authority := extractAuthority(ctx.config.Cluster)
	headers := [][2]string{
		{":method", "POST"},
		{":path", ctx.config.AuthorizeURL},
		{":authority", authority},
		{"authorization", fmt.Sprintf("Bearer %s", token)},
		{"content-length", "0"},
	}

	calloutID, err := proxywasm.DispatchHttpCall(
		ctx.config.Cluster,
		headers,
		nil,  // no body
		nil,  // no trailers
		ctx.config.Timeout,
		ctx.onAuthorizeResponse,
	)

	if err != nil {
		proxywasm.LogCriticalf("Failed to dispatch authorize call: %v", err)
		return sendError(503, "Authorization service unavailable")
	}

	proxywasm.LogInfof("📤 Dispatched authorize call (ID: %d) for JTI: %s", calloutID, jti)

	// 5. Pause request processing until callback completes
	return types.ActionPause
}

// onAuthorizeResponse is called when the /authorize HTTP call completes
func (ctx *httpContext) onAuthorizeResponse(numHeaders, bodySize, numTrailers int) {
	// Get response headers to check status code
	headers, err := proxywasm.GetHttpCallResponseHeaders()
	if err != nil {
		proxywasm.LogCriticalf("Failed to get response headers: %v", err)
		_ = proxywasm.SendHttpResponse(503, nil, []byte("Authorization service error\n"), -1)
		return
	}

	// Find status code in response headers
	var statusCode string
	for _, header := range headers {
		if header[0] == ":status" {
			statusCode = header[1]
			break
		}
	}

	proxywasm.LogInfof("📥 Authorization response for JTI %s: status=%s", ctx.jti, statusCode)

	switch statusCode {
	case "200":
		// Token is valid and not revoked - allow request to proceed
		proxywasm.LogInfof("✅ Token authorized: JTI=%s (status 200)", ctx.jti)
		if err := proxywasm.ResumeHttpRequest(); err != nil {
			proxywasm.LogCriticalf("Failed to resume request: %v", err)
		}

	case "401":
		// Token validation failed (invalid signature, expired, etc.)
		proxywasm.LogWarnf("🚫 Token validation failed: JTI=%s (status 401)", ctx.jti)
		_ = proxywasm.SendHttpResponse(401, nil, []byte("Unauthorized: Invalid JWT token\n"), -1)

	case "403":
		// Token is revoked - block request
		proxywasm.LogWarnf("🚫 Token revoked: JTI=%s (status 403)", ctx.jti)
		_ = proxywasm.SendHttpResponse(403, nil, []byte("Forbidden: Token has been revoked\n"), -1)

	default:
		// Unexpected status or timeout - fail closed (block request)
		proxywasm.LogErrorf("❌ Unexpected authorization response: JTI=%s, status=%s - blocking request (fail-closed)", ctx.jti, statusCode)
		_ = proxywasm.SendHttpResponse(503, nil, []byte("Authorization service unavailable\n"), -1)
	}
}

func sendError(statusCode uint32, body string) types.Action {
	err := proxywasm.SendHttpResponse(statusCode, nil, []byte(body+"\n"), -1)
	if err != nil {
		proxywasm.LogErrorf("Failed to send response: %v", err)
		panic(err)
	}
	return types.ActionPause
}

// extractAuthority extracts the hostname from Envoy cluster name
// Cluster format: "outbound|PORT||HOSTNAME"
// Example: "outbound|8080||jwt-token-service.api-x-devnet.svc.cluster.local" -> "jwt-token-service.api-x-devnet.svc.cluster.local"
func extractAuthority(cluster string) string {
	// Split by "||" and take the second part (hostname)
	parts := strings.Split(cluster, "||")
	if len(parts) >= 2 {
		return parts[1]
	}
	// Fallback: return cluster as-is if format doesn't match
	proxywasm.LogWarnf("Unexpected cluster format: %s, using as-is for authority", cluster)
	return cluster
}
