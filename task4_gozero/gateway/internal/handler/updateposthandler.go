package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gotask/task4_gozero/gateway/internal/logic"
	"gotask/task4_gozero/gateway/internal/svc"
	"gotask/task4_gozero/gateway/internal/types"
)

func UpdatePostHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdatePostReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUpdatePostLogic(r.Context(), svcCtx)
		resp, err := l.UpdatePost(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
