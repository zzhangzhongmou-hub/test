package handler

import (
	"strconv"
	"test/pkg/response"
	"test/service"

	"github.com/gin-gonic/gin"
)

type CreateHomeworkRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Department  string `json:"department" binding:"required"`
	Deadline    string `json:"deadline" binding:"required"`
	AllowLate   bool   `json:"allow_late"`
}

func CreateHomework(c *gin.Context) {
	var req CreateHomeworkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 10001, "参数错误："+err.Error())
		return
	}
	creatorID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, 10002, "未获取到用户信息")
		return
	}

	err := service.CreateHomework(service.CreateHomeworkRequest{
		Title:       req.Title,
		Description: req.Description,
		Department:  req.Department,
		CreatorID:   creatorID.(uint),
		Deadline:    req.Deadline,
		AllowLate:   req.AllowLate,
	})
	if err != nil {
		response.Error(c, 10002, err.Error())
		return
	}

	response.Success(c, nil)
}

func GetHomeworkList(c *gin.Context) {
	department := c.Query("department")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	homeworks, total, err := service.GetHomeworkList(department, page, pageSize)
	if err != nil {
		response.Error(c, 10003, err.Error())
		return
	}
	list := make([]gin.H, 0, len(homeworks))
	for _, hw := range homeworks {
		list = append(list, gin.H{
			"id":               hw.ID,
			"title":            hw.Title,
			"department":       hw.Department,
			"department_label": getDeptLabel(hw.Department),
			"deadline":         hw.Deadline.Format("2006-01-02 15:04:05"),
			"allow_late":       hw.AllowLate,
			"creator_id":       hw.CreatorID,
			"version":          hw.Version,
		})
	}

	response.Success(c, gin.H{
		"list":      list,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetHomeworkDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, 10001, "ID格式错误")
		return
	}
	homework, err := service.GetHomeworkDetail(uint(id))
	if err != nil {
		response.Error(c, 10004, "作业不存在")
		return
	}
	response.Success(c, gin.H{
		"id":               homework.ID,
		"title":            homework.Title,
		"description":      homework.Description,
		"department":       homework.Department,
		"department_label": getDeptLabel(homework.Department),
		"deadline":         homework.Deadline.Format("2025-02-05 15:04:05"),
		"allow_late":       homework.AllowLate,
		"creator_id":       homework.CreatorID,
		"version":          homework.Version,
		"created_at":       homework.CreatedAt,
	})
}

type UpdateHomeworkRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Deadline    string `json:"deadline"`
	AllowLate   bool   `json:"allow_late"`
	Version     int    `json:"version" binding:"required"`
}

func UpdateHomework(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, 10001, "ID格式错误")
		return
	}
	var req UpdateHomeworkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 10001, "参数错误："+err.Error())
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, 10002, "未获取到用户信息")
		return
	}
	dept, exists := c.Get("department")
	if !exists {
		response.Error(c, 10002, "未获取到用户部门信息")
		return
	}
	err = service.UpdateHomework(service.UpdateHomeworkRequest{
		ID:          uint(id),
		Title:       req.Title,
		Description: req.Description,
		Deadline:    req.Deadline,
		AllowLate:   req.AllowLate,
		Version:     req.Version,
		UpdaterID:   userID.(uint),
		UpdaterDept: dept.(string),
	})
	if err != nil {
		response.Error(c, 10005, err.Error())
		return
	}

	response.Success(c, nil)
}
func DeleteHomework(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, 10001, "ID格式错误")
		return
	}
	dept, exists := c.Get("department")
	if !exists {
		response.Error(c, 10002, "未获取到用户部门信息")
		return
	}

	err = service.DeleteHomework(uint(id), dept.(string))
	if err != nil {
		response.Error(c, 10006, err.Error())
		return
	}

	response.Success(c, nil)
}
