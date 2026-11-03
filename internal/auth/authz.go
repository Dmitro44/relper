package auth

import (
	"context"
	"net/http"
	"slices"

	"github.com/danielgtaylor/huma/v2"
)

func RequireRole(api huma.API, roles ...string) func(huma.Context, func(huma.Context)) {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(ctx huma.Context, next func(huma.Context)) {
		for _, role := range Roles(ctx.Context()) {
			if _, ok := allowed[role]; ok {
				next(ctx)
				return
			}
		}
		_ = huma.WriteErr(api, ctx, http.StatusForbidden, "insufficient role")
	}
}

func HasRole(ctx context.Context, role string) bool {
	return slices.Contains(Roles(ctx), role)
}
