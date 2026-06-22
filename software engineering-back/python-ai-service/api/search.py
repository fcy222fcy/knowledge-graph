"""
搜索 API 路由模块

本模块提供知识库的语义搜索和智能问答功能，支持三种模式：
1. 基本语义搜索 - 仅返回相关文档片段
2. 搜索+LLM回答 - 检索相关文档后使用大语言模型生成答案
3. 知识图谱增强搜索 - 结合知识图谱进行更精准的检索和回答
"""

from fastapi import APIRouter, Request
from fastapi.responses import JSONResponse
from models.schemas import SearchRequest, SearchResponse, SearchResult
from services.vector_service import vector_service  # 向量检索服务
from services.answer_service import answer_service, graph_qa_service  # 问答服务
import json
import logging

logger = logging.getLogger(__name__)
router = APIRouter()


@router.post("/search", response_model=SearchResponse)
async def search(request: Request):
    """
    基本语义搜索端点

    使用向量检索服务在知识库中搜索与查询语义最相关的文档片段。

    请求参数:
        query (str): 用户查询文本
        top_k (int): 返回结果数量，默认为 3

    返回:
        SearchResponse: 包含搜索结果列表，每个结果包含：
            - chunk_text: 文档片段文本
            - score: 相似度分数（越高越相关）
            - document_id: 文档ID
            - knowledge_point_ids: 关联的知识点ID列表
    """
    # 解析请求体
    body = await request.body()
    data = json.loads(body.decode('utf-8'))
    query = data.get('query', '')
    top_k = data.get('top_k', 3)

    # 执行向量检索
    results = vector_service.search(query, top_k)

    # 将检索结果转换为响应格式
    search_results = []
    for meta, score in results:
        search_results.append(SearchResult(
            chunk_text=meta.get("text", ""),
            score=score,
            document_id=meta.get("document_id", 0),
            knowledge_point_ids=meta.get("knowledge_point_ids", [])
        ))

    return SearchResponse(results=search_results)


@router.post("/search_and_answer")
async def search_and_answer(request: Request):
    """
    语义检索 + LLM 智能问答端点

    工作流程：
    1. 接收用户查询
    2. 使用向量检索获取相关文档片段
    3. 将检索结果作为上下文发送给大语言模型
    4. LLM 基于上下文生成自然语言回答

    请求参数:
        query (str): 用户查询文本
        top_k (int): 检索返回的文档片段数量，默认为 3
        history (list, optional): 对话历史记录，用于多轮对话

    返回:
        包含 LLM 生成的回答和相关来源

    错误处理:
        - 解码失败：尝试 UTF-8、GBK、Latin-1 三种编码
        - 服务异常：返回 400 错误码和错误信息
    """
    try:
        body = await request.body()

        # 多编码兼容：处理不同编码格式的请求体
        try:
            text = body.decode('utf-8')
        except UnicodeDecodeError:
            try:
                text = body.decode('gbk')
            except UnicodeDecodeError:
                text = body.decode('latin-1')

        data = json.loads(text)
        query = data.get('query', '')
        top_k = data.get('top_k', 3)
        history = data.get('history')

        # 调用问答服务：检索 + LLM 生成回答
        return answer_service.search_and_answer(query, top_k, history)
    except Exception as e:
        logger.error(f"search_and_answer error: {e}")
        return JSONResponse(
            status_code=400,
            content={"error": f"请求处理失败: {str(e)}"}
        )


@router.post("/search_and_answer_with_graph")
async def search_and_answer_with_graph(request: Request):
    """
    知识图谱增强的语义检索 + LLM 智能问答端点

    工作流程：
    1. 接收用户查询
    2. 使用向量检索获取相关文档片段
    3. 从知识图谱中提取相关的实体和关系
    4. 将向量检索结果和知识图谱信息结合作为上下文
    5. LLM 基于丰富的上下文生成更准确的回答

    优势：
    - 知识图谱提供结构化的实体关系信息
    - 能够回答需要多跳推理的复杂问题
    - 回答更具可解释性和准确性

    请求参数:
        query (str): 用户查询文本
        top_k (int): 检索返回的文档片段数量，默认为 3
        history (list, optional): 对话历史记录，用于多轮对话

    返回:
        包含 LLM 生成的回答、相关来源和知识图谱信息

    错误处理:
        - 解码失败：尝试 UTF-8、GBK、Latin-1 三种编码
        - 服务异常：返回 400 错误码和错误信息
    """
    try:
        body = await request.body()

        # 多编码兼容：处理不同编码格式的请求体
        try:
            text = body.decode('utf-8')
        except UnicodeDecodeError:
            try:
                text = body.decode('gbk')
            except UnicodeDecodeError:
                text = body.decode('latin-1')

        data = json.loads(text)
        query = data.get('query', '')
        top_k = data.get('top_k', 3)
        history = data.get('history')

        # 调用知识图谱问答服务：检索 + 图谱增强 + LLM 生成回答
        return graph_qa_service.search_and_answer_with_graph(query, top_k, history)
    except Exception as e:
        logger.error(f"search_and_answer_with_graph error: {e}")
        return JSONResponse(
            status_code=400,
            content={"error": f"请求处理失败: {str(e)}"}
        )
