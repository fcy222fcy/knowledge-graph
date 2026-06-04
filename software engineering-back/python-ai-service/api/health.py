from fastapi import APIRouter
from services.neo4j_service import neo4j_service

router = APIRouter()

@router.get("/health")
async def health_check():
    return {
        "status": "ok",
        "neo4j": "connected" if neo4j_service.is_available() else "disconnected"
    }
