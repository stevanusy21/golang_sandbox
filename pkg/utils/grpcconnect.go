package utils

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ConnectGrpc(serviceName, address string) (*grpc.ClientConn, error) {
	LogInfo(serviceName, "Mencoba konek ke "+serviceName+" (gRPC) di "+address+"...")

	grpcConn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		LogFatal(serviceName, "Gagal konek ke "+serviceName, err)
	}

	return grpcConn, nil
}
