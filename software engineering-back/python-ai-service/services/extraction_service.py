import re
import json
import httpx
import asyncio
import logging
from typing import List, Dict, Optional
from collections import defaultdict
from config import config

logger = logging.getLogger(__name__)

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


async def _call_ollama_for_extraction(prompt: str, system_prompt: str) -> str:
    """调用 Ollama 进行知识点提取"""
    try:
        async with httpx.AsyncClient(timeout=180.0) as client:
            response = await client.post(
                f"{config.ollama_base_url}/api/generate",
                json={
                    "model": config.ollama_model,
                    "prompt": prompt,
                    "system": system_prompt,
                    "stream": False,
                    "think": False,
                    "options": {
                        "temperature": 0.3,  # 降低温度以获得更稳定的结构化输出
                        "top_p": 0.9,
                        "num_predict": 2048,
                    },
                },
            )
            response.raise_for_status()
            result = response.json()
            answer = result.get("response", "").strip()
            if not answer:
                answer = result.get("thinking", "").strip()
            return answer
    except Exception as e:
        logger.error(f"Ollama extraction call failed: {e}")
        return ""


def _parse_llm_extraction_response(response: str) -> Dict:
    """解析LLM返回的JSON提取结果"""
    try:
        # 尝试提取JSON块
        json_match = re.search(r'\{[\s\S]*\}', response)
        if json_match:
            json_str = json_match.group(0)
            data = json.loads(json_str)

            # 验证结构
            if "points" in data and "relations" in data:
                return data
    except json.JSONDecodeError as e:
        logger.warning(f"Failed to parse LLM response as JSON: {e}")
    except Exception as e:
        logger.error(f"Error parsing LLM response: {e}")

    return {"points": [], "relations": []}


