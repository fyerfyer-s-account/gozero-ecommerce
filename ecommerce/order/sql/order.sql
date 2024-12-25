CREATE DATABASE IF NOT EXISTS `mall_order`;
USE `mall_order`;

-- 订单主表
CREATE TABLE `orders` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '订单ID',
    `order_no` varchar(64) NOT NULL COMMENT '订单编号',
    `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
    `total_amount` decimal(10,2) NOT NULL COMMENT '订单总金额',
    `pay_amount` decimal(10,2) NOT NULL COMMENT '应付金额',
    `freight_amount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '运费',
    `discount_amount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '优惠金额',
    `status` tinyint NOT NULL DEFAULT '1' COMMENT '订单状态 1:待支付 2:待发货 3:待收货 4:已完成 5:已取消 6:售后中',
    `address` varchar(255) NOT NULL COMMENT '收货地址',
    `receiver` varchar(32) NOT NULL COMMENT '收货人',
    `phone` varchar(20) NOT NULL COMMENT '联系电话',
    `remark` varchar(500) DEFAULT NULL COMMENT '订单备注',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_order_no` (`order_no`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单表';

-- 订单项表
CREATE TABLE `order_items` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '订单项ID',
    `order_id` bigint unsigned NOT NULL COMMENT '订单ID',
    `product_id` bigint unsigned NOT NULL COMMENT '商品ID',
    `sku_id` bigint unsigned NOT NULL COMMENT 'SKU ID',
    `product_name` varchar(100) NOT NULL COMMENT '商品名称',
    `sku_name` varchar(100) NOT NULL COMMENT 'SKU名称',
    `price` decimal(10,2) NOT NULL COMMENT '商品单价',
    `quantity` int NOT NULL COMMENT '购买数量',
    `total_amount` decimal(10,2) NOT NULL COMMENT '总金额',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_order_id` (`order_id`),
    KEY `idx_product_id` (`product_id`),
    KEY `idx_sku_id` (`sku_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单项表';

-- 支付信息表
CREATE TABLE `order_payments` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '支付ID',
    `order_id` bigint unsigned NOT NULL COMMENT '订单ID',
    `payment_no` varchar(64) NOT NULL COMMENT '支付流水号',
    `payment_method` tinyint NOT NULL COMMENT '支付方式 1:微信 2:支付宝 3:余额',
    `amount` decimal(10,2) NOT NULL COMMENT '支付金额',
    `status` tinyint NOT NULL DEFAULT '0' COMMENT '支付状态 0:未支付 1:已支付 2:已退款',
    `pay_time` timestamp NULL DEFAULT NULL COMMENT '支付时间',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_payment_no` (`payment_no`),
    KEY `idx_order_id` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='支付信息表';

-- 物流信息表
CREATE TABLE `order_shipping` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '物流ID',
    `order_id` bigint unsigned NOT NULL COMMENT '订单ID',
    `shipping_no` varchar(64) DEFAULT NULL COMMENT '物流单号',
    `company` varchar(64) DEFAULT NULL COMMENT '物流公司',
    `status` tinyint NOT NULL DEFAULT '0' COMMENT '物流状态 0:待发货 1:已发货 2:已签收',
    `ship_time` timestamp NULL DEFAULT NULL COMMENT '发货时间',
    `receive_time` timestamp NULL DEFAULT NULL COMMENT '签收时间',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_shipping_no` (`shipping_no`),
    KEY `idx_order_id` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='物流信息表';

-- 退款信息表
CREATE TABLE `order_refunds` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '退款ID',
    `order_id` bigint unsigned NOT NULL COMMENT '订单ID',
    `refund_no` varchar(64) NOT NULL COMMENT '退款编号',
    `amount` decimal(10,2) NOT NULL COMMENT '退款金额',
    `reason` varchar(500) NOT NULL COMMENT '退款原因',
    `status` tinyint NOT NULL DEFAULT '0' COMMENT '退款状态 0:待处理 1:已同意 2:已拒绝 3:已退款',
    `description` text COMMENT '问题描述',
    `images` json DEFAULT NULL COMMENT '图片凭证',
    `reply` varchar(500) DEFAULT NULL COMMENT '处理回复',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_refund_no` (`refund_no`),
    KEY `idx_order_id` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='退款信息表';