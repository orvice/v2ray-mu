package server

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/p4gefau1t/trojan-go/api/service"
	"google.golang.org/grpc"
)

type TrojanMgr struct {
	client service.TrojanServerServiceClient

	userStatusMap map[int64]*service.UserStatus
}

func newTrojanMgr(addr string) (*TrojanMgr, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := service.NewTrojanServerServiceClient(conn)

	return &TrojanMgr{
		client:        client,
		userStatusMap: make(map[int64]*service.UserStatus),
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

func (t *TrojanMgr) GetUser(ctx context.Context, password string) (*service.GetUsersResponse, error) {
	var err error

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	stream, err := t.client.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

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

func (t *TrojanMgr) RemoveUser(ctx context.Context, password string) error {
	var err error

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	stream, err := t.client.SetUsers(ctx)
	if err != nil {
		return err
	}

	err = stream.Send(&service.SetUsersRequest{
		Operation: service.SetUsersRequest_Delete,
		Status: &service.UserStatus{
			User: &service.User{
				Password: password,
			},
		},
	})

	if err != nil {
		tjLogger.Errorw("[trojan] remove user fail ",
			"error", err,
		)
		return err
	}

	resp, err := stream.Recv()
	if err != nil {
		tjLogger.Errorw("[trojan] remove user fail ",
			"error", err,
		)
		return err
	}
	tjLogger.Infow("[trojan] remove user success ",
		"resp", resp,
	)
	return nil
}

func (t *TrojanMgr) AddUser(ctx context.Context, password string) error {
	var err error

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	stream, err := t.client.SetUsers(ctx)
	if err != nil {
		return err
	}

	err = stream.Send(&service.SetUsersRequest{
		Operation: service.SetUsersRequest_Add,
		Status: &service.UserStatus{
			User: &service.User{
				Password: password,
			},
		},
	})

	if err != nil {
		tjLogger.Errorw("[trojan] add user fail ",
			"error", err,
		)
		return err
	}

	resp, err := stream.Recv()
	if err != nil {
		tjLogger.Errorw("[trojan] add user fail ",
			"error", err,
		)
		return err
	}
	tjLogger.Infow("[trojan] add user success ",
		"resp", resp,
	)
	return nil
}
