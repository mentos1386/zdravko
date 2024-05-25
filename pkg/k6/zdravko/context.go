package zdravko

import "context"

type zdravkoContextKey string

type Context struct {
	Target Target
}

func WithZdravkoContext(ctx context.Context, zdravkoContext Context) context.Context {
	return context.WithValue(ctx, zdravkoContextKey("zdravko-ctx"), zdravkoContext)
}

func GetZdravkoContext(ctx context.Context) Context {
	return ctx.Value(zdravkoContextKey("zdravko-ctx")).(Context)
}
