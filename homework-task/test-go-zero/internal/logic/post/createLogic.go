package post

import (
	"context"
	"encoding/json"
	"fmt"

	"testGoZero/common/enum"
	"testGoZero/internal/svc"
	"testGoZero/internal/types"
	"testGoZero/model/posts"

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
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()

	post := posts.Posts{
		Title:         req.Title,
		Content:       req.Content,
		UserId:        uint64(uid),
		CommentStatus: string(enum.NoComments),
	}
	if dbError := l.svcCtx.GormDB.Create(&post).Error; dbError != nil {
		return &types.PostInfoResp{
			Flag: false,
			Msg:  fmt.Sprintf("数据库操作失败:%v", dbError.Error()),
		}, nil
	}

	return &types.PostInfoResp{
		Flag: true,
		Msg:  "创建成功",
	}, nil
}
