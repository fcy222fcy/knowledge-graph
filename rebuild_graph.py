#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
重建Neo4j知识图谱
清除旧数据，从MySQL重新同步
"""

from neo4j import GraphDatabase
import mysql.connector

# Neo4j配置
NEO4J_URI = "bolt://localhost:7687"
NEO4J_USER = "neo4j"
NEO4J_PASSWORD = "12345678"

# MySQL配置
MYSQL_CONFIG = {
    'host': '127.0.0.1',
    'port': 3306,
    'user': 'root',
    'password': '123456',
    'database': 'software_qa_platform',
    'charset': 'utf8mb4'
}

def rebuild_graph():
    """重建Neo4j图谱"""
    print("=== 重建Neo4j知识图谱 ===\n")

    # 1. 连接Neo4j
    try:
        driver = GraphDatabase.driver(NEO4J_URI, auth=(NEO4J_USER, NEO4J_PASSWORD))
        driver.verify_connectivity()
        print("[OK] Neo4j连接成功")
    except Exception as e:
        print(f"[ERROR] Neo4j连接失败: {e}")
        return

    # 2. 清除所有旧数据
    with driver.session() as session:
        session.run("MATCH (n) DETACH DELETE n")
        print("[OK] 已清除所有旧节点和关系")

    # 3. 从MySQL读取知识点
    conn = mysql.connector.connect(**MYSQL_CONFIG)
    cursor = conn.cursor(dictionary=True)

    cursor.execute("""
        SELECT kp.id, kp.name, kp.description, kp.document_id, kp.category,
               d.title as document_title
        FROM knowledge_points kp
        LEFT JOIN documents d ON kp.document_id = d.id
    """)
    points = cursor.fetchall()
    print(f"[OK] 从MySQL读取了 {len(points)} 个知识点")

    # 4. 创建知识点节点
    with driver.session() as session:
        for point in points:
            session.run("""
                CREATE (n:KnowledgePoint {
                    id: $id,
                    name: $name,
                    description: $description,
                    document_id: $document_id,
                    category: $category
                })
            """, {
                'id': point['id'],
                'name': point['name'],
                'description': point['description'] or '',
                'document_id': point['document_id'] or 0,
                'category': point['category'] or ''
            })
    print(f"[OK] 创建了 {len(points)} 个知识点节点")

    # 5. 从MySQL读取关系
    cursor.execute("""
        SELECT kr.id, kr.source_id, kr.target_id, kr.relation_type, kr.description
        FROM knowledge_relations kr
        WHERE kr.source_id IN (SELECT id FROM knowledge_points)
          AND kr.target_id IN (SELECT id FROM knowledge_points)
    """)
    relations = cursor.fetchall()
    print(f"[OK] 从MySQL读取了 {len(relations)} 条关系")

    # 6. 创建关系
    with driver.session() as session:
        for rel in relations:
            # 使用实际的关系类型作为Neo4j关系类型
            rel_type = rel['relation_type'] or 'RELATED'
            session.run(f"""
                MATCH (a:KnowledgePoint {{id: $source_id}})
                MATCH (b:KnowledgePoint {{id: $target_id}})
                CREATE (a)-[r:{rel_type} {{
                    id: $id,
                    description: $description
                }}]->(b)
            """, {
                'id': rel['id'],
                'source_id': rel['source_id'],
                'target_id': rel['target_id'],
                'description': rel['description'] or ''
            })
    print(f"[OK] 创建了 {len(relations)} 条关系")

    # 7. 验证结果
    with driver.session() as session:
        result = session.run("MATCH (n:KnowledgePoint) RETURN count(n) as count")
        node_count = result.single()['count']

        result = session.run("MATCH ()-[r]->() RETURN count(r) as count")
        rel_count = result.single()['count']

    print(f"\n=== 重建完成 ===")
    print(f"  节点数: {node_count}")
    print(f"  关系数: {rel_count}")

    cursor.close()
    conn.close()
    driver.close()

if __name__ == '__main__':
    rebuild_graph()
