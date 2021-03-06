package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/go-goim/core/pkg/mid"
	"github.com/go-goim/core/pkg/router"
	"github.com/go-goim/core/pkg/web/response"

	"github.com/go-goim/push-service/internal/service"
)

type ConnRouter struct {
	router.Router
	upgrader websocket.Upgrader
}

func NewConnRouter() *ConnRouter {
	return &ConnRouter{
		Router: &router.BaseRouter{},
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

func (r *ConnRouter) Load(g *gin.RouterGroup) {
	g.GET("/ws", mid.AuthJwt, r.wsHandler)
}

// @Summary websocket
// @Description websocket 长连接
// @Tags Conn
// @Accept json
// @Produce json
// @Success 200 {string} string "success"
// @Failure 401 {object} response.Response "invalid jwt cookie"
// @Failure 400 {object} response.Response "invalid request"
// @Router /push/v1/conn/ws [get]
func (r *ConnRouter) wsHandler(c *gin.Context) {
	conn, err := r.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		response.ErrorRespWithStatus(c, http.StatusBadRequest, err)
		return
	}

	service.HandleWsConn(mid.GetContext(c), conn, mid.GetUID(c))
}
