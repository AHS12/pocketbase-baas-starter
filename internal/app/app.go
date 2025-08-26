package app

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	"ims-pocketbase-baas-starter/internal/apidoc"
	"ims-pocketbase-baas-starter/internal/commands"
	"ims-pocketbase-baas-starter/internal/crons"
	_ "ims-pocketbase-baas-starter/internal/database/migrations" //side effect migration load(from pocketbase)
	"ims-pocketbase-baas-starter/internal/hooks"
	"ims-pocketbase-baas-starter/internal/jobs"
	"ims-pocketbase-baas-starter/internal/middlewares"
	"ims-pocketbase-baas-starter/internal/routes"
	"ims-pocketbase-baas-starter/pkg/logger"
	"ims-pocketbase-baas-starter/pkg/metrics"
)

// NewApp creates and configures a new PocketBase app instance
func NewApp() *pocketbase.PocketBase {
	app := pocketbase.New()

	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate:  isGoRun, // auto-create migration files only in dev
		TemplateLang: migratecmd.TemplateLangGo,
	})

	logger := logger.GetLogger(app)
	logger.SetStoreLogs(true) // Enable storing logs in DB

	metricsProvider := metrics.GetInstance()
	logger.Info("Metrics provider initialized", "provider", metricsProvider != nil)

	jobManager := jobs.GetJobManager()
	// Only initialize if not already initialized
	if jobManager.GetProcessor() == nil {
		if err := jobManager.Initialize(app); err != nil {
			log.Fatalf("Failed to initialize job manager: %v", err)
		}
	}

	logger.Info("Registering job handlers")
	if err := jobs.RegisterJobs(app); err != nil {
		log.Fatalf("Failed to register job handlers: %v", err)
	}

	logger.Info("Registering scheduled cron jobs")
	if err := crons.RegisterCrons(app); err != nil {
		log.Fatalf("Failed to register schedule crons: %v", err)
	}

	logger.Info("Registering custom console commands")
	if err := commands.RegisterCommands(app); err != nil {
		log.Fatalf("Failed to register commands: %v", err)
	}

	logger.Info("Registering custom event hooks")
	if err := hooks.RegisterHooks(app); err != nil {
		log.Fatalf("Failed to register hooks: %v", err)
	}

	app.OnTerminate().BindFunc(func(te *core.TerminateEvent) error {
		if metricsProvider != nil {
			logger.Info("Shutting down metrics provider")
			if err := metricsProvider.Shutdown(context.Background()); err != nil {
				logger.Error("Failed to shutdown metrics provider", "error", err)
			}
		}
		return te.Next()
	})

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		if handler := metricsProvider.GetHandler(); handler != nil {
			se.Router.GET("/metrics", func(e *core.RequestEvent) error {
				handler.ServeHTTP(e.Response, e.Request)
				return nil
			})
			logger.Info("Metrics endpoint registered", "path", "/metrics")
		}

		generator := apidoc.InitializeGenerator(app)

		logger.Info("Registering middlewares")
		if err := middlewares.RegisterMiddlewares(se); err != nil {
			log.Fatalf("Failed to register middlewares: %v", err)
		}

		se.Router.GET("/{path...}", apis.Static(os.DirFS("./pb_public"), false))

		if err := routes.RegisterCustom(se); err != nil {
			log.Fatalf("Failed to register routes: %v", err)
		}

		apidoc.RegisterEndpoints(se, generator)

		return se.Next()
	})

	return app
}

func Run() {
	app := NewApp()

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
