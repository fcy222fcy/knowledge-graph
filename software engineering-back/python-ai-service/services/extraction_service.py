import re
from typing import List, Tuple, Dict

# 软件工程课程术语词典
KNOWLEDGE_TERMS = {
    "需求相关": ["需求分析", "需求规格", "需求验证", "需求管理", "用例", "用户故事", "功能需求", "非功能需求"],
    "设计相关": ["架构设计", "详细设计", "模块设计", "接口设计", "数据库设计", "UI设计"],
    "编码相关": ["编码规范", "代码审查", "单元测试", "集成测试", "调试"],
    "测试相关": ["测试计划", "测试用例", "黑盒测试", "白盒测试", "回归测试", "性能测试"],
    "项目管理": ["项目计划", "进度管理", "风险管理", "质量管理", "配置管理"],
    "过程模型": ["瀑布模型", "增量模型", "螺旋模型", "敏捷开发", "Scrum", "看板"],
}

def extract_knowledge_points(content: str, document_id: int) -> List[Dict]:
    """从文档内容中抽取知识点"""
    points = []
    point_id = 1

    for category, terms in KNOWLEDGE_TERMS.items():
        for term in terms:
            if term in content:
                # 查找术语所在的句子
                sentences = re.split(r'[。！？\n]', content)
                description = ""
                for sent in sentences:
                    if term in sent:
                        description = sent.strip()[:200]
                        break

                points.append({
                    "id": point_id,
                    "name": term,
                    "description": description or f"文档中提到的{term}",
                    "category": category,
                    "document_id": document_id
                })
                point_id += 1

    return points

def extract_relations(points: List[Dict]) -> List[Dict]:
    """基于共现和类别推断关系"""
    relations = []
    rel_id = 1

    for i, p1 in enumerate(points):
        for j, p2 in enumerate(points):
            if i >= j:
                continue

            # 同类别知识点建立 RELATED 关系
            if p1["category"] == p2["category"]:
                relations.append({
                    "id": rel_id,
                    "source_id": p1["id"],
                    "target_id": p2["id"],
                    "relation_type": "RELATED",
                    "description": f"{p1['name']}与{p2['name']}相关"
                })
                rel_id += 1

            # 不同类别但有依赖关系的
            elif p1["category"] == "需求相关" and p2["category"] == "设计相关":
                relations.append({
                    "id": rel_id,
                    "source_id": p1["id"],
                    "target_id": p2["id"],
                    "relation_type": "DEPENDS_ON",
                    "description": f"{p2['name']}依赖于{p1['name']}"
                })
                rel_id += 1

    return relations

def chunk_text(content: str, chunk_size: int = 500) -> List[str]:
    """将文本分块"""
    chunks = []
    sentences = re.split(r'[。！？\n]', content)
    current_chunk = ""

    for sent in sentences:
        sent = sent.strip()
        if not sent:
            continue
        if len(current_chunk) + len(sent) > chunk_size:
            if current_chunk:
                chunks.append(current_chunk)
            current_chunk = sent
        else:
            current_chunk += "。" + sent if current_chunk else sent

    if current_chunk:
        chunks.append(current_chunk)

    return chunks
