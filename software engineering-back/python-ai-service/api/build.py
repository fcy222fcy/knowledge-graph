from fastapi import APIRouter
from models.schemas import BuildRequest, BuildResponse, GraphNode, GraphEdge
from services.extraction_service import (
    extract_knowledge_points, extract_relations, chunk_text,
    extract_with_llm
)
from services.vector_service import vector_service
from services.answer_service import graph_qa_service
import logging

logger = logging.getLogger(__name__)
router = APIRouter()

@router.post("/build", response_model=BuildResponse)
async def build_graph(req: BuildRequest):
    # 1. 优先使用LLM提取知识点和关系
    llm_result = extract_with_llm(req.content, req.document_id)
    points = llm_result.get("points", [])
    relations = llm_result.get("relations", [])

    # 2. 如果LLM提取失败或结果过少，降级到正则提取
    if len(points) < 3:
        logger.info("LLM extraction insufficient, falling back to regex extraction")
        points = extract_knowledge_points(req.content, req.document_id)
        relations = extract_relations(points)

    # 3. 文本分块并向量化
    chunks = chunk_text(req.content)
    metadata = [{"text": c, "document_id": req.document_id, "document_title": req.title} for c in chunks]
    vector_service.add_texts(chunks, metadata)
    vector_service.save_index()

    # 4. 缓存图谱数据到内存，用于图谱问答
    try:
        logger.info(f"Caching graph for document {req.document_id}: {len(points)} points, {len(relations)} relations")
        graph_qa_service.update_graph_cache(req.document_id, points, relations)
        logger.info(f"Graph cache updated successfully")
    except Exception as e:
        logger.error(f"Failed to update graph cache: {e}")

    # 5. 转换为响应格式返回给 Go 后端，由 Go 统一写入 MySQL 和 Neo4j
    graph_nodes = [GraphNode(id=p["id"], name=p["name"], description=p["description"], category=p["category"], document_id=p["document_id"]) for p in points]
    graph_edges = [GraphEdge(source=r["source_id"], target=r["target_id"], relation_type=r["relation_type"], description=r["description"]) for r in relations]

    return BuildResponse(
        document_id=req.document_id,
        document_title=req.title,
        created_points=len(points),
        created_relations=len(relations),
        chunk_count=len(chunks),
        vector_count=len(chunks),
        status="completed",
        message="知识图谱构建完成",
        points=graph_nodes,
        relations=graph_edges
    )
