package controller

// 定义业务状态码

type MyCode int64

const (
	CodeSuccess         MyCode = 1000
	CodeInvalidParams   MyCode = 1001
	CodeUserExist       MyCode = 1002
	CodeUserNotExist    MyCode = 1003
	CodeInvalidPassword MyCode = 1004
	CodeServerBusy      MyCode = 1005

	CodeInvalidToken      MyCode = 1006
	CodeInvalidAuthFormat MyCode = 1007
	CodeNotLogin          MyCode = 1008
	CodeVoteRepeated      MyCode = 1009
	CodeVoteTimeExpire    MyCode = 1010
	CodeNoPermission      MyCode = 1011

	CodeFileUploadFailed MyCode = 1012
	CodeFileSizeExceeded MyCode = 1013
	CodeInvalidFileType  MyCode = 1014

	CodeCommunityExist    MyCode = 1015
	CodeCommunityNotExist MyCode = 1016
	CodeCommunityHasPost  MyCode = 1017
)

var msgFlags = map[MyCode]string{
	CodeSuccess:         "success",
	CodeInvalidParams:   "请求参数错误",
	CodeUserExist:       "用户名重复",
	CodeUserNotExist:    "用户不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "服务繁忙",

	CodeInvalidToken:      "无效的Token",
	CodeInvalidAuthFormat: "认证格式有误",
	CodeNotLogin:          "未登录",
	CodeVoteRepeated:      "不允许重复投票",
	CodeVoteTimeExpire:    "投票时间已过",
	CodeNoPermission:      "无操作权限",
	CodeFileUploadFailed:  "文件上传失败",
	CodeFileSizeExceeded:  "文件大小超出限制",
	CodeInvalidFileType:   "不支持的文件类型",

	CodeCommunityExist:    "社区名称已存在",
	CodeCommunityNotExist: "社区不存在",
	CodeCommunityHasPost:  "该社区下还有帖子，无法删除",
}

func (c MyCode) Msg() string {
	msg, ok := msgFlags[c]
	if ok {
		return msg
	}
	return msgFlags[CodeServerBusy]
}
