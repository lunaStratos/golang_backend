package procedure

import (
	"context"
	"fmt"
	"strconv"

	".server/common"
	pb ".server/protocol/target/go/message"
	".server/service/message.service/entity"
	".server/service/message.service/mapper"
	"gorm.io/gorm"
)

type MessageProcedure struct {
	pb.UnimplementedMessageProcedureServer
}

// 메시지 리스트 얻기
func (p *MessageProcedure) GetMessageList(
	ctx context.Context,
	rq *pb.GetMessageListRequest) (*pb.GetMessageListReply, error) {

	rdb := common.GetRDB()

	// 리턴 변수 선언
	var messageEntityList []entity.MessageEntity

	println(uint32(common.GetSession(ctx).Id))
	search := rq.GetSearch()
	userId := uint32(common.GetSession(ctx).Id)
	pageNum := rq.GetPageNum()
	maxLenInPage := int(rq.GetMaxLenInPage())
	categoryType := rq.GetCategoryType()
	mode := rq.GetMode()

	query := rdb.Table("t_message").Where("user_id = ?", userId)

	// 검색어 있을때
	if search != "" {
		filter := fmt.Sprintf("%s%s%s", "%", search, "%")
		query = query.Where("message LIKE ?", filter)
	}

	// categoryType, ALL 아니면 동작
	if categoryType != common.MESSAGE_STATUS_ALL {
		query = query.Where("message_type = ?", categoryType)
	}

	// 헤더의 메시지는 신규만 보여준다.
	if mode == "header" {
		query = query.Where("read_yn", "N")
	}

	// 전체 카운트
	var messageCount int64
	query.Where("deleted_at is null").Count(&messageCount)

	getPageCal := common.BoardPaging(int32(messageCount), int32(pageNum), int32(maxLenInPage))

	// 페이지 있을때
	if pageNum != 0 {
		query = query.Limit(maxLenInPage).Offset(int(getPageCal.Start))
	}

	query.Order("id desc").Find(&messageEntityList)

	var newMessageCount int64
	query.Where("read_yn", "N").Count(&newMessageCount)

	messageList := mapper.MessageList(messageEntityList)

	// 리턴 값 set
	res := &pb.GetMessageListReply{
		Code:      common.API_STATUS_CODE_200,
		Message:   "success!",
		Msg:       messageList,
		AllCount:  uint32(messageCount),
		NewCount:  uint32(newMessageCount),
		StartPage: uint32(getPageCal.StartPage),
		EndPage:   uint32(getPageCal.EndPage),
		LastPage:  uint32(getPageCal.LastPage),
		NowPage:   uint32(getPageCal.NowPage),
	}
	return res, nil
}

// 판매 메시지 생성
func (p *MessageProcedure) CreateSaleMessage(ctx context.Context, rq *pb.CreateSaleMessageRequest) (*pb.CreateMessageReply, error) {

	userId := rq.GetUserId()
	messageSubType := rq.GetMessageSubType()
	region := rq.GetRegion()
	country := rq.GetCountry()
	buyerName := rq.GetBuyerName()

	msgObject := rq.GetMessageContent()

	title := msgObject.GetTitle()
	countrys := msgObject.GetCountry()
	stateProvince := msgObject.GetStateProvince()
	blocks := msgObject.GetTileCount()
	level := msgObject.GetLevel()
	from := msgObject.GetFrom()
	price := msgObject.GetPrice()

	var makeMessage = ""
	var makeMessageContents = ""

	switch messageSubType {
	case common.MESSAGE_STATUS_SUB_001: // 내 LAND 판매 완료
		makeMessage = fmt.Sprintf("%s/%s LAND was purchased by %s. Check out WALLET.", region, country, buyerName)
		makeMessageContents = fmt.Sprintf("* Title : %s \n* Country : %s \n* State/Province : %s \n* Blocks : %d \n* Level : %d \n* From : %s \n* Price : %g",
			title, countrys, stateProvince, blocks, level, from, price)
		break
	}

	rdb := common.GetRDB()

	messageData := &entity.MessageEntity{
		MessageType:     common.MESSAGE_STATUS_SUB_002,
		MessageSubType:  messageSubType,
		MessageContents: makeMessageContents,
		Message:         makeMessage,
		UserId:          userId,
	}

	//생성
	rdb.Create(messageData)

	res := &pb.CreateMessageReply{
		Code:    common.API_STATUS_CODE_200,
		Message: "success!",
	}

	return res, nil
}

