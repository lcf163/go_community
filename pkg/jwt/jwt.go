package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// MyClaims 自定义声明结构体并内嵌 jwt.StandardClaims
// jwt 包自带的 jwt.StandardClaims 只包含了官方字段
// 我们这里需要额外记录一个 UserID 字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

// mySecret 用于加密的字符串
var mySecret = []byte("夏天夏天悄悄过去")

func keyFunc(_ *jwt.Token) (i interface{}, err error) {
	return mySecret, nil
}

// AccessTokenExpireDuration JWT access_token 过期时间
const AccessTokenExpireDuration = time.Second * 24

// RefreshTokenExpireDuration JWT refresh_token 过期时间
const RefreshTokenExpireDuration = time.Hour * 24 * 7

// GenToken 生成JWT
func GenToken(userID int64) (aToken, rToken string, err error) {
	// 创建一个自己声明的数据
	c := MyClaims{
		userID, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(AccessTokenExpireDuration).Unix(), // 过期时间
			Issuer:    "go-community",                                   // 签发人
		},
	}
	// 加密并获得完整的编码后的字符串token
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(mySecret)

	// refresh token 不需要存任何自定义数据
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(RefreshTokenExpireDuration).Unix(), // 过期时间
		Issuer:    "go-community",                                    // 签发人
	}).SignedString(mySecret)

	// 使用指定的secret签名并获得完整的编码后的字符串token
	return
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (claims *MyClaims, err error) {
	// 解析token
	var token *jwt.Token
	claims = new(MyClaims)
	token, err = jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		return
	}
	if !token.Valid { // 校验token
		err = errors.New("invalid token")
	}
	return
}

// RefreshToken 刷新AccessToken
func RefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
	// refresh token 无效直接返回
	if _, err = jwt.Parse(rToken, keyFunc); err != nil {
		return
	}

	// 从旧 access token 中解析出 claims 数据
	var claims MyClaims
	_, err = jwt.ParseWithClaims(aToken, &claims, keyFunc)
	v, _ := err.(*jwt.ValidationError)

	// 当 access token 是过期错误并且 refresh token 没有过期时，就创建一个新的 access token
	if v.Errors == jwt.ValidationErrorExpired {
		return GenToken(claims.UserID)
	}
	return
}
