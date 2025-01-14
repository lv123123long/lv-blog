package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenAndParseToken(t *testing.T) {
	secret := "secret"
	issuer := "issuer"
	expire := 10

	token, err := GenToken(secret, issuer, expire, 1, []int{1, 2})
	// 断言生成过程中，没有错误
	assert.Nil(t, err)
	// 断言生成的JWT不为空
	assert.NotEmpty(t, token)

	mc, err := ParseToken(secret, token)
	assert.Nil(t, err)
	// 断言解析后的用户ID为1
	assert.Equal(t, 1, mc.UserId)
	// 断言解析后的角色数组长度为2
	assert.Len(t, mc.RoleIds, 2)
}

func TestParseTokenError(t *testing.T) {
	tokenString := "tokenString"

	_, err := ParseToken("secret", tokenString)
	assert.ErrorIs(t, err, ErrTokenMalformed)
}
