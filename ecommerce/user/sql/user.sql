CREATE DATABASE IF NOT EXISTS `mall_user`;
USE `mall_user`;

-- 用户表
CREATE TABLE `users` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `username` varchar(32) NOT NULL COMMENT '用户名',
    `password` varchar(128) NOT NULL COMMENT '密码',
    `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
    `email` varchar(50) DEFAULT NULL COMMENT '邮箱',
    `nickname` varchar(32) DEFAULT NULL COMMENT '昵称',
    `avatar` varchar(255) DEFAULT NULL COMMENT '头像URL',
    `gender` varchar(10) NOT NULL DEFAULT 'unset' COMMENT '性别',
    `member_level` tinyint DEFAULT '0' COMMENT '会员等级',
    `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态 0:禁用 1:启用',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_username` (`username`),
    UNIQUE KEY `idx_phone` (`phone`),
    UNIQUE KEY `idx_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 用户地址表
CREATE TABLE `user_addresses` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '地址ID',
    `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
    `receiver_name` varchar(32) NOT NULL COMMENT '收货人姓名',
    `receiver_phone` varchar(20) NOT NULL COMMENT '收货人电话',
    `province` varchar(32) NOT NULL COMMENT '省份',
    `city` varchar(32) NOT NULL COMMENT '城市',
    `district` varchar(32) NOT NULL COMMENT '区/县',
    `detail_address` varchar(200) NOT NULL COMMENT '详细地址',
    `is_default` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否默认地址',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户地址表';

-- 钱包账户表
CREATE TABLE `wallet_accounts` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '钱包ID',
    `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
    `balance` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '账户余额',
    `frozen_amount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '冻结金额',
    `pay_password` varchar(128) DEFAULT NULL COMMENT '支付密码',
    `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态 0:冻结 1:正常',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='钱包账户表';

-- 钱包交易记录表
CREATE TABLE `wallet_transactions` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '交易ID',
    `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
    `order_id` varchar(64) NOT NULL COMMENT '订单号',
    `amount` decimal(10,2) NOT NULL COMMENT '交易金额',
    `type` tinyint NOT NULL COMMENT '交易类型 1:充值 2:提现 3:消费 4:退款',
    `status` tinyint NOT NULL DEFAULT '0' COMMENT '交易状态 0:处理中 1:成功 2:失败',
    `remark` varchar(255) DEFAULT NULL COMMENT '交易备注',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    UNIQUE KEY `idx_order_id` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='钱包交易记录表';

-- 登录记录表
CREATE TABLE `login_records` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '记录ID',
    `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
    `login_ip` varchar(50) NOT NULL COMMENT '登录IP',
    `login_location` varchar(100) DEFAULT NULL COMMENT '登录地点',
    `device_type` varchar(50) DEFAULT NULL COMMENT '设备类型',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='登录记录表';