package main

import (
	"bufio"
	"context"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// Make host and connect to peers
	h, err := makeHost(ctx)
	if err != nil {
		runtime.LogErrorf(ctx, "Error making host: %s\n", err.Error())
		return
	}

	handler := NewHandler(ctx)
	h.SetStreamHandler(protocol.ID(PROTOCOL), handler.handleStream)
	runtime.LogInfof(ctx, "\n[*] Your multiaddress is: /ip4/%s/tcp/%d/p2p/%s\n", HOST, PORT, h.ID().Pretty())

	peerChan, err := initMDNS(h, RENDEZVOUS)
	if err != nil {
		runtime.LogErrorf(ctx, "Error setting up MDNS: %s\n", err.Error())
		return
	}

	// Listen in goroutine so that the app startup process can continue
	go func() {
		for p := range peerChan {
			// Run in goroutine so that app can connect to other peers while doing the alias handshake
			go func(p peer.AddrInfo) {
				runtime.LogInfof(ctx, "Found peer: %s\n", p.ID.Pretty())
				if err := h.Connect(ctx, p); err != nil {
					runtime.LogErrorf(ctx, "Connection failed: %s\n", err.Error())
					return
				}
				// Open a stream with new peer to ask for alias
				stream, err := h.NewStream(ctx, p.ID, protocol.ID(PROTOCOL))
				if err != nil {
					runtime.LogErrorf(ctx, "Stream open failed: %s\n", err.Error())
					return
				}
				// defer stream.Close()
				rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
				// Send alias ping
				msg := ALIAS_REQUEST + "nicki:" + h.ID().String()
				if err = writeOnce(ctx, rw, msg); err != nil {
					runtime.LogErrorf(ctx, "Buffer error: %s\n", err.Error())
					return
				}
				// Wait for alias pong
				for {
					str, err := readOnce(ctx, rw)
					if err != nil {
						runtime.LogErrorf(ctx, "Error reading from buffer: %s\n", err.Error())
						break
					}
					if str[:10] == ALIAS_RESPONSE {
						info := Peer{ID: p.ID.String(), Alias: str[10:]}
						runtime.EventsEmit(ctx, "data_backend", info)
						break
					}
				}
			}(p)
		}
	}()
}
