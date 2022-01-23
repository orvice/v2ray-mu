package server

import (
	"context"
	"io"
	"log"

	"github.com/p4gefau1t/trojan-go/api/service"
	"google.golang.org/grpc"
)

type TrojanMgr struct {
	client service.TrojanServerServiceClient
}

func newTrojanMgr(addr string) (*TrojanMgr, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := service.NewTrojanServerServiceClient(conn)

	return &TrojanMgr{
		client: client,
	}, nil
}

func (t *TrojanMgr) ListUsers() ([]*service.UserStatus, error) {
	stream, err := t.client.ListUsers(context.Background(), &service.ListUsersRequest{})
	if err != nil {
		log.Printf("failed to call grpc command: %v", err)
	}
	var out = make([]*service.UserStatus, 0)
	for {
		reply, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("faild to recv: %v", err)
			return out, nil
		}
		out = append(out, reply.Status)
	}
	return out, nil
}

func (t *TrojanMgr) GetUser(ctx context.Context, stream service.TrojanServerService_GetUsersClient, password string) (*service.GetUsersResponse, error) {
	var err error
	err = stream.Send(&service.GetUsersRequest{
		User: &service.User{
			Password: password,
		},
	})

	if err != nil {
		tjLogger.Errorw("[trojan] get user fail ",
			"error", err,
		)
		return nil, err
	}

	resp, err := stream.Recv()
	if err != nil {
		tjLogger.Errorw("[trojan] get user fail ",
			"error", err,
		)
		return nil, err
	}
	return resp, nil
}
