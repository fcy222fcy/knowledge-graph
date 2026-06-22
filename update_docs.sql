-- 更新文档3: SQL 语法示例
UPDATE documents SET
content = '# SQL 语法示例

## 目录

1. DDL（数据定义语言）
2. DML（数据操作语言）
3. DQL（数据查询语言）
4. DCL（数据控制语言）

---

## 1. DDL（数据定义语言）

DDL用于定义数据库结构，主要包括：

### 1.1 数据库操作

创建数据库: CREATE DATABASE database_name;
查看数据库: SHOW DATABASES;
使用数据库: USE database_name;
删除数据库: DROP DATABASE database_name;

### 1.2 表操作

创建表: CREATE TABLE table_name (column1 datatype constraints, column2 datatype constraints);
查看表: SHOW TABLES;
查看表结构: DESC table_name;
修改表: ALTER TABLE table_name ADD column_name datatype;
删除表: DROP TABLE table_name;

### 1.3 数据类型

- INT: 整数类型
- VARCHAR: 变长字符串
- TEXT: 长文本
- DATE: 日期类型
- DATETIME: 日期时间类型
- DECIMAL: 精确数值类型

---

## 2. DML（数据操作语言）

DML用于操作数据，主要包括：

### 2.1 插入数据

插入单条记录: INSERT INTO table_name (column1, column2) VALUES (val1, val2);
插入多条记录: INSERT INTO table_name (column1, column2) VALUES (val1, val2), (val3, val4);

### 2.2 更新数据

更新记录: UPDATE table_name SET column1 = val1 WHERE condition;
注意：不加WHERE会更新所有记录！

### 2.3 删除数据

删除记录: DELETE FROM table_name WHERE condition;
清空表: TRUNCATE TABLE table_name;

---

## 3. DQL（数据查询语言）

DQL用于查询数据，是最常用的操作。

### 3.1 基本查询

查询所有列: SELECT * FROM table_name;
查询指定列: SELECT column1, column2 FROM table_name;
使用别名: SELECT column1 AS alias1 FROM table_name;

### 3.2 条件查询

WHERE子句: SELECT * FROM table_name WHERE condition;
比较运算符: SELECT * FROM table_name WHERE column1 > 10;
逻辑运算符: SELECT * FROM table_name WHERE column1 > 10 AND column2 < 20;
IN运算符: SELECT * FROM table_name WHERE column1 IN (val1, val2);
LIKE运算符: SELECT * FROM table_name WHERE column1 LIKE '%pattern%';

### 3.3 排序和分页

排序: SELECT * FROM table_name ORDER BY column1 ASC;
分页: SELECT * FROM table_name LIMIT 10;

### 3.4 分组查询

分组: SELECT column1, COUNT(*) FROM table_name GROUP BY column1;
HAVING子句: SELECT column1, COUNT(*) as cnt FROM table_name GROUP BY column1 HAVING cnt > 5;

---

## 4. DCL（数据控制语言）

DCL用于控制数据库访问权限。

### 4.1 用户管理

创建用户: CREATE USER 'username'@'host' IDENTIFIED BY 'password';
删除用户: DROP USER 'username'@'host';

### 4.2 权限管理

授予权限: GRANT privilege ON database.table TO 'username'@'host';
撤销权限: REVOKE privilege ON database.table FROM 'username'@'host';

---

## 总结

- DDL：定义数据库结构（CREATE, ALTER, DROP）
- DML：操作数据（INSERT, UPDATE, DELETE）
- DQL：查询数据（SELECT）
- DCL：控制权限（GRANT, REVOKE）

掌握这些SQL语法是数据库操作的基础。
',
file_size = 4500
WHERE id = 3;

-- 更新文档4: Redis 示例资料
UPDATE documents SET
content = '# Redis 示例资料

## 目录

1. Redis 简介
2. Redis 数据类型
3. Redis 常用命令
4. Redis 应用场景

---

## 1. Redis 简介

### 1.1 什么是 Redis

Redis（Remote Dictionary Server）是一个开源的、基于内存的键值对数据库。

### 1.2 Redis 特点

- 高性能：数据存储在内存中，读写速度快
- 丰富的数据结构：支持字符串、列表、集合等多种数据类型
- 持久化：支持 RDB 和 AOF 两种持久化方式
- 原子性操作：所有操作都是原子性的

### 1.3 Redis 适用场景

- 缓存
- 会话管理
- 消息队列
- 实时排行榜
- 分布式锁

---

## 2. Redis 数据类型

### 2.1 String（字符串）

设置值: SET key value
获取值: GET key
设置过期时间: SETEX key seconds value
自增: INCR key

### 2.2 List（列表）

左侧插入: LPUSH key value
右侧插入: RPUSH key value
获取范围: LRANGE key start stop
弹出左侧元素: LPOP key

