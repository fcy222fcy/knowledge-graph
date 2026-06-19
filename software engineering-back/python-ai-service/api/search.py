from fastapi import APIRouter, Request
from fastapi.responses import JSONResponse
from models.schemas import SearchRequest, SearchResponse, SearchResult
from services.vector_service import vector_service
from services.answer_service import answer_service, graph_qa_service
import json
import logging

logger = logging.getLogger(__name__)
router = APIRouter()

def _fix_encoding(query: str) -> str:
    """尝试修复编码问题：如果包含 GBK 编码的字节，尝试解码"""
    try:
        # 尝试将字符串编码为 latin-1（保持字节不变），然后解码为 GBK
        raw_bytes = query.encode('latin-1')
        decoded = raw_bytes.decode('gbk')
        logger.info(f"Fixed encoding: {query} -> {decoded}")
        return decoded
    except (UnicodeDecodeError, UnicodeEncodeError):
        return query

def _parse_request(body: bytes) -> dict:
    """解析请求体，处理 GBK/UTF-8 混合编码"""
    logger.info(f"Raw body bytes: {body[:100]}")

    # 直接从原始字节中提取 query 字段的字节
    query_start = body.find(b'"query":"') + len(b'"query":"')
    query_end = body.find(b'"', query_start)
    query_bytes = body[query_start:query_end]

    # 过滤掉 UTF-8 替换字符 (0xEF 0xBF 0xBD)，只保留有效的 GBK 字节
    filtered_bytes = bytearray()
    i = 0
    while i < len(query_bytes):
        # 检查是否是 UTF-8 替换字符
        if i + 2 < len(query_bytes) and query_bytes[i:i+3] == b'\xef\xbf\xbd':
            i += 3  # 跳过替换字符
            continue
        filtered_bytes.append(query_bytes[i])
        i += 1

    # 用 GBK 解码过滤后的字节
    try:
        query_str = bytes(filtered_bytes).decode('gbk')
        logger.info(f"Decoded query as GBK (filtered): {query_str}")
    except UnicodeDecodeError:
        query_str = bytes(filtered_bytes).decode('latin-1')
        logger.info(f"Decoded query as latin-1 (filtered): {query_str}")

    # 解析完整 JSON
    try:
        data = json.loads(body.decode('utf-8'))
    except UnicodeDecodeError:
        data = json.loads(body.decode('latin-1'))

    data['query'] = query_str
    return data

@router.post("/search", response_model=SearchResponse)
async def search(request: Request):
    body = await request.body()
    data = _parse_request(body)
    query = _fix_encoding(data.get('query', ''))
    top_k = data.get('top_k', 3)

    results = vector_service.search(query, top_k)

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
    """语义检索并使用 Ollama 生成智能回答"""
    body = await request.body()
    data = _parse_request(body)
    query = _fix_encoding(data.get('query', ''))
    top_k = data.get('top_k', 3)
    history = data.get('history')

    return answer_service.search_and_answer(query, top_k, history)

@router.post("/search_and_answer_with_graph")
async def search_and_answer_with_graph(request: Request):
    """基于知识图谱的语义检索和智能回答"""
    body = await request.body()
    data = _parse_request(body)
    query = _fix_encoding(data.get('query', ''))
    top_k = data.get('top_k', 3)
    history = data.get('history')

    return graph_qa_service.search_and_answer_with_graph(query, top_k, history)
