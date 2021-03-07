package server

import (
	"context"
	"fmt"
	"time"

	"auctioneer/app/api"
	"auctioneer/app/blizz"
	"auctioneer/app/conf"
	logging "auctioneer/app/logger"
	"auctioneer/app/middleware"
	"auctioneer/app/router"

	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sirupsen/logrus"
)

const (
	readTimeOut = 3 * time.Second
)

type Auctioneer struct {
	Fib         *fiber.App
	log         *logrus.Logger
	cfg         *conf.Config
	ctx         context.Context
	BaseHandler api.Handler
}

func NewApp(ctx context.Context, cfg *conf.Config) (*Auctioneer, error) {
	app := new(Auctioneer)
	app.Fib = fiber.New(fiber.Config{
		ErrorHandler:          middleware.ErrorHandler,
		DisableStartupMessage: true,
		ReadTimeout:           readTimeOut,
	})
	app.Fib.Use(fiberLogger.New())

	logger, err := logging.NewLogger(cfg.LogLvl)
	if err != nil {
		return nil, fmt.Errorf("ERROR SETTING UP API'S LOGGER: %v", err)
	}

	app.log = logger
	app.ctx = ctx
	app.cfg = cfg

	return app, nil
}

func (a *Auctioneer) Setup() {
	blizzClient := blizz.NewClient(a.log, &a.cfg.BlizzApiCfg)
	go blizzClient.BlizzAuthRoutine() // Сервис переавторизовывается в апи близзарда раз в 24 часа.

	a.BaseHandler = api.NewBasehandler(blizzClient)

	router.SetupRoutes(a.Fib, a.BaseHandler)
}

func (a *Auctioneer) MakeBlizzAuth() error {
	return a.BaseHandler.V1MakeBlizzAuth()
}

func (a *Auctioneer) GetRealmList() error {
	return a.BaseHandler.V1GetBlizzRealms()
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
