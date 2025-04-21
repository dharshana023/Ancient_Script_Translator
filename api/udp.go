package api

import (
        "context"
        "encoding/json"
        "fmt"
        "net"
        "sync"
        "time"

        "ancient-script-decoder/services"
        "ancient-script-decoder/utils"
)

// UDPConfig contains configuration for the UDP server
type UDPConfig struct {
        Enabled    bool   `yaml:"enabled"`
        Port       int    `yaml:"port"`
        Host       string `yaml:"host"`
        BufferSize int    `yaml:"bufferSize"`
}

// UDPServer represents a UDP server for the ancient script decoder
type UDPServer struct {
        config         UDPConfig
        serviceHandler *services.ServiceHandler
        logger         *utils.Logger
        conn           *net.UDPConn
        ctx            context.Context
        cancel         context.CancelFunc
        wg             sync.WaitGroup
}

// NewUDPServer creates a new UDP server
func NewUDPServer(config UDPConfig, serviceHandler *services.ServiceHandler, logger *utils.Logger) *UDPServer {
        ctx, cancel := context.WithCancel(context.Background())
        
        if config.Host == "" {
                config.Host = "0.0.0.0" // Default to all interfaces
        }
        
        if config.BufferSize <= 0 {
                config.BufferSize = 4096 // Default buffer size
        }
        
        return &UDPServer{
                config:         config,
                serviceHandler: serviceHandler,
                logger:         logger,
                ctx:            ctx,
                cancel:         cancel,
        }
}

// Start starts the UDP server
func (s *UDPServer) Start() error {
        if !s.config.Enabled {
                s.logger.Info("UDP server is disabled")
                return nil
        }
        
        // Resolve the UDP address
        addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", s.config.Host, s.config.Port))
        if err != nil {
                return fmt.Errorf("failed to resolve UDP address: %v", err)
        }
        
        // Create a UDP connection
        conn, err := net.ListenUDP("udp", addr)
        if err != nil {
                return fmt.Errorf("failed to start UDP server: %v", err)
        }
        
        s.conn = conn
        s.logger.Info("UDP server started", "address", addr.String())
        
        // Start receiving packets in a goroutine
        s.wg.Add(1)
        go s.receivePackets()
        
        return nil
}

// receivePackets receives UDP packets
func (s *UDPServer) receivePackets() {
        defer s.wg.Done()
        
        buffer := make([]byte, s.config.BufferSize)
        
        for {
                // Check if context is canceled
                select {
                case <-s.ctx.Done():
                        return
                default:
                        // Continue receiving packets
                }
                
                // Set a read deadline to avoid blocking forever
                s.conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
                
                // Read a packet
                n, addr, err := s.conn.ReadFromUDP(buffer)
                if err != nil {
                        if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
                                // Timeout, just continue
                                continue
                        }
                        
                        select {
                        case <-s.ctx.Done():
                                // Server is shutting down, ignore the error
                                return
                        default:
                                s.logger.Error("Failed to read UDP packet", "error", err)
                                continue
                        }
                }
                
                // Handle the packet in a goroutine to allow receiving more packets
                go s.handlePacket(buffer[:n], addr)
        }
}

// handlePacket handles a UDP packet
func (s *UDPServer) handlePacket(packet []byte, addr *net.UDPAddr) {
        s.logger.Info("Received UDP packet", "address", addr.String(), "size", len(packet))
        
        // Process the packet
        response, err := s.processPacket(packet)
        if err != nil {
                s.logger.Error("Failed to process UDP packet", "address", addr.String(), "error", err)
                // Send error response
                errorResponse := []byte(fmt.Sprintf("ERROR: %v", err))
                s.conn.WriteToUDP(errorResponse, addr)
                return
        }
        
        // Marshal the response
        responseData, err := json.Marshal(response)
        if err != nil {
                s.logger.Error("Failed to marshal UDP response", "address", addr.String(), "error", err)
                // Send error response
                errorResponse := []byte("ERROR: Failed to marshal response")
                s.conn.WriteToUDP(errorResponse, addr)
                return
        }
        
        // Send the response
        _, err = s.conn.WriteToUDP(responseData, addr)
        if err != nil {
                s.logger.Error("Failed to send UDP response", "address", addr.String(), "error", err)
                return
        }
        
        s.logger.Info("Sent UDP response", "address", addr.String(), "size", len(responseData))
}

// processPacket processes a UDP packet
func (s *UDPServer) processPacket(packet []byte) (interface{}, error) {
        // Parse the request
        var request map[string]interface{}
        if err := json.Unmarshal(packet, &request); err != nil {
                return nil, fmt.Errorf("failed to parse request: %v", err)
        }
        
        // Get the request type
        requestType, ok := request["type"].(string)
        if !ok {
                return nil, fmt.Errorf("missing or invalid request type")
        }
        
        // Process the request based on the type
        switch requestType {
        case "ping":
                return s.handlePingRequest(request)
        case "health":
                return s.handleHealthRequest()
        default:
                return nil, fmt.Errorf("unsupported request type: %s", requestType)
        }
}

// handlePingRequest handles a ping request
func (s *UDPServer) handlePingRequest(request map[string]interface{}) (map[string]interface{}, error) {
        // Get the timestamp from the request (optional)
        var timestamp string
        if ts, ok := request["timestamp"].(string); ok {
                timestamp = ts
        } else {
                timestamp = time.Now().Format(time.RFC3339)
        }
        
        // Return the response
        return map[string]interface{}{
                "type":      "pong",
                "timestamp": timestamp,
                "serverTime": time.Now().Format(time.RFC3339),
        }, nil
}

// handleHealthRequest handles a health check request
func (s *UDPServer) handleHealthRequest() (map[string]interface{}, error) {
        // Return the health status
        return map[string]interface{}{
                "type":   "health",
                "status": "healthy",
                "time":   time.Now().Format(time.RFC3339),
        }, nil
}

// Stop stops the UDP server
func (s *UDPServer) Stop() {
        s.cancel() // Cancel the context to signal shutdown
        
        // Close the connection if it exists
        if s.conn != nil {
                s.conn.Close()
        }
        
        // Wait for the receiving goroutine to finish
        s.wg.Wait()
        
        s.logger.Info("UDP server stopped")
}