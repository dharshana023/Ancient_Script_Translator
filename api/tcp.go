package api

import (
        "bufio"
        "context"
        "encoding/json"
        "fmt"
        "net"
        "sync"
        "time"

        "ancient-script-decoder/models"
        "ancient-script-decoder/services"
        "ancient-script-decoder/utils"
)

// TCPConfig contains configuration for the TCP server
type TCPConfig struct {
        Enabled        bool   `yaml:"enabled"`
        Port           int    `yaml:"port"`
        Host           string `yaml:"host"`
        MaxConnections int    `yaml:"maxConnections"`
        Timeout        int    `yaml:"timeout"`
}

// TCPServer represents a TCP server for the ancient script decoder
type TCPServer struct {
        config         TCPConfig
        serviceHandler *services.ServiceHandler
        logger         *utils.Logger
        listener       net.Listener
        connections    map[string]net.Conn
        connectionsMu  sync.Mutex
        ctx            context.Context
        cancel         context.CancelFunc
}

// NewTCPServer creates a new TCP server
func NewTCPServer(config TCPConfig, serviceHandler *services.ServiceHandler, logger *utils.Logger) *TCPServer {
        ctx, cancel := context.WithCancel(context.Background())
        
        if config.Host == "" {
                config.Host = "0.0.0.0" // Default to all interfaces
        }
        
        return &TCPServer{
                config:         config,
                serviceHandler: serviceHandler,
                logger:         logger,
                connections:    make(map[string]net.Conn),
                ctx:            ctx,
                cancel:         cancel,
        }
}

// Start starts the TCP server
func (s *TCPServer) Start() error {
        if !s.config.Enabled {
                s.logger.Info("TCP server is disabled")
                return nil
        }
        
        // Start listening for connections
        addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
        listener, err := net.Listen("tcp", addr)
        if err != nil {
                return fmt.Errorf("failed to start TCP server: %v", err)
        }
        
        s.listener = listener
        s.logger.Info("TCP server started", "address", addr)
        
        // Start accepting connections in a goroutine
        go s.acceptConnections()
        
        return nil
}

// acceptConnections accepts new connections
func (s *TCPServer) acceptConnections() {
        for {
                // Check if context is canceled
                select {
                case <-s.ctx.Done():
                        return
                default:
                        // Continue accepting connections
                }
                
                // Accept a new connection
                conn, err := s.listener.Accept()
                if err != nil {
                        select {
                        case <-s.ctx.Done():
                                // Server is shutting down, ignore the error
                                return
                        default:
                                s.logger.Error("Failed to accept connection", "error", err)
                                // Sleep briefly to avoid spinning on accept errors
                                time.Sleep(100 * time.Millisecond)
                                continue
                        }
                }
                
                // Check if we have reached the maximum number of connections
                s.connectionsMu.Lock()
                if len(s.connections) >= s.config.MaxConnections {
                        s.connectionsMu.Unlock()
                        s.logger.Error("Maximum number of connections reached, rejecting connection")
                        conn.Close()
                        continue
                }
                
                // Add the connection to the map
                connID := conn.RemoteAddr().String()
                s.connections[connID] = conn
                s.connectionsMu.Unlock()
                
                // Handle the connection in a goroutine
                go s.handleConnection(conn, connID)
        }
}

