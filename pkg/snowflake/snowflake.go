package snowflake

import (
	"strconv"
	"time"

	sf "github.com/bwmarrin/snowflake"
)

var node *sf.Node

// Init 初始化雪花算法节点
func Init(startTime string, machineId int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}

	// 标准雪花算法的分配：
	// 1位符号位 + 41位时间戳 + 10位机器ID + 12位序列号 = 64位
	sf.NodeBits = uint8(10)
	sf.StepBits = uint8(12)

	// 修改雪花算法的位数分配:
	// 1位符号位 + 41位时间戳 + 5位机器ID + 6位序列号 = 53位
	// sf.NodeBits = uint8(5) // 机器ID位数
	// sf.StepBits = uint8(6) // 序列号位数

	sf.Epoch = st.UnixNano() / 1000000 // 时间起点
	node, err = sf.NewNode(machineId)
	return
}

// GetID 返回生成的id值
func GetID() int64 {
	return node.Generate().Int64()
}

// GetIDStr 返回字符串格式的id值
func GetIDStr() string {
	return strconv.FormatInt(node.Generate().Int64(), 10)
}
