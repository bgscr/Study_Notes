package logic

import (
	"context"
	"encoding/json"
	"strconv"

	"rpc/posts/internal/svc"
	"rpc/posts/posts"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/sync/singleflight"
)

type GetPostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPostLogic {
	return &GetPostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

var sf singleflight.Group

type postCache struct {
	Id            uint64 `json:"id"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	UserId        uint64 `json:"user_id"`
	CommentStatus string `json:"comment_status"`
}

func (l *GetPostLogic) GetPost(in *posts.SinglePostInfoReq) (*posts.SinglePostInfoResp, error) {
	// todo: add your logic here and delete this line

	key := "GetPost_" + strconv.Itoa(int(in.Id))
	ctx := context.Background()
	postStr, err := l.svcCtx.RedisClient.GetCtx(ctx, key)
	if err == nil && postStr != "" {
		var post posts.SinglePostInfoResp
		if jsonErr := json.Unmarshal([]byte(postStr), &post); jsonErr == nil {
			return &post, nil
		}
	}

	v, err, _ := sf.Do(key, func() (interface{}, error) {
		// get data from DB
		data := posts.SinglePostInfoResp{
			Id:            in.Id,
			Title:         "测试",
			Content:       "测试content",
			UserId:        33,
			CommentStatus: "aaa",
		}
		// 转成无锁结构体
		cacheData := postCache{
			Id:            data.Id,
			Title:         data.Title,
			Content:       data.Content,
			UserId:        data.UserId,
			CommentStatus: data.CommentStatus,
		}
		if bytes, err := json.Marshal(cacheData); err == nil {

			setNxFlag, err := l.svcCtx.RedisClient.SetnxCtx(ctx, key, string(bytes))
			if !setNxFlag {
				return nil, err
			}
			l.svcCtx.RedisClient.ExpireCtx(ctx, key, 3600)
		}
		return &data, nil
	})

	if err != nil {
		return nil, err
	}

	return v.(*posts.SinglePostInfoResp), nil
}
