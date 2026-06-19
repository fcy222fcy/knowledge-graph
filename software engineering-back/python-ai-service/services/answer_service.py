from services.vector_service import vector_service
from services.neo4j_service import neo4j_service
from config import config
from typing import Dict, List, Optional
import httpx
import json
import asyncio
import concurrent.futures
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


class GraphQAService:
    """基于知识图谱的问答服务"""

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
                answer = result.get("response", "").strip()
                if not answer:
                    answer = result.get("thinking", "").strip()
                return answer or "抱歉，无法生成回答。"
        except httpx.TimeoutException:
            return "抱歉，AI 服务响应超时，请稍后重试。"
        except Exception as e:
            logger.error(f"Ollama call failed: {e}")
            return f"抱歉，AI 服务暂时不可用：{str(e)}"

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
            # 优先从文件缓存获取
            cache = _load_graph_cache()
            logger.info(f"Loaded graph cache: {len(cache)} documents")
            all_nodes = []
            all_relations = []
            for doc_id, data in cache.items():
                nodes = data.get("nodes", [])
                relations = data.get("relations", [])
                logger.info(f"Document {doc_id}: {len(nodes)} nodes, {len(relations)} relations")
                all_nodes.extend(nodes)
                all_relations.extend(relations)

            # 如果缓存为空，尝试从Neo4j获取
            if not all_nodes:
                graph_data = neo4j_service.get_all_graph_data()
                if graph_data:
                    all_nodes = graph_data.get("nodes", [])
                    all_relations = graph_data.get("relations", [])

            if not all_nodes:
                return {"nodes": [], "relations": []}

            nodes = all_nodes
            relations = all_relations

            # 尝试修复编码问题：如果包含 GBK 编码的字节，尝试解码
            try:
                raw_bytes = query.encode('latin-1')
                query = raw_bytes.decode('gbk')
                logger.info(f"Decoded GBK query: {query}")
            except (UnicodeDecodeError, UnicodeEncodeError):
                pass

            # 搜索与查询相关的节点
            query_lower = query.lower()
            # 移除常见停用词
            stop_words = {"什么是", "哪些", "如何", "怎么", "为什么", "请", "介绍", "说明", "解释", "的", "了", "在", "是", "？", "吗", "呢"}
            # 按字符分割，而不是按字节分割
            query_keywords = []
            for i in range(len(query_lower)):
                # 尝试提取2-4字的关键词
                for length in [4, 3, 2]:
                    if i + length <= len(query_lower):
                        word = query_lower[i:i+length]
                        if word and word not in stop_words and len(word) >= 2:
                            query_keywords.append(word)
                            break

            logger.info(f"Search query: {repr(query_lower)}, keywords: {query_keywords}")

            relevant_nodes = []
            for node in nodes:
                name = node.get("name", "").lower()
                description = node.get("description", "").lower()

                # 精确匹配
                if query_lower in name or query_lower in description or name in query_lower:
                    relevant_nodes.append(node)
                    continue

                # 关键词匹配（至少匹配1个关键词）
                match_count = sum(1 for kw in query_keywords if kw in name or kw in description)
                if match_count >= 1:
                    relevant_nodes.append(node)

            logger.info(f"Found {len(relevant_nodes)} relevant nodes")

            # 获取相关节点的关系
            relevant_node_ids = {n.get("id") for n in relevant_nodes}
            relevant_relations = []
            for rel in relations:
                source_id = rel.get("source_id") or rel.get("source")
                target_id = rel.get("target_id") or rel.get("target")
                if source_id in relevant_node_ids or target_id in relevant_node_ids:
                    relevant_relations.append(rel)

            # 如果关系太少，扩展到2-hop邻居
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

            # 更新节点列表（包含新发现的邻居）
            all_node_ids = {n.get("id") for n in nodes}
            for nid in relevant_node_ids - {n.get("id") for n in relevant_nodes}:
                for node in nodes:
                    if node.get("id") == nid:
                        relevant_nodes.append(node)
                        break

            return {
                "nodes": relevant_nodes[:20],  # 限制节点数量
                "relations": relevant_relations[:30]  # 限制关系数量
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

        # 构建节点信息
        node_info = []
        for node in nodes:
            name = node.get("name", "")
            desc = node.get("description", "")
            category = node.get("category", "")
            node_info.append(f"- {name}（{category}）：{desc}")

        # 构建关系信息
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
        logger.info(f"Graph search for '{query}': {len(graph_data.get('nodes', []))} nodes, {len(graph_data.get('relations', []))} relations")
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

        # 4. 组合上下文（向量检索 + 图谱信息）
        vector_context = "\n\n".join(context_parts[:3])

        # 5. 构建prompt
        system_prompt = """你是一个软件工程课程的智能助教。请根据提供的知识库内容和知识图谱信息回答问题。

回答要求：
1. 优先使用知识图谱中的结构化信息（概念和关系）
2. 结合知识库中的详细内容进行补充
3. 如果知识图谱和知识库中都没有相关内容，请说明并给出一般性建议
4. 在回答中明确引用相关的概念和关系
5. 使用中文回答，语言要准确、专业、易于理解"""

        user_prompt = f"""问题：{query}

"""
        if graph_context:
            user_prompt += f"""{graph_context}

"""
        if vector_context:
            user_prompt += f"""知识库内容：
{vector_context}

"""
        if history:
            recent = history[-10:]
            history_text = "\n".join(
                f"{'用户' if m.get('role') == 'user' else '助手'}：{m.get('content', '')}"
                for m in recent
            )
            user_prompt += f"""对话历史：
{history_text}

"""
        user_prompt += "请基于以上信息回答问题。"

        # 6. 调用LLM生成答案
        try:
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
            t.join(timeout=90)

            if error[0]:
                raise error[0]
            answer = result[0] or "抱歉，无法生成回答。"
        except Exception as e:
            logger.error(f"Failed to get LLM answer: {e}")
            # 回退到简单回答
            answer = f"关于「{query}」的回答：\n\n"
            if graph_data.get("nodes"):
                answer += "相关概念：\n"
                for node in graph_data["nodes"][:5]:
                    answer += f"- {node.get('name', '')}：{node.get('description', '')}\n"
            for i, text in enumerate(context_parts[:2], 1):
                answer += f"\n{i}. {text[:150]}..."
            answer += "\n\n以上内容来自知识图谱和知识库检索，仅供参考。"

        # 7. 计算置信度
        avg_score = sum(score for _, score in search_results) / len(search_results) if search_results else 0
        confidence = min(avg_score * 1.2, 1.0)

        # 8. 收集相关知识点
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