// 판매 메시지 생성
func (p *MessageProcedure) CreateNoticeMessage(ctx context.Context, rq *pb.CreateNoticeMessageRequest) (*pb.CreateMessageReply, error) {

	userId := rq.GetUserId()
	messageSubType := rq.GetMessageSubType()
	msgObject := rq.GetMessageContent()

	dateOfApplication := msgObject.GetDateOfApplication()
	applicationAmount := msgObject.GetApplicationAmount()
	reviewTheResult := msgObject.GetReviewTheResult()
	withdrawFee := msgObject.GetWithdrawFee()
	title := msgObject.GetTitle()
	description := msgObject.GetDescription()

	var makeMessage = ""
	var makeMessageContents = ""

	switch messageSubType {
	case common.MESSAGE_STATUS_SUB_001: // 취소/환불 신청 결과
		makeMessage = "Cancellation/Refund\nYour requested cancellation/refund processing has been reviewed."
		makeMessageContents = fmt.Sprintf("* Date of application: %s \n* Application amount: %d \n* Review the results : %s ",
			dateOfApplication, applicationAmount, reviewTheResult)
		break
	case common.MESSAGE_STATUS_SUB_002: // 출금신청결과
		makeMessage = "Withdraw\nThe review of application for withdrawal of holdings has been completed."
		makeMessageContents = fmt.Sprintf("* Date of application: %s \n* Application amount: %d \n* Review the results : %s \n * Withdrawal fee: %g ",
			dateOfApplication, applicationAmount, reviewTheResult, withdrawFee)
		break
	case common.MESSAGE_STATUS_SUB_003: // 일반운영공지
		makeMessage = "Notice\nThere is a new operation notice. Please check."
		makeMessageContents = fmt.Sprintf("* Title : %s \n* Description : %s",
			title, description)
		break

	case common.MESSAGE_STATUS_SUB_004: // 마케팅안내
		makeMessage = "event\nThere is a new event information. Join us."
		makeMessageContents = fmt.Sprintf("* Title : %s  \n* Description : %s",
			title, description)
		break
	case common.MESSAGE_STATUS_SUB_005: // 시스템공지
		makeMessage = "System guide\nThere are major announcements related to system operation. Please check."
		makeMessageContents = fmt.Sprintf("* Title : %s \n* Description : %s",
			title, description)
		break
	}

	rdb := common.GetRDB()

	messageData := &entity.MessageEntity{
		MessageType:     common.MESSAGE_STATUS_SUB_005,
		MessageSubType:  messageSubType,
		MessageContents: makeMessageContents,
		Message:         makeMessage,
		UserId:          userId,
	}

	//생성
	rdb.Create(messageData)

	res := &pb.CreateMessageReply{
		Code:    common.API_STATUS_CODE_200,
		Message: "success!",
	}

	return res, nil
}

// 메시지 삭제 처리
func (p *MessageProcedure) DoMessageDelete(ctx context.Context, rq *pb.DoMessageRequest) (*pb.DoMessageReply, error) {

	msgId := rq.GetMsgId()
	rdb := common.GetRDB()

	deleteMessage(rdb, msgId)

	res := &pb.DoMessageReply{
		Code:    common.API_STATUS_CODE_200,
		Message: "delete success!",
	}

	return res, nil
}

// 메시지 읽음 처리
func (p *MessageProcedure) DoMessageRead(ctx context.Context, rq *pb.DoMessageRequest) (*pb.DoMessageReply, error) {

	msgId := rq.GetMsgId()
	rdb := common.GetRDB()

	readMessage(rdb, msgId)

	res := &pb.DoMessageReply{
		Code:    common.API_STATUS_CODE_200,
		Message: "read success!",
	}

	return res, nil
}

func (p *MessageProcedure) DoMessageReadMarked(ctx context.Context, rq *pb.DoMessageListRequest) (*pb.DoMessageReply, error) {

	msgIdList := rq.GetMsgIds()
	rdb := common.GetRDB()

	msgList := []string{}
	for i := 0; i < len(msgIdList); i++ {
		msgList = append(msgList, strconv.Itoa(int(msgIdList[i])))
	}

	rdb.Model(&entity.MessageEntity{}).Where("id in ?", msgList).Update("read_yn", "Y")

	res := &pb.DoMessageReply{
		Code:    common.API_STATUS_CODE_200,
		Message: "read success!",
	}

	return res, nil
}

func (p *MessageProcedure) DoMessageDeleteMarked(ctx context.Context, rq *pb.DoMessageListRequest) (*pb.DoMessageReply, error) {

	msgIdList := rq.GetMsgIds()
	rdb := common.GetRDB()

	err := rdb.Where("id in ?", msgIdList).Delete(&entity.MessageEntity{}).Error
	if err != nil {
		res := &pb.DoMessageReply{
			Code:    common.API_STATUS_CODE_500,
			Message: "error",
		}
		return res, err
	}

	res := &pb.DoMessageReply{
		Code:    common.API_STATUS_CODE_200,
		Message: "read success!",
	}

	return res, nil
}

// 삭제 DB 모듈
func deleteMessage(tx *gorm.DB, msgId uint32) error {
	err := tx.Where("id", msgId).Delete(&entity.MessageEntity{}).Error
	if err != nil {
		return err
	}

	return nil
}

// 읽음 DB 모듈
func readMessage(tx *gorm.DB, msgId uint32) error {
	err := tx.Model(&entity.MessageEntity{}).
		Where("id = ?", msgId).
		Update("read_yn", "Y").Error

	if err != nil {
		return err
	}

	return nil
}
