-- ============================================================
-- SE智图问答平台 数据库建表语句
-- 数据库: software_engineering
-- ============================================================

CREATE DATABASE IF NOT EXISTS `software_engineering`
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

USE `software_engineering`;

-- ------------------------------------------------------------
-- 1. users - 学生用户表
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `users` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `created_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME        DEFAULT NULL COMMENT '软删除时间',
  `username`   VARCHAR(50)     NOT NULL COMMENT '用户名',
  `password`   VARCHAR(255)    NOT NULL COMMENT '密码哈希',
  `email`      VARCHAR(100)    NOT NULL COMMENT '邮箱地址',
  `nickname`   VARCHAR(50)     DEFAULT '' COMMENT '用户昵称',
  `avatar`     VARCHAR(255)    DEFAULT '' COMMENT '头像URL',
  `status`     TINYINT         NOT NULL DEFAULT 1 COMMENT '状态 1=启用 0=禁用',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  UNIQUE KEY `uk_email` (`email`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='学生用户表';

-- ------------------------------------------------------------
-- 2. teachers - 教师表
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `teachers` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `created_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME        DEFAULT NULL COMMENT '软删除时间',
  `username`   VARCHAR(50)     NOT NULL COMMENT '用户名',
  `password`   VARCHAR(255)    NOT NULL COMMENT '密码哈希',
  `email`      VARCHAR(100)    NOT NULL COMMENT '邮箱地址',
  `nickname`   VARCHAR(50)     DEFAULT '' COMMENT '教师昵称',
  `avatar`     VARCHAR(255)    DEFAULT '' COMMENT '头像URL',
  `status`     TINYINT         NOT NULL DEFAULT 1 COMMENT '状态 1=启用 0=禁用',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  UNIQUE KEY `uk_email` (`email`),
  KEY `idx_teachers_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='教师表';

-- ------------------------------------------------------------
-- 3. documents - 文档表
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `documents` (
  `id`             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `created_at`     DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`     DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at`     DATETIME        DEFAULT NULL COMMENT '软删除时间',
  `user_id`        BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `title`          VARCHAR(200)    NOT NULL COMMENT '文档标题',
  `description`    VARCHAR(500)    DEFAULT '' COMMENT '文档描述',
  `filename`       VARCHAR(200)    NOT NULL COMMENT '原始文件名',
  `file_path`      VARCHAR(500)    NOT NULL COMMENT '文件存储路径',
  `file_size`      BIGINT          DEFAULT 0 COMMENT '文件大小字节',
  `file_type`      VARCHAR(20)     DEFAULT '' COMMENT '文件类型 pdf/docx/txt 等',
  `content`        LONGTEXT        COMMENT '解析后的文本内容',
  `status`         VARCHAR(20)     NOT NULL DEFAULT 'pending' COMMENT '处理状态 pending/approved/rejected/completed/failed',
  `review_comment` VARCHAR(500)    DEFAULT '' COMMENT '审核意见',
  PRIMARY KEY (`id`),
  KEY `idx_documents_user_id` (`user_id`),
  KEY `idx_documents_status` (`status`),
  KEY `idx_documents_deleted_at` (`deleted_at`),
  KEY `idx_documents_user_id_status` (`user_id`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文档表';

-- ------------------------------------------------------------
-- 4. knowledge_points - 知识点表
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `knowledge_points` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `created_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at`  DATETIME        DEFAULT NULL COMMENT '软删除时间',
  `name`        VARCHAR(100)    NOT NULL COMMENT '知识点名称',
  `description` VARCHAR(500)    DEFAULT '' COMMENT '知识点描述',
  `document_id` BIGINT UNSIGNED DEFAULT 0 COMMENT '关联的文档ID',
  `category`    VARCHAR(50)     DEFAULT '' COMMENT '知识点分类',
  PRIMARY KEY (`id`),
  KEY `idx_knowledge_points_document_id` (`document_id`),
  KEY `idx_knowledge_points_category` (`category`),
  KEY `idx_knowledge_points_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='知识点表';

