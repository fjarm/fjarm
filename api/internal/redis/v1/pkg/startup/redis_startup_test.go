package startup

import (
	"context"
	"io"
	"log/slog"
	"strings"
	"testing"
)

type stringWriterCloser struct {
	buffer string
}

func (sc *stringWriterCloser) Write(p []byte) (n int, err error) {
	sc.buffer += string(p)
	return len(p), nil
}

func (sc *stringWriterCloser) Close() error {
	return nil
}

func TestWriteNewRedisPrimaryConfigFile(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	swc := stringWriterCloser{buffer: ""}

	rss := RedisServerStarter{logger: logger}

	err := rss.WriteNewRedisPrimaryConfigFile(context.Background(), "replicauser", "youshallnotpass", &swc)
	if err != nil {
		t.Errorf("writeNewRedisPrimaryConfigFile got an unexpected error: %v", err)
	}

	if !strings.Contains(swc.buffer, "user replicauser on >") {
		t.Errorf("writeNewRedisPrimaryConfigFile failed to produce valid redis config")
	}
}
