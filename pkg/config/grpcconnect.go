package config

import (
	"github.com/stevanusy21/golang_sandbox/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ConnectGrpc(serviceName, address string) (*grpc.ClientConn, error) {
	utils.LogInfo(serviceName, "Mencoba konek ke "+serviceName+" (gRPC) di "+address+"...")

	grpcConn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		utils.LogFatal(serviceName, "Gagal konek ke "+serviceName, err)
	}

	return grpcConn, nil
}
