package handler

import (
	"test/pkg/response"
	"test/service"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Nickname   string `json:"nickname" binding:"required"`
	Department string `json:"department" binding:"required"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 10001, "参数错误")
		return
	}
	validDepts := map[string]bool{
		"backend": true, "frontend": true, "sre": true,
		"product": true, "design": true, "android": true, "ios": true,
	}
	if !validDepts[req.Department] {
		response.Error(c, 10001, "部门参数错误")
		return
	}
	err := service.Register(service.RegisterRequest{
		Username:   req.Username,
		Password:   req.Password,
		Nickname:   req.Nickname,
		Department: req.Department,
	})
	if err != nil {
		response.Error(c, 10002, err.Error())
		return
	}

	response.Success(c, nil)
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 10001, "参数错误")
		return
	}

	resp, err := service.Login(service.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		response.Error(c, 10003, err.Error())
		return
	}
	response.Success(c, gin.H{
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
		"user": gin.H{
			"id":               resp.User.ID,
			"username":         resp.User.Username,
			"nickname":         resp.User.Nickname,
			"role":             resp.User.Role,
			"department":       resp.User.Department,
			"department_label": getDeptLabel(resp.User.Department),
		},
	})
}

func getDeptLabel(dept string) string {
	labels := map[string]string{
		"backend":  "后端",
		"frontend": "前端",
		"sre":      "SRE",
		"product":  "产品",
		"design":   "视觉设计",
		"android":  "Android",
		"ios":      "iOS",
	}
	if label, ok := labels[dept]; ok {
		return label
	}
	return dept
}

func GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	user, err := service.GetUserInfo(userID.(uint))
	if err != nil {
		response.Error(c, 10004, err.Error())
		return
	}
	response.Success(c, gin.H{
		"id":               user.ID,
		"username":         user.Username,
		"nickname":         user.Nickname,
		"role":             user.Role,
		"department":       user.Department,
		"department_label": getDeptLabel(user.Department),
		"email":            user.Email,
		"created_at":       user.CreatedAt,
	})
}

type BindEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func BindEmail(c *gin.Context) {
	var req BindEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 10001, "邮箱格式错误")
		return
	}
	userID, _ := c.Get("user_id")
	err := service.BindEmail(service.BindEmailRequest{
		UserID: userID.(uint),
		Email:  req.Email,
	})
	if err != nil {
		response.Error(c, 10005, err.Error())
		return
	}
	response.Success(c, nil)
}
