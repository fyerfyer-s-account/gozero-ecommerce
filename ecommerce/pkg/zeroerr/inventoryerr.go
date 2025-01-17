package zeroerr

// Common Stock Errors
var (
    ErrStockNotFound      = NewCodeError(105001, "库存记录不存在")
    ErrStockCreateFailed  = NewCodeError(105002, "创建库存记录失败")
    ErrStockUpdateFailed  = NewCodeError(105003, "更新库存记录失败")
    ErrInsufficientStock  = NewCodeError(105004, "库存不足")
    ErrStockAlreadyExists = NewCodeError(105005, "库存记录已存在")
)

// Warehouse Errors
var (
    ErrWarehouseNotFound     = NewCodeError(105101, "仓库不存在")
    ErrWarehouseCreateFailed = NewCodeError(105102, "创建仓库失败")
    ErrWarehouseUpdateFailed = NewCodeError(105103, "更新仓库失败")
    ErrWarehouseDisabled     = NewCodeError(105104, "仓库已停用")
    ErrDuplicateWarehouse    = NewCodeError(105105, "仓库已存在")
)

// Stock Operation Errors
var (
    ErrStockInFailed    = NewCodeError(105201, "入库操作失败")
    ErrStockOutFailed   = NewCodeError(105202, "出库操作失败")
    ErrExceedMaxStock   = NewCodeError(105204, "超出最大库存限制")
    ErrNegativeStock    = NewCodeError(105205, "库存不能为负数")
)

// Stock Lock Errors
var (
    ErrStockLockFailed     = NewCodeError(105301, "库存锁定失败")
    ErrStockUnlockFailed   = NewCodeError(105302, "库存解锁失败")
    ErrStockDeductFailed   = NewCodeError(105303, "库存扣减失败")
    ErrLockNotFound        = NewCodeError(105304, "库存锁定记录不存在")
    ErrLockExpired         = NewCodeError(105305, "库存锁定已过期")
    ErrDuplicateLock       = NewCodeError(105306, "重复的库存锁定")
    ErrInsufficientLocked  = NewCodeError(105307, "锁定库存不足")
)

// Stock Record Errors
var (
    ErrRecordNotFound      = NewCodeError(105401, "库存记录不存在")
    ErrRecordCreateFailed  = NewCodeError(105402, "创建库存记录失败")
    ErrInvalidRecordType   = NewCodeError(105403, "无效的记录类型")
    ErrDuplicateRecord     = NewCodeError(105404, "重复的库存记录")
    ErrRecordQueryFailed   = NewCodeError(105405, "查询库存记录失败")
)