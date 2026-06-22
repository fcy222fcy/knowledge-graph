#!/usr/bin/env python3
"""
问答系统可信度测试脚本
测试维度：
1. 基本问答功能
2. 知识库检索准确性
3. 知识图谱增强效果
4. 置信度计算合理性
5. 边界情况处理
"""

import requests
import json
import time
from typing import Dict, List, Optional

# 配置
BASE_URL = "http://localhost:8080/api/v1"
PYTHON_AI_URL = "http://localhost:5000"

# 测试用例
TEST_CASES = [
    {
        "name": "基础概念问题",
        "question": "什么是软件工程？",
        "expected_keywords": ["软件工程", "定义", "开发"],
        "category": "basic"
    },
    {
        "name": "流程相关问题",
        "question": "软件开发的生命周期有哪些阶段？",
        "expected_keywords": ["需求", "设计", "编码", "测试", "维护"],
        "category": "process"
    },
    {
        "name": "方法论问题",
        "question": "敏捷开发和瀑布模型有什么区别？",
        "expected_keywords": ["敏捷", "瀑布", "迭代", "文档"],
        "category": "methodology"
    },
    {
        "name": "工具相关问题",
        "question": "Git是什么？如何使用？",
        "expected_keywords": ["版本控制", "Git", "仓库", "提交"],
        "category": "tools"
    },
    {
        "name": "超出知识库范围的问题",
        "question": "今天天气怎么样？",
        "expected_keywords": ["抱歉", "无法找到", "请上传"],
        "category": "out_of_scope"
    },
    {
        "name": "模糊问题测试",
        "question": "软件",
        "expected_keywords": [],  # 可能有多种回答
        "category": "vague"
    },
    {
        "name": "专业术语问题",
        "question": "什么是设计模式中的单例模式？",
        "expected_keywords": ["单例", "设计模式", "实例"],
        "category": "technical"
    },
    {
        "name": "多轮对话上下文",
        "question": "继续刚才的话题，能详细说说吗？",
        "expected_keywords": [],  # 依赖上下文
        "category": "context"
    }
]


