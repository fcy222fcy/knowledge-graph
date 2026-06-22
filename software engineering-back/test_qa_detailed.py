#!/usr/bin/env python3
"""问答可信度详细测试"""
import requests
import json
import time

PYTHON_URL = "http://localhost:5000"

TEST_CASES = [
    {"q": "什么是软件工程", "expect_kw": ["软件", "工程"]},
    {"q": "Git是什么", "expect_kw": ["Git", "版本控制"]},
    {"q": "敏捷开发", "expect_kw": ["敏捷", "迭代"]},
    {"q": "Docker是什么", "expect_kw": ["Docker", "容器"]},
    {"q": "Kubernetes是什么", "expect_kw": ["Kubernetes", "容器"]},
    {"q": "什么是单元测试", "expect_kw": ["测试"]},
    {"q": "RESTful API设计", "expect_kw": ["API", "REST"]},
]

print("=" * 70)
print("  问答系统可信度详细测试")
print("=" * 70)

for i, tc in enumerate(TEST_CASES, 1):
    q = tc["q"]
    print(f"\n--- [{i}] 问题: {q} ---")

    # 搜索测试
    t0 = time.time()
    try:
        resp = requests.post(f"{PYTHON_URL}/search", json={"query": q, "top_k": 3}, timeout=30)
        search_time = time.time() - t0
        if resp.status_code == 200:
            results = resp.json().get("results", [])
            all_text = " ".join([r.get("chunk_text", "") for r in results])
            kw_match = sum(1 for kw in tc["expect_kw"] if kw.lower() in all_text.lower())
            avg_score = sum(r.get("score", 0) for r in results) / max(1, len(results))
            print(f"  搜索: {len(results)}条结果, 平均分={avg_score:.3f}, 关键词匹配={kw_match}/{len(tc['expect_kw'])}, 耗时={search_time:.2f}s")
        else:
            print(f"  搜索失败: {resp.status_code}")
    except Exception as e:
        print(f"  搜索异常: {e}")

    # RAG 问答测试
    t0 = time.time()
    try:
        resp = requests.post(f"{PYTHON_URL}/search_and_answer", json={"query": q, "top_k": 3}, timeout=120)
        rag_time = time.time() - t0
        if resp.status_code == 200:
            data = resp.json()
            conf = data.get("confidence", 0)
            answer = data.get("answer", "")
            sources = len(data.get("sources", []))
            print(f"  RAG: 置信度={conf:.3f}, 来源={sources}, 耗时={rag_time:.2f}s")
            print(f"  回答: {answer[:200]}")
        else:
            print(f"  RAG失败: {resp.status_code}")
    except Exception as e:
        print(f"  RAG异常: {e}")

    # 图谱问答测试
    t0 = time.time()
    try:
        resp = requests.post(f"{PYTHON_URL}/search_and_answer_with_graph", json={"query": q, "top_k": 3}, timeout=120)
        graph_time = time.time() - t0
        if resp.status_code == 200:
            data = resp.json()
            conf = data.get("confidence", 0)
            nodes = data.get("graph_nodes_count", 0)
            rels = data.get("graph_relations_count", 0)
            kps = len(data.get("related_knowledge_points", []))
            answer = data.get("answer", "")
            print(f"  图谱: 置信度={conf:.3f}, 节点={nodes}, 关系={rels}, 知识点={kps}, 耗时={graph_time:.2f}s")
            print(f"  回答: {answer[:200]}")
        else:
            print(f"  图谱失败: {resp.status_code}")
    except Exception as e:
        print(f"  图谱异常: {e}")
