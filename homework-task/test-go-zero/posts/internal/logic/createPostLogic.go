package logic

import (
	"context"

	"testGoZero/posts/internal/svc"
	"testGoZero/posts/posts"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePostLogic {
	return &CreatePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreatePostLogic) CreatePost(in *posts.PostInfoReq) (*posts.PostInfoResp, error) {
	// todo: add your logic here and delete this line

	return &posts.PostInfoResp{}, nil
}
