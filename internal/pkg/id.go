package pkg

import "github.com/bwmarrin/snowflake"

var node *snowflake.Node

func init() {
	n, err := snowflake.NewNode(1)
	if err != nil {
		panic("Snowflake 发号器初始化失败")
	}
	node = n
}

func GetId() int64 {
	return node.Generate().Int64()
}
