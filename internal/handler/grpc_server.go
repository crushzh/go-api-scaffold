package handler

import (
	"context"

	"go-api-scaffold/internal/model"
	"go-api-scaffold/internal/service"
	pb "go-api-scaffold/api/proto/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ExampleGRPCServer implements the gRPC service
type ExampleGRPCServer struct {
	pb.UnimplementedExampleServiceServer
	svc *service.ExampleService
}

// NewGRPCServer creates and registers gRPC services
func NewGRPCServer(exampleSvc *service.ExampleService) *grpc.Server {
	s := grpc.NewServer()

	pb.RegisterExampleServiceServer(s, &ExampleGRPCServer{svc: exampleSvc})
	// GEN:GRPC_REGISTER - Auto-appended by code generator, do not remove

	return s
}

// GetExample returns an example by ID
func (s *ExampleGRPCServer) GetExample(ctx context.Context, req *pb.GetExampleRequest) (*pb.ExampleResponse, error) {
	item, err := s.svc.GetByID(uint(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "not found: %v", err)
	}

	return &pb.ExampleResponse{
		Id:          uint64(item.ID),
		Name:        item.Name,
		Description: item.Description,
		Status:      item.Status,
	}, nil
}

// ListExamples returns a paginated list
func (s *ExampleGRPCServer) ListExamples(ctx context.Context, req *pb.ListExamplesRequest) (*pb.ListExamplesResponse, error) {
	query := &model.QueryExampleRequest{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
		Keyword:  req.Keyword,
	}

	items, total, err := s.svc.List(query)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list failed: %v", err)
	}

	pbItems := make([]*pb.ExampleResponse, len(items))
	for i, item := range items {
		pbItems[i] = &pb.ExampleResponse{
			Id:          uint64(item.ID),
			Name:        item.Name,
			Description: item.Description,
			Status:      item.Status,
		}
	}

	return &pb.ListExamplesResponse{
		Items: pbItems,
		Total: total,
	}, nil
}

// CreateExample creates a new example
func (s *ExampleGRPCServer) CreateExample(ctx context.Context, req *pb.CreateExampleRequest) (*pb.ExampleResponse, error) {
	item, err := s.svc.Create(&model.CreateExampleRequest{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create failed: %v", err)
	}

	return &pb.ExampleResponse{
		Id:          uint64(item.ID),
		Name:        item.Name,
		Description: item.Description,
		Status:      item.Status,
	}, nil
}
