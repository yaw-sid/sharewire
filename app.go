package main

import (
	"bufio"
	"context"

	"github.com/libp2p/go-libp2p-core/network"
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

	go func() {
		for p := range peerChan {
			runtime.LogInfo(ctx, "Found peer: "+p.ID.Pretty()+", connecting")
			if err := h.Connect(ctx, p); err != nil {
				runtime.LogErrorf(ctx, "Connection failed: %s\n", err.Error())
				continue
			}

			// Open a stream and exchange aliases
			stream, err := h.NewStream(ctx, p.ID, protocol.ID(PROTOCOL))
			if err != nil {
				runtime.LogErrorf(ctx, "Stream open failed: %s\n", err.Error())
				continue
			}
			rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
			// Send alias
			// TODO: Alias command
			msg := ALIAS_REQUEST + "nicki"
			if err = writeOnce(ctx, rw, msg); err != nil {
				runtime.LogErrorf(ctx, "Buffer error: %s\n", err.Error())
				continue
			}
			// TODO: Create an infinite loop listening for alias response command.
			// Once it is found save it and close the stream
			// A new stream will be opened for subsequent communications
			go func(p peer.AddrInfo, s network.Stream) {
				for {
					str, err := readOnce(ctx, rw)
					if err != nil {
						runtime.LogErrorf(ctx, "Error reading from buffer: %s\n", err.Error())
						continue
					}
					if str[0:10] == ALIAS_RESPONSE {
						info := Peer{ID: p.ID.String(), Alias: str[10:]}
						runtime.EventsEmit(ctx, "data_backend", info)
						// Close stream
						stream.Reset()
					}
				}
			}(p, stream)
		}
	}()

	/* go func(ctx context.Context, h host.Host) {
		for peer := range peerChan {
			runtime.LogInfo(ctx, "Found peer: "+peer.ID.Pretty()+", connecting")
			if err := h.Connect(ctx, peer); err != nil {
				runtime.LogErrorf(ctx, "Connection failed: %s\n", err.Error())
			}

			stream, err := h.NewStream(ctx, peer.ID, protocol.ID(PROTOCOL))
			if err != nil {
				runtime.LogErrorf(ctx, "Stream open failed: %s\n", err.Error())
			} else {
				runtime.EventsEmit(ctx, "data_backend", peer.ID)
				hdler := NewHandler(ctx)
				rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

				go hdler.readData(rw)
				runtime.EventsOn(ctx, "data_frontend", func(optionalData ...interface{}) {
					hdler.writeData(rw, optionalData)
				})
			}
		}
	}(ctx, h) */
}
