package handlers

import (
	authn "github.com/Dmitro44/relper/internal/auth"
	audithandler "github.com/Dmitro44/relper/internal/handlers/audit"
	"github.com/Dmitro44/relper/internal/handlers/auth"
	"github.com/Dmitro44/relper/internal/handlers/comments"
	"github.com/Dmitro44/relper/internal/handlers/likes"
	"github.com/Dmitro44/relper/internal/handlers/locations"
	"github.com/Dmitro44/relper/internal/handlers/places"
	"github.com/Dmitro44/relper/internal/handlers/ratings"
	"github.com/Dmitro44/relper/internal/handlers/tags"
	"github.com/Dmitro44/relper/internal/handlers/tasks"
	"github.com/Dmitro44/relper/internal/handlers/tips"
	"github.com/Dmitro44/relper/internal/handlers/users"
	auditsvc "github.com/Dmitro44/relper/internal/services/audit"
	authsvc "github.com/Dmitro44/relper/internal/services/auth"
	commentsvc "github.com/Dmitro44/relper/internal/services/comments"
	likessvc "github.com/Dmitro44/relper/internal/services/likes"
	locationsvc "github.com/Dmitro44/relper/internal/services/locations"
	placesvc "github.com/Dmitro44/relper/internal/services/places"
	ratingsvc "github.com/Dmitro44/relper/internal/services/ratings"
	tagsvc "github.com/Dmitro44/relper/internal/services/tags"
	tasksvc "github.com/Dmitro44/relper/internal/services/tasks"
	tipsvc "github.com/Dmitro44/relper/internal/services/tips"
	usersvc "github.com/Dmitro44/relper/internal/services/users"
	"github.com/danielgtaylor/huma/v2"
)

type Deps struct {
	Users     *usersvc.UserService
	Auth      *authsvc.AuthService
	Audit     *auditsvc.AuditService
	Comments  *commentsvc.CommentService
	Locations *locationsvc.LocationService
	Places    *placesvc.PlaceService
	Tags      *tagsvc.TagService
	Tasks     *tasksvc.TaskService
	Tips      *tipsvc.TipService
	Ratings   *ratingsvc.RatingService
	Likes     *likessvc.LikeService
}

func Register(api huma.API, d Deps) {
	protected := huma.NewGroup(api)
	protected.UseMiddleware(authn.Middleware(api, d.Auth))

	admin := huma.NewGroup(protected)
	admin.UseMiddleware(authn.RequireRole(api, "admin"))

	users.NewUserHandler(d.Users).Register(protected, admin)
	comments.NewCommentHandler(d.Comments).Register(api, protected)
	locations.NewLocationHandler(d.Locations).Register(api, admin)
	places.NewPlaceHandler(d.Places).Register(api, admin)
	ratings.NewRatingHandler(d.Ratings).Register(api, protected)
	tags.NewTagHandler(d.Tags).Register(api, admin)
	tasks.NewTaskHandler(d.Tasks).Register(protected)
	tips.NewTipHandler(d.Tips).Register(api, protected)
	likes.NewLikeHandler(d.Likes).Register(api, protected)

	// auth routes stay on root api and remain public (except /auth/me)
	auth.NewAuthHandler(d.Auth, d.Users).Register(api, protected)
	audithandler.NewAuditHandler(d.Audit).Register(admin)
}

func Config() huma.Config {
	config := huma.DefaultConfig("Relper API", "0.1.0")
	config.DocsRenderer = huma.DocsRendererScalar
	config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"bearerAuth": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		},
	}
	config.Tags = []*huma.Tag{
		{Name: "Auth"},
		{Name: "Audit"},
		{Name: "Users"},
		{Name: "Comments"},
		{Name: "Locations"},
		{Name: "Places"},
		{Name: "Ratings"},
		{Name: "Likes"},
		{Name: "Tags"},
		{Name: "Tasks"},
		{Name: "Tips"},
	}

	return config
}
