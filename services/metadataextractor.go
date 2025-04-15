package services

import (
        "regexp"
        "strings"
        "time"
        "unicode"
        
        "ancient-script-decoder/models"
        "ancient-script-decoder/utils"
)

// MetadataConfig contains settings for metadata extraction
type MetadataConfig struct {
        EnableGeographicDetection bool    `yaml:"enableGeographicDetection"`
        EnablePeriodDetection     bool    `yaml:"enablePeriodDetection"`
        EnableCultureDetection    bool    `yaml:"enableCultureDetection"`
        ContextSensitivity        float64 `yaml:"contextSensitivity"`
        ReferenceDatabase         string  `yaml:"referenceDatabase"`
}

// MetadataExtractor handles extraction of historical context from manuscripts
type MetadataExtractor struct {
        config MetadataConfig
        logger *utils.Logger
        
        // Maps for geo-temporal context detection
        periodKeywords  map[string][]models.TimePeriod
        regionKeywords  map[string][]models.Region
        cultureKeywords map[string][]string
}

// NewMetadataExtractor creates a new metadata extractor service
func NewMetadataExtractor(config MetadataConfig, logger *utils.Logger) *MetadataExtractor {
        extractor := &MetadataExtractor{
                config: config,
                logger: logger,
        }
        
        extractor.initializeKeywordMaps()
        return extractor
}

