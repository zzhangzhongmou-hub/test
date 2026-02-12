package handler

import (
	"strconv"
	"test/pkg/response"
	"test/service"
	"time"

	"github.com/gin-gonic/gin"
)

func AIReview(c *gin.Context) {
	idStr := c.Param("id")
	submissionID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, 10001, "ID格式错误")
		return
	}
	resultChan := service.EvaluateSubmissionAsync(uint(submissionID))
	select {
	case result := <-resultChan:
		if result.Error != nil && result.Comment == "" {
			response.Error(c, 10010, "AI评价失败: "+result.Error.Error())
			return
		}
		
		response.Success(c, gin.H{
			"ai_comment":      result.Comment,
			"suggested_score": result.Score,
		})

	case <-time.After(8 * time.Second):
		response.Success(c, gin.H{
			"ai_comment":      "AI评价超时，建议人工评价",
			"suggested_score": 0,
		})
	}
}
