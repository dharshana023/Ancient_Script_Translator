package api

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "ancient-script-decoder/proto"
	"ancient-script-decoder/services"
	"ancient-script-decoder/utils"
)

// GRPCServer represents the gRPC server
type GRPCServer struct {
	port           int
	serviceHandler *services.ServiceHandler
	server         *grpc.Server
	logger         *utils.Logger
}

// NewGRPCServer creates a new gRPC server
func NewGRPCServer(port int, serviceHandler *services.ServiceHandler, logger *utils.Logger) *GRPCServer {
	return &GRPCServer{
		port:           port,
		serviceHandler: serviceHandler,
		logger:         logger,
	}
}

// Start starts the gRPC server
func (s *GRPCServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// Create gRPC server
	s.server = grpc.NewServer()
	
	// Register services
	pb.RegisterTranslatorServiceServer(s.server, s)
	
	// Register reflection service for grpcurl
	reflection.Register(s.server)

	// Start the server
	return s.server.Serve(lis)
}

// Stop stops the gRPC server
func (s *GRPCServer) Stop() {
	if s.server != nil {
		s.server.GracefulStop()
	}
}

// TranslateManuscript handles the translation request via gRPC
func (s *GRPCServer) TranslateManuscript(ctx context.Context, req *pb.TranslateRequest) (*pb.TranslateResponse, error) {
	s.logger.Info("Received gRPC translation request", "scriptType", req.ScriptType)

	// Process and translate the manuscript
	translatedText, err := s.serviceHandler.ProcessAndTranslate(req.ManuscriptImage, req.ScriptType)
	if err != nil {
		s.logger.Error("Failed to process and translate manuscript", "error", err)
		return nil, fmt.Errorf("failed to process and translate manuscript: %v", err)
	}

	// Generate summary for the translated text
	summary, err := s.serviceHandler.SummarizeText(translatedText)
	if err != nil {
		s.logger.Error("Failed to generate summary", "error", err)
		return nil, fmt.Errorf("failed to generate summary: %v", err)
	}

	// Create response
	return &pb.TranslateResponse{
		OriginalScript: req.ScriptType,
		TranslatedText: translatedText,
		Summary:        summary,
	}, nil
}

// SummarizeText handles the summarization request via gRPC
func (s *GRPCServer) SummarizeText(ctx context.Context, req *pb.SummarizeRequest) (*pb.SummarizeResponse, error) {
	s.logger.Info("Received gRPC summarization request", "textLength", len(req.Text))

	// Validate request
	if req.Text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}

	// Generate summary
	summary, err := s.serviceHandler.SummarizeText(req.Text)
	if err != nil {
		s.logger.Error("Failed to generate summary", "error", err)
		return nil, fmt.Errorf("failed to generate summary: %v", err)
	}

	// Create response
	return &pb.SummarizeResponse{
		Summary:    summary,
		TextLength: int32(len(req.Text)),
	}, nil
}
