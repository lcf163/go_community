package file

import (
	"fmt"
	"math/rand"
	"path"
	"time"
)

const (
	// AvatarBaseURL 头像存储的基础路径（云存储或本地存储的基础URL）
	//AvatarBaseURL = "static/img/avatar/"
	AvatarBaseURL = ""
	// DefaultAvatarCount 默认头像的数量
	DefaultAvatarCount = 6
)

// GetRandomDefaultAvatar 获取随机默认头像
func GetRandomDefaultAvatar() string {
	rand.Seed(time.Now().UnixNano())
	avatarIndex := rand.Intn(DefaultAvatarCount) + 1 // 1-6
	return path.Join(AvatarBaseURL, fmt.Sprintf("default_%02d.png", avatarIndex))
}

// GetAvatarPath 获取头像的完整相对路径
func GetAvatarPath(filename string) string {
	if filename == "" {
		return GetRandomDefaultAvatar()
	}
	return path.Join(AvatarBaseURL, filename)
}

// GenerateAvatarFilename 生成头像文件名
func GenerateAvatarFilename(userId int64, fileExt string) string {
	return fmt.Sprintf("avatar_%d_%d%s", userId, time.Now().Unix(), fileExt)
}
