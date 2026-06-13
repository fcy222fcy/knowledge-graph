from fastapi import APIRouter, HTTPException
from models.schemas import GraphResponse, GraphNode, GraphEdge
from services.neo4j_service import neo4j_service

router = APIRouter()

@router.get("/graph", response_model=GraphResponse)
async def get_graph():
    if not neo4j_service.is_available():
        raise HTTPException(status_code=503, detail="Neo4j not available")

    nodes_data, edges_data = neo4j_service.get_all_graph_data()

    nodes = [GraphNode(**n) for n in nodes_data]
    edges = [GraphEdge(**e) for e in edges_data]

    return GraphResponse(nodes=nodes, edges=edges)
