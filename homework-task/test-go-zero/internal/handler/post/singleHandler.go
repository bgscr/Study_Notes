package post

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"testGoZero/internal/logic/post"
	"testGoZero/internal/svc"
	"testGoZero/internal/types"
)

func SingleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SinglePostInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := post.NewSingleLogic(r.Context(), svcCtx)
		resp, err := l.Single(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
