import numpy as np
import hashlib
import json
import os
import logging
from typing import List, Tuple
from config import config

logger = logging.getLogger(__name__)


class VectorService:
    def __init__(self):
        self.dimension = 768  # nomic-embed-text 输出维度
        self.vectors = []
        self.metadata = []
        self.index_path = config.vector_index_path
        self.ollama_url = f"{config.ollama_base_url}/api/embeddings"
        self.embedding_model = config.ollama_embedding_model
        self._use_ollama = True  # 尝试 Ollama，失败则降级到 MD5
        os.makedirs(os.path.dirname(self.index_path), exist_ok=True)
        self.load_index()

    def _embed_with_ollama(self, text: str) -> np.ndarray:
        """调用 Ollama embedding API"""
        import httpx
        resp = httpx.post(
            self.ollama_url,
            json={"model": self.embedding_model, "prompt": text},
            timeout=30.0,
        )
        resp.raise_for_status()
        embedding = resp.json()["embedding"]
        vec = np.array(embedding, dtype=np.float32)
        vec = vec / np.linalg.norm(vec)
        return vec

    def _embed_with_md5(self, text: str) -> np.ndarray:
        """MD5 降级方案（Ollama 不可用时）"""
        hash_obj = hashlib.md5(text.encode())
        hash_bytes = hash_obj.digest()
        extended = hash_bytes * (self.dimension // len(hash_bytes) + 1)
        vector = np.frombuffer(extended[:self.dimension], dtype=np.uint8).astype(np.float32)
        vector = vector / np.linalg.norm(vector)
        return vector

    def _text_to_vector(self, text: str) -> np.ndarray:
        """将文本转换为向量"""
        if self._use_ollama:
            try:
                return self._embed_with_ollama(text)
            except Exception as e:
                logger.warning(f"Ollama embedding failed, falling back to MD5: {e}")
                self._use_ollama = False
        return self._embed_with_md5(text)

    def add_texts(self, texts: List[str], metadata: List[dict]):
        """添加文本到索引（逐条 embedding）"""
        for text, meta in zip(texts, metadata):
            vector = self._text_to_vector(text)
            self.vectors.append(vector)
            self.metadata.append(meta)

    def search(self, query: str, top_k: int = 3) -> List[Tuple[dict, float]]:
        """搜索最相似的文本"""
        if not self.vectors:
            return []

        query_vector = self._text_to_vector(query)
        vectors_array = np.array(self.vectors)

        # 计算余弦相似度
        similarities = np.dot(vectors_array, query_vector)

        # 获取 top_k 结果
        top_indices = np.argsort(similarities)[::-1][:top_k]

        results = []
        for idx in top_indices:
            results.append((self.metadata[idx], float(similarities[idx])))

        return results

    def save_index(self):
        """保存索引到磁盘"""
        data = {
            "dimension": self.dimension,
            "vectors": [v.tolist() for v in self.vectors],
            "metadata": self.metadata,
        }
        with open(self.index_path + ".json", "w") as f:
            json.dump(data, f)

    def load_index(self):
        """从磁盘加载索引"""
        path = self.index_path + ".json"
        if os.path.exists(path):
            with open(path, "r") as f:
                data = json.load(f)
                # 检查维度是否匹配，不匹配则清空重建
                saved_dim = data.get("dimension", 128)
                if saved_dim != self.dimension:
                    logger.warning(
                        f"Vector dimension mismatch: index={saved_dim}, expected={self.dimension}. "
                        "Clearing index for rebuild."
                    )
                    self.vectors = []
                    self.metadata = []
                    return
                self.vectors = [np.array(v) for v in data["vectors"]]
                self.metadata = data["metadata"]


vector_service = VectorService()
