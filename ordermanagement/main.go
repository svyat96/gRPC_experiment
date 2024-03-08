package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	pb "ordermanagement/output"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterOrderManagementServer(s, &server{})

	log.Printf("Starting gRPC " + port)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

type server struct {
	orderMap map[string]*pb.Order
}

func (s *server) GetOrder(ctx context.Context, orderId *wrapperspb.StringValue) (*pb.Order, error) {
	return s.orderMap[orderId.Value], nil
}

func (s *server) SearchOrders(searchQuery *wrapperspb.StringValue, stream pb.OrderManagement_SearchOrdersServer) error {
	fmt.Println("Send query: ", searchQuery)
	stream.Send(&pb.Order{Id: "200", Items: []string{"a", "b"}, Description: "Что то нашел!", Price: 320, Destination: "В гугле конечно!"})
	stream.Send(&pb.Order{Id: "201", Items: []string{"a", "b"}, Description: "Что то нашел!", Price: 320, Destination: "В гугле конечно!"})
	stream.Send(&pb.Order{Id: "202", Items: []string{"a", "b"}, Description: "Что то нашел!", Price: 320, Destination: "В гугле конечно!"})
	stream.Send(&pb.Order{Id: "203", Items: []string{"a", "b"}, Description: "Что то нашел!", Price: 320, Destination: "В гугле конечно!"})
	for key, order := range s.orderMap {
		log.Print(key, order)
		for _, itemStr := range order.Items {
			log.Print(itemStr)
			if strings.Contains(
				itemStr, searchQuery.Value,
			) {
				err := stream.Send(order)
				if err != nil {
					return fmt.Errorf("Error sending message: %v", err)
				}
				log.Print("Matching Order Found: " + key)
				break
			}
		}
	}
	return nil
}
