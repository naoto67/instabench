package instabench

import "context"

type Executor interface {
	Setup(context.Context) (context.Context, error)
	ExecContext(context.Context) (context.Context, error)
}

type Preparer interface {
	PrepareContext(context.Context) (context.Context, error)
}
