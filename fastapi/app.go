package fastapi

import (
	"log"
	"net"
	"runtime/debug"
	"time"

	"git.tube/funny/link"
	"git.tube/funny/pprof"
	"git.tube/funny/slab"
)

type Handler interface {
	InitSession(*link.Session) error
	Transaction(*link.Session, Message, func())
	ErrorCallback(error)
}

type App struct {
	serviceTypes []*ServiceType
	services     [256]Provider
	timeRecoder  *pprof.TimeRecorder

	Pool         slab.Pool
	ReadBufSize  int
	SendChanSize int
	MaxRecvSize  int
	MaxSendSize  int
	RecvTimeout  time.Duration
	SendTimeout  time.Duration
}

func New() *App {
	return &App{
		timeRecoder:  pprof.NewTimeRecorder(),
		Pool:         &slab.NoPool{},
		ReadBufSize:  1024,
		SendChanSize: 1024,
		MaxRecvSize:  64 * 1024,
		MaxSendSize:  64 * 1024,
	}
}

func (app *App) Listen(network, address string, handler Handler) (*link.Server, error) {
	listener, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}
	return app.NewServer(listener, handler), nil
}

func (app *App) NewServer(listener net.Listener, handler Handler) *link.Server {
	if handler == nil {
		handler = &noHandler{}
	}
	return link.NewServer(
		listener, link.ProtocolFunc(app.newServerCodec), app.SendChanSize,
		link.HandlerFunc(func(session *link.Session) {
			app.handleSession(session, handler)
		}),
	)
}

type noHandler struct {
}

func (t *noHandler) InitSession(session *link.Session) error {
	return nil
}

func (t *noHandler) ErrorCallback(err error) {

}

func (t *noHandler) Transaction(session *link.Session, req Message, work func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("fastapi: unhandled panic when processing %s - %s", req.Identity(), err)
			log.Println(string(debug.Stack()))
		}
	}()
	work()
}
