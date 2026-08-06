package handlers

import (
	"github.com/Dmitro44/relper/internal/handlers/comments"
	"github.com/Dmitro44/relper/internal/handlers/locations"
	"github.com/Dmitro44/relper/internal/handlers/places"
	"github.com/Dmitro44/relper/internal/handlers/tags"
	"github.com/Dmitro44/relper/internal/handlers/tasks"
	"github.com/Dmitro44/relper/internal/handlers/tips"
	"github.com/Dmitro44/relper/internal/handlers/users"
	commentsvc "github.com/Dmitro44/relper/internal/services/comments"
	locationsvc "github.com/Dmitro44/relper/internal/services/locations"
	placesvc "github.com/Dmitro44/relper/internal/services/places"
	tagsvc "github.com/Dmitro44/relper/internal/services/tags"
	tasksvc "github.com/Dmitro44/relper/internal/services/tasks"
	tipsvc "github.com/Dmitro44/relper/internal/services/tips"
	usersvc "github.com/Dmitro44/relper/internal/services/users"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v3"
)

type Deps struct {
	Users     *usersvc.UserService
	Comments  *commentsvc.CommentService
	Locations *locationsvc.LocationService
	Places    *placesvc.PlaceService
	Tags      *tagsvc.TagService
	Tasks     *tasksvc.TaskService
	Tips      *tipsvc.TipService
}

func Register(app *fiber.App, d Deps) {
	config := huma.DefaultConfig("Relper API", "0.1.0")
	config.DocsRenderer = huma.DocsRendererScalar
	config.OpenAPI.Tags = []*huma.Tag{
		{Name: "Users"},
		{Name: "Comments"},
		{Name: "Locations"},
		{Name: "Places"},
		{Name: "Tags"},
		{Name: "Tasks"},
		{Name: "Tips"},
	}

	api := humafiber.New(app, config)

	users.NewUserHandler(d.Users).Register(api)
	comments.NewCommentHandler(d.Comments).Register(api)
	locations.NewLocationHandler(d.Locations).Register(api)
	places.NewPlaceHandler(d.Places).Register(api)
	tags.NewTagHandler(d.Tags).Register(api)
	tasks.NewTaskHandler(d.Tasks).Register(api)
	tips.NewTipHandler(d.Tips).Register(api)
}
