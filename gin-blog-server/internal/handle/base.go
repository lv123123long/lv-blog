package handle

import (
	"errors"
	"log/slog"
	"lv-blog/internal/global"
	"lv-blog/internal/model"
	"net/http"

	"github.com/gin-contrib/sessions"     // 会话管理
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

/*
响应式设计方案： 不使用HTTP码来表示业务状态，采用业务状态码的方式
- 只要能达到后端的请求，HTTP状态码都为 200
- 业务状态码为0 表示成功，其他表示失败
-当后端发生panic并且被gin中间件捕获时，才会返回HTTP 500状态码
*/

// 响应结构体  泛型结构体
type Response[T any] struct {
	Code    int    `json:"code"`    // 业务状态码
	Message string `json:"message"` // 响应消息
	Data    T      `json:"data"` // T表示任何参数，响应数据
	// 泛型 允许我们使用占位符T，在实例化结构体的时候指定具体的类型
}

/*
gin框架的c.JSON() 方法将响应写入HTTP响应流
*/

// HTTP状态码 + 业务码 + 消息 + 数据
// 定义一个函数，用于返回HTTP响应，包含HTTP状态码、业务状态码、消息和数据。
func ReturnHttpResponse(c *gin.Context, httpCode, code int, msg string, data any) {
	c.JSON(httpCode, Response[any]{
		Code:    code,
		Message: msg,
		Data:    data,
	})
}

// 业务码 + 数据
// 定义一个函数，用于返回业务状态码和数据，默认HTTP状态码为200。
func ReturnResponse(c *gin.Context, r global.Result, data any) {
	ReturnHttpResponse(c, http.StatusOK, r.Code(), r.Msg(), data)
}

// 成功业务码 + 数据
// 定义一个函数，用于返回成功业务状态码和数据
func ReturnSuccess(c *gin.Context, data any) {
	ReturnResponse(c, global.OkResult, data)
}

// 所有可预料的错误 = 业务错误 + 系统错误, 在业务层面处理, 返回 HTTP 200 状态码
// 对于不可预料的错误, 会触发 panic, 由 gin 中间件捕获, 并返回 HTTP 500 状态码
// err 是业务错误, data 是错误数据 (可以是 error 或 string)
// 定义一个函数，用于返回业务错误和错误数据，默认HTTP状态码为200。如果发生不可预料的错误，会触发panic，由Gin中间件捕获并返回HTTP 500状态码。
func ReturnError(c *gin.Context, r global.Result, data any) {
	slog.Info("[Func-ReturnError] " + r.Msg())

	var val string = r.Msg()

	if data != nil {
		switch v := data.(type) {
		case error:
			val = v.Error()
		case string:
			val = v
		}
		slog.Error(val) // 错误日志
	}

	c.AbortWithStatusJSON(
		http.StatusOK,
		Response[any]{
			Code:    r.Code(),
			Message: r.Msg(),
			Data:    val,
		},
	)
}

// 分页获取数据
type PageQuery struct {
	Page    int    `from:"page_num"`
	Size    int    `from:"page_size"`
	Keyword string `from:"keyword"`
}

// 分页响应数据
type PageResult[T any] struct {
	Page  int   `json:"page_num"`  // 每页条数
	Size  int   `json:"page_size"` // 上次页数
	Total int64 `json:"total"`     // 总条数
	List  []T   `json:"page_data"` // 分页数据
}

// 获取 *gorm.DB
// 定义一个函数，用于从Gin上下文中获取GORM数据库实例。
func GetDB(c *gin.Context) *gorm.DB {
	return c.MustGet(global.CTX_DB).(*gorm.DB)
	// MustGet方法会返回一个interface{}类型的值，表示从上下文中获取的数据。并将其转换为*gorm.DB类型。
}

// 获取*redis.Client
// 定义一个函数，用于从Gin上下文中获取Redis客户端实例。
func GetRDB(c *gin.Context) *redis.Client {
	return c.MustGet(global.CTX_RDB).(*redis.Client)
}

/*
获取当前登录用户信息
1. 能从 gin Context 上获取到 user 对象, 说明本次请求链路中获取过了
2. 从 session 中获取到 uid
3. 根据 uid 获取用户信息, 并设置到 gin Context 上
*/
// 定义一个函数，用于获取当前登录用户信息。首先尝试从Gin上下文中获取用户信息，
// 如果获取不到，则从会话中获取用户ID，再根据用户ID从数据库中获取用户信息，并设置到Gin上下文中。
func CurrentUserAuth(c *gin.Context) (*model.UserAuth, error) {
	key := global.CTX_USER_AUTH

	//1
	if cache, exist := c.Get(key); exist && cache != nil {
		slog.Debug("[Func-CurrentUserAuth] get from cache: " + cache.(*model.UserAuth).Username)
		return cache.(*model.UserAuth), nil
	}

	// 2
	session := sessions.Default(c)
	id := session.Get(key)
	if id == nil {
		return nil, errors.New("session 中没有 user_auth_id")
	}

	//3
	db := GetDB(c)
	user, err := model.GetUserAuthInfoById(db, id.(int))
	if err != nil {
		return nil, err
	}

	c.Set(key, user)
	return user, nil
}