### 2.3 Set（集合）

添加元素: SADD key value
获取所有元素: SMEMBERS key
判断元素是否存在: SISMEMBER key value
交集: SINTER key1 key2

### 2.4 Hash（哈希）

设置字段值: HSET key field value
获取字段值: HGET key field
获取所有字段: HGETALL key

### 2.5 Sorted Set（有序集合）

添加元素: ZADD key score value
获取范围: ZRANGE key start stop
获取分数: ZSCORE key value

---

## 3. Redis 常用命令

### 3.1 连接命令

连接 Redis: redis-cli
测试连接: PING

### 3.2 键命令

查看所有键: KEYS *
删除键: DEL key
判断键是否存在: EXISTS key
设置过期时间: EXPIRE key seconds
查看键的类型: TYPE key

### 3.3 服务器命令

查看信息: INFO
清空数据库: FLUSHDB
选择数据库: SELECT index

---

## 4. Redis 应用场景

### 4.1 缓存

设置缓存: SET user:1001 '{"name":"张三","age":25}' EX 3600
获取缓存: GET user:1001

### 4.2 会话管理

存储会话: SET session:abc123 '{"user_id":1001}' EX 1800
获取会话: GET session:abc123

### 4.3 计数器

文章阅读量: INCR article:1001:views
获取阅读量: GET article:1001:views

### 4.4 排行榜

更新分数: ZADD leaderboard 100 player1
获取前10名: ZREVRANGE leaderboard 0 9 WITHSCORES

---

## 总结

Redis 是一个功能强大的内存数据库，适用于各种需要高性能读写的场景。掌握 Redis 的基本数据类型和命令是使用 Redis 的基础。
',
file_size = 3500
WHERE id = 4;

-- 更新文档5: AI 名词解释示例
UPDATE documents SET
content = '# AI 名词解释示例

## 目录

1. 基础概念
2. 机器学习
3. 深度学习
4. 自然语言处理
5. 计算机视觉

---

## 1. 基础概念

### 1.1 人工智能（AI）

人工智能是计算机科学的一个分支，致力于创建能够模拟人类智能的系统。

### 1.2 机器学习（ML）

机器学习是AI的一个子领域，使计算机能够从数据中学习，而无需显式编程。

### 1.3 深度学习（DL）

深度学习是机器学习的一个子领域，使用多层神经网络来学习数据的层次化表示。

---

## 2. 机器学习

### 2.1 监督学习

使用标记的训练数据来学习输入到输出的映射。

- 分类：预测离散标签（如垃圾邮件检测）
- 回归：预测连续值（如房价预测）

### 2.2 无监督学习

从未标记的数据中发现隐藏的模式。

- 聚类：将数据分组（如客户细分）
- 降维：减少数据维度（如可视化）

### 2.3 强化学习

通过与环境交互来学习最优策略。

- 探索：尝试新动作
- 利用：使用已知的最佳动作

---

## 3. 深度学习

### 3.1 神经网络

由多层节点（神经元）组成的网络结构。

- 输入层：接收原始数据
- 隐藏层：处理数据
- 输出层：产生结果

### 3.2 卷积神经网络（CNN）

专门用于处理网格结构数据（如图像）的神经网络。

应用：
- 图像分类
- 目标检测
- 图像分割

### 3.3 循环神经网络（RNN）

专门用于处理序列数据的神经网络。

应用：
- 语音识别
- 自然语言处理
- 时间序列预测

---

## 4. 自然语言处理（NLP）

### 4.1 什么是 NLP

NLP 是 AI 的一个分支，专注于计算机和人类语言之间的交互。

### 4.2 常见 NLP 任务

- 文本分类：将文本分配到预定义类别
- 情感分析：判断文本的情感倾向
- 命名实体识别：识别文本中的实体
- 机器翻译：将文本从一种语言翻译成另一种语言
- 问答系统：根据问题生成答案

### 4.3 预训练模型

- BERT：双向编码器表示
- GPT：生成式预训练变换器
- T5：文本到文本转换器

---

## 5. 计算机视觉（CV）

### 5.1 什么是 CV

CV 是 AI 的一个分支，专注于让计算机理解和处理视觉信息。

### 5.2 常见 CV 任务

- 图像分类：识别图像中的主要对象
- 目标检测：定位并识别图像中的多个对象
- 语义分割：为图像中的每个像素分配标签
- 实例分割：区分同一类别的不同实例

### 5.3 常用模型

- ResNet：残差网络
- YOLO：实时目标检测
- U-Net：医学图像分割

---

## 总结

AI 是一个快速发展的领域，涵盖了从基础的机器学习到复杂的深度学习等多个子领域。了解这些基本概念对于理解和应用 AI 技术至关重要。
',
file_size = 3000
WHERE id = 5;
