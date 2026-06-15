from pydantic import BaseModel
from typing import List, Optional

class BuildRequest(BaseModel):
    document_id: int
    title: str
    content: str
    source: Optional[str] = ""

class ChatMessage(BaseModel):
    role: str
    content: str

class SearchRequest(BaseModel):
    query: str
    top_k: int = 3
    history: Optional[List[ChatMessage]] = None

class SearchResult(BaseModel):
    chunk_text: str
    score: float
    document_id: int
    knowledge_point_ids: List[int]

class SearchResponse(BaseModel):
    results: List[SearchResult]

class GraphNode(BaseModel):
    id: int
    name: str
    description: str
    category: str
    document_id: int

class GraphEdge(BaseModel):
    source: int
    target: int
    relation_type: str
    description: str

class GraphResponse(BaseModel):
    nodes: List[GraphNode]
    edges: List[GraphEdge]

class BuildResponse(BaseModel):
    document_id: int
    document_title: str
    created_points: int
    created_relations: int
    chunk_count: int
    vector_count: int
    status: str
    message: str
    points: Optional[List[GraphNode]] = []
    relations: Optional[List[GraphEdge]] = []
