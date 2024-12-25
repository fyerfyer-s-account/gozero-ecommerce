CREATE DATABASE IF NOT EXISTS `mall_cart`;
USE `mall_cart`;

-- 购物车商品表
CREATE TABLE `cart_items` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '购物车项ID',
    `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
    `product_id` bigint unsigned NOT NULL COMMENT '商品ID',
    `sku_id` bigint unsigned NOT NULL COMMENT 'SKU ID',
    `product_name` varchar(100) NOT NULL COMMENT '商品名称',
    `sku_name` varchar(100) NOT NULL COMMENT 'SKU名称',
    `image` varchar(255) DEFAULT NULL COMMENT '商品图片',
    `price` decimal(10,2) NOT NULL COMMENT '商品单价',
    `quantity` int NOT NULL DEFAULT '1' COMMENT '商品数量',
    `selected` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否选中',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_user_sku` (`user_id`, `sku_id`),
    KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='购物车商品表';

-- 购物车统计表（可选，用于存储购物车汇总信息）
CREATE TABLE `cart_statistics` (
    `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
    `total_quantity` int NOT NULL DEFAULT '0' COMMENT '商品总数量',
    `selected_quantity` int NOT NULL DEFAULT '0' COMMENT '已选商品数量',
    `total_amount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '商品总金额',
    `selected_amount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '已选商品金额',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='购物车统计表';