// initializeKeywordMaps sets up the reference data for context extraction
func (m *MetadataExtractor) initializeKeywordMaps() {
        // Initialize period detection keywords
        m.periodKeywords = map[string][]models.TimePeriod{
                "ancient": {
                        {
                                Name:        "Ancient Period",
                                StartYear:   -3000,
                                EndYear:     500,
                                Description: "The time period of ancient civilizations through the fall of the Western Roman Empire",
                        },
                },
                "medieval": {
                        {
                                Name:        "Medieval Period",
                                StartYear:   500,
                                EndYear:     1500,
                                Description: "The middle ages from the fall of Rome to the Renaissance",
                        },
                },
                "renaissance": {
                        {
                                Name:        "Renaissance",
                                StartYear:   1300,
                                EndYear:     1700,
                                Description: "Period of cultural rebirth and renewed interest in classical learning",
                        },
                },
                "bronze age": {
                        {
                                Name:        "Bronze Age",
                                StartYear:   -3000,
                                EndYear:     -1200,
                                Description: "Period characterized by the use of bronze and early writing systems",
                        },
                },
                "iron age": {
                        {
                                Name:        "Iron Age",
                                StartYear:   -1200,
                                EndYear:     -500,
                                Description: "Period characterized by the widespread use of iron",
                        },
                },
                "classical": {
                        {
                                Name:        "Classical Period",
                                StartYear:   -800,
                                EndYear:     500,
                                Description: "Greco-Roman classical period of literature and art",
                        },
                },
                "hellenistic": {
                        {
                                Name:        "Hellenistic Period",
                                StartYear:   -323,
                                EndYear:     -31,
                                Description: "Period between Alexander the Great and the rise of the Roman Empire",
                        },
                },
        }
        
        // Initialize region detection keywords
        m.regionKeywords = map[string][]models.Region{
                "mesopotamia": {
                        {
                                Name:        "Mesopotamia",
                                ModernAreas: []string{"Iraq", "Syria", "Turkey", "Iran"},
                                Description: "Area between the Tigris and Euphrates rivers, cradle of civilization",
                        },
                },
                "egypt": {
                        {
                                Name:        "Ancient Egypt",
                                ModernAreas: []string{"Egypt"},
                                Description: "Civilization along the lower Nile River, known for pyramids and hieroglyphs",
                        },
                },
                "greece": {
                        {
                                Name:        "Ancient Greece",
                                ModernAreas: []string{"Greece", "Turkey", "Italy", "Libya", "Egypt"},
                                Description: "Civilization that flourished around the Mediterranean Sea, known for philosophy and democracy",
                        },
                },
                "rome": {
                        {
                                Name:        "Ancient Rome",
                                ModernAreas: []string{"Italy", "Mediterranean Basin", "Europe", "North Africa", "Middle East"},
                                Description: "Civilization that expanded from the city of Rome to control the Mediterranean region",
                        },
                },
                "persia": {
                        {
                                Name:        "Ancient Persia",
                                ModernAreas: []string{"Iran", "Iraq", "Turkey", "Central Asia", "Caucasus", "Pakistan"},
                                Description: "One of the oldest civilizations in history, centered in modern-day Iran",
                        },
                },
                "maya": {
                        {
                                Name:        "Maya Civilization",
                                ModernAreas: []string{"Mexico", "Guatemala", "Belize", "Honduras", "El Salvador"},
                                Description: "Mesoamerican civilization known for advanced writing, art, and astronomical systems",
                        },
                },
                "china": {
                        {
                                Name:        "Ancient China",
                                ModernAreas: []string{"China"},
                                Description: "One of the world's oldest continuous civilizations, known for innovations in bureaucracy, philosophy, and technology",
                        },
                },
                "india": {
                        {
                                Name:        "Ancient India",
                                ModernAreas: []string{"India", "Pakistan", "Bangladesh", "Nepal"},
                                Description: "Civilization of the Indian subcontinent, known for religious and philosophical traditions",
                        },
                },
        }
        
        // Initialize culture detection keywords
        m.cultureKeywords = map[string][]string{
                "sumerian":    {"cuneiform", "ziggurat", "city-state", "mesopotamia"},
                "egyptian":    {"pharaoh", "hieroglyph", "pyramid", "nile", "mummy"},
                "greek":       {"polis", "acropolis", "agora", "philosophy", "democracy", "olympian"},
                "roman":       {"senate", "republic", "legion", "empire", "consul", "caesar"},
                "christian":   {"church", "monastery", "bishop", "pope", "scripture", "gospel"},
                "islamic":     {"mosque", "caliph", "quran", "sultan", "hadith"},
                "persian":     {"zoroastrian", "achaemenid", "sassanid", "ahura mazda"},
                "hebrew":      {"temple", "covenant", "torah", "prophet", "synagogue"},
                "viking":      {"norse", "fjord", "saga", "rune", "drakkar", "valhalla"},
                "celtic":      {"druid", "clan", "ogham", "gaul", "tribe"},
                "babylonian":  {"hammurabi", "marduk", "babylon", "euphrates", "ishtar"},
                "assyrian":    {"ashur", "nineveh", "lamassu", "tiglath"},
                "hittite":     {"anatolia", "hattusa", "tarhun", "hattian"},
                "phoenician":  {"alphabet", "tyre", "sidon", "byblos", "carthage", "purple"},
                "etruscan":    {"tuscany", "haruspex", "rite", "lucumo"},
                "byzantine":   {"constantinople", "orthodox", "basilica", "theodosius", "justinian"},
        }
}

