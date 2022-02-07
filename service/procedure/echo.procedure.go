package procedure

import (
	"context"

	pb ".server/protocol/target/go/message"
)

type EchoProcedure struct {
	pb.UnimplementedEchoProcedureServer
}

func (p *EchoProcedure) Echo(ctx context.Context, rq *pb.EchoRequest) (*pb.EchoReply, error) {
	return &pb.EchoReply{Msg: rq.GetMsg()}, nil
}
