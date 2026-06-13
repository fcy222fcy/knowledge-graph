from fastapi import APIRouter
from models.schemas import BuildRequest, BuildResponse, GraphNode, GraphEdge
from services.extraction_service import extract_knowledge_points, extract_relations, chunk_text
from services.vector_service import vector_service

router = APIRouter()

@router.post("/build", response_model=BuildResponse)
async def build_graph(req: BuildRequest):
    # 1. 抽取知识点（只提取，不写入 Neo4j，由 Go 后端统一持久化）
    points = extract_knowledge_points(req.content, req.document_id)

    # 2. 抽取关系
    relations = extract_relations(points)

    # 3. 文本分块并向量化
    chunks = chunk_text(req.content)
    metadata = [{"text": c, "document_id": req.document_id, "document_title": req.title} for c in chunks]
    vector_service.add_texts(chunks, metadata)
    vector_service.save_index()

    # 4. 转换为响应格式返回给 Go 后端，由 Go 统一写入 MySQL 和 Neo4j
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
