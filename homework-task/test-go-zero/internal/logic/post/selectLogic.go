package post

import (
	"context"
	"math"

	"testGoZero/internal/svc"
	"testGoZero/internal/types"
	"testGoZero/model/posts"

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
	var total int64
	var list []posts.Posts

	dbError := l.svcCtx.GormDB.Model(&posts.Posts{}).Count(&total).Error
	if dbError != nil {
		return nil, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	err = l.svcCtx.GormDB.Offset(int(offset)).Limit(int(req.PageSize)).Find(&list).Error

	var data []types.SinglePostInfoResp
	for _, v := range list {
		data = append(data, types.SinglePostInfoResp{
			Id:            v.Id,
			Title:         v.Title,
			Content:       v.Content,
			UserId:        v.UserId,
			CommentStatus: v.CommentStatus,
		})
	}
	return &types.SelectPostInfoResp{
		PageResp: types.PageResp{
			CurrentPage: req.Page,
			PageSize:    req.PageSize,
			Total:       uint64(total),
			TotalPages:  uint64(math.Ceil(float64(total) / float64(req.PageSize))),
		},
		Data: data,
	}, err
}
