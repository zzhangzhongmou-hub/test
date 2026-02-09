package handler

import (
	"test/dao"
	"test/pkg/response"
	"test/service"

	"github.com/gin-gonic/gin"
)

type CreateExamRequest struct {
	Title         string `json:"title" binding:"required"`
	Description   string `json:"description" binding:"required"`
	Department    string `json:"department" binding:"required"`
	Deadline      string `json:"deadline" binding:"required"`
	Duration      int    `json:"duration"`
	ReviewerCount int    `json:"reviewer_count"`
}

func CreateExam(c *gin.Context) {
	var req CreateExamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 10001, "参数错误："+err.Error())
		return
	}

	creatorID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, 10002, "未获取到用户信息")
		return
	}
	err := service.CreateExam(service.CreateExamRequest{
		Title:         req.Title,
		Description:   req.Description,
		Department:    req.Department,
		Deadline:      req.Deadline,
		Duration:      req.Duration,
		ReviewerCount: req.ReviewerCount,
		CreatorID:     creatorID.(uint),
	})
	if err != nil {
		response.Error(c, 10003, err.Error())
		return
	}
	response.Success(c, nil)
}

func GetMyReviews(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, 10002, "未获取到用户信息")
		return
	}
	reviews, err := dao.GetPendingReviewsByReviewer(userID.(uint))
	if err != nil {
		response.Error(c, 10003, err.Error())
		return
	}
	list := make([]gin.H, 0, len(reviews))
	for _, r := range reviews {
		list = append(list, gin.H{
			"id":          r.ID,
			"exam_id":     r.ExamID,
			"student_id":  r.StudentID,
			"status":      r.Status,
			"version":     r.Version,
			"assigned_at": r.AssignedAt.Format("2025-02-02 15:04:05"),
		})
	}

	response.Success(c, gin.H{
		"list":  list,
		"total": len(list),
	})
}

type SubmitReviewRequest struct {
	ReviewID int    `json:"review_id" binding:"required"`
	Score    int    `json:"score" binding:"required,min=0,max=100"`
	Comment  string `json:"comment" binding:"required"`
	Version  int    `json:"version" binding:"required"`
}

func SubmitReview(c *gin.Context) {
	var req SubmitReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 10001, "参数错误："+err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, 10002, "未获取到用户信息")
		return
	}
	err := service.SubmitReview(service.SubmitReviewRequest{
		ReviewID:   uint(req.ReviewID),
		ReviewerID: userID.(uint),
		Score:      req.Score,
		Comment:    req.Comment,
		Version:    req.Version,
	})
	if err != nil {
		if err.Error() == "记录已被修改，请刷新后重试" ||
			err.Error() == "考卷已被其他管理员修改，请刷新后重试" {
			response.Error(c, 10009, err.Error())
			return
		}
		response.Error(c, 10008, err.Error())
		return
	}

	response.Success(c, nil)
}
