import numpy as np
import hashlib
import json
import os
from typing import List, Tuple
from config import config

class VectorService:
    def __init__(self):
        self.dimension = 128
        self.vectors = []
        self.metadata = []
        self.index_path = config.vector_index_path
        os.makedirs(os.path.dirname(self.index_path), exist_ok=True)
        self.load_index()

    def _text_to_vector(self, text: str) -> np.ndarray:
        """将文本转换为向量（使用哈希方法，可后续替换为sentence-transformers）"""
        hash_obj = hashlib.md5(text.encode())
        hash_bytes = hash_obj.digest()

        # 扩展到所需维度
        extended = hash_bytes * (self.dimension // len(hash_bytes) + 1)
        vector = np.frombuffer(extended[:self.dimension], dtype=np.uint8).astype(np.float32)
        vector = vector / np.linalg.norm(vector)  # 归一化
        return vector

    def add_texts(self, texts: List[str], metadata: List[dict]):
        """添加文本到索引"""
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
            "vectors": [v.tolist() for v in self.vectors],
            "metadata": self.metadata
        }
        with open(self.index_path + ".json", "w") as f:
            json.dump(data, f)

    def load_index(self):
        """从磁盘加载索引"""
        path = self.index_path + ".json"
        if os.path.exists(path):
            with open(path, "r") as f:
                data = json.load(f)
                self.vectors = [np.array(v) for v in data["vectors"]]
                self.metadata = data["metadata"]

vector_service = VectorService()
