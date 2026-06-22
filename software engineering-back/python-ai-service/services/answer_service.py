from services.vector_service import vector_service
from services.neo4j_service import neo4j_service
from config import config
from typing import Dict, List, Optional
import httpx
import json
import asyncio
import logging
import os

logger = logging.getLogger(__name__)

# 图谱缓存文件路径
GRAPH_CACHE_FILE = os.path.join(os.path.dirname(__file__), '..', 'data', 'graph_cache.json')


def _load_graph_cache() -> Dict:
    """从文件加载图谱缓存"""
    try:
        if os.path.exists(GRAPH_CACHE_FILE):
            with open(GRAPH_CACHE_FILE, 'r', encoding='utf-8') as f:
                return json.load(f)
    except Exception as e:
        logger.error(f"Failed to load graph cache: {e}")
    return {}


def _save_graph_cache(cache: Dict):
    """保存图谱缓存到文件"""
    try:
        os.makedirs(os.path.dirname(GRAPH_CACHE_FILE), exist_ok=True)
        with open(GRAPH_CACHE_FILE, 'w', encoding='utf-8') as f:
            json.dump(cache, f, ensure_ascii=False, indent=2)
    except Exception as e:
        logger.error(f"Failed to save graph cache: {e}")


# ──────────────────────────────────────────────
#  统一 LLM 调用：Ollama
# ──────────────────────────────────────────────

async def _call_llm_async(user_prompt: str, system_prompt: str) -> str:
    """
    统一异步 LLM 调用，使用本地 Ollama 模型
    """
    return await _call_ollama_async(user_prompt, system_prompt)


