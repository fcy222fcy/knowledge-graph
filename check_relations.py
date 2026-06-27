#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
检查Neo4j中的关系类型
"""

from neo4j import GraphDatabase

NEO4J_URI = "bolt://localhost:7687"
NEO4J_USER = "neo4j"
NEO4J_PASSWORD = "12345678"

def check_relations():
    driver = GraphDatabase.driver(NEO4J_URI, auth=(NEO4J_USER, NEO4J_PASSWORD))

    with driver.session() as session:
        # 检查关系类型
        print("=== 关系类型统计 ===")
        result = session.run("""
            MATCH ()-[r]->()
            RETURN type(r) as rel_type, count(r) as cnt
        """)
        for record in result:
            print(f"  {record['rel_type']}: {record['cnt']} 条")

        # 检查一些关系的详细信息
        print("\n=== 关系详细信息 ===")
        result = session.run("""
            MATCH (a:KnowledgePoint)-[r]->(b:KnowledgePoint)
            RETURN a.id as source_id, a.name as source_name,
                   type(r) as rel_type,
                   b.id as target_id, b.name as target_name
            LIMIT 5
        """)
        for record in result:
            print(f"  {record['source_name']} (id:{record['source_id']}) --[{record['rel_type']}]--> {record['target_name']} (id:{record['target_id']})")

        # 检查是否有空关系
        print("\n=== 空关系检查 ===")
        result = session.run("""
            MATCH (a:KnowledgePoint)-[r]->(b:KnowledgePoint)
            WHERE r.type IS NULL OR r.type = ''
            RETURN count(r) as empty_count
        """)
        empty_count = result.single()['empty_count']
        print(f"  空type的关系: {empty_count} 条")

    driver.close()

if __name__ == '__main__':
    check_relations()
