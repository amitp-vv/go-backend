package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type AuthPayload struct {
	Address string `json:"address"`
	jwt.RegisteredClaims
}

type contextKey string

const authPayloadKey = contextKey("authPayload")

// Attach AuthPayload to request context
func setAuthPayload(r *http.Request, payload *AuthPayload) *http.Request {
	ctx := context.WithValue(r.Context(), authPayloadKey, payload)
	return r.WithContext(ctx)
}

// Retrieve AuthPayload from request context
func getAuthPayload(r *http.Request) *AuthPayload {
	val := r.Context().Value(authPayloadKey)
	if payload, ok := val.(*AuthPayload); ok {
		return payload
	}
	return nil
}

// RequireAuth middleware: validates JWT and attaches payload to context
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		tokenStr := ""
		if strings.HasPrefix(header, "Bearer ") {
			tokenStr = header[7:]
		}
		if tokenStr == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "missing_token"})
			return
		}
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "devsecret"
		}
		token, err := jwt.ParseWithClaims(tokenStr, &AuthPayload{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid_token"})
			return
		}
		claims, ok := token.Claims.(*AuthPayload)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid_token"})
			return
		}
		r = setAuthPayload(r, claims)
		next.ServeHTTP(w, r)
	})
}

// RequireAdmin middleware: checks if authenticated user is admin
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := getAuthPayload(r)
		if auth == nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
			return
		}
		list := strings.ToLower(os.Getenv("ADMIN_ADDRESSES"))
		addresses := []string{}
		for _, addr := range strings.Split(list, ",") {
			addr = strings.TrimSpace(addr)
			if addr != "" {
				addresses = append(addresses, addr)
			}
		}
		found := false
		for _, addr := range addresses {
			if addr == strings.ToLower(auth.Address) {
				found = true
				break
			}
		}
		if !found {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "forbidden"})
			return
		}
		next.ServeHTTP(w, r)
	})
}
