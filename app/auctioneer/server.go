package server

import (
	"auctioneer/app/api"
	"auctioneer/app/api/v1"
	"auctioneer/app/cache"
	"auctioneer/app/conf"
	logging "auctioneer/app/logger"
	"auctioneer/app/router"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
)

type Auctioneer struct {
	Fib         *fiber.App
	log         *logging.Logger
	cfg         *conf.Config
	ctx         context.Context
	baseHandler api.Handler
}

func Setup(ctx context.Context, cfg *conf.Config) (*Auctioneer, error) {
	logger, err := logging.NewLogger(cfg.LogLvl)
	if err != nil {
		return nil, fmt.Errorf("ERROR SETTING UP API'S LOGGER: %v", err)
	}

	auctioneer := NewApp(logger, cfg)
	auctioneer.ctx = ctx
	cache := cache.NewCache()
	auctioneer.baseHandler = api.NewBasehandler(cfg, cache)

	auctioneer.SetupRoutes()

	return auctioneer, nil
}

func (a *Auctioneer) MakeBlizzAuth() error {
	if err := a.baseHandler.V1MakeBlizzAuth(); err != nil {
		return err
	}

	return nil
}

func (a *Auctioneer) GetRealmList() error {
	if err := a.baseHandler.V1GetBlizzRealms(); err != nil {
		return err
	}

	return nil
}

func (a *Auctioneer) Serve() {
	go func() {
		if err := a.Fib.Listen(a.cfg.AppPort); err != nil {
			a.log.Errorf("Error listen: %v", err)
		}
	}()
	a.log.Info("Server started")
	<-a.ctx.Done()

	a.log.Info("Server stopping")
	if err := a.Fib.Shutdown(); err != nil {
		a.log.Fatalf("server Shutdown Failed:%+s", err)
	}
}

func NewApp(logger *logging.Logger, cfg *conf.Config) *Auctioneer {
	app := new(Auctioneer)
	app.Fib = fiber.New(fiber.Config{
		ErrorHandler:          app.errorHandler,
		DisableStartupMessage: true,
	})
	app.Fib.Use(fiberLogger.New())
	app.log = logger
	app.cfg = cfg

	return app
}

func (a *Auctioneer) SetupRoutes() {
	v1 := a.Fib.Group("/api/v1")

	router.SetupV1Routes(v1, a.baseHandler.V1Handler())
}

func (a *Auctioneer) errorHandler(c *fiber.Ctx, incomingError error) error {
	code := fiber.StatusInternalServerError
	resp := v1.ResponseV1{
		Success: false,
	}

	if e, ok := incomingError.(*fiber.Error); ok {
		code = e.Code
		resp.Message = e.Message
	} else {
		resp.Message = incomingError.Error()
	}

	a.log.Error(resp.Message)
	return c.Status(code).JSON(resp)
}
