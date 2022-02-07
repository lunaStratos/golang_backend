package procedure

import (
	configPb ".server/protocol/target/go/config"
	pb ".server/protocol/target/go/message"
	"google.golang.org/grpc"
)

func Register(s *grpc.Server) {
	pb.RegisterEchoProcedureServer(s, &EchoProcedure{})
	pb.RegisterMessageProcedureServer(s, &MessageProcedure{})
	configPb.RegisterConfigCopyDBProcedureServer(s, &ConfigCopyDBProcedure{}) //방법:2

}
