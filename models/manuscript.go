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
        Metadata       Metadata  `json:"metadata,omitempty"`
        TranslatedAt   time.Time `json:"translatedAt"`
}

// TranslationResponse represents the API response for a translation request
type TranslationResponse struct {
        OriginalScript string   `json:"originalScript"`
        TranslatedText string   `json:"translatedText"`
        Summary        string   `json:"summary"`
        Metadata       Metadata `json:"metadata,omitempty"`
        ProcessedAt    string   `json:"processedAt"`
}

// TimePeriod represents a historical time period
type TimePeriod struct {
        Name        string `json:"name"`
        StartYear   int    `json:"startYear"`
        EndYear     int    `json:"endYear"`
        Description string `json:"description"`
}

// Region represents a geographical region
type Region struct {
        Name        string   `json:"name"`
        ModernAreas []string `json:"modernAreas"`
        Description string   `json:"description"`
}

// HistoricalEvent represents an event referenced in the manuscript
type HistoricalEvent struct {
        Name        string `json:"name"`
        EventType   string `json:"eventType"`
        Year        int    `json:"year,omitempty"`
        Description string `json:"description"`
}

// Metadata represents historical context for a manuscript
type Metadata struct {
        ScriptType       string           `json:"scriptType"`
        TimePeriods      []TimePeriod     `json:"timePeriods,omitempty"`
        Regions          []Region         `json:"regions,omitempty"`
        CulturalContext  []string         `json:"culturalContext,omitempty"`
        MaterialContext  []string         `json:"materialContext,omitempty"`
        HistoricalEvents []HistoricalEvent `json:"historicalEvents,omitempty"`
        ConfidenceScore  float64          `json:"confidenceScore"`
        DetectedDate     string           `json:"detectedDate"`
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
