package client

// JWT Revocation Check - Async HTTP Call Pattern (IMPLEMENTED)
//
// This WasmPlugin implements real-time JWT token revocation enforcement using the
// proxy-wasm async HTTP call pattern. The implementation is in main.go.
//
// ## How It Works:
//
// 1. OnHttpRequestHeaders extracts full JWT token from Authorization: Bearer <token>
// 2. Calls DispatchHttpCall to jwt-token-service:8080/authorize with Authorization header
// 3. Returns ActionPause to hold the request
// 4. Callback onAuthorizeResponse receives the HTTP response
// 5. Based on response status:
//    - 200 OK → ResumeHttpRequest() (allow)
//    - 403 Forbidden → SendHttpResponse(403, "Token revoked") (block)
//    - 401 Unauthorized → SendHttpResponse(401, "Invalid token") (block)
//    - Other/Timeout → SendHttpResponse(503) (fail-closed)
//
// ## Configuration:
//
// The WasmPlugin is configured via pluginConfig in the WasmPlugin CRD:
//
//   pluginConfig:
//     cluster: "outbound|8080||jwt-token-service.api-x-devnet.svc.cluster.local"
//     authorizeUrl: "/authorize"
//     timeout: 2000  # milliseconds
//
// ## Authorization Request Format:
//
// The WasmPlugin sends:
//   POST /authorize
//   Authorization: Bearer <full-jwt-token>
//
// The jwt-token-service validates the token, extracts the JTI, and checks revocation status.
//
// ## Cluster Naming Convention:
//
// Envoy cluster names follow the format:
//   outbound|PORT||SERVICE.NAMESPACE.svc.cluster.local
//
// Example:
//   outbound|8080||jwt-token-service.api-x-devnet.svc.cluster.local
//
// This tells Envoy to route HTTP calls to the jwt-token-service in the api-x-devnet namespace.
//
// ## Latency Impact:
//
// Each request requires an async HTTP call to the jwt-token-service, adding ~1-5ms latency
// depending on network conditions and service response time. The timeout is set to 2000ms
// to prevent indefinite blocking.
//
// ## Fail-Closed Behavior:
//
// If the authorization service is unavailable or times out, the request is BLOCKED (fail-closed).
// This ensures revoked tokens cannot be used even if the revocation service is down.
//
// To change to fail-open (allow requests when service unavailable), modify the callback logic
// in main.go to call ResumeHttpRequest() on timeout/error instead of SendHttpResponse(503).
