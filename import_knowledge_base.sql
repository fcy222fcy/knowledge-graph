-- 导入软件工程综合知识库到数据库
-- 使用方法: mysql -u root -p123456 software_qa_platform < import_knowledge_base.sql

-- 首先检查是否已存在，避免重复导入
-- 删除旧的综合知识库（如果存在）
DELETE FROM documents WHERE title = '软件工程综合知识库';

-- 插入综合知识库文档
INSERT INTO documents (title, file_path, content, status, created_at, updated_at)
VALUES (
    '软件工程综合知识库',
    './Knowledge/软件工程综合知识库.md',
    LOAD_FILE('e:/goCode/goFile/software engineering/Knowledge/软件工程综合知识库.md'),
    'active',
    NOW(),
    NOW()
);

-- 验证导入
SELECT id, title, file_path, LENGTH(content) as content_size, status, created_at
FROM documents
WHERE title = '软件工程综合知识库';

-- 显示所有文档
SELECT id, title, file_path, status, created_at
FROM documents
ORDER BY id;
