package services

import (
        "bytes"
        "encoding/base64"
        "fmt"
        "image"
        "image/jpeg"
        "image/png"
        _ "image/jpeg"
        _ "image/png"
        "sync"
        
        "ancient-script-decoder/utils"
)

// ImageProcessingConfig contains the configuration for image processing
type ImageProcessingConfig struct {
        EnhancementEnabled   bool    `yaml:"enhancementEnabled"`
        ContrastFactor       float64 `yaml:"contrastFactor"`
        BrightnessAdjust     float64 `yaml:"brightnessAdjust"`
        DenoiseLevel         int     `yaml:"denoiseLevel"`
        GaussianBlurSigma    float64 `yaml:"gaussianBlurSigma"`
        GaussianBlurSize     int     `yaml:"gaussianBlurSize"`
        BoxBlurSize          int     `yaml:"boxBlurSize"`
        SobelThreshold       uint8   `yaml:"sobelThreshold"`
        RotationAngle        float64 `yaml:"rotationAngle"`
        ConcurrencyLevel     int     `yaml:"concurrencyLevel"`
        UseParallelProcessing bool    `yaml:"useParallelProcessing"`
}

// ImageProcessor handles the processing of manuscript images
type ImageProcessor struct {
        config    ImageProcessingConfig
        cacheLock sync.RWMutex
        // Cache for processed images (map of operation name -> image hash -> result)
        cache     map[string]map[string]image.Image
}

// NewImageProcessor creates a new image processor
func NewImageProcessor(config ImageProcessingConfig) *ImageProcessor {
        // Set default concurrency level if not specified
        if config.ConcurrencyLevel <= 0 {
                config.ConcurrencyLevel = 4
        }
        
        return &ImageProcessor{
                config: config,
                cache:  make(map[string]map[string]image.Image),
        }
}

// ProcessImage processes the manuscript image to prepare it for OCR
// Demonstrates error handling, conditional logic, and method chaining
func (p *ImageProcessor) ProcessImage(imageData []byte) ([]byte, error) {
        // Decode the image
        img, format, err := image.Decode(bytes.NewReader(imageData))
        if err != nil {
                return nil, fmt.Errorf("failed to decode image: %v", err)
        }

        // Create a slice of image processing algorithms to apply
        // Demonstrates slices, interfaces, and polymorphism
        var algorithms []utils.ImageProcessingAlgorithm
        
        // Apply grayscale conversion (always applied for OCR preprocessing)
        algorithms = append(algorithms, utils.NewGrayscaleProcessor(
                p.config.ConcurrencyLevel, 
                p.config.UseParallelProcessing,
        ))
        
        // Apply edge detection if enhancement is enabled
        if p.config.EnhancementEnabled {
                algorithms = append(algorithms, utils.NewSobelEdgeDetector(
                        p.config.SobelThreshold,
                        p.config.ConcurrencyLevel,
                        p.config.UseParallelProcessing,
                ))
        }
        
        // Apply denoising if level is set
        if p.config.DenoiseLevel > 0 {
                // Choose between box blur and Gaussian blur based on denoise level
                if p.config.DenoiseLevel <= 3 {
                        algorithms = append(algorithms, utils.NewBoxBlurProcessor(
                                p.config.BoxBlurSize,
                                p.config.ConcurrencyLevel,
                                p.config.UseParallelProcessing,
                        ))
                } else {
                        algorithms = append(algorithms, utils.NewGaussianBlurProcessor(
                                p.config.GaussianBlurSigma,
                                p.config.GaussianBlurSize,
                                p.config.ConcurrencyLevel,
                                p.config.UseParallelProcessing,
                        ))
                }
        }
        
        // Process the image through the pipeline
        processedImg := utils.ProcessImagePipeline(img, algorithms)
        
        // Encode the processed image based on the original format
        var buf bytes.Buffer
        switch format {
        case "jpeg":
                if err := jpeg.Encode(&buf, processedImg, nil); err != nil {
                        return nil, fmt.Errorf("failed to encode JPEG: %v", err)
                }
        case "png":
                if err := png.Encode(&buf, processedImg); err != nil {
                        return nil, fmt.Errorf("failed to encode PNG: %v", err)
                }
        default:
                return nil, fmt.Errorf("unsupported image format: %s", format)
        }
        
        return buf.Bytes(), nil
}

