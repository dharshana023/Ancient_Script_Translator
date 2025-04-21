package utils

import (
	"image"
	"image/color"
	"math"
	"sync"
)

// ImageProcessingAlgorithm is an interface that all image processing algorithms must implement
// Demonstrates interface and polymorphism
type ImageProcessingAlgorithm interface {
	Process(img image.Image) image.Image
	GetName() string
}

// Common struct that holds configuration for image processing algorithms
type ImageProcessor struct {
	name           string
	concurrency    int
	useParallel    bool
	processingLock sync.Mutex // Mutex for thread-safe operations
}

// Base constructor that all image processors will use
func NewImageProcessor(name string, concurrency int, useParallel bool) ImageProcessor {
	return ImageProcessor{
		name:        name,
		concurrency: concurrency,
		useParallel: useParallel,
	}
}

// GetName implements the ImageProcessingAlgorithm interface
func (ip *ImageProcessor) GetName() string {
	return ip.name
}

// UpsideDownProcessor flips an image upside down
// Demonstrates loops, arrays and slice operations
type UpsideDownProcessor struct {
	ImageProcessor
}

func NewUpsideDownProcessor(concurrency int, useParallel bool) *UpsideDownProcessor {
	return &UpsideDownProcessor{
		ImageProcessor: NewImageProcessor("Upside Down", concurrency, useParallel),
	}
}

func (p *UpsideDownProcessor) Process(img image.Image) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	result := image.NewRGBA(bounds)

	// Using a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	if p.useParallel {
		// Process the image in parallel using goroutines (concurrency)
		wg.Add(height)
		for y := 0; y < height; y++ {
			go func(y int) {
				defer wg.Done()
				for x := 0; x < width; x++ {
					// Calculate the position in the flipped image
					result.Set(x, height-1-y, img.At(x, y))
				}
			}(y)
		}
		wg.Wait()
	} else {
		// Single-threaded processing
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				result.Set(x, height-1-y, img.At(x, y))
			}
		}
	}

	return result
}

// RotateProcessor rotates an image by a given angle in degrees
// Demonstrates arithmetic operations, error handling
type RotateProcessor struct {
	ImageProcessor
	angle float64 // in degrees
}

func NewRotateProcessor(angle float64, concurrency int, useParallel bool) *RotateProcessor {
	return &RotateProcessor{
		ImageProcessor: NewImageProcessor("Rotate", concurrency, useParallel),
		angle:          angle,
	}
}

func (p *RotateProcessor) Process(img image.Image) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	
	// Convert angle to radians
	angleRad := p.angle * math.Pi / 180.0
	
	// Calculate the dimensions of the rotated image
	// Using absolute values to handle negative angles
	sinA, cosA := math.Abs(math.Sin(angleRad)), math.Abs(math.Cos(angleRad))
	newWidth := int(float64(width)*cosA + float64(height)*sinA)
	newHeight := int(float64(width)*sinA + float64(height)*cosA)
	
	// Create a new image with the calculated dimensions
	result := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	
	// Calculate the center of the original and new images
	origCenterX, origCenterY := float64(width)/2, float64(height)/2
	newCenterX, newCenterY := float64(newWidth)/2, float64(newHeight)/2
	
	// Rotation using a channel for work distribution
	var wg sync.WaitGroup
	pixelChan := make(chan [2]int, p.concurrency*10) // Buffered channel for work distribution
	
	// Start worker goroutines
	if p.useParallel {
		// Launch worker goroutines
		for i := 0; i < p.concurrency; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for pixel := range pixelChan {
					newX, newY := pixel[0], pixel[1]
					
					// Calculate the corresponding position in the original image
					// Translate to origin, rotate, then translate back
					origX := cosA*(float64(newX)-newCenterX) + sinA*(float64(newY)-newCenterY) + origCenterX
					origY := -sinA*(float64(newX)-newCenterX) + cosA*(float64(newY)-newCenterY) + origCenterY
					
					// Check if the pixel is within the bounds of the original image
					if origX >= 0 && origX < float64(width) && origY >= 0 && origY < float64(height) {
						// Get the color at the original position
						result.Set(newX, newY, img.At(int(origX), int(origY)))
					}
				}
			}()
		}
		
		// Send work to the channel
		for newY := 0; newY < newHeight; newY++ {
			for newX := 0; newX < newWidth; newX++ {
				pixelChan <- [2]int{newX, newY}
			}
		}
		close(pixelChan)
		wg.Wait()
	} else {
		// Single-threaded processing
		for newY := 0; newY < newHeight; newY++ {
			for newX := 0; newX < newWidth; newX++ {
				// Calculate the corresponding position in the original image
				origX := cosA*(float64(newX)-newCenterX) + sinA*(float64(newY)-newCenterY) + origCenterX
				origY := -sinA*(float64(newX)-newCenterX) + cosA*(float64(newY)-newCenterY) + origCenterY
				
				// Check if the pixel is within the bounds of the original image
				if origX >= 0 && origX < float64(width) && origY >= 0 && origY < float64(height) {
					// Get the color at the original position
					result.Set(newX, newY, img.At(int(origX), int(origY)))
				}
			}
		}
	}
	
	return result
}

