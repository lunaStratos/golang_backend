package mapper

import (
	"fmt"

	pb ".server/protocol/target/go/message"
	".server/service/message.service/entity"
)

func MessageList(entities []entity.MessageEntity) []*pb.Msg {
	length := len(entities)
	msgList := make([]*pb.Msg, length)

	for j := 0; j < length; j++ {
		ent := entities[j]
		MsgObj := &pb.Msg{
			Id:              uint32(ent.ID),
			Message:         ent.Message,
			MessageType:     ent.MessageType,
			MessageSubType:  ent.MessageSubType,
			MessageContents: ent.MessageContents,
			ReadYn:          ent.ReadYn,
			DateTime:        fmt.Sprintf("%v", ent.CreatedAt.Unix()),
		}
		msgList[j] = MsgObj
	}
	return msgList
}