def extract_with_llm(content: str, document_id: int) -> Dict:
    """使用LLM从文档内容中智能抽取知识点和关系

    Args:
        content: 文档内容
        document_id: 文档ID

    Returns:
        Dict: 包含 points 和 relations 的字典
    """
    # 限制输入长度，避免超出LLM上下文窗口
    max_content_length = 6000
    if len(content) > max_content_length:
        content = content[:max_content_length] + "\n\n...(文档已截断)..."

    system_prompt = """你是一个软件工程领域的知识图谱构建专家。你的任务是从文档中提取核心知识点和它们之间的关系，构建高质量的知识图谱。

提取原则：
1. 知识点应该是核心概念、术语、方法、技术、工具等
2. 每个知识点名称应该简洁准确（2-20个字符）
3. 描述应该清晰说明该知识点的含义（10-100个字符）
4. 分类应该准确反映知识点所属领域
5. 关系应该反映知识点之间的真实语义联系
6. 关系类型必须是以下之一：RELATED, DEPENDS_ON, PART_OF, IS_A, EXAMPLE_OF, USES, IMPLEMENTS

分类选项（category）：
- 软件生命周期：需求分析、设计、编码、测试、维护等阶段
- 软件设计：架构、模式、接口、数据库设计等
- 质量保障：测试方法、质量标准、缺陷管理等
- 过程模型：瀑布、敏捷、Scrum、DevOps等
- 数据库：SQL、NoSQL、数据建模等
- 编程语言：Java、Python、Go等语言和技术
- Web开发：前端、后端、框架、协议等
- AI/ML：机器学习、深度学习、大模型等
- 工具与平台：IDE、版本控制、CI/CD等
- 其他：不属于以上分类的知识点"""

    user_prompt = f"""请从以下软件工程课程文档中提取知识点和关系。

文档内容：
{content}

请严格按照以下JSON格式返回结果（不要添加任何其他文字）：
{{
  "points": [
    {{
      "name": "知识点名称",
      "description": "简要描述该知识点",
      "category": "分类"
    }}
  ],
  "relations": [
    {{
      "source": "源知识点名称",
      "target": "目标知识点名称",
      "relation_type": "关系类型",
      "description": "关系描述"
    }}
  ]
}}

注意：
- 提取10-30个核心知识点
- 提取5-20个有意义的关系
- 关系类型必须是：RELATED, DEPENDS_ON, PART_OF, IS_A, EXAMPLE_OF, USES, IMPLEMENTS
- 确保source和target都是提取出的知识点名称"""

    try:
        import threading

        result = [None]

        def _run():
            result[0] = asyncio.run(_call_ollama_for_extraction(user_prompt, system_prompt))

        t = threading.Thread(target=_run)
        t.start()
        t.join(timeout=120)

        if result[0]:
            data = _parse_llm_extraction_response(result[0])

            # 为每个知识点添加ID和document_id
            for i, point in enumerate(data.get("points", []), 1):
                point["id"] = i
                point["document_id"] = document_id
                # 确保必填字段存在
                if "name" not in point or not point["name"]:
                    point["name"] = f"未知知识点{i}"
                if "description" not in point:
                    point["description"] = f"{point['name']}相关概念"
                if "category" not in point:
                    point["category"] = "其他"

            # 为关系添加ID并验证节点存在
            point_names = {p["name"] for p in data.get("points", [])}
            valid_relations = []
            for i, rel in enumerate(data.get("relations", []), 1):
                if (rel.get("source") in point_names and
                    rel.get("target") in point_names and
                    rel.get("relation_type") in ["RELATED", "DEPENDS_ON", "PART_OF",
                                                   "IS_A", "EXAMPLE_OF", "USES", "IMPLEMENTS"]):
                    rel["id"] = i
                    valid_relations.append(rel)

            data["relations"] = valid_relations

            logger.info(f"LLM extraction: {len(data['points'])} points, {len(data['relations'])} relations")
            return data

    except Exception as e:
        logger.error(f"LLM extraction failed: {e}")

    return {"points": [], "relations": []}


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
        # 英文技术术语（支持多单词组合如 Spring Boot）
        (r'\b([A-Z][a-z]+(?:\s[A-Z][a-z]+)*)\b', "技术术语"),
        # 缩写（2-6个大写字母）
        (r'\b([A-Z]{2,6})\b', "缩写"),
        # 中文术语 + 后缀（更宽泛的匹配）
        (r'([一-龥]{2,10}(?:模型|模式|系统|框架|工具|技术|方法|测试|设计|开发|架构|接口|数据库|服务器|容器|平台|服务|协议|标准|规范|流程|工程))', "技术概念"),
        # 中文术语 + 动词/名词后缀
        (r'([一-龥]{2,8}(?:分析|管理|优化|集成|部署|配置|实现|验证|评估|监控|维护|演进|迭代|重构|抽象|封装|继承|多态))', "技术概念"),
        # 带版本号的术语
        (r'([A-Za-z]+\s*\d+(?:\.\d+)*)', "版本号"),
        # 点分隔的技术术语（如 Spring Boot、Node.js）
        (r'\b([A-Z][a-z]+(?:\.[A-Za-z][a-z]+)+)\b', "技术术语"),
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

    # 1. 同类别内的关系（基于文档顺序和语义）
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

    # 3. 基于名称相似性的关系（包含关系）
    for i, p1 in enumerate(points):
        for j, p2 in enumerate(points):
            if i >= j:
                continue
            key = (p1["id"], p2["id"])
            if key in seen:
                continue

            # 检查包含关系
            if p1["name"] in p2["name"] or p2["name"] in p1["name"]:
                # 短的被长的包含
                if len(p1["name"]) < len(p2["name"]):
                    source, target = p1, p2
                else:
                    source, target = p2, p1

                relations.append({
                    "id": rel_id,
                    "source_id": source["id"],
                    "target_id": target["id"],
                    "relation_type": "PART_OF",
                    "description": f"{source['name']} 是 {target['name']} 的一部分"
                })
                rel_id += 1
                seen.add(key)

    # 4. 同一文档内的相关关系（如果关系太少）
    if len(relations) < len(points) // 2:
        for i, p1 in enumerate(points):
            for j, p2 in enumerate(points):
                if i >= j:
                    continue
                key = (p1["id"], p2["id"])
                if key in seen:
                    continue

                # 同一类别内的节点建立RELATED关系
                if p1["category"] == p2["category"]:
                    relations.append({
                        "id": rel_id,
                        "source_id": p1["id"],
                        "target_id": p2["id"],
                        "relation_type": "RELATED",
                        "description": f"{p1['name']} 与 {p2['name']} 相关"
                    })
                    rel_id += 1
                    seen.add(key)

                    # 限制每个节点最多3个关系
                    source_count = sum(1 for r in relations if r["source_id"] == p1["id"])
                    if source_count >= 3:
                        break

    return relations


def _infer_relation_type(name1: str, name2: str) -> str:
    """推断两个术语之间的关系类型"""
    # 演化关系（DEPENDS_ON）
    evolution_pairs = [
        ("需求分析", "概要设计"), ("概要设计", "详细设计"),
        ("瀑布模型", "增量模型"), ("增量模型", "敏捷开发"),
        ("单元测试", "集成测试"), ("集成测试", "系统测试"),
        ("设计", "编码"), ("编码", "测试"), ("测试", "部署")
    ]
    for src, tgt in evolution_pairs:
        if src in name1 and tgt in name2:
            return "DEPENDS_ON"

    # 包含关系（PART_OF）
    if len(name1) > len(name2) and name2 in name1:
        return "PART_OF"
    if len(name2) > len(name1) and name1 in name2:
        return "PART_OF"

    # 使用关系（USES）
    usage_keywords = ["使用", "采用", "基于", "利用"]
    for kw in usage_keywords:
        if kw in name1 or kw in name2:
            return "USES"

    # 实现关系（IMPLEMENTS）
    impl_keywords = ["实现", "执行", "运行"]
    for kw in impl_keywords:
        if kw in name1 or kw in name2:
            return "IMPLEMENTS"

    # 示例关系（EXAMPLE_OF）
    example_keywords = ["示例", "例子", "案例", "实例"]
    for kw in example_keywords:
        if kw in name1 or kw in name2:
            return "EXAMPLE_OF"

    # 默认关系（RELATED）
    return "RELATED"


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
