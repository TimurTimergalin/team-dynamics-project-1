package auth_sdk

import "net/http"

func ExtractAuthFromHTTPRequest(r *http.Request) *IncomingAuth {
	token := r.Header.Get(headerToken)
	if token == "" {
		return nil
	}
	return &IncomingAuth{Token: token}
}
