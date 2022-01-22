package server

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/p4gefau1t/trojan-go/api/service"
	"google.golang.org/grpc"
)

type TrojanMgr struct {
	client         service.TrojanServerServiceClient
	setUserStream  service.TrojanServerService_SetUsersClient
	listUserStream service.TrojanServerService_ListUsersClient
}

func newTrojanMgr(addr string) (*TrojanMgr, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := service.NewTrojanServerServiceClient(conn)

	setUserStream, err := client.SetUsers(context.Background())
	if err != nil {
		return nil, err
	}
	listStream, err := client.ListUsers(context.Background(), &service.ListUsersRequest{})
	if err != nil {
		return nil, err
	}

	return &TrojanMgr{
		client:         client,
		setUserStream:  setUserStream,
		listUserStream: listStream,
	}, nil
}

func (t *TrojanMgr) ListUsers() {
	stream, err := t.client.ListUsers(context.Background(), &service.ListUsersRequest{})
	if err != nil {
		log.Printf("failed to call grpc command: %v", err)
	}
	for {
		reply, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("faild to recv: %v", err)
		}
		fmt.Println(reply)
		break
	}
}
