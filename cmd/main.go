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
	auditrepo "github.com/Dmitro44/relper/internal/repositories/audit/pgx"
	commentsrepo "github.com/Dmitro44/relper/internal/repositories/comments/pgx"
	likesrepo "github.com/Dmitro44/relper/internal/repositories/likes/pgx"
	locationsrepo "github.com/Dmitro44/relper/internal/repositories/locations/pgx"
	placesrepo "github.com/Dmitro44/relper/internal/repositories/places/pgx"
	ratingsrepo "github.com/Dmitro44/relper/internal/repositories/ratings/pgx"
	reftokrepo "github.com/Dmitro44/relper/internal/repositories/refreshtokens/pgx"
	rolesrepo "github.com/Dmitro44/relper/internal/repositories/roles/pgx"
	tagsrepo "github.com/Dmitro44/relper/internal/repositories/tags/pgx"
	tasksrepo "github.com/Dmitro44/relper/internal/repositories/tasks/pgx"
	tipsrepo "github.com/Dmitro44/relper/internal/repositories/tips/pgx"
	usersrepo "github.com/Dmitro44/relper/internal/repositories/users/pgx"
	auditsvc "github.com/Dmitro44/relper/internal/services/audit"
	"github.com/Dmitro44/relper/internal/services/auth"
	"github.com/Dmitro44/relper/internal/services/comments"
	"github.com/Dmitro44/relper/internal/services/likes"
	"github.com/Dmitro44/relper/internal/services/locations"
	"github.com/Dmitro44/relper/internal/services/places"
	"github.com/Dmitro44/relper/internal/services/ratings"
	"github.com/Dmitro44/relper/internal/services/tags"
	"github.com/Dmitro44/relper/internal/services/tasks"
	"github.com/Dmitro44/relper/internal/services/tips"
	userservice "github.com/Dmitro44/relper/internal/services/users"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
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
	auditService := auditsvc.NewAuditService(auditrepo.NewAuditRepo(pool))
	userService := userservice.NewUserService(usersrepo.NewUsersRepo(pool), rolesrepo.NewRolesRepo(pool), auditService)
	authService := auth.NewAuthService(userService, reftokrepo.NewRefreshTokenRepo(pool), cfg, auditService)
	placeService := places.NewPlaceService(placesrepo.NewPlacesRepo(pool), locationService)
	tipService := tips.NewTipService(tipsrepo.NewTipsRepo(pool), placeService, locationService)
	commentService := comments.NewCommmentService(commentsrepo.NewCommentsRepo(pool), placeService, tipService)
	ratingService := ratings.NewRatingService(ratingsrepo.NewRatingsRepo(pool), placeService)
	likeService := likes.NewLikeService(likesrepo.NewLikesRepo(pool), tipService, commentService)

	app := fiber.New()
	app.Use(middlewares.RequestLogger())

	api := humafiber.New(app, handlers.Config())

	handlers.Register(api, handlers.Deps{
		Users:     userService,
		Auth:      authService,
		Audit:     auditService,
		Comments:  commentService,
		Locations: locationService,
		Places:    placeService,
		Tags:      tags.NewTagService(tagsrepo.NewTagsRepo(pool)),
		Tasks:     tasks.NewTaskService(tasksrepo.NewTasksRepo(pool)),
		Tips:      tipService,
		Ratings:   ratingService,
		Likes:     likeService,
	})

	log.Fatal(app.Listen(":3000"))
}
