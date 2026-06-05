package seed

import (
	"log"
	"software_engineering/pkg/database"
	"software_engineering/internal/model/entity"
	"software_engineering/internal/repository"
)

func SeedAll() {
	seedUsers()
	seedKnowledgePoints()
	seedQuestions()
	seedDocuments()
}

func seedUsers() {
	var count int64
	database.DB.Model(&entity.User{}).Count(&count)
	if count > 0 {
		return
	}

	users := []entity.User{
		{Username: "student001", Password: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", Email: "student001@example.com", Nickname: "张三", Status: 1},
		{Username: "student002", Password: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", Email: "student002@example.com", Nickname: "李四", Status: 1},
		{Username: "teacher001", Password: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", Email: "teacher001@example.com", Nickname: "王老师", Status: 1},
	}

	if err := database.DB.Create(&users).Error; err != nil {
		log.Printf("seed users failed: %v", err)
		return
	}
	log.Println("seeded 3 demo users")
}

func seedKnowledgePoints() {
	var count int64
	database.DB.Model(&entity.KnowledgePoint{}).Count(&count)
	if count > 0 {
		return
	}

	points := []entity.KnowledgePoint{
		{Name: "需求分析", Description: "识别和确认用户需求的过程", Category: "需求相关"},
		{Name: "软件测试", Description: "验证软件是否满足需求的过程", Category: "测试相关"},
		{Name: "软件生命周期", Description: "软件从提出到废弃的整个过程", Category: "基础概念"},
		{Name: "编码实现", Description: "将设计转化为可执行代码的过程", Category: "开发相关"},
		{Name: "项目管理", Description: "对软件项目进行计划、组织和控制", Category: "管理相关"},
	}

	if err := database.DB.Create(&points).Error; err != nil {
		log.Printf("seed knowledge points failed: %v", err)
		return
	}

	// Sync knowledge points to Neo4j
	for i := range points {
		if err := repository.CreateKnowledgePointInNeo4j(&points[i]); err != nil {
			log.Printf("warning: neo4j seed knowledge point %q failed: %v", points[i].Name, err)
		}
	}

	relations := []entity.KnowledgeRelation{
		{SourceID: 1, TargetID: 2, RelationType: "DEPENDS_ON", Description: "需求分析是软件测试的前置环节"},
		{SourceID: 1, TargetID: 4, RelationType: "DEPENDS_ON", Description: "需求分析完成后进入编码实现"},
		{SourceID: 3, TargetID: 1, RelationType: "PART_OF", Description: "需求分析是软件生命周期的一个阶段"},
		{SourceID: 5, TargetID: 3, RelationType: "RELATED", Description: "项目管理贯穿整个软件生命周期"},
	}

	if err := database.DB.Create(&relations).Error; err != nil {
		log.Printf("seed knowledge relations failed: %v", err)
		return
	}

	// Sync relations to Neo4j
	for i := range relations {
		if err := repository.CreateRelationInNeo4j(&relations[i]); err != nil {
			log.Printf("warning: neo4j seed relation %d->%d failed: %v", relations[i].SourceID, relations[i].TargetID, err)
		}
	}

	log.Println("seeded 5 knowledge points and 4 relations")
}

func seedQuestions() {
	var count int64
	database.DB.Model(&entity.Question{}).Count(&count)
	if count > 0 {
		return
	}

	questions := []entity.Question{
		{
			Title:            "以下哪个不是需求分析的活动？",
			Type:             "single",
			Difficulty:       "easy",
			KnowledgePointID: 1,
			Options:          `[{"key":"A","value":"需求获取"},{"key":"B","value":"需求分析"},{"key":"C","value":"代码编写"},{"key":"D","value":"需求验证"}]`,
			Answer:           "C",
			Explanation:      "代码编写属于编码阶段，不是需求分析的活动",
		},
		{
			Title:            "黑盒测试的主要关注点是什么？",
			Type:             "single",
			Difficulty:       "easy",
			KnowledgePointID: 2,
			Options:          `[{"key":"A","value":"程序内部逻辑"},{"key":"B","value":"程序外部功能"},{"key":"C","value":"代码覆盖率"},{"key":"D","value":"算法效率"}]`,
			Answer:           "B",
			Explanation:      "黑盒测试关注程序的外部功能，不关心内部实现细节",
		},
		{
			Title:            "软件生命周期的第一个阶段是？",
			Type:             "single",
			Difficulty:       "easy",
			KnowledgePointID: 3,
			Options:          `[{"key":"A","value":"编码"},{"key":"B","value":"测试"},{"key":"C","value":"需求分析"},{"key":"D","value":"维护"}]`,
			Answer:           "C",
			Explanation:      "软件生命周期的第一个阶段是需求分析",
		},
	}

	if err := database.DB.Create(&questions).Error; err != nil {
		log.Printf("seed questions failed: %v", err)
		return
	}
	log.Println("seeded 3 demo questions")
}

func seedDocuments() {
	var count int64
	database.DB.Model(&entity.Document{}).Count(&count)
	if count > 0 {
		return
	}

	// 示例文档内容 - 软件工程基础知识
	documents := []entity.Document{
		{
			UserID:      1, // student001
			Title:       "软件工程基础概念",
			Description: "软件工程的基本概念和原理介绍",
			Filename:    "software_engineering_basics.md",
			FilePath:    "./uploads/software_engineering_basics.md",
			FileSize:    0,
			FileType:    ".md",
			Content: `# 软件工程基础概念

## 什么是软件工程？

软件工程是一门研究用工程化方法构建和维护有效的、实用的和高质量的软件的学科。它涉及程序设计语言、数据库、软件开发工具、系统平台、标准、设计模式等方面。

## 软件生命周期

软件生命周期是指软件从提出到废弃的整个过程，通常包括以下阶段：

1. **需求分析阶段**：识别和确认用户需求
2. **设计阶段**：根据需求设计软件架构和详细设计
3. **编码实现阶段**：将设计转化为可执行代码
4. **测试阶段**：验证软件是否满足需求
5. **部署阶段**：将软件部署到生产环境
6. **维护阶段**：对软件进行改进和修复

## 需求分析

需求分析是软件开发的第一步，主要包括：

- **需求获取**：通过访谈、问卷等方式收集用户需求
- **需求分析**：分析需求的可行性和完整性
- **需求规格说明**：编写需求规格说明书
- **需求验证**：确认需求是否正确反映用户意图

## 软件测试

软件测试是验证软件是否满足需求的过程，主要方法包括：

- **黑盒测试**：关注程序的外部功能，不关心内部实现
- **白盒测试**：关注程序的内部逻辑结构
- **单元测试**：测试单个模块的功能
- **集成测试**：测试模块间的交互
- **系统测试**：测试整个系统的功能
- **验收测试**：用户验证软件是否满足需求

## 项目管理

软件项目管理是对软件项目进行计划、组织和控制的过程，主要包括：

- 项目计划制定
- 进度管理
- 成本管理
- 质量管理
- 风险管理
- 团队管理`,
			Status: "completed",
		},
		{
			UserID:      1, // student001
			Title:       "软件开发方法论",
			Description: "常见的软件开发方法论介绍",
			Filename:    "development_methodologies.md",
			FilePath:    "./uploads/development_methodologies.md",
			FileSize:    0,
			FileType:    ".md",
			Content: `# 软件开发方法论

## 瀑布模型

瀑布模型是一种线性的软件开发方法，按照以下顺序进行：

1. 需求分析
2. 系统设计
3. 实现
4. 测试
5. 部署
6. 维护

**优点**：简单易懂，易于管理
**缺点**：不适应需求变化，风险后置

## 敏捷开发

敏捷开发是一种迭代式的软件开发方法，强调：

- 个体和互动 高于 流程和工具
- 工作的软件 高于 详尽的文档
- 客户合作 高于 合同谈判
- 响应变化 高于 遵循计划

常见框架：Scrum、XP、Kanban

## DevOps

DevOps是开发和运维的结合，目标是：

- 缩短开发周期
- 提高部署频率
- 提高系统稳定性
- 快速修复问题

实践：持续集成、持续部署、基础设施即代码`,
			Status: "completed",
		},
		{
			UserID:      1, // student001
			Title:       "软件质量保证",
			Description: "软件质量保证的方法和最佳实践",
			Filename:    "quality_assurance.md",
			FilePath:    "./uploads/quality_assurance.md",
			FileSize:    0,
			FileType:    ".md",
			Content: `# 软件质量保证

## 软件质量的维度

软件质量通常从以下维度评估：

1. **功能性**：软件是否满足功能需求
2. **可靠性**：软件在规定条件下完成规定功能的能力
3. **易用性**：用户使用软件的难易程度
4. **效率**：软件的性能表现
5. **可维护性**：软件被修改的难易程度
6. **可移植性**：软件从一个环境迁移到另一个环境的难易程度

## 代码质量

提高代码质量的方法：

- **代码审查**：通过同行评审发现潜在问题
- **静态分析**：使用工具自动检测代码问题
- **测试驱动开发**：先写测试再写代码
- **重构**：在不改变外部行为的前提下改善代码结构
- **编码规范**：遵循统一的编码标准

## 持续集成与持续交付

**持续集成（CI）**：
- 频繁地将代码集成到主干
- 每次集成都通过自动化构建验证
- 快速发现和修复集成问题

**持续交付（CD）**：
- 确保软件随时可以部署到生产环境
- 自动化部署流程
- 灰度发布和回滚机制`,
			Status: "completed",
		},
	}

	if err := database.DB.Create(&documents).Error; err != nil {
		log.Printf("seed documents failed: %v", err)
		return
	}
	log.Println("seeded 3 demo documents")
}
