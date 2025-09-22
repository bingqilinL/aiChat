package user

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// User 用户结构体（内部使用，包含密码）
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"` // 存储明文密码（简化版本）
}

// UserResponse 用户响应结构体（不包含密码）
type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// UserManager 用户管理器
type UserManager struct {
	redisClient *redis.Client
}

// NewUserManager 创建用户管理器
func NewUserManager(redisClient *redis.Client) *UserManager {
	return &UserManager{
		redisClient: redisClient,
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	Message string        `json:"message"`
	User    *UserResponse `json:"user,omitempty"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Message string        `json:"message"`
	User    *UserResponse `json:"user"`
	Token   string        `json:"token"`
}

// LogoutResponse 退出登录响应
type LogoutResponse struct {
	Message string `json:"message"`
}

// generateUserID 生成用户ID
func generateUserID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// hashPassword 加密密码
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// checkPassword 验证密码
func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// toUserResponse 将User转换为UserResponse
func toUserResponse(user *User) *UserResponse {
	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
}

// Register 用户注册
func (um *UserManager) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	// 检查用户名是否已存在
	userKey := fmt.Sprintf("user:%s", req.Username)
	exists, err := um.redisClient.Exists(ctx, userKey).Result()
	if err != nil {
		log.Printf("检查用户是否存在失败: %s, 错误: %v", req.Username, err)
		return nil, fmt.Errorf("检查用户是否存在失败: %v", err)
	}
	if exists > 0 {
		log.Printf("用户注册失败: 用户名已存在 - %s", req.Username)
		return nil, fmt.Errorf("用户名已存在")
	}

	// 生成用户ID
	userID := generateUserID()
	// 加密密码
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		log.Printf("用户注册失败: 密码加密失败 - %s, 错误: %v", req.Username, err)
		return nil, fmt.Errorf("密码加密失败: %v", err)
	}
	// 创建用户对象
	user := &User{
		ID:       userID,
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	// 将用户信息存储到Redis
	userData, err := json.Marshal(user)
	if err != nil {
		log.Printf("用户注册失败: 用户数据序列化失败 - %s, 错误: %v", req.Username, err)
		return nil, fmt.Errorf("用户数据序列化失败: %v", err)
	}

	// 存储用户信息（key: user:用户名, value: 用户对象）
	err = um.redisClient.Set(ctx, userKey, userData, 0).Err()
	if err != nil {
		log.Printf("用户注册失败: 存储用户信息失败 - %s, 错误: %v", req.Username, err)
		return nil, fmt.Errorf("存储用户信息失败: %v", err)
	}
	//// 创建用户名到用户ID的映射
	//usernameToIDKey := fmt.Sprintf("username_to_id:%s", req.Username)
	//err = um.redisClient.Set(ctx, usernameToIDKey, userID, 0).Err()
	//if err != nil {
	//	return nil, fmt.Errorf("存储用户名映射失败: %v", err)
	//}
	//
	//// 创建用户ID到用户名的映射
	//userIDToUsernameKey := fmt.Sprintf("user_id_to_username:%s", userID)
	//err = um.redisClient.Set(ctx, userIDToUsernameKey, req.Username, 0).Err()
	//if err != nil {
	//	return nil, fmt.Errorf("存储用户ID映射失败: %v", err)
	//}

	log.Printf("用户注册成功: %s", req.Username)

	return &RegisterResponse{
		Message: "注册成功",
		User:    toUserResponse(user),
	}, nil
}

// Login 用户登录
func (um *UserManager) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// 获取用户信息
	userKey := fmt.Sprintf("user:%s", req.Username)
	userData, err := um.redisClient.Get(ctx, userKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Printf("用户登录失败: 用户不存在 - %s", req.Username)
			return nil, fmt.Errorf("用户名或密码错误")
		}
		log.Printf("用户登录失败: 获取用户信息失败 - %s, 错误: %v", req.Username, err)
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}

	// 反序列化用户信息
	var user User
	err = json.Unmarshal([]byte(userData), &user)
	if err != nil {
		log.Printf("用户登录失败: 用户数据反序列化失败 - %s, 错误: %v", req.Username, err)
		return nil, fmt.Errorf("用户数据反序列化失败: %v", err)
	}

	// 验证密码
	if !checkPassword(req.Password, user.Password) {
		log.Printf("用户登录失败: 密码错误 - %s", req.Username)
		return nil, fmt.Errorf("用户名或密码错误")
	}
	// 生成登录token
	//token := generateUserID()
	////存储token到用户ID的映射（设置24小时过期）
	//tokenKey := fmt.Sprintf("token:%s", token)
	//err = um.redisClient.Set(ctx, tokenKey, user.ID, 24*time.Hour).Err()
	//if err != nil {
	//	return nil, fmt.Errorf("存储登录token失败: %v", err)
	//}

	//// 存储用户ID到token的映射
	//userTokenKey := fmt.Sprintf("user_token:%s", user.ID)
	//err = um.redisClient.Set(ctx, userTokenKey, token, 24*time.Hour).Err()
	//if err != nil {
	//	return nil, fmt.Errorf("存储用户token失败: %v", err)
	//}
	log.Printf("用户登录成功: %s", req.Username)

	return &LoginResponse{
		Message: "登录成功",
		User:    toUserResponse(&user),
		Token:   "token",
	}, nil
}

// Logout 用户退出登录
func (um *UserManager) Logout(ctx context.Context, token string) (*LogoutResponse, error) {
	// 获取token对应的用户ID
	//tokenKey := fmt.Sprintf("token:%s", token)
	//userID, err := um.redisClient.Get(ctx, tokenKey).Result()
	//if err != nil {
	//	if err == redis.Nil {
	//		log.Printf("用户退出登录失败: 无效的token - %s", token)
	//		return nil, fmt.Errorf("无效的token")
	//	}
	//	log.Printf("用户退出登录失败: 获取token信息失败 - %s, 错误: %v", token, err)
	//	return nil, fmt.Errorf("获取token信息失败: %v", err)
	//}
	//
	//// 删除token映射
	//err = um.redisClient.Del(ctx, tokenKey).Err()
	//if err != nil {
	//	log.Printf("用户退出登录失败: 删除token失败 - 用户ID: %s, 错误: %v", userID, err)
	//	return nil, fmt.Errorf("删除token失败: %v", err)
	//}

	//// 删除用户token映射
	//userTokenKey := fmt.Sprintf("user_token:%s", userID)
	//err = um.redisClient.Del(ctx, userTokenKey).Err()
	//if err != nil {
	//	log.Printf("用户退出登录失败: 删除用户token失败 - 用户ID: %s, 错误: %v", userID, err)
	//	return nil, fmt.Errorf("删除用户token失败: %v", err)
	//}
	// 简化版本：直接返回成功，不需要token验证
	log.Printf("用户退出登录成功")

	return &LogoutResponse{
		Message: "退出登录成功",
	}, nil
}

// GetUserByToken 根据token获取用户信息（简化版本，暂时不支持）
func (um *UserManager) GetUserByToken(ctx context.Context, token string) (*User, error) {
	// 简化版本：暂时不支持token验证
	return nil, fmt.Errorf("简化版本不支持token验证")
}

// GetUserByUsername 根据用户名获取用户信息
func (um *UserManager) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	userKey := fmt.Sprintf("user:%s", username)
	userData, err := um.redisClient.Get(ctx, userKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}

	// 反序列化用户信息
	var user User
	err = json.Unmarshal([]byte(userData), &user)
	if err != nil {
		return nil, fmt.Errorf("用户数据反序列化失败: %v", err)
	}

	return &user, nil
}
