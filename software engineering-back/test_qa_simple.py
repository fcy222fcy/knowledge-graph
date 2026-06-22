#!/usr/bin/env python3
"""问答系统快速测试"""
import requests
import json

PYTHON_URL = "http://localhost:5000"
BACKEND_URL = "http://localhost:8080/api/v1"

def test_search():
    print("[1] 测试搜索服务...")
    try:
        resp = requests.post(f"{PYTHON_URL}/search", json={"query": "软件工程", "top_k": 3}, timeout=30)
        if resp.status_code == 200:
            data = resp.json()
            results = data.get("results", [])
            print(f"    搜索返回 {len(results)} 条结果")
            for i, r in enumerate(results[:2]):
                print(f"    [{i+1}] score={r.get('score',0):.3f} text={r.get('chunk_text','')[:80]}...")
            return True
        else:
            print(f"    搜索失败: {resp.status_code}")
    except Exception as e:
        print(f"    搜索异常: {e}")
    return False

def test_rag_answer():
    print("[2] 测试 RAG 问答...")
    try:
        resp = requests.post(f"{PYTHON_URL}/search_and_answer", json={"query": "什么是软件工程", "top_k": 3}, timeout=120)
        if resp.status_code == 200:
            data = resp.json()
            print(f"    置信度: {data.get('confidence', 0):.3f}")
            print(f"    来源数: {len(data.get('sources', []))}")
            answer = data.get("answer", "")
            print(f"    回答预览: {answer[:150]}...")
            return True
        else:
            print(f"    问答失败: {resp.status_code}")
    except Exception as e:
        print(f"    问答异常: {e}")
    return False

def test_graph_answer():
    print("[3] 测试知识图谱问答...")
    try:
        resp = requests.post(f"{PYTHON_URL}/search_and_answer_with_graph", json={"query": "什么是软件工程", "top_k": 3}, timeout=120)
        if resp.status_code == 200:
            data = resp.json()
            print(f"    置信度: {data.get('confidence', 0):.3f}")
            print(f"    图谱节点: {data.get('graph_nodes_count', 0)}")
            print(f"    图谱关系: {data.get('graph_relations_count', 0)}")
            print(f"    相关知识点: {len(data.get('related_knowledge_points', []))}")
            answer = data.get("answer", "")
            print(f"    回答预览: {answer[:150]}...")
            return True
        else:
            print(f"    图谱问答失败: {resp.status_code}")
    except Exception as e:
        print(f"    图谱问答异常: {e}")
    return False

def test_backend_api():
    print("[4] 测试后端 API...")
    try:
        resp = requests.post(f"{BACKEND_URL}/ask", json={"question": "什么是软件工程"}, timeout=120)
        if resp.status_code == 200:
            data = resp.json()
            d = data.get("data", {})
            print(f"    置信度: {d.get('confidence', 0):.3f}")
            print(f"    来源数: {len(d.get('sources', []))}")
            answer = d.get("answer", "")
            print(f"    回答预览: {answer[:150]}...")
            return True
        else:
            print(f"    API 失败: {resp.status_code} {resp.text[:100]}")
    except Exception as e:
        print(f"    API 异常: {e}")
    return False

if __name__ == "__main__":
    print("=" * 60)
    print("  问答系统可信度快速测试")
    print("=" * 60)

    results = {}
    results["search"] = test_search()
    results["rag"] = test_rag_answer()
    results["graph"] = test_graph_answer()
    results["api"] = test_backend_api()

    print("\n" + "=" * 60)
    print("  测试总结")
    print("=" * 60)
    for k, v in results.items():
        status = "[PASS]" if v else "[FAIL]"
        print(f"  {k:12s} {status}")
    passed = sum(1 for v in results.values() if v)
    total = len(results)
    print(f"\n  通过率: {passed}/{total} ({passed/total*100:.0f}%)")
