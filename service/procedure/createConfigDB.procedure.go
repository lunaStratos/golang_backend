package procedure

import (
	"context"

	".server/common"
	configPb ".server/protocol/target/go/config"
)

type ConfigCopyDBProcedure struct {
	configPb.UnimplementedConfigCopyDBProcedureServer
}

func (h *ConfigCopyDBProcedure) CreateCommonCodeDB(ctx context.Context, req *configPb.CreateCommonCodeDBRequest) (*configPb.CreateCommonCodeDBReply, error) {
	rdb := common.GetRDB()
	at := common.ApplyTable

	rdb.Migrator().DropTable(&common.CommonCodeEntity{})

	if at(&common.CommonCodeEntity{}, rdb) {
		for i := 0; i < len(req.CommonCodes); i++ {
			rdb.Create(&common.CommonCodeEntity{
				ParentCode:   req.CommonCodes[i].ParentCode,
				ParentCodeNm: req.CommonCodes[i].ParentCode,
				Code:         req.CommonCodes[i].Code,
				CodeNm:       req.CommonCodes[i].CodeNm,
				Description:  req.CommonCodes[i].Description,
				Is_use:       req.CommonCodes[i].IsUse,
			})
		}
	}

	reply := &configPb.CreateCommonCodeDBReply{
		Code: common.API_STATUS_CODE_200,
	}
	return reply, nil
}
