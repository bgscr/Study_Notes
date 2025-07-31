package login

import (
	"context"
	"fmt"
	"strings"

	"testGoZero/common/cryptx"
	"testGoZero/internal/svc"
	"testGoZero/internal/types"
	"testGoZero/model/users"

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
	defer func() {
		if r := recover(); r != nil {
			// 修改命名返回值
			resp = &types.RegisterUserInfoResp{
				Flag: false,
				Msg:  "新增用户数据异常",
			}
			// 将 panic 转换为 error
			switch x := r.(type) {
			case error:
				err = x
			default:
				err = fmt.Errorf("%v", x)
			}
		}
	}()

	user := users.Users{
		Username: strings.TrimSpace(req.UsernName),
		Password: cryptx.PasswordEncrypt(l.svcCtx.Config.Salt, strings.TrimSpace(req.Password)),
		Email:    strings.TrimSpace(req.Email),
	}

	if dbError := l.svcCtx.GormDB.Create(&user).Error; dbError != nil {
		return &types.RegisterUserInfoResp{
			Flag: false,
			Msg:  fmt.Sprintf("数据库操作失败:%v", dbError.Error()),
		}, nil
	}

	return &types.RegisterUserInfoResp{
		Flag: true,
		Msg:  "注册成功",
	}, nil
}
