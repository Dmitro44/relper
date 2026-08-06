package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/Dmitro44/relper/internal/config"
	"github.com/Dmitro44/relper/internal/database"
	"github.com/Dmitro44/relper/internal/handlers"
	commentsrepo "github.com/Dmitro44/relper/internal/repositories/comments/pgx"
	locationsrepo "github.com/Dmitro44/relper/internal/repositories/locations/pgx"
	placesrepo "github.com/Dmitro44/relper/internal/repositories/places/pgx"
	tagsrepo "github.com/Dmitro44/relper/internal/repositories/tags/pgx"
	tasksrepo "github.com/Dmitro44/relper/internal/repositories/tasks/pgx"
	tipsrepo "github.com/Dmitro44/relper/internal/repositories/tips/pgx"
	usersrepo "github.com/Dmitro44/relper/internal/repositories/users/pgx"
	"github.com/Dmitro44/relper/internal/services/comments"
	"github.com/Dmitro44/relper/internal/services/locations"
	"github.com/Dmitro44/relper/internal/services/places"
	"github.com/Dmitro44/relper/internal/services/tags"
	"github.com/Dmitro44/relper/internal/services/tasks"
	"github.com/Dmitro44/relper/internal/services/tips"
	userservice "github.com/Dmitro44/relper/internal/services/users"
	"github.com/gofiber/fiber/v3"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("unable to load config", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()

	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("unable to connect to database", "error", err)
		os.Exit(1)
	}

	locationService := locations.NewLocationService(locationsrepo.NewLocationsRepo(pool))
	app := fiber.New()
	handlers.Register(app, handlers.Deps{
		Users:     userservice.NewUserService(usersrepo.NewUsersRepo(pool)),
		Comments:  comments.NewCommmentService(commentsrepo.NewCommentsRepo(pool)),
		Locations: locationService,
		Places:    places.NewPlaceService(placesrepo.NewPlacesRepo(pool), *locationService),
		Tags:      tags.NewTagService(tagsrepo.NewTagsRepo(pool)),
		Tasks:     tasks.NewTaskService(tasksrepo.NewTasksRepo(pool)),
		Tips:      tips.NewTipService(tipsrepo.NewTipsRepo(pool)),
	})

	log.Fatal(app.Listen(":3000"))
}
