create database `db_jike`;

CREATE TABLE `db_jike`.`account`(
                                    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
                                    `mobile` varchar(20) NOT NULL DEFAULT '' COMMENT '手机号',
                                    `email` varchar(100) NOT NULL DEFAULT '' COMMENT '邮箱',
                                    `password` varchar(255) NOT NULL DEFAULT '' COMMENT '密码',
                                    `nickname` varchar(50) NOT NULL DEFAULT '' COMMENT '昵称',
                                    `avatar` varchar(500) NOT NULL DEFAULT '' COMMENT '头像URL',
                                    `gender` tinyint(1) NOT NULL DEFAULT 0 COMMENT '性别: 0-未知 1-男 2-女',
                                    `birthday` date NOT NULL DEFAULT '1970-01-01' COMMENT '生日',
                                    `bio` varchar(500) NOT NULL DEFAULT '' COMMENT '个人简介',
                                    `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态: 0-禁用 1-正常 2-冻结',
                                    `last_login_at` datetime NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '最后登录时间',
                                    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                    PRIMARY KEY (`id`),
                                    UNIQUE KEY `uk_mobile` (`mobile`),
                                    KEY `idx_status` (`status`),
                                    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户账户表';

CREATE TABLE `article_channel` (
     `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '频道Id',
     `name` VARCHAR(100) NOT NULL COMMENT '频道名称',
     `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
     `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
     PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章频道表';

-- 插入文章频道数据
INSERT INTO `article_channel` (`name`) VALUES
                                           ('技术'),
                                           ('生活'),
                                           ('旅游'),
                                           ('美食'),
                                           ('科技'),
                                           ('编程');

CREATE TABLE `article` (
                           `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '文章Id',
                           `title` varchar(200) NOT NULL COMMENT '文章标题',
                           `content` longtext NOT NULL COMMENT '文章内容',
                           `cover_type` tinyint(4) NOT NULL DEFAULT 1 COMMENT '封面类型：1-无图，2-单图，3-多图',
                           `cover_images` varchar(200) DEFAULT NULL COMMENT '封面图片地址数组',
                           `channel_id` bigint(20) NOT NULL COMMENT '频道Id',
                           `status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '文章状态：0-草稿，1-已发布，2-已删除',
                           `author_id` bigint(20) NOT NULL COMMENT '作者Id',
                           `view_count` int(11) NOT NULL DEFAULT 0 COMMENT '阅读量',
                           `like_count` int(11) NOT NULL DEFAULT 0 COMMENT '点赞数',
                           `comment_count` int(11) NOT NULL DEFAULT 0 COMMENT '评论数',
                           `is_top` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否置顶：0-否，1-是',
                           `is_recommend` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否推荐：0-否，1-是',
                           `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                           `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                           PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章表';
