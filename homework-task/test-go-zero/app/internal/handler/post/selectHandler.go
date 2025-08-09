package post

import (
	"net/http"

	"blog-api/internal/logic/post"
	"blog-api/internal/svc"
	"blog-api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func SelectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := post.NewSelectLogic(r.Context(), svcCtx)
		resp, err := l.Select(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
