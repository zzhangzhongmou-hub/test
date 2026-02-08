package service

import (
	"errors"
	"fmt"
	"test/dao"
	"test/models"
	"test/pkg/hash"
	"test/pkg/jwt"
)

type RegisterRequest struct {
	Username   string
	Password   string
	Nickname   string
	Department string
}

func Register(req RegisterRequest) error {
	existingUser, _ := dao.GetUserByUsername(req.Username)
	if existingUser != nil {
		return errors.New("用户名已存在")
	}

	hashedPassword, err := hash.Encrypt(req.Password)
	if err != nil {
		return errors.New("加密失败")
	}
	user := &models.User{
		Username:   req.Username,
		Password:   hashedPassword,
		Nickname:   req.Nickname,
		Role:       models.RoleStudent,
		Department: req.Department,
	}
	return dao.CreateUser(user)
}

type LoginRequest struct {
	Username string
	Password string
}

type LoginResponse struct {
	AccessToken  string
	RefreshToken string
	User         *models.User
}

func Login(req LoginRequest) (*LoginResponse, error) {
	user, err := dao.GetUserByUsername(req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}
	if !hash.Check(req.Password, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}
	accessToken, refreshToken, err := jwt.GenerateToken(user.ID, user.Role, user.Department)
	if err != nil {
		return nil, errors.New("Token生成失败")
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func GetUserInfo(userID uint) (*models.User, error) {
	user, err := dao.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	return user, nil
}

func DeleteUser(userID uint, password string) error {
	user, err := dao.GetUserByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}
	if !hash.Check(password, user.Password) {
		return errors.New("密码错误")
	}
	return dao.DeleteUser(userID)
}

type BindEmailRequest struct {
	UserID uint
	Email  string
}

func BindEmail(req BindEmailRequest) error {
	existing, err := dao.GetUserByEmail(req.Email)
	if err != nil {
		return err
	}
	if existing != nil && existing.ID != req.UserID {
		return fmt.Errorf("该邮箱已被其他账号绑定")
	}
	return dao.UpdateUserEmail(req.UserID, req.Email)
}