// ShearRotateProcessor rotates an image using three shear matrices
// Demonstrates structure, arithmetic operations, and error handling
type ShearRotateProcessor struct {
	ImageProcessor
	angle float64 // in degrees
}

func NewShearRotateProcessor(angle float64, concurrency int, useParallel bool) *ShearRotateProcessor {
	return &ShearRotateProcessor{
		ImageProcessor: NewImageProcessor("Shear Rotate", concurrency, useParallel),
		angle:          angle,
	}
}

func (p *ShearRotateProcessor) Process(img image.Image) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	
	// Convert angle to radians
	angleRad := p.angle * math.Pi / 180.0
	
	// Calculate shear factors
	tanHalfAngle := math.Tan(angleRad / 2)
	sinAngle := math.Sin(angleRad)
	
	// Create a map to store the intermediate results
	// Demonstrates use of maps
	intermediateResults := make(map[string]image.Image)
	
	// First shear (horizontal)
	intermediate1 := image.NewRGBA(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			newX := int(float64(x) - float64(y)*tanHalfAngle)
			if newX >= 0 && newX < width {
				intermediate1.Set(newX, y, img.At(x, y))
			}
		}
	}
	intermediateResults["shear1"] = intermediate1
	
	// Second shear (vertical)
	intermediate2 := image.NewRGBA(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			newY := int(float64(y) + float64(x)*sinAngle)
			if newY >= 0 && newY < height {
				intermediate2.Set(x, newY, intermediate1.At(x, y))
			}
		}
	}
	intermediateResults["shear2"] = intermediate2
	
	// Third shear (horizontal again)
	result := image.NewRGBA(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			newX := int(float64(x) - float64(y)*tanHalfAngle)
			if newX >= 0 && newX < width {
				result.Set(newX, y, intermediate2.At(x, y))
			}
		}
	}
	intermediateResults["shear3"] = result
	
	return result
}

// GrayscaleProcessor converts an image to grayscale
// Demonstrates arithmetic operations, structure, slice
type GrayscaleProcessor struct {
	ImageProcessor
	weights [3]float64 // RGB weights for grayscale conversion
}

func NewGrayscaleProcessor(concurrency int, useParallel bool) *GrayscaleProcessor {
	// Standard weights for RGB to grayscale conversion
	return &GrayscaleProcessor{
		ImageProcessor: NewImageProcessor("Grayscale", concurrency, useParallel),
		weights:        [3]float64{0.299, 0.587, 0.114}, // Standard RGB to grayscale weights
	}
}

func (p *GrayscaleProcessor) Process(img image.Image) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	result := image.NewGray(bounds)
	
	// Process image with goroutines for parallel execution
	var wg sync.WaitGroup
	
	if p.useParallel {
		// Calculate the size of each chunk to process
		chunkSize := height / p.concurrency
		if chunkSize < 1 {
			chunkSize = 1
		}
		
		// Process the image in chunks
		for i := 0; i < p.concurrency; i++ {
			wg.Add(1)
			startY := i * chunkSize
			endY := (i + 1) * chunkSize
			if i == p.concurrency-1 {
				endY = height // Ensure we process all rows
			}
			
			go func(startY, endY int) {
				defer wg.Done()
				for y := startY; y < endY; y++ {
					for x := 0; x < width; x++ {
						// Get RGB values
						r, g, b, _ := img.At(x, y).RGBA()
						
						// Convert to 8-bit color values
						r8 := float64(r >> 8)
						g8 := float64(g >> 8)
						b8 := float64(b >> 8)
						
						// Calculate grayscale value
						gray := uint8(p.weights[0]*r8 + p.weights[1]*g8 + p.weights[2]*b8)
						
						result.SetGray(x, y, color.Gray{Y: gray})
					}
				}
			}(startY, endY)
		}
		
		wg.Wait()
	} else {
		// Single-threaded processing
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				// Get RGB values
				r, g, b, _ := img.At(x, y).RGBA()
				
				// Convert to 8-bit color values
				r8 := float64(r >> 8)
				g8 := float64(g >> 8)
				b8 := float64(b >> 8)
				
				// Calculate grayscale value
				gray := uint8(p.weights[0]*r8 + p.weights[1]*g8 + p.weights[2]*b8)
				
				result.SetGray(x, y, color.Gray{Y: gray})
			}
		}
	}
	
	return result
}

