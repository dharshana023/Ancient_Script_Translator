package api

import (
        "context"
        "encoding/json"
        "fmt"
        "io"
        "net/http"
        "time"

        "ancient-script-decoder/models"
        "ancient-script-decoder/services"
        "ancient-script-decoder/utils"
)

// RESTServer represents the REST API server
type RESTServer struct {
        port          int
        serviceHandler *services.ServiceHandler
        server        *http.Server
        logger        *utils.Logger
}

// NewRESTServer creates a new REST API server
func NewRESTServer(port int, serviceHandler *services.ServiceHandler, logger *utils.Logger) *RESTServer {
        return &RESTServer{
                port:          port,
                serviceHandler: serviceHandler,
                logger:        logger,
        }
}

// Start starts the REST API server
func (s *RESTServer) Start() error {
        mux := http.NewServeMux()

        // Register API endpoints
        mux.HandleFunc("/api/translate", s.handleTranslate)
        mux.HandleFunc("/api/translate/text", s.handleTranslateText)
        mux.HandleFunc("/api/summarize", s.handleSummarize)
        mux.HandleFunc("/api/health", s.handleHealth)
        
        // Serve static files
        fs := http.FileServer(http.Dir("static"))
        mux.Handle("/", fs)

        // Create HTTP server
        s.server = &http.Server{
                Addr:         fmt.Sprintf("0.0.0.0:%d", s.port),
                Handler:      mux,
                ReadTimeout:  60 * time.Second,
                WriteTimeout: 60 * time.Second,
                IdleTimeout:  120 * time.Second,
        }

        // Start the server
        return s.server.ListenAndServe()
}

// Stop stops the REST API server
func (s *RESTServer) Stop(ctx context.Context) error {
        return s.server.Shutdown(ctx)
}

// handleTranslate handles the translation request
func (s *RESTServer) handleTranslate(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
                http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                return
        }

        // Parse multipart form data with 10MB limit
        err := r.ParseMultipartForm(10 << 20)
        if err != nil {
                s.logger.Error("Failed to parse form", "error", err)
                http.Error(w, "Failed to parse form", http.StatusBadRequest)
                return
        }

        // Get the uploaded file
        file, handler, err := r.FormFile("manuscript")
        if err != nil {
                s.logger.Error("Failed to get file from form", "error", err)
                http.Error(w, "Failed to get file from form", http.StatusBadRequest)
                return
        }
        defer file.Close()

        s.logger.Info("Received manuscript", "filename", handler.Filename, "size", handler.Size)

        // Read the file content
        fileBytes, err := io.ReadAll(file)
        if err != nil {
                s.logger.Error("Failed to read file", "error", err)
                http.Error(w, "Failed to read file", http.StatusInternalServerError)
                return
        }

        // Get script type from form
        scriptType := r.FormValue("scriptType")
        if scriptType == "" {
                scriptType = "auto" // Default to auto-detection
        }

        // Process, translate the manuscript, and extract metadata
        processedText, metadata, err := s.serviceHandler.ProcessTranslateWithMetadata(fileBytes, scriptType)
        if err != nil {
                s.logger.Error("Failed to process and translate manuscript", "error", err)
                http.Error(w, fmt.Sprintf("Failed to process and translate manuscript: %v", err), http.StatusInternalServerError)
                return
        }

        // Generate summary for the translated text
        summary, err := s.serviceHandler.SummarizeText(processedText)
        if err != nil {
                s.logger.Error("Failed to generate summary", "error", err)
                http.Error(w, fmt.Sprintf("Failed to generate summary: %v", err), http.StatusInternalServerError)
                return
        }

        // Create response
        response := models.TranslationResponse{
                OriginalScript: scriptType,
                TranslatedText: processedText,
                Summary:        summary,
                Metadata:       metadata,
                ProcessedAt:    time.Now().Format(time.RFC3339),
        }

        // Send JSON response
        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(response); err != nil {
                s.logger.Error("Failed to encode response", "error", err)
                http.Error(w, "Failed to encode response", http.StatusInternalServerError)
                return
        }
}

// handleSummarize handles the summarization request
func (s *RESTServer) handleSummarize(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
                http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                return
        }

        // Parse JSON request
        var request models.SummarizeRequest
        if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
                s.logger.Error("Failed to parse request", "error", err)
                http.Error(w, "Failed to parse request", http.StatusBadRequest)
                return
        }

        // Validate request
        if request.Text == "" {
                http.Error(w, "Text cannot be empty", http.StatusBadRequest)
                return
        }

        // Generate summary with specified algorithm if provided
        summary, err := s.serviceHandler.SummarizeTextWithAlgorithm(request.Text, request.Algorithm)
        if err != nil {
                s.logger.Error("Failed to generate summary", "error", err, "algorithm", request.Algorithm)
                http.Error(w, fmt.Sprintf("Failed to generate summary: %v", err), http.StatusInternalServerError)
                return
        }

        // Create response
        response := models.SummarizeResponse{
                Summary:     summary,
                TextLength:  len(request.Text),
                ProcessedAt: time.Now().Format(time.RFC3339),
        }

        // Send JSON response
        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(response); err != nil {
                s.logger.Error("Failed to encode response", "error", err)
                http.Error(w, "Failed to encode response", http.StatusInternalServerError)
                return
        }
}

// handleHealth handles the health check request
func (s *RESTServer) handleHealth(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
                http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                return
        }

        // Create response
        response := map[string]string{
                "status": "healthy",
                "time":   time.Now().Format(time.RFC3339),
        }

        // Send JSON response
        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(response); err != nil {
                s.logger.Error("Failed to encode response", "error", err)
                http.Error(w, "Failed to encode response", http.StatusInternalServerError)
                return
        }
}

// handleTranslateText handles the direct text translation and metadata extraction request
func (s *RESTServer) handleTranslateText(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
                http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                return
        }

        // Parse JSON request
        var request struct {
                OriginalText string `json:"originalText"`
                ScriptType   string `json:"scriptType"`
        }
        
        if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
                s.logger.Error("Failed to parse request", "error", err)
                http.Error(w, "Failed to parse request", http.StatusBadRequest)
                return
        }

        // Validate request
        if request.OriginalText == "" {
                http.Error(w, "Text cannot be empty", http.StatusBadRequest)
                return
        }

        if request.ScriptType == "" {
                request.ScriptType = "auto" // Default to auto-detection
        }

        // For direct text input, we skip the image processing step
        // and extract metadata directly from the text
        metadata, err := s.serviceHandler.ExtractMetadata(request.OriginalText, request.ScriptType)
        if err != nil {
                s.logger.Error("Failed to extract metadata", "error", err)
                http.Error(w, fmt.Sprintf("Failed to extract metadata: %v", err), http.StatusInternalServerError)
                return
        }

        // Generate summary for the text
        summary, err := s.serviceHandler.SummarizeText(request.OriginalText)
        if err != nil {
                s.logger.Error("Failed to generate summary", "error", err)
                http.Error(w, fmt.Sprintf("Failed to generate summary: %v", err), http.StatusInternalServerError)
                return
        }

        // Create response
        response := models.TranslationResponse{
                OriginalScript: request.ScriptType,
                TranslatedText: request.OriginalText, // For text input, we don't translate
                Summary:        summary,
                Metadata:       metadata,
                ProcessedAt:    time.Now().Format(time.RFC3339),
        }

        // Send JSON response
        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(response); err != nil {
                s.logger.Error("Failed to encode response", "error", err)
                http.Error(w, "Failed to encode response", http.StatusInternalServerError)
                return
        }
}
