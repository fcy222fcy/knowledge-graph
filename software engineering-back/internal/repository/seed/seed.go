package seed

import (
	"log"
	"software_engineering/internal/database"
	"software_engineering/internal/model/entity"
)

func SeedAll() {
	seedUsers()
	seedKnowledgePoints()
	seedQuestions()
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
