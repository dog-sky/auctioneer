package server

import (
	api "auctioneer/app/api"
	"auctioneer/app/api/v1"
	"auctioneer/app/conf"
	logging "auctioneer/app/logger"
	router "auctioneer/app/router"
	"context"
	"fmt"
	fiber "github.com/gofiber/fiber/v2"
)

type Auctioneer struct {
	Fib         *fiber.App
	log         *logging.Logger
	cfg         *conf.Config
	ctx         context.Context
	baseHandler *api.BaseHandler
}

func Setup(ctx context.Context, cfg *conf.Config) (*Auctioneer, error) {
	logger, err := logging.NewLogger(cfg.LogLvl)
	if err != nil {
		return nil, fmt.Errorf("ERROR SETTING UP API'S LOGGER: %v", err)
	}

	auctioneer := NewApp(logger, cfg)
	auctioneer.ctx = ctx
	auctioneer.baseHandler = api.NewBasehandler(cfg).(*api.BaseHandler)

	auctioneer.SetupRoutes()

	return auctioneer, nil
}

func (a *Auctioneer) MakeBlizzAuth() error {
	if err := a.baseHandler.V1.MakeBlizzAuth(); err != nil {
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
	app.log = logger
	app.cfg = cfg

	return app
}

func (a *Auctioneer) SetupRoutes() {
	v1 := a.Fib.Group("/api/v1")

	router.SetupV1Routes(v1, a.baseHandler.V1)
}

func (a *Auctioneer) errorHandler(c *fiber.Ctx, incomingError error) error {
	code := fiber.StatusInternalServerError
	resp := v1.ResponseV1{
		Success: false,
	}

	if e, ok := incomingError.(*fiber.Error); ok {
		code = e.Code
		resp.Result = e.Message
	} else {
		resp.Result = incomingError.Error()
	}

	a.log.Error(resp.Result)
	return c.Status(code).JSON(resp)
}
