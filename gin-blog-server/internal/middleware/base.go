package middleware

import (
	"errors"
	"log/slog"
	"lv-blog/internal/global"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// WithRedisDB 将 redis.Client 注入到 gin.Context
// handler 中通过 c.MustGet(g.CTX_RDB).(*redis.Client) 来使用
// 将Redis客户端注入到Gin的上下文中，以便在处理请求时可以访问Redis
// 使用ctx.Set() 方法将Redis客户端存储在Gin的上下文，键为 global.CTX_RDB
func WithRedisDB(rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(global.CTX_RDB, rdb)
		ctx.Next()
	}
}

// 将Gorm数据库实例注入到Gin的上下文中，以便在处理请求时可以访问数据库
func WithGormDB(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(global.CTX_DB, db)
		ctx.Next()
	}
}

// CORS 跨域请求
func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"PUT", "POST", "GET", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders: []string{"Origin", "Authorization", "Content-Type", "X-Requested-With"},
		ExposeHeaders: []string{"Content-type"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 24 * time.Hour,
	})
}

// WithCookieStore 基于 cookie 的 session
func WithCookieStore(name, secret string) gin.HandlerFunc {
	// 使用cookie存储会话数据
	store := cookie.NewStore([]byte(secret))
	// 设置会话的路径和最大年龄
	store.Options(sessions.Options{Path: "/", MaxAge: 600})
	return sessions.Sessions(name, store)
}

// WithMemStore 基于内存的session
func WithMemStore(name, secret string) gin.HandlerFunc {
	// 使用内存存储会话数据
	store := memstore.NewStore([]byte(secret))
	store.Options(sessions.Options{Path: "/", MaxAge: 600})
	return sessions.Sessions(name, store)
}

// Logger 日志记录
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		cost := time.Since(start)

		slog.Info("[GIN]",
				slog.String("path", c.Request.URL.Path),
				slog.String("query", c.Request.URL.RawQuery),
				slog.Int("status", c.Writer.Status()),
				slog.String("method", c.Request.Method),
				slog.String("ip", c.ClientIP()),
				slog.Int("size", c.Writer.Size()),
				slog.Duration("cost", cost),
	)
	}
}

// Recovery 恢复中间件
func Recovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					var se *os.SyscallError
					if errors.As(ne, &se) {
						seStr := strings.ToLower(se.Error())
						if strings.Contains(seStr, "broken pipe") ||
							strings.Contains(seStr, "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
			}
			handle.ReturnHttpResponse(c, http.StatusInternalServerError, global.FAIL, global.GetMsg(global.FAIL), err)

			httpRequest, _ := httputil.DumpRequest(c.Request, false)
			if brokenPipe {
				slog.Error(c.Request.URL.Path,
					slog.Any("error", err),
					slog.String("request", string(httpRequest)),
				)
				
				_ = c.Errors(err.(error))
				c.Abort()
				return
			}

			if stack {
				slog.Error("[Recovery from panic]",
					slog.Any("error", err),
					slog.String("request", string(httpRequest)),
					slog.String("stack", string(debug.Stack())),
				)
			} else {
				slog.Error("[Recovery from panic]",
					slog.Any("error", err),
					slog.String("request", string(httpRequest)),
				)
			}
			c.AbortWithStatus(http.StatusInternalServerError)
		} ()
		c.Next()
	}
	
}