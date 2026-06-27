#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
为知识点添加多文档来源关系
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

def add_multi_sources():
    """为常见知识点添加多文档来源"""
    conn = mysql.connector.connect(**DB_CONFIG)
    cursor = conn.cursor(dictionary=True)

    try:
        # 定义需要添加多来源的知识点和对应的文档
        # 格式: {知识点名称: [文档ID列表]}
        multi_source_points = {
            '编码': [4, 5, 6, 9, 11],        # 来自多个文档
            '代码': [4, 5, 6, 9],
            '开发': [4, 5, 6, 9, 11],
            '质量': [4, 5, 6, 9, 11],
            '测试': [6, 9, 11],
            '设计': [5, 9, 11],
            '需求': [4, 9, 11],
            '维护': [5, 6, 9, 11],
            '版本控制': [9, 11],
            '项目管理': [4, 5, 6, 9, 11],
            '风险': [5, 6, 9],
            '成本': [4, 6, 9],
            '进度': [6, 9],
            '评审': [4, 9],
            '团队': [9],
            '编码实现': [9, 11],
            '需求分析': [4, 9, 11],
            '系统设计': [5, 9, 11],
            '部署': [9, 11],
            '敏捷开发方法论': [9, 11],
            'Scrum框架': [9, 11],
        }

        total_added = 0

        for point_name, doc_ids in multi_source_points.items():
            # 查找知识点
            cursor.execute(
                "SELECT id, name FROM knowledge_points WHERE name = %s",
                (point_name,)
            )
            point = cursor.fetchone()

            if not point:
                print(f"[SKIP] 知识点不存在: {point_name}")
                continue

            point_id = point['id']
            added_count = 0

            for doc_id in doc_ids:
                try:
                    cursor.execute("""
                        INSERT IGNORE INTO knowledge_point_documents
                        (knowledge_point_id, document_id, created_at, updated_at)
                        VALUES (%s, %s, NOW(), NOW())
                    """, (point_id, doc_id))
                    if cursor.rowcount > 0:
                        added_count += 1
                except Exception as e:
                    print(f"  [ERROR] 添加失败: {e}")

            if added_count > 0:
                print(f"[OK] {point_name}: 添加了 {added_count} 个文档来源")
                total_added += added_count

        conn.commit()

        # 验证结果
        print("\n" + "=" * 50)
        print(f"总共添加了 {total_added} 条关联记录")

        # 统计每个知识点的文档数
        print("\n[STATS] 知识点来源统计:")
        cursor.execute("""
            SELECT kp.name, COUNT(kpd.document_id) as doc_count
            FROM knowledge_points kp
            LEFT JOIN knowledge_point_documents kpd ON kp.id = kpd.knowledge_point_id
            GROUP BY kp.id, kp.name
            HAVING doc_count > 1
            ORDER BY doc_count DESC
        """)

        for row in cursor.fetchall():
            print(f"  {row['name']}: {row['doc_count']} 个来源")

    except Exception as e:
        print(f"错误: {e}")
        import traceback
        traceback.print_exc()
        conn.rollback()
    finally:
        cursor.close()
        conn.close()

if __name__ == '__main__':
    add_multi_sources()
