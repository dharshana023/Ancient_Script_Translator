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
        lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
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

        // Process, translate the manuscript, and extract metadata
        translatedText, metadata, err := s.serviceHandler.ProcessTranslateWithMetadata(req.ManuscriptImage, req.ScriptType)
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

        // Convert Go metadata to protobuf metadata
        metadataProto := &pb.MetadataResponse{
                ScriptType:      metadata.ScriptType,
                ConfidenceScore: metadata.ConfidenceScore,
                DetectedDate:    metadata.DetectedDate,
        }

        // Add time periods
        for _, period := range metadata.TimePeriods {
                metadataProto.TimePeriods = append(metadataProto.TimePeriods, &pb.TimePeriod{
                        Name:        period.Name,
                        StartYear:   int32(period.StartYear),
                        EndYear:     int32(period.EndYear),
                        Description: period.Description,
                })
        }

        // Add regions
        for _, region := range metadata.Regions {
                metadataProto.Regions = append(metadataProto.Regions, &pb.Region{
                        Name:        region.Name,
                        ModernAreas: region.ModernAreas,
                        Description: region.Description,
                })
        }

        // Add cultural context
        metadataProto.CulturalContext = metadata.CulturalContext
        metadataProto.MaterialContext = metadata.MaterialContext

        // Add historical events
        for _, event := range metadata.HistoricalEvents {
                metadataProto.HistoricalEvents = append(metadataProto.HistoricalEvents, &pb.HistoricalEvent{
                        Name:        event.Name,
                        EventType:   event.EventType,
                        Year:        int32(event.Year),
                        Description: event.Description,
                })
        }

        // Create response
        return &pb.TranslateResponse{
                OriginalScript: req.ScriptType,
                TranslatedText: translatedText,
                Summary:        summary,
                Metadata:       metadataProto,
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
