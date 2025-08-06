package post

import (
	"context"

	"testGoZero/app/internal/svc"
	"testGoZero/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SelectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSelectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SelectLogic {
	return &SelectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SelectLogic) Select(req *types.PageReq) (resp *types.SelectPostInfoResp, err error) {
	// todo: add your logic here and delete this line

	return
}
