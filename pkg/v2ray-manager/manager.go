package v2raymanager

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/v2fly/v2ray-core/v5/app/proxyman/command"
	statscmd "github.com/v2fly/v2ray-core/v5/app/stats/command"
	"github.com/v2fly/v2ray-core/v5/common/protocol"
	"github.com/v2fly/v2ray-core/v5/common/serial"
	"github.com/v2fly/v2ray-core/v5/proxy/vmess"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Manager struct {
	client      command.HandlerServiceClient
	statsClient statscmd.StatsServiceClient

	inBoundTag string
	logger     *slog.Logger
}

const (
	UplinkFormat   = "user>>>%s>>>traffic>>>uplink"
	DownlinkFormat = "user>>>%s>>>traffic>>>downlink"
)

type TrafficInfo struct {
	Up, Down int64
}

func NewManager(addr, tag string, l *slog.Logger) (*Manager, error) {
	cc, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := command.NewHandlerServiceClient(cc)
	statsClient := statscmd.NewStatsServiceClient(cc)
	m := &Manager{
		client:      client,
		statsClient: statsClient,
		inBoundTag:  tag,
		logger:      l,
	}
	if m.logger == nil {
		m.logger = slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{}))
	}
	return m, nil
}

func (m *Manager) SetLogger(l *slog.Logger) {
	m.logger = l
}

// return is exist,and error
func (m *Manager) AddUser(ctx context.Context, u User) (bool, error) {
	resp, err := m.client.AlterInbound(ctx, &command.AlterInboundRequest{
		Tag: m.inBoundTag,
		Operation: serial.ToTypedMessage(&command.AddUserOperation{
			User: &protocol.User{
				Level: u.GetLevel(),
				Email: u.GetEmail(),
				Account: serial.ToTypedMessage(&vmess.Account{
					Id:               u.GetUUID(),
					AlterId:          u.GetAlterID(),
					SecuritySettings: &protocol.SecurityConfig{Type: protocol.SecurityType_AUTO},
				}),
			},
		}),
	})
	if err != nil && !IsAlreadyExistsError(err) {
		m.logger.Error("failed to call add user",
			"resp", resp,
			"error", err,
		)
		return false, err
	}
	return IsAlreadyExistsError(err), nil
}

func (m *Manager) RemoveUser(ctx context.Context, u User) error {
	resp, err := m.client.AlterInbound(ctx, &command.AlterInboundRequest{
		Tag: m.inBoundTag,
		Operation: serial.ToTypedMessage(&command.RemoveUserOperation{
			Email: u.GetEmail(),
		}),
	})
	if err != nil {
		m.logger.Error("failed to call remove user : ", "error", err)
		return TODOErr
	}
	m.logger.Debug("call remove user resp: ", "resp", resp)

	return nil
}

// @todo error handle
func (m *Manager) GetTrafficAndReset(ctx context.Context, u User) (TrafficInfo, error) {
	ti := TrafficInfo{}
	up, err := m.statsClient.GetStats(ctx, &statscmd.GetStatsRequest{
		Name:   fmt.Sprintf(UplinkFormat, u.GetEmail()),
		Reset_: true,
	})
	if err != nil && !IsNotFoundError(err) {
		m.logger.Error("get traffic user ", "u", u, "error", err)
		return ti, err
	}

	down, err := m.statsClient.GetStats(ctx, &statscmd.GetStatsRequest{
		Name:   fmt.Sprintf(DownlinkFormat, u.GetEmail()),
		Reset_: true,
	})
	if err != nil && !IsNotFoundError(err) {
		m.logger.Error("get traffic user fail",
			"user", u,
			"error", err)
		return ti, nil
	}

	if up != nil {
		ti.Up = up.Stat.Value
	}
	if down != nil {
		ti.Down = down.Stat.Value
	}
	return ti, nil
}

type UserData struct {
	User        User
	TrafficInfo TrafficInfo
}

func (m *Manager) GetUserList(ctx context.Context, reset bool) ([]UserData, error) {
	resp, err := m.statsClient.QueryStats(ctx, &statscmd.QueryStatsRequest{
		Reset_: reset,
	})
	if err != nil {
		return nil, err
	}

	var users = make(map[string]UserData)

	for _, v := range resp.Stat {

		email := getEmailFromStatName(v.GetName())
		uuid := getUUDIFromEmail(email)

		if _, ok := users[uuid]; !ok {
			users[uuid] = UserData{
				TrafficInfo: TrafficInfo{},
			}
		}

		u := user{
			email: email,
			uuid:  uuid,
		}
		ti := users[uuid].TrafficInfo

		if strings.Contains(v.GetName(), "downlink") {
			ti.Down = v.Value
		} else {
			ti.Up = v.Value
		}

		users[uuid] = UserData{
			User:        u,
			TrafficInfo: ti,
		}

	}

	var data = make([]UserData, 0, len(users))
	for _, v := range users {
		data = append(data, v)
	}

	return data, nil
}

func getEmailFromStatName(s string) string {
	arr := strings.Split(s, ">>>")
	if len(arr) > 1 {
		return arr[1]
	}
	return s
}

func getUUDIFromEmail(s string) string {
	arr := strings.Split(s, "@")
	if len(arr) > 0 {
		return arr[0]
	}
	return s
}
