package bootstrap

import "context"

func (a *App) initWorkerContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	a.workerCtx = ctx
	a.workerCancel = cancel

	return a.workerCtx, a.workerCancel
}