// ApplyImageTransformations applies various geometric transformations to an image
// Demonstrates method chaining, error handling, and interface usage
func (p *ImageProcessor) ApplyImageTransformations(imageData []byte, transformations map[string]interface{}) ([]byte, error) {
        // Decode the image
        img, format, err := image.Decode(bytes.NewReader(imageData))
        if err != nil {
                return nil, fmt.Errorf("failed to decode image: %v", err)
        }
        
        // Create a slice of image processing algorithms to apply
        var algorithms []utils.ImageProcessingAlgorithm
        
        // Apply transformations based on the provided options
        if val, ok := transformations["upsideDown"]; ok && val.(bool) {
                algorithms = append(algorithms, utils.NewUpsideDownProcessor(
                        p.config.ConcurrencyLevel,
                        p.config.UseParallelProcessing,
                ))
        }
        
        if val, ok := transformations["rotationAngle"]; ok {
                angle := val.(float64)
                // Use shear rotation for specific angles for better performance
                if angle == 90 || angle == 180 || angle == 270 {
                        algorithms = append(algorithms, utils.NewShearRotateProcessor(
                                angle,
                                p.config.ConcurrencyLevel,
                                p.config.UseParallelProcessing,
                        ))
                } else {
                        algorithms = append(algorithms, utils.NewRotateProcessor(
                                angle,
                                p.config.ConcurrencyLevel,
                                p.config.UseParallelProcessing,
                        ))
                }
        }
        
        if val, ok := transformations["grayscale"]; ok && val.(bool) {
                algorithms = append(algorithms, utils.NewGrayscaleProcessor(
                        p.config.ConcurrencyLevel,
                        p.config.UseParallelProcessing,
                ))
        }
        
        if val, ok := transformations["boxBlur"]; ok && val.(bool) {
                size := p.config.BoxBlurSize
                if sizeVal, sizeOk := transformations["boxBlurSize"]; sizeOk {
                        size = sizeVal.(int)
                }
                algorithms = append(algorithms, utils.NewBoxBlurProcessor(
                        size,
                        p.config.ConcurrencyLevel,
                        p.config.UseParallelProcessing,
                ))
        }
        
        if val, ok := transformations["gaussianBlur"]; ok && val.(bool) {
                sigma := p.config.GaussianBlurSigma
                size := p.config.GaussianBlurSize
                if sigmaVal, sigmaOk := transformations["gaussianBlurSigma"]; sigmaOk {
                        sigma = sigmaVal.(float64)
                }
                if sizeVal, sizeOk := transformations["gaussianBlurSize"]; sizeOk {
                        size = sizeVal.(int)
                }
                algorithms = append(algorithms, utils.NewGaussianBlurProcessor(
                        sigma,
                        size,
                        p.config.ConcurrencyLevel,
                        p.config.UseParallelProcessing,
                ))
        }
        
        if val, ok := transformations["edgeDetection"]; ok && val.(bool) {
                threshold := p.config.SobelThreshold
                if threshVal, threshOk := transformations["edgeThreshold"]; threshOk {
                        threshold = uint8(threshVal.(int))
                }
                algorithms = append(algorithms, utils.NewSobelEdgeDetector(
                        threshold,
                        p.config.ConcurrencyLevel,
                        p.config.UseParallelProcessing,
                ))
        }
        
        // Process the image through the pipeline
        processedImg := utils.ProcessImagePipeline(img, algorithms)
        
        // Encode the processed image based on the original format
        var buf bytes.Buffer
        switch format {
        case "jpeg":
                if err := jpeg.Encode(&buf, processedImg, nil); err != nil {
                        return nil, fmt.Errorf("failed to encode JPEG: %v", err)
                }
        case "png":
                if err := png.Encode(&buf, processedImg); err != nil {
                        return nil, fmt.Errorf("failed to encode PNG: %v", err)
                }
        default:
                return nil, fmt.Errorf("unsupported image format: %s", format)
        }
        
        return buf.Bytes(), nil
}

// GetImageBase64 converts an image to base64 for web display
func (p *ImageProcessor) GetImageBase64(imageData []byte) (string, error) {
        return base64.StdEncoding.EncodeToString(imageData), nil
}

// ExtractTextFromImage extracts text from the processed image
// In a real implementation, this would use OCR specific to ancient scripts
func (p *ImageProcessor) ExtractTextFromImage(processedImageData []byte, scriptType string) (string, error) {
        // Different script types would use different OCR algorithms in a real implementation
        scriptHandlers := map[string]func([]byte) (string, error){
                "latin": p.processLatinScript,
                "greek": p.processGreekScript,
                "cuneiform": p.processCuneiformScript,
                "hieroglyphic": p.processHieroglyphicScript,
                "runic": p.processRunicScript,
                "auto": p.autoDetectScript,
        }
        
        // Get the appropriate handler function for the script type
        handler, ok := scriptHandlers[scriptType]
        if !ok {
                // Default to auto-detection if script type is not recognized
                handler = p.autoDetectScript
        }
        
        // Process the image with the selected handler
        return handler(processedImageData)
}

// Script processing functions - in a real implementation, these would use
// specialized OCR algorithms for different ancient scripts

func (p *ImageProcessor) processLatinScript(imageData []byte) (string, error) {
        return "Example Latin text extracted from manuscript: Senatus Populusque Romanus", nil
}

func (p *ImageProcessor) processGreekScript(imageData []byte) (string, error) {
        return "Example Greek text extracted from manuscript: Ἐν ἀρχῇ ἦν ὁ λόγος", nil
}

func (p *ImageProcessor) processCuneiformScript(imageData []byte) (string, error) {
        return "Example Cuneiform text extracted from manuscript: Laws of Hammurabi, first section", nil
}

func (p *ImageProcessor) processHieroglyphicScript(imageData []byte) (string, error) {
        return "Example Hieroglyphic text extracted from manuscript: Excerpt from Book of the Dead", nil
}

func (p *ImageProcessor) processRunicScript(imageData []byte) (string, error) {
        return "Example Runic text extracted from manuscript: Norse inscription commemorating victory", nil
}

func (p *ImageProcessor) autoDetectScript(imageData []byte) (string, error) {
        // In a real implementation, this would analyze the image to determine the script type
        // Then call the appropriate processing function
        return "Auto-detected ancient text from manuscript. Script appears to be a form of early Mediterranean writing.", nil
}
