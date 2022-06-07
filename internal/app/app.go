package app

import (
	"fmt"
	"net"

	"github.com/go-goim/core/pkg/app"
)

type Application struct {
	*app.Application
}

var (
	application *Application
)

func InitApplication() (*Application, error) {
	// do some own biz logic if needed
	a, err := app.InitApplication()
	if err != nil {
		return nil, err
	}

	application = &Application{
		Application: a,
	}

	return application, nil
}

func GetApplication() *Application {
	return application
}

func GetAgentIP() string {
	return net.JoinHostPort(application.GetHost(), fmt.Sprintf("%d", application.Config.SrvConfig.Http.Port))
}
