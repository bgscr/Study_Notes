package login

import (
	"context"

	"testGoZero/app/internal/svc"
	"testGoZero/app/internal/types"
	"testGoZero/common/jwt"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginUserInfoReq) (resp *types.LoginUserInfoResp, err error) { // todo: add your logic here and delete this line

	token, _ := jwt.GenerateToken(l.svcCtx.Config.Auth.AccessSecret,
		l.svcCtx.Config.Auth.AccessExpire, 1)
	return &types.LoginUserInfoResp{
		Flag:  true,
		Msg:   "登陆成功",
		Token: "Bearer " + token,
	}, nil
	return
}
