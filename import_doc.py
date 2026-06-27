#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
导入软件工程综合知识库到数据库
"""

import mysql.connector
import os

# 数据库配置
DB_CONFIG = {
    'host': '127.0.0.1',
    'port': 3306,
    'user': 'root',
    'password': '123456',
    'database': 'software_qa_platform',
    'charset': 'utf8mb4'
}

# 文档路径
DOC_PATH = './Knowledge/软件工程综合知识库.md'

def import_document():
    """导入文档到数据库"""
    # 读取文件内容
    script_dir = os.path.dirname(os.path.abspath(__file__))
    full_path = os.path.join(script_dir, DOC_PATH)

    if not os.path.exists(full_path):
        print(f"错误: 文件不存在 - {full_path}")
        return False

    with open(full_path, 'r', encoding='utf-8') as f:
        content = f.read()

    print(f"读取文件成功: {len(content)} 字符")

    # 连接数据库
    try:
        conn = mysql.connector.connect(**DB_CONFIG)
        cursor = conn.cursor()

        # 检查是否已存在
        cursor.execute(
            "SELECT id FROM documents WHERE title = %s",
            ('软件工程综合知识库',)
        )
        existing = cursor.fetchone()

        if existing:
            print(f"文档已存在，ID: {existing[0]}，将更新内容")
            cursor.execute(
                """UPDATE documents
                   SET content = %s, file_path = %s, filename = %s, file_size = %s, file_type = 'md', status = 'active', updated_at = NOW()
                   WHERE title = %s""",
                (content, DOC_PATH, '软件工程综合知识库.md', len(content), '软件工程综合知识库')
            )
        else:
            print("插入新文档...")
            cursor.execute(
                """INSERT INTO documents (title, filename, file_path, file_size, file_type, content, status, created_at, updated_at)
                   VALUES (%s, %s, %s, %s, 'md', %s, 'active', NOW(), NOW())""",
                ('软件工程综合知识库', '软件工程综合知识库.md', DOC_PATH, len(content), content)
            )

        conn.commit()
        print("导入成功!")

        # 验证
        cursor.execute(
            "SELECT id, title, LENGTH(content) as size, status FROM documents WHERE title = '软件工程综合知识库'"
        )
        result = cursor.fetchone()
        if result:
            print(f"\n验证结果:")
            print(f"  ID: {result[0]}")
            print(f"  标题: {result[1]}")
            print(f"  内容大小: {result[2]} 字符")
            print(f"  状态: {result[3]}")

        # 显示所有文档
        print(f"\n所有文档:")
        cursor.execute("SELECT id, title, status FROM documents ORDER BY id")
        for row in cursor.fetchall():
            print(f"  {row[0]}. {row[1]} [{row[2]}]")

        return True

    except mysql.connector.Error as e:
        print(f"数据库错误: {e}")
        return False
    finally:
        if 'conn' in locals() and conn.is_connected():
            cursor.close()
            conn.close()

if __name__ == '__main__':
    import_document()
