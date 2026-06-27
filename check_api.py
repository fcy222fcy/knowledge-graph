#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
检查图谱API返回的数据
"""

import requests

API_URL = "http://localhost:8080/api/v1/graph"
TOKEN = "your_token_here"  # 需要替换为实际的token

def check_api():
    # 先获取所有数据
    headers = {"Authorization": f"Bearer {TOKEN}"}

    try:
        # 获取所有图谱数据
        resp = requests.get(API_URL, headers=headers, timeout=5)
        if resp.status_code == 200:
            data = resp.json()
            nodes = data.get('data', {}).get('nodes', [])
            edges = data.get('data', {}).get('edges', [])

            print(f"节点数: {len(nodes)}")
            print(f"关系数: {len(edges)}")

            # 统计关系类型
            rel_types = {}
            for edge in edges:
                rel_type = edge.get('relation_type', 'unknown')
                rel_types[rel_type] = rel_types.get(rel_type, 0) + 1

            print("\n关系类型统计:")
            for rt, count in sorted(rel_types.items(), key=lambda x: -x[1]):
                print(f"  {rt}: {count}")

            # 测试筛选
            print("\n=== 测试筛选 contains ===")
            resp2 = requests.get(API_URL, params={"relation_type": "contains"}, headers=headers, timeout=5)
            if resp2.status_code == 200:
                data2 = resp2.json()
                edges2 = data2.get('data', {}).get('edges', [])
                print(f"筛选后关系数: {len(edges2)}")
            else:
                print(f"筛选请求失败: {resp2.status_code}")
        else:
            print(f"请求失败: {resp.status_code}")
            print(resp.text)
    except Exception as e:
        print(f"错误: {e}")

if __name__ == '__main__':
    check_api()
