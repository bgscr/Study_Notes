package post

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"testGoZero/app/internal/logic/post"
	"testGoZero/app/internal/svc"
	"testGoZero/app/internal/types"
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
