CREATE DATABASE IF NOT EXISTS `mall_search`;
USE `mall_search`;

-- 热门关键词表
CREATE TABLE `hot_keywords` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `keyword` varchar(100) NOT NULL COMMENT '关键词',
    `count` bigint NOT NULL DEFAULT '0' COMMENT '搜索次数',
    `is_manual` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否手动添加',
    `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态 1:显示 2:隐藏',
    `sort` int NOT NULL DEFAULT '0' COMMENT '排序',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_keyword` (`keyword`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='热门关键词表';

-- 搜索历史表
CREATE TABLE `search_histories` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
    `keyword` varchar(100) NOT NULL COMMENT '搜索关键词',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_time` (`user_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='搜索历史表';

-- 搜索统计表
CREATE TABLE `search_statistics` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `keyword` varchar(100) NOT NULL COMMENT '搜索关键词',
    `count` bigint NOT NULL DEFAULT '0' COMMENT '搜索次数',
    `result_count` bigint NOT NULL DEFAULT '0' COMMENT '结果数量',
    `date` date NOT NULL COMMENT '统计日期',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_keyword_date` (`keyword`, `date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='搜索统计表';

-- 商品索引表（用于备份）
CREATE TABLE `product_indices` (
    `product_id` bigint unsigned NOT NULL COMMENT '商品ID',
    `data` json NOT NULL COMMENT '索引数据',
    `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态 1:正常 2:删除',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品索引表';