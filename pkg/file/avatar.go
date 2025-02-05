package file

import (
	"errors"
	"fmt"
	"go_community/settings"
	"math/rand"
	"strings"
	"time"
)

const (
	// AvatarBaseURL 头像存储的基础路径（云存储或本地存储的基础URL）
	//AvatarBaseURL = "static/img/avatar/"
	AvatarBaseURL = ""
	// DefaultAvatarCount 默认头像的数量
	DefaultAvatarCount = 6
)

var (
	// ErrorFileLimit 文件大小超出限制
	ErrorFileLimit = errors.New("文件大小超出限制")
	// ErrorFileType 文件类型不支持
	ErrorFileType = errors.New("文件类型不支持")
	// ErrorFileDirectory 创建目录失败
	ErrorFileDirectory = errors.New("创建目录失败")
)

// GetRandomDefaultAvatar 获取随机默认头像文件名
func GetRandomDefaultAvatar() string {
	rand.Seed(time.Now().UnixNano())
	avatarIndex := rand.Intn(DefaultAvatarCount) + 1 // 1-6
	return fmt.Sprintf("default-avatar_%02d.png", avatarIndex)
}

// GetAvatarPath 获取头像的完整URL路径
func GetAvatarPath(filename string) string {
	if filename == "" {
		return settings.Conf.Avatar.GetDomain() + settings.Conf.Avatar.BaseURL + GetRandomDefaultAvatar()
	}
	// 如果已经是完整URL，直接返回
	if strings.HasPrefix(filename, "http") {
		return filename
	}
	// 返回完整的URL路径，filename 只包含文件名
	return settings.Conf.Avatar.GetDomain() + settings.Conf.Avatar.BaseURL + filename
}

// GenerateAvatarFilename 生成头像文件名
func GenerateAvatarFilename(userId int64, fileExt string) string {
	// 只返回文件名，不包含路径
	return fmt.Sprintf("avatar_%d_%d%s", userId, time.Now().Unix(), fileExt)
}