// BoxBlurProcessor applies a box blur filter to an image
// Demonstrates arrays, slice operations, arithmetic operations
type BoxBlurProcessor struct {
	ImageProcessor
	kernelSize int // Size of the kernel (must be odd)
}

func NewBoxBlurProcessor(kernelSize int, concurrency int, useParallel bool) *BoxBlurProcessor {
	// Ensure kernel size is odd
	if kernelSize%2 == 0 {
		kernelSize++
	}
	
	return &BoxBlurProcessor{
		ImageProcessor: NewImageProcessor("Box Blur", concurrency, useParallel),
		kernelSize:     kernelSize,
	}
}

func (p *BoxBlurProcessor) Process(img image.Image) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	result := image.NewRGBA(bounds)
	
	// Calculate the radius of the kernel
	radius := p.kernelSize / 2
	
	// Create channels for parallel processing
	type workItem struct {
		x, y int
	}
	
	var wg sync.WaitGroup
	workChan := make(chan workItem, width*height)
	
	if p.useParallel {
		// Start worker goroutines
		for i := 0; i < p.concurrency; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for work := range workChan {
					x, y := work.x, work.y
					
					// Apply box blur
					var sumR, sumG, sumB, count uint32
					
					for ky := -radius; ky <= radius; ky++ {
						for kx := -radius; kx <= radius; kx++ {
							// Calculate neighbor position
							nx, ny := x+kx, y+ky
							
							// Check if the neighbor is within bounds
							if nx >= 0 && nx < width && ny >= 0 && ny < height {
								r, g, b, _ := img.At(nx, ny).RGBA()
								sumR += r
								sumG += g
								sumB += b
								count++
							}
						}
					}
					
					// Calculate average
					if count > 0 {
						avgR := sumR / count
						avgG := sumG / count
						avgB := sumB / count
						
						result.Set(x, y, color.RGBA{
							R: uint8(avgR >> 8),
							G: uint8(avgG >> 8),
							B: uint8(avgB >> 8),
							A: 255,
						})
					}
				}
			}()
		}
		
		// Send work items
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				workChan <- workItem{x, y}
			}
		}
		
		close(workChan)
		wg.Wait()
	} else {
		// Single-threaded processing
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				// Apply box blur
				var sumR, sumG, sumB, count uint32
				
				for ky := -radius; ky <= radius; ky++ {
					for kx := -radius; kx <= radius; kx++ {
						// Calculate neighbor position
						nx, ny := x+kx, y+ky
						
						// Check if the neighbor is within bounds
						if nx >= 0 && nx < width && ny >= 0 && ny < height {
							r, g, b, _ := img.At(nx, ny).RGBA()
							sumR += r
							sumG += g
							sumB += b
							count++
						}
					}
				}
				
				// Calculate average
				if count > 0 {
					avgR := sumR / count
					avgG := sumG / count
					avgB := sumB / count
					
					result.Set(x, y, color.RGBA{
						R: uint8(avgR >> 8),
						G: uint8(avgG >> 8),
						B: uint8(avgB >> 8),
						A: 255,
					})
				}
			}
		}
	}
	
	return result
}

// GaussianBlurProcessor applies a Gaussian blur filter to an image
// Demonstrates 2D arrays, nested loops, error handling
type GaussianBlurProcessor struct {
	ImageProcessor
	sigma      float64   // Standard deviation for Gaussian
	kernelSize int       // Size of the kernel (must be odd)
	kernel     [][]float64 // Precalculated kernel
}

func NewGaussianBlurProcessor(sigma float64, kernelSize int, concurrency int, useParallel bool) *GaussianBlurProcessor {
	// Ensure kernel size is odd
	if kernelSize%2 == 0 {
		kernelSize++
	}
	
	processor := &GaussianBlurProcessor{
		ImageProcessor: NewImageProcessor("Gaussian Blur", concurrency, useParallel),
		sigma:          sigma,
		kernelSize:     kernelSize,
	}
	
	// Generate the Gaussian kernel
	processor.kernel = processor.generateGaussianKernel()
	
	return processor
}

