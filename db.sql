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