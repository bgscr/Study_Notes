package login

import (
	"context"

	"testGoZero/app/internal/svc"
	"testGoZero/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterUserInfoReq) (resp *types.RegisterUserInfoResp, err error) {
	// todo: add your logic here and delete this line

	return
}
