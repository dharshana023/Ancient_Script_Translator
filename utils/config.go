package utils

import (
        "fmt"
        "os"
        "path/filepath"

        "gopkg.in/yaml.v2"
)

// Config represents the application configuration
type Config struct {
        REST struct {
                Port int `yaml:"port"`
        } `yaml:"rest"`
        GRPC struct {
                Port int `yaml:"port"`
        } `yaml:"grpc"`
        ImageProcessing struct {
                EnhancementEnabled bool    `yaml:"enhancementEnabled"`
                ContrastFactor     float64 `yaml:"contrastFactor"`
                BrightnessAdjust   float64 `yaml:"brightnessAdjust"`
                DenoiseLevel       int     `yaml:"denoiseLevel"`
        } `yaml:"imageProcessing"`
        Translation struct {
                DefaultTargetLanguage string   `yaml:"defaultTargetLanguage"`
                SupportedScripts      []string `yaml:"supportedScripts"`
                UseExternalAPI        bool     `yaml:"useExternalAPI"`
                APIEndpoint           string   `yaml:"apiEndpoint"`
                APIKey                string   `yaml:"apiKey"`
        } `yaml:"translation"`
        Summarization struct {
                MaxSummaryLength    int     `yaml:"maxSummaryLength"`
                Algorithm           string  `yaml:"algorithm"`
                ModelPath           string  `yaml:"modelPath"`
                SentenceImportance  float64 `yaml:"sentenceImportance"`
                KeywordImportance   float64 `yaml:"keywordImportance"`
                ContextImportance   float64 `yaml:"contextImportance"`
                ContextWindowSize   int     `yaml:"contextWindowSize"`
                KeyphrasesPerDoc    int     `yaml:"keyphrasesPerDoc"`
                MinSentenceLength   int     `yaml:"minSentenceLength"`
                EnableEntityRecog   bool    `yaml:"enableEntityRecognition"`
        } `yaml:"summarization"`
        Metadata struct {
                EnableGeographicDetection bool    `yaml:"enableGeographicDetection"`
                EnablePeriodDetection     bool    `yaml:"enablePeriodDetection"`
                EnableCultureDetection    bool    `yaml:"enableCultureDetection"`
                ContextSensitivity        float64 `yaml:"contextSensitivity"`
                ReferenceDatabase         string  `yaml:"referenceDatabase"`
        } `yaml:"metadata"`
}

// LoadConfig loads the configuration from a file
func LoadConfig(path string) (*Config, error) {
        // Create default configuration
        config := &Config{}
        config.REST.Port = 5000
        config.GRPC.Port = 8000
        
        // Default image processing settings
        config.ImageProcessing.EnhancementEnabled = true
        config.ImageProcessing.ContrastFactor = 1.5
        config.ImageProcessing.BrightnessAdjust = 0.1
        config.ImageProcessing.DenoiseLevel = 2
        
        // Default translation settings
        config.Translation.DefaultTargetLanguage = "en"
        config.Translation.SupportedScripts = []string{"latin", "greek", "cuneiform", "hieroglyphic", "runic"}
        config.Translation.UseExternalAPI = false
        config.Translation.APIEndpoint = os.Getenv("TRANSLATION_API_ENDPOINT")
        config.Translation.APIKey = os.Getenv("TRANSLATION_API_KEY")
        
        // Default summarization settings
        config.Summarization.MaxSummaryLength = 500
        config.Summarization.Algorithm = "textrank"
        config.Summarization.ModelPath = "models/word2vec.bin"
        config.Summarization.SentenceImportance = 0.3
        config.Summarization.KeywordImportance = 0.4
        config.Summarization.ContextImportance = 0.3
        config.Summarization.ContextWindowSize = 2
        config.Summarization.KeyphrasesPerDoc = 5
        config.Summarization.MinSentenceLength = 3
        config.Summarization.EnableEntityRecog = true
        
        // Default metadata settings
        config.Metadata.EnableGeographicDetection = true
        config.Metadata.EnablePeriodDetection = true
        config.Metadata.EnableCultureDetection = true
        config.Metadata.ContextSensitivity = 0.7
        config.Metadata.ReferenceDatabase = "models/historical_reference.db"
        
        // If configuration file exists, load it
        if _, err := os.Stat(path); err == nil {
                // Read the YAML file
                data, err := os.ReadFile(filepath.Clean(path))
                if err != nil {
                        return nil, fmt.Errorf("failed to read configuration file: %v", err)
                }
                
                // Parse the YAML file
                if err := yaml.Unmarshal(data, config); err != nil {
                        return nil, fmt.Errorf("failed to parse configuration file: %v", err)
                }
        }
        
        return config, nil
}
