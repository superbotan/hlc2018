package main

import (
	"go_v1/common"
	"time"
)

var (
	accounts = AccountStorageCreate()

	cities     = common.SimpleDict16Create()
	countries  = common.SimpleDict8Create()
	fnames     = common.SimpleDict16Create()
	snames     = common.SimpleDict16Create()
	domains    = common.SimpleDict8Create()
	interests  = common.SimpleDict8Create()
	phonecodes = common.SimpleDict8Create()

	index = GlobalIndexCreate()

	//likeindex = LikeIndexCreate()

	//likeindex2 = LikesindexV2Create()

	groupindex = GroupIndexCreate()

	nowtime = int32(0)
)

func GetTimeNow() int32 {
	if nowtime == 0 {
		return int32(time.Now().Unix())
	}

	return nowtime
}
