package server

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"syscall"

	"k8sman/internal/config"
	"k8sman/internal/handlers"
	"k8sman/internal/providers"
	"k8sman/pkg/server"

	"github.com/gin-gonic/gin"
)

type RuntimeContext struct {
	Cfg         *config.ServiceConfig
	ShutdownCtx context.Context
	cancelFunc  context.CancelFunc
}

func NewRuntimeContext() *RuntimeContext {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("unable to load service configuration: %v", err)
	}

	// logger := slog.Default()
	// logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// logger.Setup(cfg.General.AppName, cfg.General.Version)
	// notifications.NotifyServiceRelease(cfg.General)

	shutdownCtx, cancelFunc := context.WithCancel(context.Background())

	return &RuntimeContext{
		Cfg:         cfg,
		ShutdownCtx: shutdownCtx,
		cancelFunc:  cancelFunc,
	}
}

func (rCtx *RuntimeContext) Run() {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		done()
		if r := recover(); r != nil {
			log.Fatalf("application panicked: %v", r)
		}
	}()

	srv, err := server.New(rCtx.Cfg.HTTPServer.Port)
	if err != nil {
		log.Fatalf("server.New: %v", err)
	}
	slog.Info("server running", "port", rCtx.Cfg.HTTPServer.Port)

	err = srv.ServeHTTPHandler(ctx, rCtx.SetupRoutes(ctx))
	if err != nil {
		log.Fatalf("ServeHTTPHandler: %v", err)
	}
	done()

	slog.Info("successful shutdown")
}

func (rCtx *RuntimeContext) SetupRoutes(ctx context.Context) http.Handler {
	if slices.Contains([]string{"dev", "local"}, rCtx.Cfg.General.Stage) {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// mux := gin.New()
	mux := gin.Default()

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	mux.Use(gin.Recovery())

	k8sProvider := providers.NewK8SProvider(rCtx.Cfg.General.Stage)

	httpHandler := handlers.NewHandler(
		rCtx.Cfg,
		k8sProvider,
	)

	unSecuredRoutes := mux.Group("")
	handlers.AddUnsecuredEndpoints(unSecuredRoutes, httpHandler)

	mux.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Page Not Found",
		})
	})
	return mux
}