// ExtractMetadata analyzes text content to extract historical metadata
// If imageData is nil, extraction will be based only on text content
func (m *MetadataExtractor) ExtractMetadata(text string, scriptType string, imageData []byte) (models.Metadata, error) {
        metadata := models.Metadata{
                ConfidenceScore: 0.0,
                ScriptType:      scriptType,
                DetectedDate:    time.Now().Format(time.RFC3339),
        }
        
        // Extract time periods
        periods := m.extractTimePeriods(text)
        if len(periods) > 0 {
                metadata.TimePeriods = periods
                metadata.ConfidenceScore += 0.2
        }
        
        // Extract geographical regions
        regions := m.extractRegions(text)
        if len(regions) > 0 {
                metadata.Regions = regions
                metadata.ConfidenceScore += 0.2
        }
        
        // Extract cultural context
        cultures := m.extractCulturalContext(text)
        if len(cultures) > 0 {
                metadata.CulturalContext = cultures
                metadata.ConfidenceScore += 0.2
        }
        
        // Analyze script type to increase confidence
        if scriptType != "auto" && scriptType != "" {
                metadata.ConfidenceScore += 0.2
                
                // Associate script type with likely cultures
                switch scriptType {
                case "cuneiform":
                        if !contains(metadata.CulturalContext, "Sumerian") {
                                metadata.CulturalContext = append(metadata.CulturalContext, "Sumerian")
                        }
                        if !contains(metadata.CulturalContext, "Babylonian") {
                                metadata.CulturalContext = append(metadata.CulturalContext, "Babylonian")
                        }
                case "hieroglyphic":
                        if !contains(metadata.CulturalContext, "Egyptian") {
                                metadata.CulturalContext = append(metadata.CulturalContext, "Egyptian")
                        }
                case "greek":
                        if !contains(metadata.CulturalContext, "Greek") {
                                metadata.CulturalContext = append(metadata.CulturalContext, "Greek")
                        }
                case "latin":
                        if !contains(metadata.CulturalContext, "Roman") {
                                metadata.CulturalContext = append(metadata.CulturalContext, "Roman")
                        }
                case "runic":
                        if !contains(metadata.CulturalContext, "Norse") {
                                metadata.CulturalContext = append(metadata.CulturalContext, "Norse")
                        }
                }
        }
        
        // Add material analysis based on script type
        metadata.MaterialContext = m.analyzeMaterial(scriptType)
        
        // Extract potential historical events based on text content
        events := m.extractHistoricalEvents(text)
        if len(events) > 0 {
                metadata.HistoricalEvents = events
                metadata.ConfidenceScore += 0.2
        }
        
        // Cap confidence at 1.0
        if metadata.ConfidenceScore > 1.0 {
                metadata.ConfidenceScore = 1.0
        }
        
        return metadata, nil
}

// extractTimePeriods identifies time periods mentioned in the text
func (m *MetadataExtractor) extractTimePeriods(text string) []models.TimePeriod {
        text = strings.ToLower(text)
        var periods []models.TimePeriod
        periodSet := make(map[string]bool)
        
        for keyword, possiblePeriods := range m.periodKeywords {
                if strings.Contains(text, keyword) {
                        for _, period := range possiblePeriods {
                                if !periodSet[period.Name] {
                                        periods = append(periods, period)
                                        periodSet[period.Name] = true
                                }
                        }
                }
        }
        
        // Try to detect century-based references
        centuryRegex := regexp.MustCompile(`(\d+)(st|nd|rd|th) century`)
        matches := centuryRegex.FindAllStringSubmatch(text, -1)
        
        for _, match := range matches {
                if len(match) >= 2 {
                        century, err := utils.StringToInt(match[1])
                        if err == nil {
                                startYear := (century - 1) * 100
                                endYear := century * 100
                                period := models.TimePeriod{
                                        Name:        match[0],
                                        StartYear:   startYear,
                                        EndYear:     endYear,
                                        Description: "Period spanning the " + match[0],
                                }
                                
                                if !periodSet[period.Name] {
                                        periods = append(periods, period)
                                        periodSet[period.Name] = true
                                }
                        }
                }
        }
        
        return periods
}

// extractRegions identifies geographical regions mentioned in the text
func (m *MetadataExtractor) extractRegions(text string) []models.Region {
        text = strings.ToLower(text)
        var regions []models.Region
        regionSet := make(map[string]bool)
        
        for keyword, possibleRegions := range m.regionKeywords {
                if strings.Contains(text, keyword) {
                        for _, region := range possibleRegions {
                                if !regionSet[region.Name] {
                                        regions = append(regions, region)
                                        regionSet[region.Name] = true
                                }
                        }
                }
        }
        
        return regions
}

