package handler

import (
	"strconv"
	"test/dao"
	"test/pkg/response"
	"test/service"

	"github.com/gin-gonic/gin"
)

func Submit(c *gin.Context) {
	var req service.SubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 10001, "参数错误："+err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, 10002, "未获取到用户信息")
		return
	}
	req.StudentID = userID.(uint)

	if err := service.Submit(req); err != nil {
		response.Error(c, 10005, err.Error())
		return
	}

	response.Success(c, nil)
}

func GetMySubmissions(c *gin.Context) {
	userID, _ := c.Get("user_id")

	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	subs, total, err := service.GetMySubmissions(userID.(uint), page, pageSize)
	if err != nil {
		response.Error(c, 10003, err.Error())
		return
	}

	list := make([]gin.H, 0, len(subs))
	for _, sub := range subs {
		list = append(list, gin.H{
			"id":           sub.ID,
			"homework_id":  sub.HomeworkID,
			"content":      sub.Content,
			"file_url":     sub.FileURL,
			"is_late":      sub.IsLate,
			"score":        sub.Score,
			"comment":      sub.Comment,
			"is_excellent": sub.IsExcellent,
			"status":       sub.Status,
			"created_at":   sub.CreatedAt,
		})
	}

	response.Success(c, gin.H{
		"list":      list,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

type ReviewRequest struct {
	Score   int    `json:"score" binding:"required,min=0,max=100"`
	Comment string `json:"comment" binding:"required"`
}

func GetSubmissionsByHomework(c *gin.Context) {
	homeworkIDStr := c.Param("homework_id")
	homeworkID, err := strconv.ParseUint(homeworkIDStr, 10, 64)
	if err != nil {
		response.Error(c, 10001, "作业ID格式错误")
		return
	}

	dept, _ := c.Get("department")
	homework, _ := dao.GetHomeworkByID(uint(homeworkID))
	if homework.Department != dept.(string) {
		response.Error(c, 10003, "只能查看本部门作业的提交")
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	subs, total, err := service.GetSubmissionsByHomework(uint(homeworkID), page, pageSize)
	if err != nil {
		response.Error(c, 10003, err.Error())
		return
	}

	list := make([]gin.H, 0, len(subs))
	for _, sub := range subs {
		item := gin.H{
			"id":           sub.ID,
			"student_id":   sub.StudentID,
			"content":      sub.Content,
			"file_url":     sub.FileURL,
			"is_late":      sub.IsLate,
			"score":        sub.Score,
			"comment":      sub.Comment,
			"is_excellent": sub.IsExcellent,
			"status":       sub.Status,
			"created_at":   sub.CreatedAt,
		}

		if sub.Student.ID != 0 {
			item["student_name"] = sub.Student.Nickname
		}

		list = append(list, item)
	}

	response.Success(c, gin.H{
		"list":      list,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func Review(c *gin.Context) {
	idStr := c.Param("id")
	submissionID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, 10001, "ID格式错误")
		return
	}

	var req ReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 10001, "参数错误："+err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	dept, _ := c.Get("department")

	err = service.Review(service.ReviewRequest{
		SubmissionID: uint(submissionID),
		ReviewerID:   userID.(uint),
		ReviewerDept: dept.(string),
		Score:        req.Score,
		Comment:      req.Comment,
	})
	if err != nil {
		response.Error(c, 10008, err.Error())
		return
	}

	response.Success(c, nil)
}

func MarkExcellent(c *gin.Context) {
	idStr := c.Param("id")
	submissionID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, 10001, "ID格式错误")
		return
	}

	var req struct {
		IsExcellent bool `json:"is_excellent" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 10001, "参数错误："+err.Error())
		return
	}

	dept, _ := c.Get("department")

	err = service.MarkExcellent(service.MarkExcellentRequest{
		SubmissionID: uint(submissionID),
		IsExcellent:  req.IsExcellent,
		ReviewerDept: dept.(string),
	})
	if err != nil {
		response.Error(c, 10008, err.Error())
		return
	}

	response.Success(c, nil)
}

func GetExcellentSubmissions(c *gin.Context) {
	department := c.Query("department")

	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	subs, total, err := service.GetExcellentSubmissions(department, page, pageSize)
	if err != nil {
		response.Error(c, 10003, err.Error())
		return
	}

	list := make([]gin.H, 0, len(subs))
	for _, sub := range subs {
		list = append(list, gin.H{
			"id":             sub.ID,
			"homework_title": sub.Homework.Title,
			"student_id":     sub.StudentID,

			"student_nickname": sub.Student.Nickname,
			"department":       sub.Homework.Department,
			"department_label": getDeptLabel(sub.Homework.Department),
			"score":            sub.Score,
			"comment":          sub.Comment,
			"content":          sub.Content,
			"created_at":       sub.CreatedAt,
		})
	}

	response.Success(c, gin.H{
		"list":      list,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
