package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
	"time"

	"github.com/Dmitro44/relper/internal/config"
	"github.com/Dmitro44/relper/internal/services/auth"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func signToken(t *testing.T, secret string, id uuid.UUID, roles ...string) string {
	t.Helper()
	now := time.Now()
	claims := auth.Claims{
		Subject:   id.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
		Roles:     roles,
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	return token
}

func TestRequireRole(t *testing.T) {
	const secret = "0234qwerty3940ytrewq538430"
	authService := auth.NewAuthService(nil, nil, &config.Config{JWTSecret: secret})
	app := fiber.New()
	api := humafiber.New(app, huma.DefaultConfig("test", "1.0"))

	protected := huma.NewGroup(api)
	protected.UseMiddleware(Middleware(api, authService))

	admin := huma.NewGroup(protected)
	admin.UseMiddleware(RequireRole(api, "admin"))

	huma.Register(admin, huma.Operation{
		OperationID:   "admin-test",
		Method:        http.MethodGet,
		Path:          "/admin/test",
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, _ *struct{}) (*struct{}, error) { return &struct{}{}, nil })

	cases := []struct {
		name  string
		token string
		want  int
	}{
		{"empty", "", http.StatusUnauthorized},
		{"garbage", "abc", http.StatusUnauthorized},
		{"user", signToken(t, secret, uuid.New(), "user"), http.StatusForbidden},
		{"no role", signToken(t, secret, uuid.New()), http.StatusForbidden},
		{"admin", signToken(t, secret, uuid.New(), "admin"), http.StatusOK},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/admin/test", nil)
			if c.token != "" {
				req.Header.Set("Authorization", "Bearer "+c.token)
			}
			resp, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}
			if resp.StatusCode != c.want {
				t.Fatalf("got %d, want %d", resp.StatusCode, c.want)
			}
		})
	}
}

func TestHasRole(t *testing.T) {
	ctx := context.WithValue(context.Background(), rolesKey, []string{"user", "admin"})

	if !HasRole(ctx, "user") {
		t.Fatal("expected user role")
	}
	if !HasRole(ctx, "admin") {
		t.Fatal("expected admin role")
	}
	if HasRole(ctx, "moderator") {
		t.Fatal("unepected moderator role")
	}
	if HasRole(context.Background(), "admin") {
		t.Fatal("empty context must not grant roles")
	}
}

func TestRoles(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		name    string
		context context.Context
		want    []string
	}{
		{"missing", ctx, nil},
		{"wrong type", context.WithValue(ctx, rolesKey, "admin"), nil},
		{"empty", context.WithValue(ctx, rolesKey, []string{}), []string{}},
		{"present", context.WithValue(ctx, rolesKey, []string{"admin"}), []string{"admin"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Roles(c.context)
			if slices.Compare(got, c.want) != 0 {
				t.Fatalf("got: %v, want: %v", got, c.want)
			}
		})
	}
}
