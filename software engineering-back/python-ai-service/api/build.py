from fastapi import APIRouter
from models.schemas import BuildRequest, BuildResponse
from services.extraction_service import extract_knowledge_points, extract_relations, chunk_text
from services.neo4j_service import neo4j_service
from services.vector_service import vector_service

router = APIRouter()

@router.post("/build", response_model=BuildResponse)
async def build_graph(req: BuildRequest):
    # 1. 抽取知识点
    points = extract_knowledge_points(req.content, req.document_id)

    # 2. 抽取关系
    relations = extract_relations(points)

    # 3. 写入 Neo4j
    for p in points:
        neo4j_service.create_knowledge_point(
            mysql_id=p["id"],
            name=p["name"],
            description=p["description"],
            category=p["category"],
            document_id=p["document_id"]
        )

    for r in relations:
        neo4j_service.create_relation(
            source_id=r["source_id"],
            target_id=r["target_id"],
            relation_type=r["relation_type"],
            description=r["description"]
        )

    # 4. 文本分块并向量化
    chunks = chunk_text(req.content)
    metadata = [{"text": c, "document_id": req.document_id, "document_title": req.title} for c in chunks]
    vector_service.add_texts(chunks, metadata)
    vector_service.save_index()

    return BuildResponse(
        document_id=req.document_id,
        document_title=req.title,
        created_points=len(points),
        created_relations=len(relations),
        chunk_count=len(chunks),
        vector_count=len(chunks),
        status="completed",
        message="知识图谱构建完成"
    )
