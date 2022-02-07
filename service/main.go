package main

import (
	"server/common"
	"server/service/message.service/entity"
	"server/service/message.service/procedure"
)

func main() {

	// env load 및 argument 처리
	mode := common.LoadEnv("message")

	entity.Migrate()

	if mode == "grpc" {
		common.GrpcServe(procedure.Register)
	}

}