// Generate the Gaussian kernel
func (p *GaussianBlurProcessor) generateGaussianKernel() [][]float64 {
	radius := p.kernelSize / 2
	kernel := make([][]float64, p.kernelSize)
	
	// Gaussian function: G(x,y) = (1/(2*pi*sigma^2)) * e^(-(x^2+y^2)/(2*sigma^2))
	// We'll skip the normalization factor (1/(2*pi*sigma^2)) and normalize at the end
	twoSigmaSquared := 2 * p.sigma * p.sigma
	
	// Calculate kernel values
	sum := 0.0
	for y := -radius; y <= radius; y++ {
		kernel[y+radius] = make([]float64, p.kernelSize)
		for x := -radius; x <= radius; x++ {
			// Calculate Gaussian value
			exponent := -(float64(x*x+y*y) / twoSigmaSquared)
			value := math.Exp(exponent)
			kernel[y+radius][x+radius] = value
			sum += value
		}
	}
	
	// Normalize the kernel so it sums to 1
	for y := 0; y < p.kernelSize; y++ {
		for x := 0; x < p.kernelSize; x++ {
			kernel[y][x] /= sum
		}
	}
	
	return kernel
}

func (p *GaussianBlurProcessor) Process(img image.Image) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	result := image.NewRGBA(bounds)
	
	// Calculate the radius of the kernel
	radius := p.kernelSize / 2
	
	// Create a worker pool for parallel processing
	type pixelWork struct {
		x, y int
	}
	
	var wg sync.WaitGroup
	workChan := make(chan pixelWork, width*height)
	
	if p.useParallel {
		// Start worker goroutines
		for i := 0; i < p.concurrency; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for work := range workChan {
					x, y := work.x, work.y
					
					// Apply Gaussian blur
					var sumR, sumG, sumB float64
					
					for ky := 0; ky < p.kernelSize; ky++ {
						for kx := 0; kx < p.kernelSize; kx++ {
							// Calculate source pixel position
							sx := x + (kx - radius)
							sy := y + (ky - radius)
							
							// Check if source pixel is within bounds
							if sx >= 0 && sx < width && sy >= 0 && sy < height {
								// Get color
								r, g, b, _ := img.At(sx, sy).RGBA()
								
								// Convert to float and apply kernel weight
								weight := p.kernel[ky][kx]
								sumR += float64(r) * weight
								sumG += float64(g) * weight
								sumB += float64(b) * weight
							}
						}
					}
					
					// Set the result pixel
					result.Set(x, y, color.RGBA{
						R: uint8(sumR / 65535),
						G: uint8(sumG / 65535),
						B: uint8(sumB / 65535),
						A: 255,
					})
				}
			}()
		}
		
		// Send work items
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				workChan <- pixelWork{x, y}
			}
		}
		
		close(workChan)
		wg.Wait()
	} else {
		// Single-threaded processing
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				// Apply Gaussian blur
				var sumR, sumG, sumB float64
				
				for ky := 0; ky < p.kernelSize; ky++ {
					for kx := 0; kx < p.kernelSize; kx++ {
						// Calculate source pixel position
						sx := x + (kx - radius)
						sy := y + (ky - radius)
						
						// Check if source pixel is within bounds
						if sx >= 0 && sx < width && sy >= 0 && sy < height {
							// Get color
							r, g, b, _ := img.At(sx, sy).RGBA()
							
							// Convert to float and apply kernel weight
							weight := p.kernel[ky][kx]
							sumR += float64(r) * weight
							sumG += float64(g) * weight
							sumB += float64(b) * weight
						}
					}
				}
				
				// Set the result pixel
				result.Set(x, y, color.RGBA{
					R: uint8(sumR / 65535),
					G: uint8(sumG / 65535),
					B: uint8(sumB / 65535),
					A: 255,
				})
			}
		}
	}
	
	return result
}

// SobelEdgeDetector detects edges using Sobel operator
// Demonstrates 2D arrays, image processing, error handling
type SobelEdgeDetector struct {
	ImageProcessor
	threshold uint8 // Threshold for edge detection
}

