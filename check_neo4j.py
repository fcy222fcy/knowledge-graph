#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
检查Neo4j中的图谱数据
"""

from neo4j import GraphDatabase

NEO4J_URI = "bolt://localhost:7687"
NEO4J_USER = "neo4j"
NEO4J_PASSWORD = "12345678"

def check_neo4j():
    driver = GraphDatabase.driver(NEO4J_URI, auth=(NEO4J_USER, NEO4J_PASSWORD))

    with driver.session() as session:
        # 检查节点数
        result = session.run("MATCH (n:KnowledgePoint) RETURN count(n) as count")
        node_count = result.single()['count']
        print(f"节点数: {node_count}")

        # 检查关系数
        result = session.run("MATCH ()-[r]->() RETURN count(r) as count")
        rel_count = result.single()['count']
        print(f"关系数: {rel_count}")

        # 查看一些关系
        print("\n=== 关系示例 ===")
        result = session.run("""
            MATCH (a:KnowledgePoint)-[r]->(b:KnowledgePoint)
            RETURN a.name as source, type(r) as rel_type, b.name as target
            LIMIT 10
        """)
        for record in result:
            print(f"  {record['source']} --[{record['rel_type']}]--> {record['target']}")

        # 检查是否有重复节点
        print("\n=== 重复节点检查 ===")
        result = session.run("""
            MATCH (n:KnowledgePoint)
            WITH n.name as name, count(n) as cnt
            WHERE cnt > 1
            RETURN name, cnt
        """)
        duplicates = list(result)
        if duplicates:
            for d in duplicates:
                print(f"  {d['name']}: {d['cnt']} 个")
        else:
            print("  没有重复节点")

    driver.close()

if __name__ == '__main__':
    check_neo4j()
