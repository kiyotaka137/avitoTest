package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"avitoTest/internal/config"
	transport "avitoTest/internal/http"
	"avitoTest/internal/repository"
	"avitoTest/internal/repository/postgres"
	"avitoTest/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

type App struct {
	engine *gin.Engine
	server *http.Server
	log    *slog.Logger
	config config.Config
	pool   *pgxpool.Pool
}

func New(log *slog.Logger, cfg *config.Config) *App {
	dsn := buildDSN(cfg)
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		log.Error("db connect failed", "error", err)
		panic(err)
	}

	tx := repository.NewTxManager(pool)

	teamRepo := postgres.NewTeamRepo(pool, log)
	userRepo := postgres.NewUserRepo(pool, log)
	prRepo := postgres.NewPRRepo(pool, log)

	teamSvc := service.NewTeamService(tx, teamRepo)
	userSvc := service.NewUserService(tx, userRepo, prRepo)
	prSvc := service.NewPRService(tx, prRepo, userRepo, teamRepo)

	teamHandler := transport.NewTeamHandler(teamSvc)
	userHandler := transport.NewUserHandler(userSvc)
	prHandler := transport.NewPullRequestHandler(prSvc)

	statsRepo := postgres.NewStatsRepo(pool, log)
	statsSvc := service.NewStatsService(tx, statsRepo)
	statsH := transport.NewStatsHandler(statsSvc)

	engine := transport.NewRouter(teamHandler, userHandler, prHandler, statsH)

	srv := &http.Server{
		Addr:         cfg.Http.Address,
		Handler:      engine,
		ReadTimeout:  cfg.Http.Timeout,
		WriteTimeout: cfg.Http.Timeout,
		IdleTimeout:  cfg.Http.IdleTimeout,
	}

	return &App{
		engine: engine,
		server: srv,
		log:    log,
		config: *cfg,
		pool:   pool,
	}
}

func (a *App) Run() error {
	a.log.Info("starting http server", "addr", a.config.Http.Address)
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	a.log.Info("shutting down http server")
	defer a.pool.Close()
	return a.server.Shutdown(ctx)
}

func buildDSN(cfg *config.Config) string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.Database.User, cfg.Database.Password),
		Host:   fmt.Sprintf("%s:%d", cfg.Database.Host, cfg.Database.Port),
		Path:   cfg.Database.DBName,
	}
	q := u.Query()
	q.Set("sslmode", cfg.Database.SSLMode)
	u.RawQuery = q.Encode()
	return u.String()
}
