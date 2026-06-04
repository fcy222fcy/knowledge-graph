from fastapi import APIRouter
from models.schemas import SearchRequest, SearchResponse, SearchResult
from services.vector_service import vector_service

router = APIRouter()

@router.post("/search", response_model=SearchResponse)
async def search(req: SearchRequest):
    results = vector_service.search(req.query, req.top_k)

    search_results = []
    for meta, score in results:
        search_results.append(SearchResult(
            chunk_text=meta.get("text", ""),
            score=score,
            document_id=meta.get("document_id", 0),
            knowledge_point_ids=meta.get("knowledge_point_ids", [])
        ))

    return SearchResponse(results=search_results)
