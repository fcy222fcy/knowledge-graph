-- 知识点多对多关系迁移脚本
-- 功能：创建中间表，迁移现有数据，删除旧的document_id字段

-- 1. 创建中间表
CREATE TABLE IF NOT EXISTS knowledge_point_documents (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    created_at DATETIME(3) DEFAULT NULL,
    updated_at DATETIME(3) DEFAULT NULL,
    knowledge_point_id BIGINT UNSIGNED NOT NULL COMMENT '知识点ID',
    document_id BIGINT UNSIGNED NOT NULL COMMENT '文档ID',
    PRIMARY KEY (id),
    UNIQUE KEY uk_point_document (knowledge_point_id, document_id),
    KEY idx_knowledge_point_id (knowledge_point_id),
    KEY idx_document_id (document_id),
    CONSTRAINT fk_kpd_knowledge_point FOREIGN KEY (knowledge_point_id) REFERENCES knowledge_points(id) ON DELETE CASCADE,
    CONSTRAINT fk_kpd_document FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='知识点-文档关联表';

-- 2. 迁移现有数据
INSERT IGNORE INTO knowledge_point_documents (knowledge_point_id, document_id, created_at, updated_at)
SELECT id, document_id, created_at, updated_at
FROM knowledge_points
WHERE document_id > 0;

-- 3. 验证迁移结果
SELECT '=== 迁移统计 ===' as info;
SELECT COUNT(*) as total_records FROM knowledge_point_documents;

SELECT '=== 每个知识点关联的文档数 ===' as info;
SELECT kp.name, COUNT(kpd.document_id) as doc_count
FROM knowledge_points kp
LEFT JOIN knowledge_point_documents kpd ON kp.id = kpd.knowledge_point_id
GROUP BY kp.id, kp.name
HAVING doc_count > 1
ORDER BY doc_count DESC
LIMIT 10;

-- 4. 注意：暂时保留document_id字段，等代码修改完成后再删除
-- ALTER TABLE knowledge_points DROP COLUMN document_id;
