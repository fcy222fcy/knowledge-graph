import re
from typing import List, Dict
from collections import defaultdict

# 软件工程领域知识本体 - 用于分类和关系推断
ONTOLOGY = {
    "软件生命周期": {
        "keywords": ["需求分析", "概要设计", "详细设计", "编码", "测试", "维护", "部署", "运维"],
        "relations": [("需求分析", "概要设计", "演化为"), ("概要设计", "详细设计", "细化为"), ("详细设计", "编码", "实现为")]
    },
    "软件设计": {
        "keywords": ["架构设计", "模块设计", "接口设计", "数据库设计", "UI设计", "设计模式", "MVC", "MVP", "MVVM"],
        "relations": [("架构设计", "模块设计", "分解为"), ("模块设计", "接口设计", "定义接口")]
    },
    "质量保障": {
        "keywords": ["软件测试", "单元测试", "集成测试", "系统测试", "验收测试", "回归测试", "黑盒测试", "白盒测试"],
        "relations": [("单元测试", "集成测试", "组合为"), ("集成测试", "系统测试", "扩展为")]
    },
    "过程模型": {
        "keywords": ["瀑布模型", "增量模型", "螺旋模型", "敏捷开发", "Scrum", "看板", "迭代开发"],
        "relations": [("瀑布模型", "增量模型", "改进为"), ("增量模型", "敏捷开发", "演变为")]
    },
    "数据库": {
        "keywords": ["SQL", "MySQL", "PostgreSQL", "Oracle", "关系型数据库", "NoSQL", "Redis", "MongoDB"],
        "relations": [("关系型数据库", "SQL", "使用"), ("Redis", "缓存", "用于")]
    },
    "编程语言": {
        "keywords": ["Java", "Python", "Go", "JavaScript", "TypeScript", "C++", "Rust"],
        "relations": []
    },
    "Web开发": {
        "keywords": ["HTTP", "HTTPS", "REST", "API", "前端", "后端", "全栈", "Docker", "Nginx", "Tomcat"],
        "relations": [("前端", "后端", "调用"), ("Docker", "Nginx", "部署")]
    },
    "AI/ML": {
        "keywords": ["机器学习", "深度学习", "神经网络", "LLM", "大语言模型", "RAG", "Agent", "向量数据库"],
        "relations": [("机器学习", "深度学习", "包含"), ("LLM", "RAG", "结合")]
    }
}

# 关系类型定义
RELATION_TYPES = {
    "演化": "实体A随时间或技术发展演变为实体B",
    "包含": "实体A包含实体B作为子组件",
    "使用": "实体A使用实体B来实现功能",
    "依赖": "实体A依赖实体B才能运行",
    "测试": "实体A用于测试实体B",
    "部署": "实体A用于部署实体B",
    "关联": "实体A和实体B在概念上相关"
}


def extract_knowledge_points(content: str, document_id: int) -> List[Dict]:
    """从文档内容中智能抽取知识点 - 基于文档结构和内容分析"""
    points = []
    point_id = 1
    seen_names = set()

    # 按段落和标题分割文档
    sections = re.split(r'\n#{1,3}\s+|\n\n+|[。！？]', content)

    for section in sections:
        section = section.strip()
        if len(section) < 5:
            continue

        # 尝试从每个段落提取实体
        extracted = _extract_entities_from_text(section)
        for name, category, description in extracted:
            if name not in seen_names and len(name) >= 2:
                seen_names.add(name)
                points.append({
                    "id": point_id,
                    "name": name,
                    "description": description[:200] if description else f"{name}相关概念",
                    "category": category,
                    "document_id": document_id
                })
                point_id += 1

    # 如果提取的点太少，使用标题作为知识点
    if len(points) < 3:
        titles = re.findall(r'#{1,3}\s+(.+)', content)
        for title in titles:
            title = title.strip()
            if title and title not in seen_names and len(title) >= 2:
                seen_names.add(title)
                category = _classify_term(title)
                points.append({
                    "id": point_id,
                    "name": title,
                    "description": f"文档标题：{title}",
                    "category": category,
                    "document_id": document_id
                })
                point_id += 1

    return points[:50]  # 限制最多50个知识点


