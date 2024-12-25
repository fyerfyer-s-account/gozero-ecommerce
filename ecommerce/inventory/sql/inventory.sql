CREATE DATABASE IF NOT EXISTS `mall_inventory`;
USE `mall_inventory`;

-- 库存表
CREATE TABLE `stocks` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '库存ID',
    `sku_id` bigint unsigned NOT NULL COMMENT 'SKU ID',
    `warehouse_id` bigint unsigned NOT NULL COMMENT '仓库ID',
    `available` int NOT NULL DEFAULT '0' COMMENT '可用库存',
    `locked` int NOT NULL DEFAULT '0' COMMENT '锁定库存',
    `total` int NOT NULL DEFAULT '0' COMMENT '总库存',
    `alert_quantity` int NOT NULL DEFAULT '0' COMMENT '库存预警数量',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_sku_warehouse` (`sku_id`, `warehouse_id`),
    KEY `idx_sku_id` (`sku_id`),
    KEY `idx_warehouse_id` (`warehouse_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='库存表';

-- 仓库表
CREATE TABLE `warehouses` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '仓库ID',
    `name` varchar(100) NOT NULL COMMENT '仓库名称',
    `address` varchar(255) NOT NULL COMMENT '仓库地址',
    `contact` varchar(50) DEFAULT NULL COMMENT '联系人',
    `phone` varchar(20) DEFAULT NULL COMMENT '联系电话',
    `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态 1:正常 2:停用',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='仓库表';

-- 库存记录表
CREATE TABLE `stock_records` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '记录ID',
    `sku_id` bigint unsigned NOT NULL COMMENT 'SKU ID',
    `warehouse_id` bigint unsigned NOT NULL COMMENT '仓库ID',
    `type` tinyint NOT NULL COMMENT '类型 1:入库 2:出库 3:锁定 4:解锁',
    `quantity` int NOT NULL COMMENT '数量',
    `order_no` varchar(64) DEFAULT NULL COMMENT '订单号',
    `remark` varchar(255) DEFAULT NULL COMMENT '备注',
    `operator` varchar(50) DEFAULT NULL COMMENT '操作人',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_sku_id` (`sku_id`),
    KEY `idx_warehouse_id` (`warehouse_id`),
    KEY `idx_order_no` (`order_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='库存记录表';

-- 库存锁定表
CREATE TABLE `stock_locks` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '锁定ID',
    `order_no` varchar(64) NOT NULL COMMENT '订单号',
    `sku_id` bigint unsigned NOT NULL COMMENT 'SKU ID',
    `warehouse_id` bigint unsigned NOT NULL COMMENT '仓库ID',
    `quantity` int NOT NULL COMMENT '锁定数量',
    `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态 1:锁定 2:已解锁 3:已扣减',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_order_sku` (`order_no`, `sku_id`),
    KEY `idx_sku_id` (`sku_id`),
    KEY `idx_warehouse_id` (`warehouse_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='库存锁定表';