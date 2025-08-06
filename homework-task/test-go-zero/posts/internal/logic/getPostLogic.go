package logic

import (
	"context"

	"testGoZero/posts/internal/svc"
	"testGoZero/posts/posts"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPostLogic {
	return &GetPostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPostLogic) GetPost(in *posts.SinglePostInfoReq) (*posts.SinglePostInfoResp, error) {
	// todo: add your logic here and delete this line

	return &posts.SinglePostInfoResp{
		Id:            1,
		Title:         "测试",
		Content:       "测试content",
		UserId:        33,
		CommentStatus: "aaa",
	}, nil
}
