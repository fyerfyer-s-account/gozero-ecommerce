package logic

import (
	"context"
	"log"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCategoriesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCategoriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCategoriesLogic {
	return &GetCategoriesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCategoriesLogic) GetCategories(in *product.Empty) (*product.GetCategoriesResponse, error) {
    log.Println("Starting GetCategories")
    
    categories, err := l.svcCtx.CategoriesModel.FindAll(l.ctx)
    if err != nil {
        log.Printf("Failed to get categories from DB: %v", err)
        return nil, err
    }
    log.Printf("Found %d categories in DB", len(categories))

    var resp []*product.Category
    for _, c := range categories {
        l.Logger.Infof("Processing category: ID=%d, Name=%s", c.Id, c.Name)
        resp = append(resp, &product.Category{
            Id:       int64(c.Id),
            Name:     c.Name,
            ParentId: c.ParentId.Int64,
            Level:    c.Level,
            Sort:     c.Sort,
            Icon:     c.Icon.String,
            CreatedAt: c.CreatedAt.Unix(),
        })
    }

    log.Println("GetCategories completed successfully")
    return &product.GetCategoriesResponse{Categories: resp}, nil
}