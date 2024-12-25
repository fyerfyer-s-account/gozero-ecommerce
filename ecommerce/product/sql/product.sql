CREATE DATABASE IF NOT EXISTS `mall_product`;
USE `mall_product`;

-- 商品分类表
CREATE TABLE `categories` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '分类ID',
    `name` varchar(50) NOT NULL COMMENT '分类名称',
    `parent_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '父分类ID',
    `level` tinyint NOT NULL DEFAULT '1' COMMENT '层级',
    `sort` int NOT NULL DEFAULT '0' COMMENT '排序',
    `icon` varchar(255) DEFAULT NULL COMMENT '图标URL',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_parent_id` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品分类表';

-- 商品表
CREATE TABLE `products` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '商品ID',
    `name` varchar(100) NOT NULL COMMENT '商品名称',
    `description` text COMMENT '商品描述',
    `category_id` bigint unsigned NOT NULL COMMENT '分类ID',
    `brand` varchar(50) DEFAULT NULL COMMENT '品牌',
    `images` json DEFAULT NULL COMMENT '商品图片列表',
    `price` decimal(10,2) NOT NULL COMMENT '商品价格',
    `sales` int NOT NULL DEFAULT '0' COMMENT '销量',
    `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态 1:上架 2:下架',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_category_id` (`category_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品表';

-- SKU表
CREATE TABLE `skus` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'SKU ID',
    `product_id` bigint unsigned NOT NULL COMMENT '商品ID',
    `sku_code` varchar(64) NOT NULL COMMENT 'SKU编码',
    `attributes` json NOT NULL COMMENT 'SKU属性',
    `price` decimal(10,2) NOT NULL COMMENT 'SKU价格',
    `stock` int NOT NULL DEFAULT '0' COMMENT '库存',
    `sales` int NOT NULL DEFAULT '0' COMMENT '销量',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_sku_code` (`sku_code`),
    KEY `idx_product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品SKU表';

-- 商品评价表
CREATE TABLE `product_reviews` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '评价ID',
    `product_id` bigint unsigned NOT NULL COMMENT '商品ID',
    `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
    `order_id` bigint unsigned NOT NULL COMMENT '订单ID',
    `rating` tinyint NOT NULL DEFAULT '5' COMMENT '评分 1-5',
    `content` text COMMENT '评价内容',
    `images` json DEFAULT NULL COMMENT '评价图片',
    `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态 0:待审核 1:已通过 2:已拒绝',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_product_id` (`product_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_order_id` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品评价表';