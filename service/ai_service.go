package service

import (
	"context"
	"fmt"
	"test/dao"
	"test/pkg/ai"
	"time"
)

type AIResult struct {
	Comment string
	Score   int
	Error   error
}

func EvaluateSubmissionAsync(submissionID uint) <-chan AIResult {
	resultChan := make(chan AIResult, 1)

	go func() {
		defer close(resultChan)

		submission, err := dao.GetSubmissionByID(submissionID)
		if err != nil {
			resultChan <- AIResult{Error: fmt.Errorf("提交记录不存在")}
			return
		}

		client := ai.NewClient()
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()
		comment, score, err := client.EvaluateCode(ctx, submission.Content)

		if err != nil {
			resultChan <- AIResult{
				Comment: "AI评价服务暂时不可用，建议人工评价。原因：" + err.Error(),
				Score:   0,
				Error:   err,
			}
			return
		}
		resultChan <- AIResult{
			Comment: comment,
			Score:   score,
			Error:   nil,
		}
	}()

	return resultChan
}