def _extract_entities_from_text(text: str) -> List[tuple]:
    """从文本中提取实体（名称，类别，描述）"""
    entities = []

    # 匹配技术术语模式
    patterns = [
        # 英文技术术语
        (r'\b([A-Z][a-z]+(?:\s[A-Z][a-z]+)*)\b', "技术术语"),
        # 缩写
        (r'\b([A-Z]{2,6})\b', "缩写"),
        # 中文术语（2-8个字）
        (r'([一-龥]{2,8}(?:模型|模式|系统|框架|工具|技术|方法|测试|设计|开发|架构|接口|数据库|服务器|容器))', "技术概念"),
        # 带版本号的术语
        (r'([A-Za-z]+\s*\d+(?:\.\d+)*)', "版本号"),
    ]

    for pattern, category in patterns:
        matches = re.findall(pattern, text)
        for match in matches:
            if len(match) >= 2:
                # 获取匹配词所在的句子作为描述
                description = _get_context_sentence(text, match)
                entities.append((match, category, description))

    # 检查本体中的关键词
    for cat_name, cat_info in ONTOLOGY.items():
        for keyword in cat_info["keywords"]:
            if keyword in text:
                description = _get_context_sentence(text, keyword)
                entities.append((keyword, cat_name, description))

    return entities


def _get_context_sentence(text: str, term: str) -> str:
    """获取术语所在的上下文句子"""
    sentences = re.split(r'[。！？\n]', text)
    for sent in sentences:
        if term in sent:
            return sent.strip()[:200]
    return ""


def _classify_term(term: str) -> str:
    """根据本体对术语进行分类"""
    term_lower = term.lower()
    for cat_name, cat_info in ONTOLOGY.items():
        for keyword in cat_info["keywords"]:
            if keyword in term or term in keyword:
                return cat_name
    return "其他"


def extract_relations(points: List[Dict]) -> List[Dict]:
    """基于语义和本体提取有意义的关系"""
    relations = []
    rel_id = 1
    seen = set()

    # 按类别分组
    category_points = defaultdict(list)
    for p in points:
        category_points[p["category"]].append(p)

    # 1. 同类别内的链式关系（基于文档顺序）
    for cat, cat_pts in category_points.items():
        if len(cat_pts) <= 1:
            continue
        # 每个节点最多与前后2个节点建立关系
        for i, p in enumerate(cat_pts):
            # 与前一个节点
            if i > 0:
                prev = cat_pts[i - 1]
                key = (prev["id"], p["id"])
                if key not in seen:
                    seen.add(key)
                    rel_type = _infer_relation_type(prev["name"], p["name"])
                    relations.append({
                        "id": rel_id,
                        "source_id": prev["id"],
                        "target_id": p["id"],
                        "relation_type": rel_type,
                        "description": f"{prev['name']} {rel_type} {p['name']}"
                    })
                    rel_id += 1

    # 2. 跨类别的语义关系（基于本体定义）
    for cat_name, cat_info in ONTOLOGY.items():
        for src_name, tgt_name, rel_desc in cat_info.get("relations", []):
            src_pts = [p for p in points if src_name in p["name"]]
            tgt_pts = [p for p in points if tgt_name in p["name"]]
            if src_pts and tgt_pts:
                key = (src_pts[0]["id"], tgt_pts[0]["id"])
                if key not in seen:
                    seen.add(key)
                    relations.append({
                        "id": rel_id,
                        "source_id": src_pts[0]["id"],
                        "target_id": tgt_pts[0]["id"],
                        "relation_type": rel_desc,
                        "description": f"{src_pts[0]['name']} {rel_desc} {tgt_pts[0]['name']}"
                    })
                    rel_id += 1

    return relations


def _infer_relation_type(name1: str, name2: str) -> str:
    """推断两个术语之间的关系类型"""
    # 演化关系
    evolution_pairs = [
        ("需求分析", "概要设计"), ("概要设计", "详细设计"),
        ("瀑布模型", "增量模型"), ("增量模型", "敏捷开发"),
        ("单元测试", "集成测试"), ("集成测试", "系统测试")
    ]
    for src, tgt in evolution_pairs:
        if src in name1 and tgt in name2:
            return "演化"

    # 包含关系
    if len(name1) > len(name2) and name2 in name1:
        return "包含"
    if len(name2) > len(name1) and name1 in name2:
        return "包含"

    # 默认关系
    return "关联"


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
