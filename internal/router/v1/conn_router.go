package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/go-goim/core/pkg/mid"
	"github.com/go-goim/core/pkg/response"
	"github.com/go-goim/core/pkg/router"

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
			CheckOrigin: func(r *http.Request) bool {
				// todo use check uid/token middleware before this handler
				return true
			},
		},
	}
}

func (r *ConnRouter) Load(g *gin.RouterGroup) {
	g.GET("/ws", mid.AuthJwtCookie, r.wsHandler)
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
	// todo use check uid/token middleware before this handler
	uid := c.GetHeader("uid")
	if uid == "" {
		response.ErrorRespWithStatus(c, http.StatusUnauthorized, fmt.Errorf("invalid uid"))
		return
	}

	if !r.upgrader.CheckOrigin(c.Request) {
		response.ErrorRespWithStatus(c, http.StatusForbidden, fmt.Errorf("invalid origin"))
		return
	}

	conn, err := r.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		response.ErrorRespWithStatus(c, http.StatusBadRequest, err)
		return
	}

	service.HandleWsConn(mid.GetContext(c), conn, uid)
}
