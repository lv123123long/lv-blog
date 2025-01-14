package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// 创建自定义的错误
var (
	ErrTokenExpired     = errors.New("token 已过期，请重新登录")
	ErrTokenNotValidYet = errors.New("token 无效， 请重新登录")
	ErrTokenMalformed   = errors.New("token 不正确，请重新登录")
	ErrTokenInvalid     = errors.New("这不是一个token，请重新登录")
)

// 自定义的声明结构体，用来存储JWT中的自定义数据
type MyClaims struct {
	UserId  int   `json:"user_id"`        // 用户ID
	RoleIds []int `json:"role_ids"`       // 角色ID
	jwt.RegisteredClaims                  // 存储JWT的注册声明（发行者，过期时间等）
}

// 生成JWT的函数
// secret: 用于签名的密钥
// issuer：JWT的发行者               JWT的过期时间（小时） 用户ID  用户的角色ID列表
func GenToken(secret, issuer string, expireHour, userId int, roleIds []int) (string, error) {
	claims := MyClaims{
		UserId:  userId,
		RoleIds: roleIds,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireHour) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// jwt的New 创建一个新的JWT;   使用HS256加密算法进行签名
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用给定的密钥对JWT进行签名，并返回签名后的JWT字符串
	return token.SignedString([]byte(secret))
}


// secret：用于验证签名的密钥
// token：要解析的token
func ParseToken(secret, token string) (*MyClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &MyClaims{},
		// 提供一个函数来获取用于验证签名的密钥
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
	if err != nil {
		switch vError, ok := err.(*jwt.ValidationError); ok {
		case vError.Errors&jwt.ValidationErrorMalformed != 0:
			return nil, ErrTokenMalformed
		case vError.Errors&jwt.ValidationErrorExpired != 0:
			return nil, ErrTokenExpired
		case vError.Errors&jwt.ValidationErrorNotValidYet != 0:
			return nil, ErrTokenNotValidYet
		default:
			return nil, ErrTokenInvalid
		}
	}

	if claims, ok := jwtToken.Claims.(*MyClaims); ok && jwtToken.Valid {
		return claims, nil
	}
	// 使用类型断言 判断 jwtToken.Claims 是否为 *MyClaims 类型，
	// 类型断言的语法是value, ok := typeAssertion，其中value是断言后的值，ok是一个布尔值，表示断言是否成功。
	// 当类型断言成功时，jwtToken.Claims的值会变成*MyClaims类型，并且我们可以安全地使用claims变量。

	return nil, ErrTokenInvalid
}
