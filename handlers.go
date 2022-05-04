package main

import (
	"bufio"
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Handler struct {
	ctx context.Context
}

func NewHandler(ctx context.Context) *Handler {
	return &Handler{ctx: ctx}
}

func (h *Handler) handleStream(stream network.Stream) {
	runtime.LogInfo(h.ctx, "Got a new stream!")

	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	go h.readData(rw)
	runtime.EventsOn(h.ctx, "data_frontend", func(optionalData ...interface{}) {
		h.writeData(rw, optionalData)
	})
}

func (h *Handler) readData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			runtime.LogErrorf(h.ctx, "Error reading from buffer: %s", err.Error())
			return
		}

		if str == "" {
			return
		}

		runtime.EventsEmit(h.ctx, "data_backend", str)
	}
}

func (h *Handler) writeData(rw *bufio.ReadWriter, data []interface{}) {
	for {
		_, err := rw.WriteString(fmt.Sprintf("%s\n", data[0]))
		if err != nil {
			runtime.LogErrorf(h.ctx, "Error writing to buffer: %s", err.Error())
			return
		}
		err = rw.Flush()
		if err != nil {
			runtime.LogErrorf(h.ctx, "Error flushing buffer: %s", err.Error())
			return
		}
	}
}

func writeOnce(ctx context.Context, rw *bufio.ReadWriter, data string) error {
	_, err := rw.WriteString(fmt.Sprintf("%s\n", data))
	if err != nil {
		return err
	}
	err = rw.Flush()
	if err != nil {
		return err
	}
	return nil
}

func readOnce(ctx context.Context, rw *bufio.ReadWriter) (string, error) {
	str, err := rw.ReadString('\n')
	if err != nil {
		return "", err
	}
	if str == "" {
		return "", err
	}
	return str, nil
}
