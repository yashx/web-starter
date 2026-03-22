package ws

import (
	"context"
	"errors"
	"net/http"
	"web-starter/foundation"
	"web-starter/foundation/appError"
	"web-starter/foundation/httpHelper"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type SubRouter struct {
	app *foundation.App
}

func NewSubRouter(app *foundation.App) *SubRouter {
	return &SubRouter{app: app}
}

func (sR SubRouter) BuildHandler() (string, http.Handler) {
	router := chi.NewRouter()
	router.Get("/", sR.getWebSocketHandlerFn)
	return "/ws", router
}

func (sR SubRouter) getWebSocketHandlerFn(writer http.ResponseWriter, httpRequest *http.Request) {
	conn, err := websocket.Accept(writer, httpRequest, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		httpHelper.JsonErrorResponse(sR.app, writer, httpRequest, appError.InternalServerErrorWithCause(err))
		return
	}
	defer conn.CloseNow()

	ctx := httpRequest.Context()
	msgCh := make(chan map[string]any, 16)

	readerErrCh := sR.startReader(ctx, conn, msgCh)
	writerErrCh := sR.startWriter(ctx, conn, msgCh)

	sR.manage(ctx, conn, readerErrCh, writerErrCh)
}

func (sR SubRouter) startReader(ctx context.Context, conn *websocket.Conn, msgCh chan<- map[string]any) <-chan error {
	errCh := make(chan error, 1)
	go func() {
		defer close(msgCh)
		readCtx := context.WithoutCancel(ctx)
		for {
			var msg map[string]any
			if err := wsjson.Read(readCtx, conn, &msg); err != nil {
				errCh <- err
				return
			}
			sR.app.Logger.Info("received message", zap.Any("message", msg))
			msgCh <- msg
		}
	}()
	return errCh
}

func (sR SubRouter) startWriter(ctx context.Context, conn *websocket.Conn, msgCh <-chan map[string]any) <-chan error {
	errCh := make(chan error, 1)
	go func() {
		writeCtx := context.WithoutCancel(ctx)
		for msg := range msgCh {
			if err := wsjson.Write(writeCtx, conn, msg); err != nil {
				errCh <- err
				return
			}
		}
		errCh <- nil
	}()
	return errCh
}

func (sR SubRouter) manage(ctx context.Context, conn *websocket.Conn, readerErrCh, writerErrCh <-chan error) {
	select {
	case <-ctx.Done():
		sR.closeOnShutdown(conn)
		<-readerErrCh
		<-writerErrCh
	case err := <-readerErrCh:
		sR.handleReadError(ctx, conn, err)
		<-writerErrCh
	case err := <-writerErrCh:
		sR.handleWriteError(ctx, conn, err)
		<-readerErrCh
	}
}

func (sR SubRouter) closeOnShutdown(conn *websocket.Conn) {
	sR.app.Logger.Info("context cancelled, closing websocket connection")
	if err := conn.Close(websocket.StatusGoingAway, "server shutting down"); err != nil {
		sR.app.Logger.Warn("error closing websocket on shutdown", zap.Error(err))
	}
}

func (sR SubRouter) handleReadError(ctx context.Context, conn *websocket.Conn, err error) {
	if closeErr, ok := errors.AsType[websocket.CloseError](err); ok {
		sR.app.Logger.Info("websocket closed by client",
			zap.Int("code", int(closeErr.Code)),
			zap.String("reason", closeErr.Reason),
		)
		return
	}
	if ctx.Err() != nil {
		return
	}
	sR.app.Logger.Warn("error reading message", zap.Error(err))
	if err := conn.Close(websocket.StatusNormalClosure, ""); err != nil {
		sR.app.Logger.Warn("error closing websocket connection", zap.Error(err))
	}
}

func (sR SubRouter) handleWriteError(ctx context.Context, conn *websocket.Conn, err error) {
	if err == nil || ctx.Err() != nil {
		return
	}
	sR.app.Logger.Warn("error writing message", zap.Error(err))
	if err := conn.Close(websocket.StatusNormalClosure, ""); err != nil {
		sR.app.Logger.Warn("error closing websocket connection", zap.Error(err))
	}
}
