package snowflake

import (
	"time"

	sf "github.com/bwmarrin/snowflake"
)

var node *sf.Node

// Init 需传入当前的机器ID
func Init(startTime string, machineId int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}

	sf.Epoch = st.UnixNano() / 1000000
	node, err = sf.NewNode(machineId)
	return
}

// GetID 返回生成的id值
func GetID() int64 {
	return node.Generate().Int64()
}