// extractCulturalContext identifies cultural contexts mentioned in the text
func (m *MetadataExtractor) extractCulturalContext(text string) []string {
        text = strings.ToLower(text)
        var cultures []string
        cultureSet := make(map[string]bool)
        
        // Look for direct culture names
        for culture, keywords := range m.cultureKeywords {
                if strings.Contains(text, culture) {
                        capitalizedCulture := capitalize(culture)
                        if !cultureSet[capitalizedCulture] {
                                cultures = append(cultures, capitalizedCulture)
                                cultureSet[capitalizedCulture] = true
                        }
                } else {
                        // Check related keywords
                        keywordCount := 0
                        for _, keyword := range keywords {
                                if strings.Contains(text, keyword) {
                                        keywordCount++
                                }
                        }
                        
                        // If we found multiple keywords, it's likely this culture
                        if keywordCount >= 2 {
                                capitalizedCulture := capitalize(culture)
                                if !cultureSet[capitalizedCulture] {
                                        cultures = append(cultures, capitalizedCulture)
                                        cultureSet[capitalizedCulture] = true
                                }
                        }
                }
        }
        
        return cultures
}

// analyzeMaterial determines likely material based on script type
func (m *MetadataExtractor) analyzeMaterial(scriptType string) []string {
        var materials []string
        
        switch scriptType {
        case "cuneiform":
                materials = append(materials, "Clay tablet", "Stone")
        case "hieroglyphic":
                materials = append(materials, "Papyrus", "Stone", "Wood")
        case "latin", "greek":
                materials = append(materials, "Parchment", "Papyrus", "Wax tablet")
        case "runic":
                materials = append(materials, "Stone", "Wood", "Bone", "Metal")
        default:
                materials = append(materials, "Parchment", "Paper", "Stone")
        }
        
        return materials
}

// extractHistoricalEvents identifies potential historical events in the text
func (m *MetadataExtractor) extractHistoricalEvents(text string) []models.HistoricalEvent {
        var events []models.HistoricalEvent
        
        // Look for event patterns like battles, treaties, reigns, etc.
        battleRegex := regexp.MustCompile(`(?i)battle of ([a-zA-Z\s]+)`)
        treatyRegex := regexp.MustCompile(`(?i)treaty of ([a-zA-Z\s]+)`)
        reignRegex := regexp.MustCompile(`(?i)reign of ([a-zA-Z\s]+)`)
        
        // Find battles
        battleMatches := battleRegex.FindAllStringSubmatch(text, -1)
        for _, match := range battleMatches {
                if len(match) >= 2 {
                        event := models.HistoricalEvent{
                                Name:        "Battle of " + match[1],
                                EventType:   "Battle",
                                Description: "A military conflict that took place at " + match[1],
                        }
                        events = append(events, event)
                }
        }
        
        // Find treaties
        treatyMatches := treatyRegex.FindAllStringSubmatch(text, -1)
        for _, match := range treatyMatches {
                if len(match) >= 2 {
                        event := models.HistoricalEvent{
                                Name:        "Treaty of " + match[1],
                                EventType:   "Treaty",
                                Description: "A formal agreement between political entities signed at " + match[1],
                        }
                        events = append(events, event)
                }
        }
        
        // Find reigns
        reignMatches := reignRegex.FindAllStringSubmatch(text, -1)
        for _, match := range reignMatches {
                if len(match) >= 2 {
                        event := models.HistoricalEvent{
                                Name:        "Reign of " + match[1],
                                EventType:   "Reign",
                                Description: "The period during which " + match[1] + " held sovereign power",
                        }
                        events = append(events, event)
                }
        }
        
        return events
}

// Helper function to check if a string is in a slice
func contains(slice []string, item string) bool {
        for _, s := range slice {
                if strings.EqualFold(s, item) {
                        return true
                }
        }
        return false
}

// capitalize returns a properly capitalized version of a word or term
func capitalize(s string) string {
        if s == "" {
                return ""
        }
        
        words := strings.Split(s, " ")
        for i, word := range words {
                if len(word) == 0 {
                        continue
                }
                r := []rune(word)
                words[i] = string(append([]rune{unicode.ToUpper(r[0])}, r[1:]...))
        }
        
        return strings.Join(words, " ")
}