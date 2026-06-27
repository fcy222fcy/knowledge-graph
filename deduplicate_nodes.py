#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
知识图谱节点去重脚本
"""

import mysql.connector

DB_CONFIG = {
    'host': '127.0.0.1',
    'port': 3306,
    'user': 'root',
    'password': '123456',
    'database': 'software_qa_platform',
    'charset': 'utf8mb4'
}

def deduplicate_nodes():
    """去重知识图谱节点"""
    conn = mysql.connector.connect(**DB_CONFIG)
    cursor = conn.cursor(dictionary=True)

    try:
        # 1. 查找重复的节点
        cursor.execute("""
            SELECT name, COUNT(*) as count, GROUP_CONCAT(id) as ids
            FROM knowledge_points
            GROUP BY name
            HAVING COUNT(*) > 1
        """)
        duplicates = cursor.fetchall()

        print(f"找到 {len(duplicates)} 组重复节点\n")

        total_deleted = 0
        total_relations_updated = 0

        for dup in duplicates:
            name = dup['name']
            ids = [int(x) for x in dup['ids'].split(',')]
            keep_id = min(ids)  # 保留ID最小的节点
            delete_ids = [x for x in ids if x != keep_id]

            print(f"节点: {name}")
            print(f"  保留: ID {keep_id}")
            print(f"  删除: ID {delete_ids}")

            # 2. 更新关系表：将指向重复节点的关系重定向到主节点
            for delete_id in delete_ids:
                # 更新source_id
                cursor.execute("""
                    UPDATE knowledge_relations
                    SET source_id = %s
                    WHERE source_id = %s
                """, (keep_id, delete_id))
                total_relations_updated += cursor.rowcount

                # 更新target_id
                cursor.execute("""
                    UPDATE knowledge_relations
                    SET target_id = %s
                    WHERE target_id = %s
                """, (keep_id, delete_id))
                total_relations_updated += cursor.rowcount

            # 3. 删除重复的关系（保留主节点之间的关系）
            cursor.execute("""
                DELETE kr1 FROM knowledge_relations kr1
                INNER JOIN knowledge_relations kr2
                ON kr1.source_id = kr2.source_id
                AND kr1.target_id = kr2.target_id
                AND kr1.id > kr2.id
            """)

            # 4. 删除重复的节点
            cursor.execute("""
                DELETE FROM knowledge_points
                WHERE id IN (%s)
            """ % ','.join(['%s'] * len(delete_ids)), delete_ids)
            total_deleted += cursor.rowcount

            print(f"  已删除 {cursor.rowcount} 个节点\n")

        conn.commit()

        # 5. 统计结果
        cursor.execute("SELECT COUNT(*) as count FROM knowledge_points")
        node_count = cursor.fetchone()['count']

        cursor.execute("SELECT COUNT(*) as count FROM knowledge_relations")
        relation_count = cursor.fetchone()['count']

        print("=" * 50)
        print("去重完成！")
        print(f"  - 关系更新: {total_relations_updated} 条")
        print(f"  - 节点删除: {total_deleted} 个")
        print(f"  - 剩余节点: {node_count} 个")
        print(f"  - 剩余关系: {relation_count} 条")

        # 6. 检查是否还有重复
        cursor.execute("""
            SELECT name, COUNT(*) as count
            FROM knowledge_points
            GROUP BY name
            HAVING COUNT(*) > 1
        """)
        remaining_duplicates = cursor.fetchall()

        if remaining_duplicates:
            print(f"\n警告: 还有 {len(remaining_duplicates)} 组重复节点")
        else:
            print("\n✓ 没有重复节点了")

    except Exception as e:
        print(f"错误: {e}")
        conn.rollback()
    finally:
        cursor.close()
        conn.close()

if __name__ == '__main__':
    deduplicate_nodes()
