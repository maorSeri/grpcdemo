package main

import (
	"errors"
	"fmt"
	"google.golang.org/grpc/metadata"
	"grpcdemo/pb"
	"io"
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const port = ":9000"

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("first")
		log.Fatal(err)
	}
	creds, err := credentials.NewServerTLSFromFile("cert.pem", "key.pem")
	if err != nil {
		fmt.Println("second")
		log.Fatal(err)
	}

	opts := []grpc.ServerOption{grpc.Creds(creds)}
	s := grpc.NewServer(opts...)
	pb.RegisterEmployeeServiceServer(s, new(employeeService))
	log.Println("Starting server on port: " + port)
	s.Serve(lis)
}

type employeeService struct{}

func (s *employeeService) GetAllEmployees(req *pb.GetAllRequest, streem pb.EmployeeService_GetAllEmployeesServer) error {
	for _, e := range employees {
		streem.Send(&pb.EmployeeResponse{Employee: &e})
	}
	return nil
}

func (s *employeeService) GetEmployeeByBadgeNumber(ctx context.Context, req *pb.GetByBadgeNumberRequest) (*pb.EmployeeResponse, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		fmt.Printf("Metadata received: %v\n", md)
	}
	for _, e := range employees {
		if req.BadgeNumber == e.BadgeNumber {
			return &pb.EmployeeResponse{Employee: &e}, nil
		}
	}
	return nil, errors.New("Employee not found")
}

func (s *employeeService) SaveEmployee(ctx context.Context, req *pb.EmployeeRequest) (*pb.EmployeeResponse, error) {
	employees = append(employees, *req.Employee)
	return &pb.EmployeeResponse{Employee: req.Employee}, nil
}

func (s *employeeService) SaveAll(stream pb.EmployeeService_SaveAllServer) error {
	for {
		emp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		employees = append(employees, *emp.Employee)
		stream.Send(&pb.EmployeeResponse{Employee: emp.Employee})
	}
	for _, e := range employees {
		fmt.Println(e)
	}
	return nil
}

func (s *employeeService) AddPhoto(stream pb.EmployeeService_AddPhotoServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if ok {
		fmt.Printf("Recveiving photo for badge number %v\n", md["badgenumber"][0])
	}

	var imgData []byte
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("File received with length: %v\n", len(imgData))
			return stream.SendAndClose(&pb.AddPhotoResponse{IsOk: true})
		}
		if err != nil {
			return err
		}
		fmt.Printf("Received %v bytes\n", len(data.Data))
		imgData = append(imgData, data.Data...)
	}
	return nil
}
