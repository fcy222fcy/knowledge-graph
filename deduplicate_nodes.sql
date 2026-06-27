-- 知识图谱节点去重脚本
-- 功能：合并重复的节点，保留一个主节点，其他节点关联到主节点

-- 1. 创建临时表，存储需要保留的主节点
CREATE TEMPORARY TABLE IF NOT EXISTS nodes_to_keep AS
SELECT MIN(id) as keep_id, name
FROM knowledge_points
GROUP BY name
HAVING COUNT(*) > 1;

-- 2. 创建临时表，存储需要删除的重复节点
CREATE TEMPORARY TABLE IF NOT EXISTS nodes_to_delete AS
SELECT kp.id as delete_id, kp.name, ntk.keep_id
FROM knowledge_points kp
JOIN nodes_to_keep ntk ON kp.name = ntk.name
WHERE kp.id != ntk.keep_id;

-- 3. 查看将要处理的数据
SELECT '=== 将要保留的主节点 ===' as info;
SELECT keep_id, name FROM nodes_to_keep;

SELECT '=== 将要删除的重复节点 ===' as info;
SELECT delete_id, name, keep_id FROM nodes_to_delete;

-- 4. 更新关系表：将指向重复节点的关系重定向到主节点
UPDATE knowledge_relations kr
JOIN nodes_to_delete nd ON kr.source_id = nd.delete_id
SET kr.source_id = nd.keep_id
WHERE kr.source_id IN (SELECT delete_id FROM nodes_to_delete);

UPDATE knowledge_relations kr
JOIN nodes_to_delete nd ON kr.target_id = nd.delete_id
SET kr.target_id = nd.keep_id
WHERE kr.target_id IN (SELECT delete_id FROM nodes_to_delete);

-- 5. 删除重复的关系（保留主节点之间的关系）
DELETE kr1 FROM knowledge_relations kr1
JOIN knowledge_relations kr2
ON kr1.source_id = kr2.source_id
AND kr1.target_id = kr2.target_id
AND kr1.id > kr2.id;

-- 6. 删除重复的节点
DELETE FROM knowledge_points
WHERE id IN (SELECT delete_id FROM nodes_to_delete);

-- 7. 验证结果
SELECT '=== 去重后的节点统计 ===' as info;
SELECT COUNT(*) as total_nodes FROM knowledge_points;

SELECT '=== 去重后的重复检查 ===' as info;
SELECT name, COUNT(*) as count
FROM knowledge_points
GROUP BY name
HAVING COUNT(*) > 1;

-- 清理临时表
DROP TEMPORARY TABLE IF EXISTS nodes_to_keep;
DROP TEMPORARY TABLE IF EXISTS nodes_to_delete;
