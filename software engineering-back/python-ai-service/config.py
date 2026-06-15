import os
from dataclasses import dataclass

@dataclass
class Config:
    neo4j_uri: str = os.getenv("NEO4J_URI", "bolt://localhost:7687")
    neo4j_user: str = os.getenv("NEO4J_USER", "neo4j")
    neo4j_password: str = os.getenv("NEO4J_PASSWORD", "password")
    vector_index_path: str = os.getenv("VECTOR_INDEX_PATH", "./data/vector_index")
    host: str = os.getenv("AI_SERVICE_HOST", "0.0.0.0")
    port: int = int(os.getenv("AI_SERVICE_PORT", "5000"))
    # Ollama 配置
    ollama_base_url: str = os.getenv("OLLAMA_BASE_URL", "http://localhost:11434")
    ollama_model: str = os.getenv("OLLAMA_MODEL", "qwen3:8b")
    ollama_embedding_model: str = os.getenv("OLLAMA_EMBEDDING_MODEL", "nomic-embed-text")

config = Config()
