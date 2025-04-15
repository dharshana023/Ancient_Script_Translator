package services

import (
	"fmt"
	"strings"
)

// TranslationConfig contains the configuration for translation
type TranslationConfig struct {
	DefaultTargetLanguage string   `yaml:"defaultTargetLanguage"`
	SupportedScripts      []string `yaml:"supportedScripts"`
	UseExternalAPI        bool     `yaml:"useExternalAPI"`
	APIEndpoint           string   `yaml:"apiEndpoint"`
	APIKey                string   `yaml:"apiKey"`
}

// Translator handles the translation of ancient scripts
type Translator struct {
	config TranslationConfig
}

// NewTranslator creates a new translator
func NewTranslator(config TranslationConfig) *Translator {
	return &Translator{
		config: config,
	}
}

// TranslateText translates the extracted text to the target language
func (t *Translator) TranslateText(text, scriptType string) (string, error) {
	// Validate script type
	if scriptType != "auto" && !t.isScriptSupported(scriptType) {
		return "", fmt.Errorf("unsupported script type: %s", scriptType)
	}

	// If script type is auto, attempt to detect it
	if scriptType == "auto" {
		var err error
		scriptType, err = t.detectScriptType(text)
		if err != nil {
			return "", fmt.Errorf("failed to detect script type: %v", err)
		}
	}

	// If external API is enabled, use it for translation
	if t.config.UseExternalAPI {
		return t.translateWithExternalAPI(text, scriptType)
	}

	// Otherwise use internal translation logic
	return t.translateWithInternalLogic(text, scriptType)
}

// isScriptSupported checks if the script type is supported
func (t *Translator) isScriptSupported(scriptType string) bool {
	for _, supported := range t.config.SupportedScripts {
		if strings.EqualFold(scriptType, supported) {
			return true
		}
	}
	return false
}

// detectScriptType attempts to automatically detect the script type
func (t *Translator) detectScriptType(text string) (string, error) {
	// In a real implementation, this would analyze the text and patterns
	// to determine the most likely script type
	// For this simulation, we'll return a default script type
	return t.config.SupportedScripts[0], nil
}

// translateWithExternalAPI translates the text using an external API
func (t *Translator) translateWithExternalAPI(text, scriptType string) (string, error) {
	// In a real implementation, this would make an API call to an external
	// translation service with the text and script type
	// For this simulation, we'll return a placeholder translated text
	return fmt.Sprintf("Translated from %s script: Lorem ipsum dolor sit amet", scriptType), nil
}

// translateWithInternalLogic translates the text using internal logic
func (t *Translator) translateWithInternalLogic(text, scriptType string) (string, error) {
	// In a real implementation, this would apply specific translation rules
	// based on the script type and language patterns
	// For this simulation, we'll return a placeholder translated text
	return fmt.Sprintf("Translated from %s script: Lorem ipsum dolor sit amet", scriptType), nil
}
