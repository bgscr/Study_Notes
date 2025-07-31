package post

import (
	"context"

	"testGoZero/internal/svc"
	"testGoZero/internal/types"
	"testGoZero/model/posts"

	"github.com/zeromicro/go-zero/core/logx"
)

type SingleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSingleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SingleLogic {
	return &SingleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SingleLogic) Single(req *types.SinglePostInfoReq) (resp *types.SinglePostInfoResp, err error) {
	var post posts.Posts
	l.svcCtx.GormDB.Where(&posts.Posts{Id: req.Id}).First(&post)

	return &types.SinglePostInfoResp{
		Id:            post.Id,
		Title:         post.Title,
		Content:       post.Content,
		UserId:        uint64(post.UserId),
		CommentStatus: post.CommentStatus,
	}, nil
}
