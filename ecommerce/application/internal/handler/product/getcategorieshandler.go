package product

import (
	"net/http"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/logic/product"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetCategoriesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := product.NewGetCategoriesLogic(r.Context(), svcCtx)
		resp, err := l.GetCategories()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
