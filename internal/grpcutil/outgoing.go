package grpcutil

import (
	"context"
	"go.uber.org/fx"
	"google.golang.org/grpc/metadata"
	"strconv"
)

type ctxKey string

const CtxUserID ctxKey = "uid"
const CtxRoles ctxKey = "roles"

func WithUser(ctx context.Context, uid int64, roles []string) context.Context {
	ctx = context.WithValue(ctx, CtxUserID, uid)
	ctx = context.WithValue(ctx, CtxRoles, roles)
	return ctx
}

func OutgoingWithUser(ctx context.Context) context.Context {
	v := ctx.Value(CtxUserID)
	uid, _ := v.(int64)
	md := metadata.Pairs("x-user-id", strconv.FormatInt(uid, 10))
	return metadata.NewOutgoingContext(ctx, md)
}

var Module = fx.Provide(func() struct{} { return struct{}{} })
