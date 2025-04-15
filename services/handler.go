package services

import (
        "ancient-script-decoder/models"
        "ancient-script-decoder/utils"
)

// ServiceHandler coordinates the different services
type ServiceHandler struct {
        imageProcessor   *ImageProcessor
        translator       *Translator
        summarizer       *Summarizer
        metadataExtractor *MetadataExtractor
        logger           *utils.Logger
}

// NewServiceHandler creates a new service handler
func NewServiceHandler(imageProcessor *ImageProcessor, translator *Translator, summarizer *Summarizer, metadataExtractor *MetadataExtractor, logger *utils.Logger) *ServiceHandler {
        return &ServiceHandler{
                imageProcessor:   imageProcessor,
                translator:       translator,
                summarizer:       summarizer,
                metadataExtractor: metadataExtractor,
                logger:           logger,
        }
}

// ProcessAndTranslate processes an image and translates the extracted text
func (h *ServiceHandler) ProcessAndTranslate(imageData []byte, scriptType string) (string, error) {
        // Process the image
        h.logger.Info("Processing manuscript image")
        processedImage, err := h.imageProcessor.ProcessImage(imageData)
        if err != nil {
                h.logger.Error("Failed to process image", "error", err)
                return "", err
        }

        // Extract text from the processed image
        h.logger.Info("Extracting text from processed image", "scriptType", scriptType)
        extractedText, err := h.imageProcessor.ExtractTextFromImage(processedImage, scriptType)
        if err != nil {
                h.logger.Error("Failed to extract text", "error", err)
                return "", err
        }

        // Translate the extracted text
        h.logger.Info("Translating extracted text", "scriptType", scriptType)
        translatedText, err := h.translator.TranslateText(extractedText, scriptType)
        if err != nil {
                h.logger.Error("Failed to translate text", "error", err)
                return "", err
        }

        return translatedText, nil
}

// SummarizeText summarizes the translated text
func (h *ServiceHandler) SummarizeText(text string) (string, error) {
        return h.SummarizeTextWithAlgorithm(text, "")
}

// SummarizeTextWithAlgorithm summarizes text using the specified algorithm
func (h *ServiceHandler) SummarizeTextWithAlgorithm(text string, algorithm string) (string, error) {
        h.logger.Info("Generating summary", "textLength", len(text), "algorithm", algorithm)
        
        // Store current algorithm
        var originalAlgorithm string
        if algorithm != "" {
                // Temporarily change the algorithm if specified
                originalAlgorithm = h.summarizer.config.Algorithm
                h.summarizer.config.Algorithm = algorithm
                defer func() {
                        // Restore original algorithm after summarization
                        h.summarizer.config.Algorithm = originalAlgorithm
                }()
        }
        
        // Generate summary
        summary, err := h.summarizer.SummarizeText(text)
        if err != nil {
                h.logger.Error("Failed to generate summary", "error", err)
                return "", err
        }
        
        h.logger.Info("Summary generated successfully", "summaryLength", len(summary), "algorithm", algorithm)
        return summary, nil
}

// ExtractMetadata extracts historical context metadata from translated text and original manuscript
// If imageData is nil, metadata will be extracted from text only (for direct text input)
func (h *ServiceHandler) ExtractMetadata(translatedText string, scriptType string, imageData ...[]byte) (models.Metadata, error) {
        h.logger.Info("Extracting historical metadata", "scriptType", scriptType, "textLength", len(translatedText))
        
        // For direct text input without an image
        var imgData []byte
        if len(imageData) > 0 {
                imgData = imageData[0]
        }
        
        metadata, err := h.metadataExtractor.ExtractMetadata(translatedText, scriptType, imgData)
        if err != nil {
                h.logger.Error("Failed to extract metadata", "error", err)
                return models.Metadata{}, err
        }
        
        h.logger.Info("Metadata extraction successful", 
                "timePeriods", len(metadata.TimePeriods), 
                "regions", len(metadata.Regions),
                "cultures", len(metadata.CulturalContext),
                "confidence", metadata.ConfidenceScore)
        
        return metadata, nil
}

// ProcessTranslateWithMetadata processes, translates, and extracts metadata in one operation
func (h *ServiceHandler) ProcessTranslateWithMetadata(imageData []byte, scriptType string) (string, models.Metadata, error) {
        // First translate the text
        translatedText, err := h.ProcessAndTranslate(imageData, scriptType)
        if err != nil {
                return "", models.Metadata{}, err
        }
        
        // Extract metadata from translated text and original image
        metadata, err := h.ExtractMetadata(translatedText, scriptType, imageData)
        if err != nil {
                // Don't fail the whole operation if metadata extraction fails
                h.logger.Error("Metadata extraction failed, continuing with empty metadata", "error", err)
                return translatedText, models.Metadata{}, nil
        }
        
        return translatedText, metadata, nil
}