class QATestRunner:
    """问答系统测试运行器"""

    def __init__(self):
        self.session = requests.Session()
        self.auth_token = None
        self.test_results = []
        self.conversation_id = None

    def setup_auth(self):
        """设置认证（假设已有测试用户）"""
        # 这里需要根据实际的认证方式调整
        # 如果有测试账号，可以先登录获取token
        try:
            resp = self.session.post(f"{BASE_URL}/auth/login", json={
                "username": "testuser",
                "password": "testpass"
            })
            if resp.status_code == 200:
                data = resp.json()
                self.auth_token = data.get("data", {}).get("token")
                if self.auth_token:
                    self.session.headers.update({"Authorization": f"Bearer {self.auth_token}"})
                    print("[OK] 认证成功")
                    return True
        except Exception as e:
            print(f"[WARN] 认证失败: {e}")

        # 尝试注册测试用户
        try:
            resp = self.session.post(f"{BASE_URL}/auth/register", json={
                "username": "testuser",
                "password": "testpass",
                "email": "test@example.com"
            })
            if resp.status_code == 200:
                print("[OK] 测试用户注册成功")
                return self.setup_auth()
        except Exception as e:
            print(f"[WARN] 注册失败: {e}")

        return False

    def test_python_ai_service(self) -> Dict:
        """测试 Python AI 服务是否可用"""
        print("\n" + "="*60)
        print("[TEST] 测试 Python AI 服务可用性")
        print("="*60)

        result = {"service": "python-ai", "status": "unknown", "details": {}}

        try:
            # 测试搜索服务
            resp = self.session.post(f"{PYTHON_AI_URL}/search", json={
                "query": "软件工程",
                "top_k": 3
            })
            if resp.status_code == 200:
                data = resp.json()
                result["status"] = "available"
                result["details"]["search"] = {
                    "status": "ok",
                    "results_count": len(data.get("results", []))
                }
                print(f"[OK] 搜索服务可用，返回 {len(data.get('results', []))} 条结果")
            else:
                result["status"] = "partial"
                result["details"]["search"] = {"status": "error", "code": resp.status_code}
                print(f"[WARN]  搜索服务异常: {resp.status_code}")
        except Exception as e:
            result["status"] = "unavailable"
            result["details"]["error"] = str(e)
            print(f"[FAIL] 搜索服务不可用: {e}")

        try:
            # 测试问答服务
            resp = self.session.post(f"{PYTHON_AI_URL}/search_and_answer", json={
                "query": "什么是软件工程",
                "top_k": 3
            })
            if resp.status_code == 200:
                data = resp.json()
                result["details"]["answer"] = {
                    "status": "ok",
                    "has_answer": bool(data.get("answer")),
                    "confidence": data.get("confidence", 0)
                }
                print(f"[OK] 问答服务可用，置信度: {data.get('confidence', 0):.2f}")
            else:
                result["details"]["answer"] = {"status": "error", "code": resp.status_code}
                print(f"[WARN]  问答服务异常: {resp.status_code}")
        except Exception as e:
            result["details"]["answer"] = {"status": "error", "error": str(e)}
            print(f"[FAIL] 问答服务调用失败: {e}")

        return result

    def test_search_accuracy(self, test_case: Dict) -> Dict:
        """测试搜索准确性"""
        question = test_case["question"]
        expected = test_case["expected_keywords"]

        result = {
            "question": question,
            "category": test_case["category"],
            "search_results": [],
            "keyword_matches": 0,
            "accuracy_score": 0.0
        }

        try:
            resp = self.session.post(f"{PYTHON_AI_URL}/search", json={
                "query": question,
                "top_k": 5
            })

            if resp.status_code == 200:
                data = resp.json()
                results = data.get("results", [])

                # 提取搜索结果文本
                all_text = " ".join([r.get("chunk_text", "") for r in results])
                result["search_results"] = [
                    {
                        "text": r.get("chunk_text", "")[:100],
                        "score": r.get("score", 0),
                        "doc_id": r.get("document_id", 0)
                    }
                    for r in results[:3]
                ]

                # 检查关键词匹配
                if expected:
                    matches = sum(1 for kw in expected if kw in all_text)
                    result["keyword_matches"] = matches
                    result["accuracy_score"] = matches / len(expected)
                else:
                    result["accuracy_score"] = 1.0  # 没有期望关键词时默认满分

        except Exception as e:
            result["error"] = str(e)

        return result

    def test_qa_confidence(self, test_case: Dict) -> Dict:
        """测试问答置信度"""
        question = test_case["question"]

        result = {
            "question": question,
            "category": test_case["category"],
            "answer_preview": "",
            "confidence": 0.0,
            "sources_count": 0,
            "confidence_valid": False
        }

        try:
            resp = self.session.post(f"{PYTHON_AI_URL}/search_and_answer", json={
                "query": question,
                "top_k": 3
            })

            if resp.status_code == 200:
                data = resp.json()
                result["answer_preview"] = data.get("answer", "")[:200]
                result["confidence"] = data.get("confidence", 0)
                result["sources_count"] = len(data.get("sources", []))

                # 验证置信度范围
                result["confidence_valid"] = 0.0 <= result["confidence"] <= 1.0

        except Exception as e:
            result["error"] = str(e)

        return result

    def test_graph_qa(self, test_case: Dict) -> Dict:
        """测试知识图谱增强问答"""
        question = test_case["question"]

        result = {
            "question": question,
            "category": test_case["category"],
            "answer_preview": "",
            "confidence": 0.0,
            "graph_nodes": 0,
            "graph_relations": 0,
            "related_points": 0
        }

        try:
            resp = self.session.post(f"{PYTHON_AI_URL}/search_and_answer_with_graph", json={
                "query": question,
                "top_k": 3
            })

            if resp.status_code == 200:
                data = resp.json()
                result["answer_preview"] = data.get("answer", "")[:200]
                result["confidence"] = data.get("confidence", 0)
                result["graph_nodes"] = data.get("graph_nodes_count", 0)
                result["graph_relations"] = data.get("graph_relations_count", 0)
                result["related_points"] = len(data.get("related_knowledge_points", []))

        except Exception as e:
            result["error"] = str(e)

        return result

    def test_api_endpoint(self, test_case: Dict) -> Dict:
        """测试后端 API 端点"""
        question = test_case["question"]

        result = {
            "question": question,
            "category": test_case["category"],
            "api_status": "unknown",
            "answer_preview": "",
            "confidence": 0.0,
            "response_time": 0.0
        }

        try:
            start_time = time.time()
            resp = self.session.post(f"{BASE_URL}/ask", json={
                "question": question
            })
            result["response_time"] = time.time() - start_time

            if resp.status_code == 200:
                data = resp.json()
                result["api_status"] = "success"
                result["answer_preview"] = data.get("data", {}).get("answer", "")[:200]
                result["confidence"] = data.get("data", {}).get("confidence", 0)
                result["conversation_id"] = data.get("data", {}).get("conversation_id")
            else:
                result["api_status"] = f"error_{resp.status_code}"
                result["error_detail"] = resp.text[:200]

        except Exception as e:
            result["api_status"] = "exception"
            result["error"] = str(e)

        return result

    def run_all_tests(self):
        """运行所有测试"""
        print("\n" + "="*60)
        print("[START] 开始问答系统可信度测试")
        print("="*60)

        # 1. 测试 AI 服务可用性
        ai_service_result = self.test_python_ai_service()
        self.test_results.append(ai_service_result)

        # 2. 运行测试用例
        print("\n" + "="*60)
        print("[INFO] 运行测试用例")
        print("="*60)

        for i, test_case in enumerate(TEST_CASES, 1):
            print(f"\n--- 测试 {i}/{len(TEST_CASES)}: {test_case['name']} ---")

            # 搜索准确性测试
            search_result = self.test_search_accuracy(test_case)
            print(f"  搜索准确度: {search_result['accuracy_score']:.2%}")

            # 置信度测试
            confidence_result = self.test_qa_confidence(test_case)
            print(f"  置信度: {confidence_result['confidence']:.2f}")
            print(f"  置信度有效: {'[OK]' if confidence_result['confidence_valid'] else '[FAIL]'}")

            # 知识图谱测试
            graph_result = self.test_graph_qa(test_case)
            print(f"  图谱节点: {graph_result['graph_nodes']}, 关系: {graph_result['graph_relations']}")

            # API 测试
            api_result = self.test_api_endpoint(test_case)
            print(f"  API 状态: {api_result['api_status']}")
            print(f"  响应时间: {api_result['response_time']:.2f}s")

            self.test_results.append({
                "test_case": test_case,
                "search": search_result,
                "confidence": confidence_result,
                "graph": graph_result,
                "api": api_result
            })

        # 3. 生成测试报告
        self.generate_report()

    def generate_report(self):
        """生成测试报告"""
        print("\n" + "="*60)
        print("[REPORT] 测试报告")
        print("="*60)

        total_tests = len(TEST_CASES)
        passed_tests = 0
        failed_tests = 0

        for result in self.test_results:
            if isinstance(result, dict) and "test_case" in result:
                # 统计通过/失败
                all_valid = (
                    result["search"]["accuracy_score"] >= 0.5 and
                    result["confidence"]["confidence_valid"] and
                    result["api"]["api_status"] == "success"
                )
                if all_valid:
                    passed_tests += 1
                else:
                    failed_tests += 1

        print(f"\n[STATS] 总体统计:")
        print(f"  测试用例数: {total_tests}")
        print(f"  通过: {passed_tests}")
        print(f"  失败: {failed_tests}")
        print(f"  通过率: {passed_tests/total_tests*100:.1f}%")

        print(f"\n[DETAIL] 详细结果:")
        for result in self.test_results:
            if isinstance(result, dict) and "test_case" in result:
                tc = result["test_case"]
                print(f"\n  [{tc['category']}] {tc['name']}")
                print(f"    搜索准确度: {result['search']['accuracy_score']:.2%}")
                print(f"    置信度: {result['confidence']['confidence']:.2f}")
                print(f"    图谱增强: 节点={result['graph']['graph_nodes']}, 关系={result['graph']['graph_relations']}")
                print(f"    API 响应: {result['api']['api_status']} ({result['api']['response_time']:.2f}s)")

        # 可信度评估
        print(f"\n[ASSESS] 可信度评估:")
        avg_confidence = sum(
            r["confidence"]["confidence"]
            for r in self.test_results
            if isinstance(r, dict) and "confidence" in r
        ) / max(1, len([r for r in self.test_results if isinstance(r, dict) and "confidence" in r]))

        avg_accuracy = sum(
            r["search"]["accuracy_score"]
            for r in self.test_results
            if isinstance(r, dict) and "search" in r
        ) / max(1, len([r for r in self.test_results if isinstance(r, dict) and "search" in r]))

        print(f"  平均置信度: {avg_confidence:.2f}")
        print(f"  平均搜索准确度: {avg_accuracy:.2%}")

        if avg_confidence >= 0.8 and avg_accuracy >= 0.7:
            print("  评估结果: [OK] 高可信度")
        elif avg_confidence >= 0.6 and avg_accuracy >= 0.5:
            print("  评估结果: [WARN]  中等可信度")
        else:
            print("  评估结果: [FAIL] 低可信度，建议优化")

        # 保存报告
        report_path = "qa_system_test_report.json"
        with open(report_path, "w", encoding="utf-8") as f:
            json.dump(self.test_results, f, ensure_ascii=False, indent=2)
        print(f"\n[SAVE] 详细报告已保存到: {report_path}")


if __name__ == "__main__":
    runner = QATestRunner()
    runner.setup_auth()
    runner.run_all_tests()
