package grpcutil

import (
	"context"
	"go.uber.org/fx"
	"google.golang.org/grpc/metadata"
	"strings"
)

type ctxKey string

const CtxUserID ctxKey = "uid"
const CtxRoles ctxKey = "roles"

func WithUser(ctx context.Context, uid string, roles []string) context.Context {
	ctx = context.WithValue(ctx, CtxUserID, uid)
	ctx = context.WithValue(ctx, CtxRoles, roles)
	return ctx
}

func OutgoingWithUser(ctx context.Context) context.Context {
	v := ctx.Value(CtxUserID)
	uid, _ := v.(string)
	uid = strings.TrimSpace(uid)
	md := metadata.Pairs("x-user-id", uid)
	return metadata.NewOutgoingContext(ctx, md)
}

var Module = fx.Provide(func() struct{} { return struct{}{} })
