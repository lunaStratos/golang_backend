package entity

import "server/common"

// 메시지
type MessageEntity struct {
	common.Model
	UserId          uint32 `gorm:"type:int;not null"`
	Message         string `gorm:"type:varchar(100);not null"`
	MessageContents string `gorm:"type:varchar(200);not null"`
	MessageType     string `gorm:"type:varchar(3);not null"`
	MessageSubType  string `gorm:"type:varchar(3);not null"`
	ReadYn          string `gorm:"type:varchar(1);DEFAULT:'N'"`
}
