package logic

import (
	"context"

	"rpc/posts/internal/svc"
	"rpc/posts/posts"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdatePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePostLogic {
	return &UpdatePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdatePostLogic) UpdatePost(in *posts.PostInfoReq) (*posts.PostInfoResp, error) {
	// todo: add your logic here and delete this line

	return &posts.PostInfoResp{}, nil
}
