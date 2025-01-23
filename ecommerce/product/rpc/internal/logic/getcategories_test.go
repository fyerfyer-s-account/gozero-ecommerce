package logic

import (
    "context"
    "database/sql"
    "flag"
    "testing"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/model"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"
    "github.com/stretchr/testify/assert"
    "github.com/zeromicro/go-zero/core/conf"
)

func TestGetCategoriesLogic_GetCategories(t *testing.T) {
    t.Log("Starting GetCategories test")
    
    configFile := flag.String("f", "../../etc/product.yaml", "config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    ctx := svc.NewServiceContext(c)
    
    t.Log("Creating test categories")
    testCategories := []*model.Categories{
        {
            Name:     "Category 1",
            ParentId: 0,
            Level:    1,
            Sort:     1,
            Icon:     sql.NullString{String: "icon1.jpg", Valid: true},
        },
        {
            Name:     "Category 2",
            ParentId: 0,
            Level:    1,
            Sort:     2,
            Icon:     sql.NullString{String: "icon2.jpg", Valid: true},
        },
    }

    var categoryIds []uint64
    for i, category := range testCategories {
        t.Logf("Inserting test category %d: %s", i+1, category.Name)
        result, err := ctx.CategoriesModel.Insert(context.Background(), category)
        if err != nil {
            t.Fatalf("Failed to insert category %d: %v", i+1, err)
        }
        id, err := result.LastInsertId()
        if err != nil {
            t.Fatalf("Failed to get last insert id for category %d: %v", i+1, err)
        }
        categoryIds = append(categoryIds, uint64(id))
        t.Logf("Created category with ID: %d", id)
    }

    defer func() {
        t.Log("Cleaning up test categories")
        for _, id := range categoryIds {
            if err := ctx.CategoriesModel.Delete(context.Background(), id); err != nil {
                t.Errorf("Failed to delete test category %d: %v", id, err)
            }
        }
    }()

    t.Run("Get all categories", func(t *testing.T) {
        t.Log("Testing GetCategories")
        logic := NewGetCategoriesLogic(context.Background(), ctx)
        
        t.Log("Calling GetCategories")
        resp, err := logic.GetCategories(&product.Empty{})
        
        if err != nil {
            t.Logf("GetCategories error: %v", err)
        }
        if resp != nil {
            t.Logf("GetCategories response: %+v", resp)
        }
        
        assert.NoError(t, err)
        assert.NotNil(t, resp)
        assert.Equal(t, 2, len(resp.Categories))
    })
}