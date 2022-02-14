package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"grpcdemo/pb"
	"log"
)

const port = ":9000"

func main() {
	option := flag.Int("o", 1, "Command to run")
	creds, err := credentials.NewClientTLSFromFile("cert.pem", "")
	if err != nil {
		fmt.Println("first")
		log.Fatal(err)
	}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	conn, err := grpc.Dial("localhost"+port, opts...)
	if err != nil {
		fmt.Println("second")
		log.Fatal(err)
	}

	defer conn.Close()
	client := pb.NewEmployeeServiceClient(conn)

	switch *option {
	case 1:
		SendMetadata(client)
	}
}

func SendMetadata(client pb.EmployeeServiceClient) {
	md := metadata.MD{}
	md["user"] = []string{"mvansickle"}
	md["password"] = []string{"password1"}
	ctx := context.Background()
	ctx = metadata.NewIncomingContext(ctx, md)
	client.GetEmployeeByBadgeNumber(ctx, &pb.GetByBadgeNumberRequest{})
}