// handleConnection handles a TCP client connection
func (s *TCPServer) handleConnection(conn net.Conn, connID string) {
        defer func() {
                conn.Close()
                s.connectionsMu.Lock()
                delete(s.connections, connID)
                s.connectionsMu.Unlock()
        }()
        
        // Set a read deadline based on the timeout
        if s.config.Timeout > 0 {
                deadline := time.Now().Add(time.Duration(s.config.Timeout) * time.Second)
                conn.SetDeadline(deadline)
        }
        
        s.logger.Info("New TCP connection established", "client", connID)
        
        // Initialize a scanner to read requests
        scanner := bufio.NewScanner(conn)
        
        // Read and process requests until the connection is closed
        for scanner.Scan() {
                // Reset the deadline for each request
                if s.config.Timeout > 0 {
                        deadline := time.Now().Add(time.Duration(s.config.Timeout) * time.Second)
                        conn.SetDeadline(deadline)
                }
                
                // Get the request data
                data := scanner.Text()
                s.logger.Info("Received TCP request", "client", connID, "dataLength", len(data))
                
                // Process the request
                response, err := s.processRequest(data)
                if err != nil {
                        s.logger.Error("Failed to process TCP request", "client", connID, "error", err)
                        errorResponse := fmt.Sprintf("ERROR: %v\n", err)
                        conn.Write([]byte(errorResponse))
                        continue
                }
                
                // Send the response
                responseData, err := json.Marshal(response)
                if err != nil {
                        s.logger.Error("Failed to marshal TCP response", "client", connID, "error", err)
                        conn.Write([]byte("ERROR: Failed to marshal response\n"))
                        continue
                }
                
                // Append a newline to mark the end of the response
                responseData = append(responseData, '\n')
                
                _, err = conn.Write(responseData)
                if err != nil {
                        s.logger.Error("Failed to send TCP response", "client", connID, "error", err)
                        return
                }
                
                s.logger.Info("Sent TCP response", "client", connID, "responseLength", len(responseData))
        }
        
        if err := scanner.Err(); err != nil {
                s.logger.Error("Scanner error", "client", connID, "error", err)
        }
}

// processRequest processes a client request
func (s *TCPServer) processRequest(data string) (interface{}, error) {
        // Parse the request
        var request map[string]interface{}
        if err := json.Unmarshal([]byte(data), &request); err != nil {
                return nil, fmt.Errorf("failed to parse request: %v", err)
        }
        
        // Get the request type
        requestType, ok := request["type"].(string)
        if !ok {
                return nil, fmt.Errorf("missing or invalid request type")
        }
        
        // Process the request based on the type
        switch requestType {
        case "summarize":
                return s.handleSummarizeRequest(request)
        case "metadata":
                return s.handleMetadataRequest(request)
        default:
                return nil, fmt.Errorf("unsupported request type: %s", requestType)
        }
}

// handleSummarizeRequest handles a summarization request
func (s *TCPServer) handleSummarizeRequest(request map[string]interface{}) (interface{}, error) {
        // Get the text to summarize
        text, ok := request["text"].(string)
        if !ok || text == "" {
                return nil, fmt.Errorf("missing or invalid text")
        }
        
        // Get the algorithm (optional)
        algorithm := ""
        if alg, ok := request["algorithm"].(string); ok {
                algorithm = alg
        }
        
        // Generate the summary
        summary, err := s.serviceHandler.SummarizeTextWithAlgorithm(text, algorithm)
        if err != nil {
                return nil, fmt.Errorf("failed to generate summary: %v", err)
        }
        
        // Return the response
        return models.SummarizeResponse{
                Summary:     summary,
                TextLength:  len(text),
                ProcessedAt: time.Now().Format(time.RFC3339),
        }, nil
}

// handleMetadataRequest handles a metadata extraction request
func (s *TCPServer) handleMetadataRequest(request map[string]interface{}) (interface{}, error) {
        // Get the text
        text, ok := request["text"].(string)
        if !ok || text == "" {
                return nil, fmt.Errorf("missing or invalid text")
        }
        
        // Get the script type (optional)
        scriptType := "auto"
        if st, ok := request["scriptType"].(string); ok && st != "" {
                scriptType = st
        }
        
        // Extract metadata
        metadata, err := s.serviceHandler.ExtractMetadata(text, scriptType)
        if err != nil {
                return nil, fmt.Errorf("failed to extract metadata: %v", err)
        }
        
        // Return the response
        return map[string]interface{}{
                "metadata":    metadata,
                "scriptType":  scriptType,
                "processedAt": time.Now().Format(time.RFC3339),
        }, nil
}

// Stop stops the TCP server
func (s *TCPServer) Stop() {
        s.cancel() // Cancel the context to signal shutdown
        
        // Close the listener if it exists
        if s.listener != nil {
                s.listener.Close()
        }
        
        // Close all client connections
        s.connectionsMu.Lock()
        for connID, conn := range s.connections {
                conn.Close()
                delete(s.connections, connID)
        }
        s.connectionsMu.Unlock()
        
        s.logger.Info("TCP server stopped")
}