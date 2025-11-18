package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
)

func (m *Middlewares) AuthenticateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(header, " ")
		if len(parts) != 2 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		tokenParts := strings.Split(token, ".")
		if len(tokenParts) != 3 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		jwtHeader := tokenParts[0]
		jwtPayload := tokenParts[1]
		signature := tokenParts[2]

		message := jwtHeader + "." + jwtPayload

		h := hmac.New(sha256.New, []byte(m.cnf.JwtSecretKey))
		h.Write([]byte(message))
		expectedSig := base64UrlEncode(h.Sum(nil))

		if expectedSig != signature {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		payloadBytes, err := base64.RawURLEncoding.DecodeString(jwtPayload)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		type jwtPayloadStruct struct {
			Sub int `json:"sub"`
		}

		var claims jwtPayloadStruct
		if err := json.Unmarshal(payloadBytes, &claims); err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", claims.Sub)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func base64UrlEncode(data []byte) string {
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(data)
}
