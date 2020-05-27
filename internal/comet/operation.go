package comet

import (
	"context"
	"time"

	model "github.com/Terry-Mao/goim/api/comet/grpc"
	logic "github.com/Terry-Mao/goim/api/logic/grpc"
	"github.com/Terry-Mao/goim/pkg/strings"
	log "github.com/golang/glog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
)

// Connect connected a connection.
func (s *Server) Connect(c context.Context, p *model.Proto, cookie string) (mid int64, key, rid string, accepts []int32, heartbeat time.Duration, err error) {

	// Body内部含有这些信息，Logic并没有分配
	//	Mid      int64   `json:"mid"`      // 用户在业务中的ID (int64)
	//	Key      string  `json:"key"`      // 客户端标识别,如果为空则自动生成UUID
	//	RoomID   string  `json:"room_id"`  // 客户端加入房间
	//	Platform string  `json:"platform"` // 客户端所在平台
	//	Accepts  []int32 `json:"accepts"`  // 监听房间
	reply, err := s.rpcClient.Connect(c, &logic.ConnectReq{
		Server: s.serverID, // Comet自己的ServerID
		Cookie: cookie,
		Token:  p.Body, // 消息内容
	})
	if err != nil {
		return
	}
	return reply.Mid, reply.Key, reply.RoomID, reply.Accepts, time.Duration(reply.Heartbeat), nil
}

// Disconnect disconnected a connection.
func (s *Server) Disconnect(c context.Context, mid int64, key string) (err error) {
	_, err = s.rpcClient.Disconnect(context.Background(), &logic.DisconnectReq{
		Server: s.serverID,
		Mid:    mid,
		Key:    key,
	})
	return
}

// Heartbeat heartbeat a connection session.
func (s *Server) Heartbeat(ctx context.Context, mid int64, key string) (err error) {
	_, err = s.rpcClient.Heartbeat(ctx, &logic.HeartbeatReq{
		Server: s.serverID,
		Mid:    mid,
		Key:    key,
	})
	return
}

// RenewOnline renew room online.
func (s *Server) RenewOnline(ctx context.Context, serverID string, roomCount map[string]int32) (allRoom map[string]int32, err error) {
	reply, err := s.rpcClient.RenewOnline(ctx, &logic.OnlineReq{
		Server:    s.serverID,
		RoomCount: roomCount,
	}, grpc.UseCompressor(gzip.Name))
	if err != nil {
		return
	}
	return reply.AllRoomCount, nil
}

// Receive receive a message.
func (s *Server) Receive(ctx context.Context, mid int64, p *model.Proto) (err error) {
	_, err = s.rpcClient.Receive(ctx, &logic.ReceiveReq{Mid: mid, Proto: p})
	return
}

// Operate operate.
func (s *Server) Operate(ctx context.Context, p *model.Proto, ch *Channel, b *Bucket) error {
	switch p.Op {
	case model.OpChangeRoom:
		if err := b.ChangeRoom(string(p.Body), ch); err != nil {
			log.Errorf("b.ChangeRoom(%s) error(%v)", p.Body, err)
		}
		p.Op = model.OpChangeRoomReply
	case model.OpSub:
		if ops, err := strings.SplitInt32s(string(p.Body), ","); err == nil {
			ch.Watch(ops...)
		}
		p.Op = model.OpSubReply
	case model.OpUnsub:
		if ops, err := strings.SplitInt32s(string(p.Body), ","); err == nil {
			ch.UnWatch(ops...)
		}
		p.Op = model.OpUnsubReply
	default:
		// TODO ack ok&failed
		// 通知mid
		if err := s.Receive(ctx, ch.Mid, p); err != nil {
			log.Errorf("s.Report(%d) op:%d error(%v)", ch.Mid, p.Op, err)
		}
		p.Body = nil
	}
	return nil
}
