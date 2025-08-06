package logic

import (
	"context"

	"testGoZero/users/internal/svc"
	"testGoZero/users/users"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *users.RegisterUserInfoReq) (*users.RegisterUserInfoResp, error) {
	// todo: add your logic here and delete this line

	return &users.RegisterUserInfoResp{}, nil
}
