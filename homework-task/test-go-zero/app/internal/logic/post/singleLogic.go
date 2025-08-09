package post

import (
	"context"

	"testGoZero/app/internal/svc"
	"testGoZero/app/internal/types"
	"testGoZero/posts/posts"

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
	// todo: add your logic here and delete this line

	rpcResp, err := l.svcCtx.PostRPC.GetPost(l.ctx, &posts.SinglePostInfoReq{
		Id: req.Id,
	})

	return &types.SinglePostInfoResp{
		Id:            rpcResp.Id,
		Title:         rpcResp.Title,
		Content:       rpcResp.Content,
		UserId:        rpcResp.UserId,
		CommentStatus: rpcResp.CommentStatus,
	}, err
}
