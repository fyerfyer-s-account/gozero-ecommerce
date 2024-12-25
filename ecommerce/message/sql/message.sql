CREATE DATABASE IF NOT EXISTS `mall_message`;
USE `mall_message`;

-- 消息表
CREATE TABLE `messages` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '消息ID',
    `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
    `title` varchar(100) NOT NULL COMMENT '消息标题',
    `content` text NOT NULL COMMENT '消息内容',
    `type` tinyint NOT NULL COMMENT '消息类型 1:系统通知 2:订单消息 3:活动消息 4:物流消息',
    `send_channel` tinyint NOT NULL COMMENT '发送渠道 1:站内信 2:短信 3:邮件 4:APP推送',
    `extra_data` json DEFAULT NULL COMMENT '额外数据',
    `is_read` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否已读',
    `read_time` timestamp NULL DEFAULT NULL COMMENT '阅读时间',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_type` (`user_id`, `type`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息表';

-- 消息模板表
CREATE TABLE `message_templates` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '模板ID',
    `code` varchar(50) NOT NULL COMMENT '模板编码',
    `name` varchar(100) NOT NULL COMMENT '模板名称',
    `title_template` varchar(200) NOT NULL COMMENT '标题模板',
    `content_template` text NOT NULL COMMENT '内容模板',
    `type` tinyint NOT NULL COMMENT '消息类型',
    `channels` json NOT NULL COMMENT '发送渠道',
    `config` json DEFAULT NULL COMMENT '渠道配置',
    `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态 1:启用 2:禁用',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息模板表';

-- 消息发送记录表
CREATE TABLE `message_sends` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '记录ID',
    `message_id` bigint unsigned NOT NULL COMMENT '消息ID',
    `template_id` bigint unsigned DEFAULT NULL COMMENT '模板ID',
    `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
    `channel` tinyint NOT NULL COMMENT '发送渠道',
    `status` tinyint NOT NULL COMMENT '发送状态 1:待发送 2:发送中 3:发送成功 4:发送失败',
    `error` varchar(500) DEFAULT NULL COMMENT '错误信息',
    `retry_count` int NOT NULL DEFAULT '0' COMMENT '重试次数',
    `next_retry_time` timestamp NULL DEFAULT NULL COMMENT '下次重试时间',
    `send_time` timestamp NULL DEFAULT NULL COMMENT '发送时间',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_message_id` (`message_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_status` (`status`),
    KEY `idx_next_retry` (`next_retry_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息发送记录表';

-- 通知设置表
CREATE TABLE `notification_settings` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '设置ID',
    `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
    `type` tinyint NOT NULL COMMENT '消息类型',
    `channel` tinyint NOT NULL COMMENT '通知渠道',
    `is_enabled` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_user_type_channel` (`user_id`, `type`, `channel`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='通知设置表';