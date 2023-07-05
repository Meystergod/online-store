package main

import (
	"context"
	"time"

	"github.com/Meystergod/online-store/internal/config"
	"github.com/Meystergod/online-store/internal/controller"
	"github.com/Meystergod/online-store/internal/delivery/http/httpecho"
	"github.com/Meystergod/online-store/internal/repository/mongo"
	"github.com/Meystergod/online-store/internal/utils"
	"github.com/Meystergod/online-store/pkg/client"
	"github.com/Meystergod/online-store/pkg/httpserver"
	"github.com/Meystergod/online-store/pkg/logging"
	"github.com/Meystergod/online-store/pkg/ossignal"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func main() {
	logger := logging.NewDefaultLogger()

	if err := Run(); err != nil {
		logger.Fatal().Err(err).Msg("running service")
	}
}

func Run() error {
	var cfg config.Config
	if err := envconfig.Process("", &cfg); err != nil {
		return errors.Wrap(err, "reading config")
	}

	ctx := context.Background()
	runner, ctx := errgroup.WithContext(ctx)

	loggerDeps := &logging.LoggerDeps{
		LogLevel: cfg.Log.LogLevel,
		LogFile:  cfg.Log.LogFile,
		LogSize:  cfg.Log.LogSize,
		LogAge:   cfg.Log.LogAge,
	}

	logger, err := logging.NewLogger(loggerDeps)
	if err != nil {
		return errors.Wrap(err, "creating logger")
	}

	ctx = logger.WithContext(ctx)

	httpServerDeps := &httpserver.ServerDeps{
		Address: cfg.HTTPServer.Address,
	}

	httpServer := httpserver.NewServer(httpServerDeps)

	httpServer.Server().Validator = utils.NewValidator()

	dbConfig := client.NewMongoConfig(
		cfg.Database.Auth,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	db, err := client.NewMongoClient(ctx, dbConfig)
	if err != nil {
		return errors.Wrap(err, "connecting database")
	}

	categoryRepository := mongo.NewCategoryRepository(db, utils.CollNameCategory)
	categoryController := controller.NewCategoryController(categoryRepository)
	httpecho.SetCategoryApiRoutes(httpServer.Server(), categoryController)

	subcategoryRepository := mongo.NewSubcategoryRepository(db, utils.CollNameSubcategory)
	subcategoryController := controller.NewSubcategoryController(subcategoryRepository)
	httpecho.SetSubcategoryApiRoutes(httpServer.Server(), subcategoryController)

	discountRepository := mongo.NewDiscountRepository(db, utils.CollNameDiscount)
	discountController := controller.NewDiscountController(discountRepository)
	httpecho.SetDiscountApiRoutes(httpServer.Server(), discountController)

	productRepository := mongo.NewProductRepository(db, utils.CollNameProduct)
	productController := controller.NewProductController(productRepository)
	httpecho.SetProductApiRoutes(httpServer.Server(), productController)

	tagRepository := mongo.NewTagRepository(db, utils.CollNameTag)
	tagController := controller.NewTagController(tagRepository)
	httpecho.SetTagApiRoutes(httpServer.Server(), tagController)

	logger.Info().Msgf("start %s %s on %s", cfg.Application.Name, cfg.Application.Version, cfg.HTTPServer.Address)

	defer logger.Info().Msg("service done")

	runner.Go(func() error {
		if err := httpServer.Start(ctx); err != nil {
			return errors.Wrap(err, "listening and starting http api")
		}

		return nil
	})

	runner.Go(func() error {
		if err := ossignal.DefaultSignalWaiter(ctx); err != nil {
			return errors.Wrap(err, "os signal waiter")
		}

		return nil
	})

	runner.Go(func() error {
		<-ctx.Done()

		ctxSignal, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		defer cancel()

		logger.Info().Msg("shutdown http server")
		if err := httpServer.Shutdown(ctxSignal); err != nil {
			logger.Error().Err(err).Msg("shutdown http server")
		}

		return nil
	})

	if err := runner.Wait(); err != nil {
		switch {
		case ossignal.IsExitSignal(err):
			logger.Info().Msg("exited by exit signal")
		default:
			return errors.Wrap(err, "exiting with error")
		}
	}

	return nil
}
