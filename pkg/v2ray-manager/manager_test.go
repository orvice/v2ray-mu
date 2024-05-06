package v2raymanager

import (
	"context"
	"io"
	"log/slog"
	"os"
	"testing"
)

func TestManager(t *testing.T) {
	addr := os.Getenv("ADDR")
	l := slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{}))
	cli, err := NewManager(addr, "api", l)

	if err != nil {
		t.Error(err)
		return
	}
	resp, err := cli.GetUserList(context.Background(), false)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("len %d", len(resp))

	for _, v := range resp {
		t.Log(v.User.GetUUID(), v.TrafficInfo.Down, v.TrafficInfo.Up)
	}

}
