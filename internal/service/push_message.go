package service

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"

	responsepb "github.com/go-goim/api/transport/response"

	"github.com/go-goim/core/pkg/log"

	messagev1 "github.com/go-goim/api/message/v1"

	"github.com/go-goim/core/pkg/conn/ws"
	"github.com/go-goim/core/pkg/graceful"
	"github.com/go-goim/core/pkg/worker"
)

type PushMessager struct {
	messagev1.UnimplementedPushMessageServiceServer
	workerPool *worker.Pool
}

var (
	pm     *PushMessager
	pmOnce sync.Once
)

func GetPushMessager() *PushMessager {
	pmOnce.Do(func() {
		pm = new(PushMessager)
		pm.workerPool = worker.NewPool(100, 20)
		graceful.Register(pm.workerPool.Shutdown)
	})

	return pm
}

func (p *PushMessager) PushMessage(ctx context.Context, req *messagev1.PushMessageReq) (resp *messagev1.PushMessageResp, err error) {
	log.Info("receive msg", "content", req.String())
	resp = &messagev1.PushMessageResp{
		Response: responsepb.Code_OK.BaseResponse(),
	}
	if len(req.GetToUsers()) == 1 && req.GetToUsers()[0] == -1 {
		// cannot use request ctx in async function.It may kill the goroutine after this request finished.
		go p.handleBroadcastAsync(context.Background(), req)
		return
	}

	for _, uid := range req.GetToUsers() {
		c := ws.Get(strconv.FormatInt(uid, 10))
		if c == nil {
			log.Info("PUSH| user conn not found", "uid", uid)
			resp.FailedUsers = append(resp.FailedUsers, uid)
			continue
		}

		err1 := PushMessage(c, req)
		if err1 != nil {
			log.Error("PUSH| push message failed", "uid", uid, "err", err1.Error())
			resp.FailedUsers = append(resp.FailedUsers, uid)
		}

	}

	return
}

func (p *PushMessager) handleBroadcastAsync(ctx context.Context, req *messagev1.PushMessageReq) {
	ch := ws.LoadAllConn()
	wf := func() error {
		for c := range ch {
			if c.Err() != nil {
				continue
			}

			if err := PushMessage(c, req); err != nil {
				log.Error("PushMessage error", "err", err.Error())
			}
		}

		return nil
	}

	result := p.workerPool.Submit(ctx, wf, 5)
	log.Info("PUSH| workerPool submit", "result", result, "status=", result.Status(), "err=", result.Err())
	if result.Status() == worker.TaskStatusQueueFull {
		log.Error("worker queue buffer is full, should set more buffer")
	}
}

func PushMessage(wc *ws.WebsocketConn, req *messagev1.PushMessageReq) error {
	b, err := json.Marshal(req.Message)
	if err != nil {
		return err
	}

	return wc.Write(b)
}