-- ------------------------------------------------------------
-- 5. knowledge_relations - 知识点关系表
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `knowledge_relations` (
  `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `created_at`    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at`    DATETIME        DEFAULT NULL COMMENT '软删除时间',
  `source_id`     BIGINT UNSIGNED NOT NULL COMMENT '源知识点ID',
  `target_id`     BIGINT UNSIGNED NOT NULL COMMENT '目标知识点ID',
  `type`          VARCHAR(60)     NOT NULL DEFAULT '' COMMENT '关系类型标识',
  `relation_type` VARCHAR(20)     NOT NULL COMMENT '关系类型 RELATED/DEPENDS_ON/PART_OF',
  `description`   VARCHAR(500)    DEFAULT '' COMMENT '关系描述',
  PRIMARY KEY (`id`),
  KEY `idx_knowledge_relations_source_id` (`source_id`),
  KEY `idx_knowledge_relations_target_id` (`target_id`),
  KEY `idx_knowledge_relations_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='知识点关系表';

-- ------------------------------------------------------------
-- 6. knowledge_builds - 知识图谱构建记录表
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `knowledge_builds` (
  `id`                BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `created_at`        DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`        DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at`        DATETIME        DEFAULT NULL COMMENT '软删除时间',
  `document_ids`      VARCHAR(500)    DEFAULT '' COMMENT '逗号分隔的文档ID列表',
  `created_points`    INT             DEFAULT 0 COMMENT '创建的知识点数量',
  `created_relations` INT             DEFAULT 0 COMMENT '创建的关系数量',
  `chunk_count`       INT             DEFAULT 0 COMMENT '文档分块数量',
  `vector_count`      INT             DEFAULT 0 COMMENT '向量数量',
  `status`            VARCHAR(20)     NOT NULL DEFAULT 'completed' COMMENT '构建状态',
  `message`           VARCHAR(500)    DEFAULT '' COMMENT '构建结果描述',
  PRIMARY KEY (`id`),
  KEY `idx_knowledge_builds_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='知识图谱构建记录表';

-- ------------------------------------------------------------
-- 7. questions - 题目表
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `questions` (
  `id`                  BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `created_at`          DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`          DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at`          DATETIME        DEFAULT NULL COMMENT '软删除时间',
  `title`               VARCHAR(500)    NOT NULL COMMENT '题目标题',
  `type`                VARCHAR(20)     NOT NULL COMMENT '题目类型 single/multiple',
  `difficulty`          VARCHAR(20)     NOT NULL COMMENT '难度 easy/medium/hard',
  `knowledge_point_id`  BIGINT UNSIGNED DEFAULT 0 COMMENT '关联的知识点ID',
  `options`             TEXT            COMMENT '选项 JSON数组',
  `answer`              VARCHAR(20)     NOT NULL COMMENT '正确答案',
  `explanation`         TEXT            COMMENT '题目解析',
  PRIMARY KEY (`id`),
  KEY `idx_questions_knowledge_point_id` (`knowledge_point_id`),
  KEY `idx_questions_type` (`type`),
  KEY `idx_questions_difficulty` (`difficulty`),
  KEY `idx_questions_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='题目表';

-- ------------------------------------------------------------
-- 8. quizzes - 答题记录表
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `quizzes` (
  `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `created_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at`  DATETIME        DEFAULT NULL COMMENT '软删除时间',
  `question_id` BIGINT UNSIGNED NOT NULL COMMENT '题目ID',
  `user_id`     BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `user_answer` VARCHAR(20)     NOT NULL COMMENT '用户提交的答案',
  `is_correct`  TINYINT(1)      DEFAULT 0 COMMENT '是否回答正确',
  PRIMARY KEY (`id`),
  KEY `idx_quizzes_question_id` (`question_id`),
  KEY `idx_quizzes_user_id` (`user_id`),
  KEY `idx_quizzes_deleted_at` (`deleted_at`),
  KEY `idx_quizzes_user_id_is_correct` (`user_id`, `is_correct`),
  KEY `idx_quizzes_user_id_created_at` (`user_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='答题记录表';

-- ------------------------------------------------------------
-- 9. ask_sessions - 问答会话表
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `ask_sessions` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `created_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME        DEFAULT NULL COMMENT '软删除时间',
  `user_id`    BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `title`      VARCHAR(200)    DEFAULT '' COMMENT '会话标题',
  PRIMARY KEY (`id`),
  KEY `idx_ask_sessions_user_id` (`user_id`),
  KEY `idx_ask_sessions_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='问答会话表';

-- ------------------------------------------------------------
-- 10. ask_messages - 问答消息表
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `ask_messages` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `created_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME        DEFAULT NULL COMMENT '软删除时间',
  `session_id` BIGINT UNSIGNED NOT NULL COMMENT '会话ID',
  `role`       VARCHAR(20)     NOT NULL COMMENT '消息角色 user/assistant',
  `content`    TEXT            COMMENT '消息内容',
  `confidence` DOUBLE          DEFAULT 0 COMMENT '回答置信度 0-1',
  PRIMARY KEY (`id`),
  KEY `idx_ask_messages_session_id` (`session_id`),
  KEY `idx_ask_messages_deleted_at` (`deleted_at`),
  KEY `idx_ask_messages_session_id_role` (`session_id`, `role`),
  KEY `idx_ask_messages_session_id_created_at` (`session_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='问答消息表';
