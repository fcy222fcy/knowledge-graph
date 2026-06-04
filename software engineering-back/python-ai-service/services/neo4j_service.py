from neo4j import GraphDatabase
from config import config
from typing import List, Tuple, Optional

class Neo4jService:
    def __init__(self):
        self.driver = None
        try:
            self.driver = GraphDatabase.driver(
                config.neo4j_uri,
                auth=(config.neo4j_user, config.neo4j_password)
            )
            self.driver.verify_connectivity()
            print("Neo4j connected successfully")
        except Exception as e:
            print(f"Warning: Neo4j connection failed: {e}")

    def close(self):
        if self.driver:
            self.driver.close()

    def is_available(self) -> bool:
        return self.driver is not None

    def create_knowledge_point(self, mysql_id: int, name: str, description: str, category: str, document_id: int):
        if not self.is_available():
            return
        with self.driver.session() as session:
            session.run(
                "CREATE (n:KnowledgePoint {id: $id, name: $name, description: $description, category: $category, document_id: $document_id})",
                id=mysql_id, name=name, description=description, category=category, document_id=document_id
            )

    def create_relation(self, source_id: int, target_id: int, relation_type: str, description: str):
        if not self.is_available():
            return
        with self.driver.session() as session:
            query = f"MATCH (a:KnowledgePoint {{id: $source_id}}), (b:KnowledgePoint {{id: $target_id}}) CREATE (a)-[r:{relation_type} {{description: $description}}]->(b)"
            session.run(query, source_id=source_id, target_id=target_id, description=description)

    def get_all_graph_data(self) -> Tuple[List[dict], List[dict]]:
        if not self.is_available():
            return [], []
        with self.driver.session() as session:
            nodes_result = session.run("MATCH (n:KnowledgePoint) RETURN n.id AS id, n.name AS name, n.description AS description, n.category AS category, n.document_id AS document_id")
            nodes = [dict(record) for record in nodes_result]

            edges_result = session.run("MATCH (a:KnowledgePoint)-[r]->(b:KnowledgePoint) RETURN a.id AS source, b.id AS target, type(r) AS relation_type, r.description AS description")
            edges = [dict(record) for record in edges_result]

            return nodes, edges

neo4j_service = Neo4jService()
