package controller

import (
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/dto"
	"software_engineering/internal/service"
	"software_engineering/internal/utils"
)

func UploadDocument(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "请上传文件")
		return
	}
	title := c.PostForm("title")
	description := c.PostForm("description")
	ext := filepath.Ext(file.Filename)

	f, err := file.Open()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "文件读取失败")
		return
	}
	defer f.Close()

	resp, err := service.UploadDocument(title, description, file.Filename, file.Size, ext, f)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, resp)
}

func GetDocument(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	resp, err := service.GetDocument(uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}
	utils.Success(c, resp)
}

func GetDocumentContent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	resp, err := service.GetDocumentContent(uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}
	utils.Success(c, resp)
}

func UpdateDocument(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req dto.UpdateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	if err := service.UpdateDocument(uint(id), req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, nil)
}

func DeleteDocument(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := service.DeleteDocument(uint(id)); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, nil)
}

func ListDocuments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	keyword := c.Query("keyword")
	status := c.Query("status")
	list, total, err := service.ListDocuments(page, size, keyword, status)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Paginated(c, list, total, page, size)
}
