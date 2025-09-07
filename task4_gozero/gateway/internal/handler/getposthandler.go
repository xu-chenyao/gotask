package handler

import (
	"net/http"
	"strconv"

	"gotask/task4_gozero/gateway/internal/logic"
	"gotask/task4_gozero/gateway/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetPostHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var pathVars map[string]string
		if err := httpx.ParsePath(r, &pathVars); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		idStr := pathVars[":id"]
		id, _ := strconv.ParseInt(idStr, 10, 64)

		l := logic.NewGetPostLogic(r.Context(), svcCtx)
		resp, err := l.GetPost(id)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
