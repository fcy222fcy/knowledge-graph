from services.vector_service import vector_service
from services.neo4j_service import neo4j_service
from typing import Dict, List

class AnswerService:
    def search_and_answer(self, query: str, top_k: int = 3) -> Dict:
        """语义检索并生成回答"""
        # 1. 向量检索
        search_results = vector_service.search(query, top_k)

        # 2. 组装回答
        if not search_results:
            return {
                "answer": f"关于「{query}」：抱歉，未找到相关知识内容。请尝试上传相关文档或调整问题描述。",
                "confidence": 0.0,
                "sources": [],
                "related_knowledge_points": []
            }

        # 3. 生成回答
        sources = []
        all_text = []
        for meta, score in search_results:
            sources.append({
                "document_id": meta.get("document_id", 0),
                "document_title": meta.get("document_title", ""),
                "content": meta.get("text", "")[:200]
            })
            all_text.append(meta.get("text", ""))

        # 简单的回答生成（可后续替换为LLM）
        answer = f"关于「{query}」的回答：\n\n"
        for i, text in enumerate(all_text[:2], 1):
            answer += f"{i}. {text[:150]}...\n\n"

        answer += "以上内容来自知识库检索，仅供参考。"

        # 4. 计算置信度
        avg_score = sum(score for _, score in search_results) / len(search_results)
        confidence = min(avg_score * 1.2, 1.0)

        return {
            "answer": answer,
            "confidence": confidence,
            "sources": sources,
            "related_knowledge_points": []
        }

answer_service = AnswerService()
