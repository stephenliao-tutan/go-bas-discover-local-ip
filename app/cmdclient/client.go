package cmdclient

import (
	"google.golang.org/grpc"

	"log"

	"context"
	"fmt"

	"github.com/BASChain/go-bas-discover-local-ip/app/cmdcommon"
	"github.com/BASChain/go-bas-discover-local-ip/app/cmdpb"
	"github.com/BASChain/go-bas-discover-local-ip/config"
	"strings"
)

type CmdConnection struct {
	c      *grpc.ClientConn
	ctx    context.Context
	cancel context.CancelFunc
}

type CmdClient struct {
	addr string
	conn *CmdConnection
}

var cmdClient *CmdClient

func NewCmdClient(addr string) *CmdClient {
	return &CmdClient{addr: addr}
}

func (cc *CmdClient) DialToCmdServer() *CmdConnection {

	if cc.addr == "" {
		cfg := config.GetBasDisCfg()
		cc.addr = cfg.CmdListenPort
	}

	conn, err := grpc.Dial(cc.addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("can not connect to rcp server:", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	cc.conn = &CmdConnection{conn, ctx, cancel}

	return cc.conn

}

func (cc *CmdClient) Close() {
	cc.conn.c.Close()
	cc.conn.cancel()
}

func DefaultCmdSend(addr string, cmd int32) {
	if addr == "" || strings.Contains(addr, "127.0.0.1") {
		if _, err := cmdcommon.IsProcessStarted(); err != nil {
			log.Println(err)
			return
		}
	}

	request := &cmdpb.DefaultRequest{}
	request.Reqid = cmd

	cc := NewCmdClient(addr)

	cc.DialToCmdServer()
	defer cc.Close()

	client := cmdpb.NewDefaultcmdsrvClient(cc.conn.c)

	if resp, err := client.DefaultCmdDo(cc.conn.ctx, request); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp.Message)
	}

}

func (cc *CmdClient) GetRpcClientConn() *grpc.ClientConn {
	return cc.conn.c
}

func (cc *CmdClient) GetRpcCnxt() *context.Context {
	return &cc.conn.ctx
}
