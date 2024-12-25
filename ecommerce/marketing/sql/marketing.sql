CREATE DATABASE IF NOT EXISTS `mall_marketing`;
USE `mall_marketing`;

-- 优惠券表
CREATE TABLE `coupons` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '优惠券ID',
    `name` varchar(100) NOT NULL COMMENT '优惠券名称',
    `code` varchar(32) NOT NULL COMMENT '优惠券码',
    `type` tinyint NOT NULL COMMENT '优惠券类型 1:满减 2:折扣 3:立减',
    `value` decimal(10,2) NOT NULL COMMENT '优惠金额或折扣率',
    `min_amount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '最低使用金额',
    `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态 0:未开始 1:进行中 2:已结束 3:已失效',
    `start_time` timestamp NULL DEFAULT NULL COMMENT '开始时间',
    `end_time` timestamp NULL DEFAULT NULL COMMENT '结束时间',
    `total` int NOT NULL DEFAULT '0' COMMENT '发行总量',
    `received` int NOT NULL DEFAULT '0' COMMENT '已领取数量',
    `used` int NOT NULL DEFAULT '0' COMMENT '已使用数量',
    `per_limit` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否限制每人领取数量',
    `user_limit` int NOT NULL DEFAULT '1' COMMENT '每人限领数量',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='优惠券表';

-- 用户优惠券表
CREATE TABLE `user_coupons` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
    `coupon_id` bigint unsigned NOT NULL COMMENT '优惠券ID',
    `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态 0:未使用 1:已使用 2:已过期',
    `used_time` timestamp NULL DEFAULT NULL COMMENT '使用时间',
    `order_no` varchar(64) DEFAULT NULL COMMENT '订单号',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_coupon_id` (`coupon_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户优惠券表';

-- 促销活动表
CREATE TABLE `promotions` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '活动ID',
    `name` varchar(100) NOT NULL COMMENT '活动名称',
    `type` tinyint NOT NULL COMMENT '活动类型 1:满减 2:折扣 3:秒杀',
    `rules` json NOT NULL COMMENT '促销规则',
    `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态 0:未开始 1:进行中 2:已结束',
    `start_time` timestamp NULL DEFAULT NULL COMMENT '开始时间',
    `end_time` timestamp NULL DEFAULT NULL COMMENT '结束时间',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='促销活动表';

-- 用户积分表
CREATE TABLE `user_points` (
    `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
    `points` bigint NOT NULL DEFAULT '0' COMMENT '积分余额',
    `total_points` bigint NOT NULL DEFAULT '0' COMMENT '累计获得积分',
    `used_points` bigint NOT NULL DEFAULT '0' COMMENT '已使用积分',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户积分表';

-- 积分记录表
CREATE TABLE `points_records` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '记录ID',
    `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
    `points` bigint NOT NULL COMMENT '积分变动数量',
    `type` tinyint NOT NULL COMMENT '类型 1:获取 2:使用',
    `source` varchar(50) NOT NULL COMMENT '来源',
    `remark` varchar(255) DEFAULT NULL COMMENT '备注',
    `order_no` varchar(64) DEFAULT NULL COMMENT '订单号',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='积分记录表';