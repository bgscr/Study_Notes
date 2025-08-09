package post

import (
	"blog-api/internal/svc"
	"blog-api/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.PostInfoReq) (resp *types.PostInfoResp, err error) {
	// todo: add your logic here and delete this line

	return
}
