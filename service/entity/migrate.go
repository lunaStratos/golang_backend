package entity

import "server/common"

// Migrate db migration
func Migrate() {

	rdb := common.GetRDB()
	at := common.ApplyTable

	//RDB 생성
	at(&MessageEntity{}, rdb)
}
