package logic

import (
	"context"

	"rpc/posts/internal/svc"
	"rpc/posts/posts"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPostsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListPostsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPostsLogic {
	return &ListPostsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListPostsLogic) ListPosts(in *posts.PageReq) (*posts.SelectPostInfoResp, error) {
	// todo: add your logic here and delete this line

	return &posts.SelectPostInfoResp{}, nil
}
