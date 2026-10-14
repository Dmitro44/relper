package authn

import (
	"context"
	"net/http"
	"strings"

	authsvc "github.com/Dmitro44/relper/internal/services/auth"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type ctxKey uint8

const userIDKey ctxKey = 0

func Middleware(api huma.API, svc *authsvc.AuthService) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		token, ok := strings.CutPrefix(ctx.Header("Authorization"), "Bearer ")
		if !ok || token == "" {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "missing bearer token")
			return
		}

		claims, err := svc.ValidateAccessToken(ctx.Context(), token)
		if err != nil {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "invalid token subject")
			return
		}

		next(huma.WithValue(ctx, userIDKey, userID))
	}
}

func UserID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDKey).(uuid.UUID)
	return id, ok
}