func NewSobelEdgeDetector(threshold uint8, concurrency int, useParallel bool) *SobelEdgeDetector {
	return &SobelEdgeDetector{
		ImageProcessor: NewImageProcessor("Sobel Edge Detection", concurrency, useParallel),
		threshold:      threshold,
	}
}

func (p *SobelEdgeDetector) Process(img image.Image) image.Image {
	// First convert the image to grayscale
	grayscaleProcessor := NewGrayscaleProcessor(p.concurrency, p.useParallel)
	grayImg := grayscaleProcessor.Process(img)
	
	bounds := grayImg.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	result := image.NewGray(bounds)
	
	// Sobel operators
	sobelX := [][]int{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1},
	}
	
	sobelY := [][]int{
		{-1, -2, -1},
		{0, 0, 0},
		{1, 2, 1},
	}
	
	// Function to safely get grayscale value
	getGray := func(x, y int) uint8 {
		if x < 0 || x >= width || y < 0 || y >= height {
			return 0
		}
		switch grayImg := grayImg.(type) {
		case *image.Gray:
			return grayImg.GrayAt(x, y).Y
		default:
			r, g, b, _ := grayImg.At(x, y).RGBA()
			return uint8((r + g + b) / 3 >> 8)
		}
	}
	
	// Process with goroutines for parallel execution
	var wg sync.WaitGroup
	
	if p.useParallel {
		// We'll skip the border pixels for simplicity
		for y := 1; y < height-1; y++ {
			wg.Add(1)
			go func(y int) {
				defer wg.Done()
				
				for x := 1; x < width-1; x++ {
					// Apply Sobel operators
					var gx, gy int
					
					for ky := -1; ky <= 1; ky++ {
						for kx := -1; kx <= 1; kx++ {
							val := int(getGray(x+kx, y+ky))
							gx += val * sobelX[ky+1][kx+1]
							gy += val * sobelY[ky+1][kx+1]
						}
					}
					
					// Calculate gradient magnitude
					magnitude := math.Sqrt(float64(gx*gx + gy*gy))
					
					// Apply threshold
					var pixel uint8
					if magnitude > float64(p.threshold) {
						pixel = 255
					}
					
					result.SetGray(x, y, color.Gray{Y: pixel})
				}
			}(y)
		}
		
		wg.Wait()
	} else {
		// Single-threaded processing
		for y := 1; y < height-1; y++ {
			for x := 1; x < width-1; x++ {
				// Apply Sobel operators
				var gx, gy int
				
				for ky := -1; ky <= 1; ky++ {
					for kx := -1; kx <= 1; kx++ {
						val := int(getGray(x+kx, y+ky))
						gx += val * sobelX[ky+1][kx+1]
						gy += val * sobelY[ky+1][kx+1]
					}
				}
				
				// Calculate gradient magnitude
				magnitude := math.Sqrt(float64(gx*gx + gy*gy))
				
				// Apply threshold
				var pixel uint8
				if magnitude > float64(p.threshold) {
					pixel = 255
				}
				
				result.SetGray(x, y, color.Gray{Y: pixel})
			}
		}
	}
	
	return result
}

// Create a pipeline of image processing algorithms
// Demonstrates the use of channels, interfaces, and goroutines
func ProcessImagePipeline(img image.Image, algorithms []ImageProcessingAlgorithm) image.Image {
	// Create channels to pass images between stages
	channels := make([]chan image.Image, len(algorithms))
	for i := range channels {
		channels[i] = make(chan image.Image)
	}
	
	// Start goroutines for each processing stage
	var wg sync.WaitGroup
	wg.Add(len(algorithms))
	
	for i, algorithm := range algorithms {
		go func(i int, alg ImageProcessingAlgorithm) {
			defer wg.Done()
			defer func() {
				if i < len(channels)-1 {
					close(channels[i])
				}
			}()
			
			// Get input from previous stage or use original image for first stage
			var input image.Image
			if i == 0 {
				input = img
			} else {
				input = <-channels[i-1]
			}
			
			// Apply algorithm
			output := alg.Process(input)
			
			// Send output to next stage or return result for last stage
			if i < len(channels)-1 {
				channels[i] <- output
			} else {
				channels[i] <- output
			}
		}(i, algorithm)
	}
	
	// Create a goroutine to close the last channel when all processing is done
	go func() {
		wg.Wait()
		if len(channels) > 0 {
			close(channels[len(channels)-1])
		}
	}()
	
	// Return the result from the last stage
	if len(algorithms) > 0 {
		return <-channels[len(channels)-1]
	}
	
	return img
}