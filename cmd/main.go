package cmd

import (
	"context"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	messagev1 "github.com/go-goim/api/message/v1"
	"github.com/go-goim/core/pkg/cmd"
	"github.com/go-goim/core/pkg/graceful"
	"github.com/go-goim/core/pkg/log"
	"github.com/go-goim/core/pkg/mid"

	"github.com/go-goim/push-service/internal/app"
	"github.com/go-goim/push-service/internal/router"
	"github.com/go-goim/push-service/internal/service"

	_ "github.com/swaggo/swag" // use for swagger doc

	_ "github.com/go-goim/push-service/docs" // use for swagger doc
)

var (
	jwtSecret string
)

func init() {
	cmd.GlobalFlagSet.StringVar(&jwtSecret, "jwt-secret", "", "jwt secret")
}

func Main() {
	if err := cmd.ParseFlags(); err != nil {
		panic(err)
	}

	if jwtSecret == "" {
		panic("jwt secret is empty")
	}
	mid.SetJwtHmacSecret(jwtSecret)

	application, err := app.InitApplication()
	if err != nil {
		log.Fatal("initApplication got err", "error", err)
	}

	// register grpc
	messagev1.RegisterPushMessagerServer(application.GrpcSrv, service.GetPushMessager())

	// register router
	g := gin.New()
	g.Use(gin.Recovery(), mid.Logger)
	router.RegisterRouter(g.Group("/push"))
	application.HTTPSrv.HandlePrefix("/", g)
	// register swagger
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err = application.Run(); err != nil {
		log.Fatal("application run got error", "error", err)
	}

	graceful.Register(application.Shutdown)
	if err = graceful.Shutdown(context.TODO()); err != nil {
		log.Info("graceful shutdown got error", "error", err)
	}
}