async def _call_ollama_async(user_prompt: str, system_prompt: str) -> str:
    """调用本地 Ollama 生成回答"""
    try:
        async with httpx.AsyncClient(timeout=180.0) as client:
            response = await client.post(
                f"{config.ollama_base_url}/api/generate",
                json={
                    "model": config.ollama_model,
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
            answer = result.get("response", "").strip()
            if not answer:
                answer = result.get("thinking", "").strip()
            return answer or "抱歉，无法生成回答。"
    except httpx.TimeoutException:
        return "抱歉，AI 服务响应超时，请稍后重试。"
    except Exception as e:
        logger.error(f"Ollama call failed: {e}")
        return f"抱歉，AI 服务暂时不可用：{str(e)}"


def _call_llm_sync(user_prompt: str, system_prompt: str, timeout: int = 120) -> str:
    """同步调用 LLM（在子线程中运行异步函数）"""
    result_holder = [None]
    error_holder = [None]

    def _run():
        try:
            result_holder[0] = asyncio.run(
                _call_llm_async(user_prompt, system_prompt)
            )
        except Exception as e:
            error_holder[0] = e

    import threading
    t = threading.Thread(target=_run, daemon=True)
    t.start()
    t.join(timeout=timeout)

    if error_holder[0]:
        raise error_holder[0]
    return result_holder[0] or "抱歉，无法生成回答。"


# ──────────────────────────────────────────────
#  GraphQAService：基于知识图谱的问答
# ──────────────────────────────────────────────

class GraphQAService:
    """基于知识图谱的问答服务"""

    def update_graph_cache(self, document_id: int, points: List[Dict], relations: List[Dict]):
        """更新图谱缓存（写入文件）"""
        cache = _load_graph_cache()
        cache[str(document_id)] = {
            "nodes": points,
            "relations": relations
        }
        _save_graph_cache(cache)
        logger.info(f"Graph cache updated for document {document_id}: {len(points)} nodes, {len(relations)} relations")

    def _search_knowledge_graph(self, query: str) -> Dict:
        """在知识图谱中搜索相关节点和关系"""
        try:
            cache = _load_graph_cache()
            all_nodes = []
            all_relations = []
            for doc_id, data in cache.items():
                all_nodes.extend(data.get("nodes", []))
                all_relations.extend(data.get("relations", []))

            if not all_nodes:
                graph_data = neo4j_service.get_all_graph_data()
                if graph_data:
                    all_nodes = graph_data.get("nodes", [])
                    all_relations = graph_data.get("relations", [])

            if not all_nodes:
                return {"nodes": [], "relations": []}

            nodes = all_nodes
            relations = all_relations

            query_lower = query.lower()
            stop_words = {"什么是", "哪些", "如何", "怎么", "为什么", "请", "介绍", "说明", "解释", "的", "了", "在", "是", "？", "吗", "呢"}
            query_keywords = []
            for i in range(len(query_lower)):
                for length in [4, 3, 2]:
                    if i + length <= len(query_lower):
                        word = query_lower[i:i+length]
                        if word and word not in stop_words and len(word) >= 2:
                            query_keywords.append(word)
                            break

            relevant_nodes = []
            for node in nodes:
                name = node.get("name", "").lower()
                description = node.get("description", "").lower()
                if query_lower in name or query_lower in description or name in query_lower:
                    relevant_nodes.append(node)
                    continue
                match_count = sum(1 for kw in query_keywords if kw in name or kw in description)
                if match_count >= 1:
                    relevant_nodes.append(node)

            relevant_node_ids = {n.get("id") for n in relevant_nodes}
            relevant_relations = []
            for rel in relations:
                source_id = rel.get("source_id") or rel.get("source")
                target_id = rel.get("target_id") or rel.get("target")
                if source_id in relevant_node_ids or target_id in relevant_node_ids:
                    relevant_relations.append(rel)

            if len(relevant_relations) < 3 and relevant_nodes:
                for rel in relations:
                    source_id = rel.get("source_id") or rel.get("source")
                    target_id = rel.get("target_id") or rel.get("target")
                    if source_id in relevant_node_ids:
                        relevant_node_ids.add(target_id)
                        if rel not in relevant_relations:
                            relevant_relations.append(rel)
                    elif target_id in relevant_node_ids:
                        relevant_node_ids.add(source_id)
                        if rel not in relevant_relations:
                            relevant_relations.append(rel)

            all_node_ids = {n.get("id") for n in nodes}
            for nid in relevant_node_ids - {n.get("id") for n in relevant_nodes}:
                for node in nodes:
                    if node.get("id") == nid:
                        relevant_nodes.append(node)
                        break

            return {
                "nodes": relevant_nodes[:20],
                "relations": relevant_relations[:30]
            }
        except Exception as e:
            logger.error(f"Knowledge graph search failed: {e}")
            return {"nodes": [], "relations": []}

    def _build_graph_context(self, graph_data: Dict) -> str:
        """构建图谱上下文用于问答"""
        nodes = graph_data.get("nodes", [])
        relations = graph_data.get("relations", [])
        if not nodes:
            return ""

        node_info = []
        for node in nodes:
            name = node.get("name", "")
            desc = node.get("description", "")
            category = node.get("category", "")
            node_info.append(f"- {name}（{category}）：{desc}")

        relation_info = []
        for rel in relations:
            source_name = rel.get("source_name", "")
            target_name = rel.get("target_name", "")
            rel_type = rel.get("relation_type", "")
            desc = rel.get("description", "")
            relation_info.append(f"- {source_name} --[{rel_type}]--> {target_name}：{desc}")

        context = "知识图谱中的相关概念和关系：\n\n"
        context += "概念：\n" + "\n".join(node_info) + "\n\n"
        if relation_info:
            context += "关系：\n" + "\n".join(relation_info)
        return context

    def search_and_answer_with_graph(
        self,
        query: str,
        top_k: int = 3,
        history: Optional[List[dict]] = None,
    ) -> Dict:
        """基于知识图谱的语义检索和问答"""
        # 1. 向量检索
        search_results = vector_service.search(query, top_k)

        # 2. 知识图谱检索
        graph_data = self._search_knowledge_graph(query)
        graph_context = self._build_graph_context(graph_data)

        # 3. 准备向量检索的上下文
        sources = []
        context_parts = []
        for meta, score in search_results:
            sources.append({
                "document_id": meta.get("document_id", 0),
                "document_title": meta.get("document_title", ""),
                "content": meta.get("text", "")[:200],
            })
            context_parts.append(meta.get("text", ""))

        vector_context = "\n\n".join(context_parts[:3])

        # 4. 构建 prompt
        system_prompt = """你是一个软件工程课程的智能助教。请根据提供的知识库内容和知识图谱信息回答问题。

回答要求：
1. 优先使用知识图谱中的结构化信息（概念和关系）
2. 结合知识库中的详细内容进行补充
3. 如果知识图谱和知识库中都没有相关内容，请说明并给出一般性建议
4. 在回答中明确引用相关的概念和关系
5. 使用中文回答，语言要准确、专业、易于理解"""

        user_prompt = f"问题：{query}\n\n"
        if graph_context:
            user_prompt += f"{graph_context}\n\n"
        if vector_context:
            user_prompt += f"知识库内容：\n{vector_context}\n\n"
        if history:
            recent = history[-10:]
            history_text = "\n".join(
                f"{'用户' if m.get('role') == 'user' else '助手'}：{m.get('content', '')}"
                for m in recent
            )
            user_prompt += f"对话历史：\n{history_text}\n\n"
        user_prompt += "请基于以上信息回答问题。"

        # 5. 调用 LLM（Ollama）
        try:
            answer = _call_llm_sync(user_prompt, system_prompt, timeout=120)
        except Exception as e:
            logger.error(f"Failed to get LLM answer: {e}")
            answer = f"关于「{query}」的回答：\n\n"
            if graph_data.get("nodes"):
                answer += "相关概念：\n"
                for node in graph_data["nodes"][:5]:
                    answer += f"- {node.get('name', '')}：{node.get('description', '')}\n"
            for i, text in enumerate(context_parts[:2], 1):
                answer += f"\n{i}. {text[:150]}..."
            answer += "\n\n以上内容来自知识图谱和知识库检索，仅供参考。"

        # 6. 计算置信度
        avg_score = sum(score for _, score in search_results) / len(search_results) if search_results else 0
        confidence = min(avg_score * 1.2, 1.0)

        # 7. 收集相关知识点
        related_points = []
        for node in graph_data.get("nodes", [])[:5]:
            related_points.append({
                "id": node.get("id", 0),
                "name": node.get("name", ""),
                "description": node.get("description", "")
            })

        return {
            "answer": answer,
            "confidence": confidence,
            "sources": sources,
            "related_knowledge_points": related_points,
            "graph_nodes_count": len(graph_data.get("nodes", [])),
            "graph_relations_count": len(graph_data.get("relations", []))
        }


graph_qa_service = GraphQAService()


# ──────────────────────────────────────────────
#  AnswerService：普通 RAG 问答
# ──────────────────────────────────────────────

class AnswerService:
    def _build_prompt(
        self, query: str, context: str, history: Optional[List[dict]] = None
    ) -> str:
        parts = []
        if history:
            recent = history[-10:]
            history_text = "\n".join(
                f"{'用户' if m.get('role') == 'user' else '助手'}：{m.get('content', '')}"
                for m in recent
            )
            parts.append(f"对话历史：\n{history_text}")
        if context:
            parts.append(f"知识库内容：\n{context}")
        parts.append(f"问题：{query}")
        return "\n\n".join(parts)

    def search_and_answer(
        self,
        query: str,
        top_k: int = 3,
        history: Optional[List[dict]] = None,
    ) -> Dict:
        """语义检索并生成回答"""
        search_results = vector_service.search(query, top_k)

        if not search_results:
            return {
                "answer": f"关于「{query}」：抱歉，未找到相关知识内容。请尝试上传相关文档或调整问题描述。",
                "confidence": 0.0,
                "sources": [],
                "related_knowledge_points": [],
            }

        sources = []
        context_parts = []
        for meta, score in search_results:
            sources.append({
                "document_id": meta.get("document_id", 0),
                "document_title": meta.get("document_title", ""),
                "content": meta.get("text", "")[:200],
            })
            context_parts.append(meta.get("text", ""))

        context = "\n\n".join(context_parts[:3])

        system_prompt = """你是一个软件工程课程的智能助教。请根据提供的知识库内容回答问题。
回答要求：
1. 准确、专业、易于理解
2. 如果知识库中有相关内容，请基于内容回答
3. 如果知识库中没有相关内容，请说明并给出一般性建议
4. 使用中文回答"""

        user_prompt = self._build_prompt(query, context, history)

        try:
            answer = _call_llm_sync(user_prompt, system_prompt, timeout=90)
        except Exception as e:
            logger.error(f"Failed to get LLM answer: {e}")
            answer = f"关于「{query}」的回答：\n\n"
            for i, text in enumerate(context_parts[:2], 1):
                answer += f"{i}. {text[:150]}...\n\n"
            answer += "以上内容来自知识库检索，仅供参考。"

        avg_score = sum(score for _, score in search_results) / len(search_results)
        confidence = min(avg_score * 1.2, 1.0)

        return {
            "answer": answer,
            "confidence": confidence,
            "sources": sources,
            "related_knowledge_points": [],
        }


answer_service = AnswerService()
