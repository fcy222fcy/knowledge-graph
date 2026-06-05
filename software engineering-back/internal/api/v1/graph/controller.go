package graph

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/model/dto/request"
	"software_engineering/internal/service"
	"software_engineering/pkg/response"
)

func GetGraph(c *gin.Context) {
	documentID, _ := strconv.Atoi(c.Query("document_id"))
	keyword := c.Query("keyword")
	relationType := c.Query("relation_type")
	resp, err := service.GetGraphData(uint(documentID), keyword, relationType)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, resp)
}

func BuildGraph(c *gin.Context) {
	var req request.BuildGraphRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	resp, err := service.BuildGraph(req.DocumentIDs)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, resp)
}

func GetLatestBuild(c *gin.Context) {
	resp, err := service.GetLatestBuildResult()
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.Success(c, resp)
}

func ListBuildHistory(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	resp, err := service.ListBuildHistory(page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, resp)
}
