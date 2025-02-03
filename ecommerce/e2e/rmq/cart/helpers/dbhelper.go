package helpers

import (
    "context"
    "fmt"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/model"
    "github.com/zeromicro/go-zero/core/stores/cache"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

type DBHelper struct {
    conn            sqlx.SqlConn
    cacheRedis      cache.ClusterConf
    cartItemsModel  model.CartItemsModel
    cartStatsModel  model.CartStatisticsModel
}

// NewDBHelper creates a new database helper with the given configuration
func NewDBHelper(mysqlDSN string, redisConf cache.ClusterConf) (*DBHelper, error) {
    conn := sqlx.NewMysql(mysqlDSN)

    return &DBHelper{
        conn:            conn,
        cacheRedis:      redisConf,
        cartItemsModel:  model.NewCartItemsModel(conn, redisConf),
        cartStatsModel:  model.NewCartStatisticsModel(conn, redisConf),
    }, nil
}

// Close closes the database connections
func (h *DBHelper) Close() error {
    return nil // sqlx doesn't provide a Close method
}

// GetCartItemsModel returns the cart items model
func (h *DBHelper) GetCartItemsModel() model.CartItemsModel {
    return h.cartItemsModel
}

// GetCartStatsModel returns the cart statistics model
func (h *DBHelper) GetCartStatsModel() model.CartStatisticsModel {
    return h.cartStatsModel
}

// CleanTestData removes test data from all relevant tables
func (h *DBHelper) CleanTestData(ctx context.Context) error {
    // Clean cart items with test data pattern
    if _, err := h.conn.ExecCtx(ctx, "DELETE FROM cart_items WHERE product_name LIKE 'Test_%'"); err != nil {
        return fmt.Errorf("failed to clean table cart_items: %w", err)
    }

    // Clean cart statistics for test users
    if _, err := h.conn.ExecCtx(ctx, "DELETE FROM cart_statistics WHERE user_id IN (SELECT DISTINCT user_id FROM cart_items WHERE product_name LIKE 'Test_%')"); err != nil {
        return fmt.Errorf("failed to clean table cart_statistics: %w", err)
    }

    return nil
}

// PrepareTestCartItem creates a test cart item
func (h *DBHelper) PrepareTestCartItem(ctx context.Context, userId uint64, productId uint64, skuId uint64, quantity int64) (*model.CartItems, error) {
    item := &model.CartItems{
        UserId:      userId,
        ProductId:   productId,
        SkuId:      skuId,
        ProductName: fmt.Sprintf("Test_Product_%d", productId),
        SkuName:    fmt.Sprintf("Test_SKU_%d", skuId),
        Price:      100.00,
        Quantity:   quantity,
        Selected:   1,
    }

    result, err := h.cartItemsModel.Insert(ctx, item)
    if err != nil {
        return nil, fmt.Errorf("failed to insert test cart item: %w", err)
    }

    id, err := result.LastInsertId()
    if err != nil {
        return nil, fmt.Errorf("failed to get inserted cart item ID: %w", err)
    }

    item.Id = uint64(id)
    return item, nil
}

// PrepareTestCartStatistics creates or updates test cart statistics
func (h *DBHelper) PrepareTestCartStatistics(ctx context.Context, userId uint64, total, selected int64, totalAmount, selectedAmount float64) error {
    stats := &model.CartStatistics{
        UserId:           userId,
        TotalQuantity:    total,
        SelectedQuantity: selected,
        TotalAmount:      totalAmount,
        SelectedAmount:   selectedAmount,
    }

    return h.cartStatsModel.Upsert(ctx, stats)
}

// GetCartItemsByUserId retrieves all cart items for a user
func (h *DBHelper) GetCartItemsByUserId(ctx context.Context, userId uint64) ([]*model.CartItems, error) {
    return h.cartItemsModel.FindByUserId(ctx, userId)
}

// GetCartStatistics retrieves cart statistics for a user
func (h *DBHelper) GetCartStatistics(ctx context.Context, userId uint64) (*model.CartStatistics, error) {
    return h.cartStatsModel.FindOne(ctx, userId)
}

// ClearUserCart removes all items from a user's cart
func (h *DBHelper) ClearUserCart(ctx context.Context, userId uint64) error {
    if err := h.cartItemsModel.DeleteByUserId(ctx, userId); err != nil {
        return fmt.Errorf("failed to clear cart items: %w", err)
    }

    stats := &model.CartStatistics{
        UserId:           userId,
        TotalQuantity:    0,
        SelectedQuantity: 0,
        TotalAmount:      0,
        SelectedAmount:   0,
    }

    if err := h.cartStatsModel.Upsert(ctx, stats); err != nil {
        return fmt.Errorf("failed to reset cart statistics: %w", err)
    }

    return nil
}