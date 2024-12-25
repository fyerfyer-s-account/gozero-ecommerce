CREATE DATABASE IF NOT EXISTS `mall_payment`;
USE `mall_payment`;

-- 支付订单表
CREATE TABLE `payment_orders` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `payment_no` varchar(64) NOT NULL COMMENT '支付单号',
    `order_no` varchar(64) NOT NULL COMMENT '订单号',
    `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
    `amount` decimal(10,2) NOT NULL COMMENT '支付金额',
    `channel` tinyint NOT NULL COMMENT '支付渠道 1:微信 2:支付宝 3:余额',
    `channel_data` json DEFAULT NULL COMMENT '支付渠道数据',
    `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态 1:待支付 2:支付中 3:已支付 4:已退款 5:已关闭',
    `notify_url` varchar(255) DEFAULT NULL COMMENT '回调地址',
    `return_url` varchar(255) DEFAULT NULL COMMENT '返回地址',
    `expire_time` timestamp NULL DEFAULT NULL COMMENT '过期时间',
    `pay_time` timestamp NULL DEFAULT NULL COMMENT '支付时间',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_payment_no` (`payment_no`),
    KEY `idx_order_no` (`order_no`),
    KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='支付订单表';

-- 退款订单表
CREATE TABLE `refund_orders` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `refund_no` varchar(64) NOT NULL COMMENT '退款单号',
    `payment_no` varchar(64) NOT NULL COMMENT '支付单号',
    `order_no` varchar(64) NOT NULL COMMENT '订单号',
    `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
    `amount` decimal(10,2) NOT NULL COMMENT '退款金额',
    `reason` varchar(255) NOT NULL COMMENT '退款原因',
    `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态 1:待处理 2:处理中 3:已退款 4:退款失败',
    `channel_data` json DEFAULT NULL COMMENT '退款渠道数据',
    `notify_url` varchar(255) DEFAULT NULL COMMENT '回调地址',
    `refund_time` timestamp NULL DEFAULT NULL COMMENT '退款时间',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_refund_no` (`refund_no`),
    KEY `idx_payment_no` (`payment_no`),
    KEY `idx_order_no` (`order_no`),
    KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='退款订单表';

-- 支付渠道表
CREATE TABLE `payment_channels` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `name` varchar(50) NOT NULL COMMENT '渠道名称',
    `channel` tinyint NOT NULL COMMENT '渠道类型 1:微信 2:支付宝',
    `config` json NOT NULL COMMENT '渠道配置',
    `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态 1:启用 2:禁用',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_channel` (`channel`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='支付渠道表';

-- 支付日志表
CREATE TABLE `payment_logs` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `payment_no` varchar(64) NOT NULL COMMENT '支付单号',
    `type` tinyint NOT NULL COMMENT '类型 1:支付 2:退款',
    `channel` tinyint NOT NULL COMMENT '支付渠道',
    `request_data` json DEFAULT NULL COMMENT '请求数据',
    `response_data` json DEFAULT NULL COMMENT '响应数据',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_payment_no` (`payment_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='支付日志表';