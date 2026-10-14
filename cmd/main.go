package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/Dmitro44/relper/internal/config"
	"github.com/Dmitro44/relper/internal/database"
	"github.com/Dmitro44/relper/internal/handlers"
	"github.com/Dmitro44/relper/internal/logger"
	"github.com/Dmitro44/relper/internal/middlewares"
	commentsrepo "github.com/Dmitro44/relper/internal/repositories/comments/pgx"
	locationsrepo "github.com/Dmitro44/relper/internal/repositories/locations/pgx"
	placesrepo "github.com/Dmitro44/relper/internal/repositories/places/pgx"
	reftokrepo "github.com/Dmitro44/relper/internal/repositories/refreshtokens/pgx"
	tagsrepo "github.com/Dmitro44/relper/internal/repositories/tags/pgx"
	tasksrepo "github.com/Dmitro44/relper/internal/repositories/tasks/pgx"
	tipsrepo "github.com/Dmitro44/relper/internal/repositories/tips/pgx"
	usersrepo "github.com/Dmitro44/relper/internal/repositories/users/pgx"
	"github.com/Dmitro44/relper/internal/services/auth"
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
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	logger := logger.New(cfg.AppEnv)
	slog.SetDefault(logger)

	ctx := context.Background()

	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	locationService := locations.NewLocationService(locationsrepo.NewLocationsRepo(pool))
	userService := userservice.NewUserService(usersrepo.NewUsersRepo(pool))
	authService := auth.NewAuthService(userService, reftokrepo.NewRefreshTokenRepo(pool), cfg)

	app := fiber.New()
	app.Use(middlewares.RequestLogger())

	app.Use(middlewares.RequireAuth(authService, "/auth"))

	handlers.Register(app, handlers.Deps{
		Users:     userService,
		Auth:      authService,
		Comments:  comments.NewCommmentService(commentsrepo.NewCommentsRepo(pool)),
		Locations: locationService,
		Places:    places.NewPlaceService(placesrepo.NewPlacesRepo(pool), *locationService),
		Tags:      tags.NewTagService(tagsrepo.NewTagsRepo(pool)),
		Tasks:     tasks.NewTaskService(tasksrepo.NewTasksRepo(pool)),
		Tips:      tips.NewTipService(tipsrepo.NewTipsRepo(pool)),
	})

	log.Fatal(app.Listen(":3000"))
}
