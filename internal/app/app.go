// Package app bootstraps and runs the application.
package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"wishlist-api/internal/config"
	"wishlist-api/internal/db"
	"wishlist-api/internal/handlers"
	"wishlist-api/internal/middleware"
	"wishlist-api/internal/repository"
	"wishlist-api/internal/service"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "wishlist-api/docs" // swagger documentation
)

const shutdownTimeout = 5 * time.Second

// App holds all application dependencies and the HTTP server.
type App struct {
	cfg    *config.Config
	dbPool *pgxpool.Pool
	router *chi.Mux
	server *http.Server
}

// New creates and configures a new App instance.
func New(ctx context.Context, cfg *config.Config) (*App, error) {
	pool, err := db.NewPostgresPool(ctx,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)
	if err != nil {
		return nil, err
	}

	userRepo := repository.NewUserRepository(pool)
	wishlistRepo := repository.NewWishlistRepository(pool)
	itemRepo := repository.NewItemRepository(pool)

	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	wishlistService := service.NewWishlistService(wishlistRepo)
	itemService := service.NewItemService(itemRepo, wishlistRepo)

	authHandler := handlers.NewAuthHandler(authService)
	wishlistHandler := handlers.NewWishlistHandler(wishlistService)
	itemHandler := handlers.NewItemHandler(itemService)
	publicHandler := handlers.NewPublicHandler(wishlistService, itemService)

	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)
	r.Get("/public/wishlists/{token}", publicHandler.GetWishlistByToken)
	r.Post("/public/wishlists/{token}/items/{item_id}/book", publicHandler.BookItem)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(authService))
		r.Route("/wishlists", func(r chi.Router) {
			r.Post("/", wishlistHandler.Create)
			r.Get("/", wishlistHandler.GetAll)
			r.Get("/{id}", wishlistHandler.GetByID)
			r.Put("/{id}", wishlistHandler.Update)
			r.Delete("/{id}", wishlistHandler.Delete)
			r.Post("/{wishlist_id}/items", itemHandler.Create)
			r.Get("/{wishlist_id}/items", itemHandler.GetAll)
		})
		r.Put("/items/{id}", itemHandler.Update)
		r.Delete("/items/{id}", itemHandler.Delete)
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	server := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: r,
	}

	return &App{
		cfg:    cfg,
		dbPool: pool,
		router: r,
		server: server,
	}, nil
}

// Run starts the HTTP server and handles graceful shutdown.
func (a *App) Run() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("server started", "port", a.cfg.AppPort)
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
			quit <- syscall.SIGTERM
		}
	}()

	<-quit
	slog.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		slog.Error("shutdown error", "error", err)
	}

	a.dbPool.Close()
	slog.Info("server stopped")
}
