package logic

import (
	"context"

	"testGoZero/posts/internal/svc"
	"testGoZero/posts/posts"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePostLogic {
	return &DeletePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeletePostLogic) DeletePost(in *posts.PostInfoReq) (*posts.PostInfoResp, error) {
	// todo: add your logic here and delete this line

	return &posts.PostInfoResp{}, nil
}
