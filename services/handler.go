package services

import (
	"ancient-script-decoder/utils"
)

// ServiceHandler coordinates the different services
type ServiceHandler struct {
	imageProcessor *ImageProcessor
	translator     *Translator
	summarizer     *Summarizer
	logger         *utils.Logger
}

// NewServiceHandler creates a new service handler
func NewServiceHandler(imageProcessor *ImageProcessor, translator *Translator, summarizer *Summarizer, logger *utils.Logger) *ServiceHandler {
	return &ServiceHandler{
		imageProcessor: imageProcessor,
		translator:     translator,
		summarizer:     summarizer,
		logger:         logger,
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
	h.logger.Info("Generating summary", "textLength", len(text))
	
	// Generate summary
	summary, err := h.summarizer.SummarizeText(text)
	if err != nil {
		h.logger.Error("Failed to generate summary", "error", err)
		return "", err
	}
	
	h.logger.Info("Summary generated successfully", "summaryLength", len(summary))
	return summary, nil
}
