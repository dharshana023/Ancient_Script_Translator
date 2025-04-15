package models

import (
        "time"
)

// Manuscript represents an ancient manuscript
type Manuscript struct {
        ID           string    `json:"id"`
        OriginalPath string    `json:"originalPath"`
        ScriptType   string    `json:"scriptType"`
        UploadedAt   time.Time `json:"uploadedAt"`
        ProcessedAt  time.Time `json:"processedAt,omitempty"`
}

// TranslationResult represents the result of a translation
type TranslationResult struct {
        ManuscriptID   string    `json:"manuscriptId"`
        OriginalScript string    `json:"originalScript"`
        TranslatedText string    `json:"translatedText"`
        Summary        string    `json:"summary"`
        TranslatedAt   time.Time `json:"translatedAt"`
}

// TranslationResponse represents the API response for a translation request
type TranslationResponse struct {
        OriginalScript string `json:"originalScript"`
        TranslatedText string `json:"translatedText"`
        Summary        string `json:"summary"`
        ProcessedAt    string `json:"processedAt"`
}

// SummarizeRequest represents a request to summarize text
type SummarizeRequest struct {
        Text      string `json:"text"`
        Algorithm string `json:"algorithm,omitempty"` // Optional algorithm specification (extractive, abstractive, hybrid)
}

// SummarizeResponse represents the API response for a summarization request
type SummarizeResponse struct {
        Summary     string `json:"summary"`
        TextLength  int    `json:"textLength"`
        ProcessedAt string `json:"processedAt"`
}
