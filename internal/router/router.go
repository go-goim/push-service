package router

import (
	"github.com/gin-gonic/gin"

	"github.com/go-goim/core/pkg/router"

	v1 "github.com/go-goim/push-service/internal/router/v1"
)

type rootRouter struct {
	router.Router
}

func newRootRouter() *rootRouter {
	r := &rootRouter{
		Router: &router.BaseRouter{},
	}

	r.init()
	return r
}

func (r *rootRouter) init() {
	r.Register("/v1", v1.NewRouter())
}

func RegisterRouter(g *gin.RouterGroup) {
	r := newRootRouter()
	r.Load(g)
}
