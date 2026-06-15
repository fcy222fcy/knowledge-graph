from services.vector_service import vector_service
from services.neo4j_service import neo4j_service
from config import config
from typing import Dict, List, Optional
import httpx
import json
import asyncio
import concurrent.futures
import logging

logger = logging.getLogger(__name__)


class AnswerService:
    def __init__(self):
        self.ollama_url = f"{config.ollama_base_url}/api/generate"
        self.model = config.ollama_model

    async def _call_ollama(self, user_prompt: str, system_prompt: str) -> str:
        """调用 Ollama 生成回答"""
        try:
            async with httpx.AsyncClient(timeout=180.0) as client:
                response = await client.post(
                    self.ollama_url,
                    json={
                        "model": self.model,
                        "prompt": user_prompt,
                        "system": system_prompt,
                        "stream": False,
                        "think": False,
                        "options": {
                            "temperature": 0.7,
                            "top_p": 0.9,
                            "num_predict": 1024,
                        },
                    },
                )
                response.raise_for_status()
                result = response.json()
                # qwen3 模型默认开启 thinking，response 可能为空，回退到 thinking 字段
                answer = result.get("response", "").strip()
                if not answer:
                    answer = result.get("thinking", "").strip()
                return answer or "抱歉，无法生成回答。"
        except httpx.TimeoutException:
            return "抱歉，AI 服务响应超时，请稍后重试。"
        except Exception as e:
            logger.error(f"Ollama call failed: {e}")
            return f"抱歉，AI 服务暂时不可用：{str(e)}"

    def _build_prompt(
        self, query: str, context: str, history: Optional[List[dict]] = None
    ) -> str:
        """构建包含历史上下文的 prompt"""
        parts = []

        # 历史对话（最近 5 轮）
        if history:
            recent = history[-10:]  # 最多 5 轮（10 条消息）
            history_text = "\n".join(
                f"{'用户' if m.get('role') == 'user' else '助手'}：{m.get('content', '')}"
                for m in recent
            )
            parts.append(f"对话历史：\n{history_text}")

        # 知识库内容
        if context:
            parts.append(f"知识库内容：\n{context}")

        # 当前问题
        parts.append(f"问题：{query}")

        return "\n\n".join(parts)

    def search_and_answer(
        self,
        query: str,
        top_k: int = 3,
        history: Optional[List[dict]] = None,
    ) -> Dict:
        """语义检索并生成回答"""
        # 1. 向量检索
        search_results = vector_service.search(query, top_k)

        # 2. 组装回答
        if not search_results:
            return {
                "answer": f"关于「{query}」：抱歉，未找到相关知识内容。请尝试上传相关文档或调整问题描述。",
                "confidence": 0.0,
                "sources": [],
                "related_knowledge_points": [],
            }

        # 3. 准备上下文
        sources = []
        context_parts = []
        for meta, score in search_results:
            sources.append(
                {
                    "document_id": meta.get("document_id", 0),
                    "document_title": meta.get("document_title", ""),
                    "content": meta.get("text", "")[:200],
                }
            )
            context_parts.append(meta.get("text", ""))

        context = "\n\n".join(context_parts[:3])

        # 4. 构建 prompt 并调用 Ollama
        system_prompt = """你是一个软件工程课程的智能助教。请根据提供的知识库内容回答问题。
回答要求：
1. 准确、专业、易于理解
2. 如果知识库中有相关内容，请基于内容回答
3. 如果知识库中没有相关内容，请说明并给出一般性建议
4. 使用中文回答"""

        user_prompt = self._build_prompt(query, context, history)

        try:
            # 在新线程中运行异步调用
            import threading

            result = [None]
            error = [None]

            def _run():
                try:
                    result[0] = asyncio.run(
                        self._call_ollama(user_prompt, system_prompt)
                    )
                except Exception as e:
                    error[0] = e

            t = threading.Thread(target=_run)
            t.start()
            t.join(timeout=60)

            if error[0]:
                raise error[0]
            answer = result[0] or "抱歉，无法生成回答。"
        except Exception as e:
            logger.error(f"Failed to get LLM answer: {e}")
            # 回退到简单回答
            answer = f"关于「{query}」的回答：\n\n"
            for i, text in enumerate(context_parts[:2], 1):
                answer += f"{i}. {text[:150]}...\n\n"
            answer += "以上内容来自知识库检索，仅供参考。"

        # 5. 计算置信度
        avg_score = sum(score for _, score in search_results) / len(search_results)
        confidence = min(avg_score * 1.2, 1.0)

        return {
            "answer": answer,
            "confidence": confidence,
            "sources": sources,
            "related_knowledge_points": [],
        }


answer_service = AnswerService()
