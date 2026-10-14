package handlers

import (
	"github.com/Dmitro44/relper/internal/authn"
	"github.com/Dmitro44/relper/internal/handlers/auth"
	"github.com/Dmitro44/relper/internal/handlers/comments"
	"github.com/Dmitro44/relper/internal/handlers/locations"
	"github.com/Dmitro44/relper/internal/handlers/places"
	"github.com/Dmitro44/relper/internal/handlers/tags"
	"github.com/Dmitro44/relper/internal/handlers/tasks"
	"github.com/Dmitro44/relper/internal/handlers/tips"
	"github.com/Dmitro44/relper/internal/handlers/users"
	authsvc "github.com/Dmitro44/relper/internal/services/auth"
	commentsvc "github.com/Dmitro44/relper/internal/services/comments"
	locationsvc "github.com/Dmitro44/relper/internal/services/locations"
	placesvc "github.com/Dmitro44/relper/internal/services/places"
	tagsvc "github.com/Dmitro44/relper/internal/services/tags"
	tasksvc "github.com/Dmitro44/relper/internal/services/tasks"
	tipsvc "github.com/Dmitro44/relper/internal/services/tips"
	usersvc "github.com/Dmitro44/relper/internal/services/users"
	"github.com/danielgtaylor/huma/v2"
)

type Deps struct {
	Users     *usersvc.UserService
	Auth      *authsvc.AuthService
	Comments  *commentsvc.CommentService
	Locations *locationsvc.LocationService
	Places    *placesvc.PlaceService
	Tags      *tagsvc.TagService
	Tasks     *tasksvc.TaskService
	Tips      *tipsvc.TipService
}

func Register(api huma.API, d Deps) {
	protected := huma.NewGroup(api)
	protected.UseMiddleware(authn.Middleware(api, d.Auth))

	users.NewUserHandler(d.Users).Register(protected)
	comments.NewCommentHandler(d.Comments).Register(protected)
	locations.NewLocationHandler(d.Locations).Register(protected)
	places.NewPlaceHandler(d.Places).Register(protected)
	tags.NewTagHandler(d.Tags).Register(protected)
	tasks.NewTaskHandler(d.Tasks).Register(protected)
	tips.NewTipHandler(d.Tips).Register(protected)

	// auth routes stay on root api and remain public
	auth.NewAuthHandler(d.Auth).Register(api)
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
		{Name: "Users"},
		{Name: "Comments"},
		{Name: "Locations"},
		{Name: "Places"},
		{Name: "Tags"},
		{Name: "Tasks"},
		{Name: "Tips"},
	}

	return config
}
