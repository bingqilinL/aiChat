package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserHandler 用户处理器
type UserHandler struct {
	userManager *UserManager
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userManager *UserManager) *UserHandler {
	return &UserHandler{
		userManager: userManager,
	}
}

// Register 处理用户注册请求
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	// 验证用户名长度
	if len(req.Username) < 3 || len(req.Username) > 20 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "用户名长度必须在3-20个字符之间",
		})
		return
	}

	// 验证密码长度
	if len(req.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "密码长度至少6个字符",
		})
		return
	}

	// 调用用户管理器进行注册
	response, err := h.userManager.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Login 处理用户登录请求
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	// 调用用户管理器进行登录
	response, err := h.userManager.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.SetSameSite(http.SameSiteNoneMode)
	// 设置HTTP-only cookie存储token
	c.SetCookie("auth_token", response.Token, 24*60*60, "/", "", false, true)

	c.JSON(http.StatusOK, response)
}

// Logout 处理用户退出登录请求
func (h *UserHandler) Logout(c *gin.Context) {
	// 调试：打印所有cookie
	//fmt.Printf("所有cookies: %v\n", c.Request.Cookies())
	// 从cookie中获取token
	//token, err := c.Cookie("auth_token")
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"error": "未找到登录token",
	//	})
	//	return
	//}

	// 调用用户管理器进行退出登录
	response, err := h.userManager.Logout(c.Request.Context(), "token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 清除cookie
	c.SetCookie("auth_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, response)
}

// GetCurrentUser 获取当前用户信息
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	// 从cookie中获取token
	token, err := c.Cookie("auth_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未找到登录token",
		})
		return
	}

	// 根据token获取用户信息
	user, err := h.userManager.GetUserByToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": toUserResponse(user),
	})
}
