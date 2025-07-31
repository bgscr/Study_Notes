package login

import (
	"context"
	"strings"

	"testGoZero/common/cryptx"
	"testGoZero/common/jwt"
	"testGoZero/internal/svc"
	"testGoZero/internal/types"
	"testGoZero/model/users"

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

func (l *LoginLogic) Login(req *types.LoginUserInfoReq) (resp *types.LoginUserInfoResp, err error) {
	var user users.Users
	l.svcCtx.GormDB.Where(&users.Users{Username: req.UsernName}).First(&user)
	if user.Id == 0 {
		return &types.LoginUserInfoResp{
			Flag: false,
			Msg:  "用户不存在",
		}, nil
	}
	if cryptx.PasswordEncrypt(l.svcCtx.Config.Salt, strings.TrimSpace(req.Password)) != user.Password {

		return &types.LoginUserInfoResp{
			Flag: false,
			Msg:  "密码不正确",
		}, nil
	}

	token, _ := jwt.GenerateToken(l.svcCtx.Config.Auth.AccessSecret,
		l.svcCtx.Config.Auth.AccessExpire, user.Id)
	return &types.LoginUserInfoResp{
		Flag:  true,
		Msg:   "登陆成功",
		Token: "Bearer " + token,
	}, nil
}
