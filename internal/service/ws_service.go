package service

import (
	"context"
	"time"

	"github.com/gorilla/websocket"

	"github.com/go-goim/core/pkg/consts"
	"github.com/go-goim/core/pkg/log"

	"github.com/go-goim/core/pkg/conn/ws"

	"github.com/go-goim/push-service/internal/app"
)

func HandleWsConn(ctx context.Context, c *websocket.Conn, uid string) {
	ww := ws.WrapWs(ctx, c, uid)
	ww.AddPingAction(func() error {
		return app.GetApplication().Redis.SetEX(context.Background(),
			consts.GetUserOnlineAgentKey(uid), app.GetAgentIP(), consts.UserOnlineAgentKeyExpire).Err()
	})
	ww.AddCloseAction(func() error {
		ctx2, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		return app.GetApplication().Redis.Del(ctx2, consts.GetUserOnlineAgentKey(uid)).Err()

	})

	ctx2, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	err := app.GetApplication().Redis.Set(ctx2, consts.GetUserOnlineAgentKey(uid), app.GetAgentIP(), consts.UserOnlineAgentKeyExpire).Err()
	if err != nil {
		log.Error("redis set error", "key", consts.GetUserOnlineAgentKey(uid), "error", err)
	}

	_ = ww.Write([]byte("success"))
}
