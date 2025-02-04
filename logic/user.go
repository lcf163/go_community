package logic

import (
	"go_community/dao/mysql"
	"go_community/models"
	pkg_file "go_community/pkg/file"
	"go_community/pkg/jwt"
	"go_community/pkg/snowflake"
	"go_community/settings"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

// 存放业务逻辑的代码

// SignUp 注册业务逻辑
func SignUp(p *models.ParamSignUp) (err error) {
	// 判断用户是否存在
	if err := mysql.CheckUserExist(p.UserName); err != nil {
		return err
	}
	// 生成 UID
	UserID := snowflake.GetID()
	// 构造一个 User 实例
	user := &models.User{
		UserID:   UserID,
		UserName: p.UserName,
		Password: p.Password,
	}
	// 保存进数据库
	return mysql.InsertUser(user)
	// redis.xxx
}

// Login 登录业务逻辑
func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		UserName: p.UserName,
		Password: p.Password,
	}
	// 用户登录，传递的是指针
	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	// 生成 JWT token
	//return jwt.GenToken(user.UserID)
	accessToken, refreshToken, err := jwt.GenToken(user.UserID)
	if err != nil {
		return
	}
	user.AccessToken = accessToken
	user.RefreshToken = refreshToken
	return
}

// GetUserInfo 获取用户信息
func GetUserInfo(UserID int64) (*models.User, error) {
	// 从数据库中查询用户信息
	user, err := mysql.GetUserById(UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, mysql.ErrorUserNotExist
	}
	return user, nil
}

// UpdateUserName 更新用户名
func UpdateUserName(UserID int64, p *models.ParamUpdateUser) error {
	// 检查用户名是否已存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		return err
	}

	// 更新用户名
	return mysql.UpdateUserName(UserID, p)
}

// UpdatePassword 修改密码
func UpdatePassword(UserID int64, p *models.ParamUpdatePassword) error {
	// 验证旧密码是否正确
	if err := mysql.CheckPassword(UserID, p.OldPassword); err != nil {
		return err
	}

	// 更新密码
	return mysql.UpdatePassword(UserID, p.NewPassword)
}

// UpdateAvatar 更新用户头像
func UpdateAvatar(UserID int64, file *multipart.FileHeader) (string, error) {
	// 检查文件大小
	if file.Size > settings.Conf.Avatar.MaxSize {
		return "", pkg_file.ErrorFileLimit
	}

	// 检查文件类型
	ext := path.Ext(file.Filename)
	if !isValidImageExt(ext) {
		return "", pkg_file.ErrorFileType
	}

	// 获取用户当前头像
	user, err := mysql.GetUserById(UserID)
	if err != nil {
		return "", err
	}

	// 生成新文件名（不包含域名和路径）
	filename := pkg_file.GenerateAvatarFilename(UserID, ext)

	// 确保上传目录存在
	uploadDir := settings.Conf.Avatar.BaseURL
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", pkg_file.ErrorFileDirectory
	}

	// 如果用户已有头像且不是外部URL，删除原头像文件
	if user.Avatar != "" && !strings.HasPrefix(user.Avatar, "http") {
		oldAvatarPath := path.Join(uploadDir, user.Avatar)
		// 忽略删除错误，因为文件可能不存在
		_ = os.Remove(oldAvatarPath)
	}

	// 保存新文件
	dst := path.Join(uploadDir, filename)
	srcFile, err := file.Open()
	if err != nil {
		return "", err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return "", err
	}

	// 更新数据库中的头像路径
	if err := mysql.UpdateUserAvatar(UserID, filename); err != nil {
		// 如果数据库更新失败，删除已上传的文件
		os.Remove(dst)
		return "", err
	}

	return filename, nil
}

// isValidImageExt 检查是否为有效的图片扩展名
func isValidImageExt(ext string) bool {
	validExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}
	return validExts[strings.ToLower(ext)]
}
