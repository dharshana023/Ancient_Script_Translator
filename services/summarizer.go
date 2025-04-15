package services

import (
        "errors"
        "fmt"
        "math"
        "regexp"
        "sort"
        "strings"
        "sync"

        "github.com/sajari/word2vec"
)

// SummarizationConfig contains the configuration for text summarization
type SummarizationConfig struct {
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
}

// Summarizer handles the summarization of translated texts
type Summarizer struct {
        config        SummarizationConfig
        model         *word2vec.Model
        modelLoaded   bool
        stopwords     map[string]bool
        mutex         sync.RWMutex
}

// NewSummarizer creates a new summarizer
func NewSummarizer(config SummarizationConfig) *Summarizer {
        s := &Summarizer{
                config:      config,
                modelLoaded: false,
                stopwords:   buildStopwordsMap(),
        }

        // Initialize the model if a path is provided
        if config.ModelPath != "" {
                s.loadModel()
        }

        return s
}

// loadModel loads the word embedding model
func (s *Summarizer) loadModel() {
        // In a real implementation, this would load the specified model
        // For this simulation, we'll set modelLoaded to false
        s.modelLoaded = false
}

// SummarizeText generates a summary for the translated text
func (s *Summarizer) SummarizeText(text string) (string, error) {
        if text == "" {
                return "", errors.New("cannot summarize empty text")
        }

        // Select the summarization algorithm based on configuration
        switch s.config.Algorithm {
        case "extractive":
                return s.extractiveSummarization(text)
        case "abstractive":
                return s.abstractiveSummarization(text)
        case "hybrid":
                return s.hybridSummarization(text)
        default:
                // Default to TextRank-based extractive summarization
                return s.textRankSummarization(text)
        }
}

// extractiveSummarization implements an extractive summarization algorithm
func (s *Summarizer) extractiveSummarization(text string) (string, error) {
        // Split text into sentences
        sentences := splitIntoSentences(text)
        if len(sentences) == 0 {
                return "", errors.New("no valid sentences found in text")
        }

        // Score each sentence
        sentenceScores := make(map[string]float64)
        for _, sentence := range sentences {
                sentenceScores[sentence] = s.scoreSentence(sentence, text)
        }

        // Extract the top sentences
        topSentences := s.extractTopSentences(sentences, sentenceScores)

        // Combine sentences in their original order
        return s.combineSentencesInOrder(topSentences, sentences), nil
}

// abstractiveSummarization implements an abstractive summarization algorithm
func (s *Summarizer) abstractiveSummarization(text string) (string, error) {
        // This would typically use a neural model to generate a summary
        // For this simulation, we'll use a simplified approach
        
        // Extract key concepts
        concepts := s.extractKeyConcepts(text)
        
        // Generate sentences based on key concepts
        summary := s.generateSummaryFromConcepts(concepts, text)
        
        // Ensure summary is within length limit
        if len(summary) > s.config.MaxSummaryLength {
                return summary[:s.config.MaxSummaryLength] + "...", nil
        }
        
        return summary, nil
}

// hybridSummarization combines extractive and abstractive approaches
func (s *Summarizer) hybridSummarization(text string) (string, error) {
        // First extract important sentences
        extractiveSummary, err := s.extractiveSummarization(text)
        if err != nil {
                return "", err
        }
        
        // Then apply abstractive techniques to refine
        concepts := s.extractKeyConcepts(extractiveSummary)
        summary := s.refineExtractiveWithAbstractive(extractiveSummary, concepts)
        
        // Ensure summary is within length limit
        if len(summary) > s.config.MaxSummaryLength {
                return summary[:s.config.MaxSummaryLength] + "...", nil
        }
        
        return summary, nil
}

// textRankSummarization implements a TextRank-inspired algorithm
func (s *Summarizer) textRankSummarization(text string) (string, error) {
        // Split text into sentences
        sentences := splitIntoSentences(text)
        if len(sentences) == 0 {
                return "", errors.New("no valid sentences found in text")
        }
        
        // Create a sentence similarity graph
        similarityGraph := s.buildSimilarityGraph(sentences)
        
        // Apply PageRank algorithm to find important sentences
        sentenceScores := s.applyPageRank(similarityGraph, 0.85, 30)
        
        // Extract top sentences
        topSentences := s.extractTopSentences(sentences, sentenceScores)
        
        // Combine sentences in their original order
        return s.combineSentencesInOrder(topSentences, sentences), nil
}

// scoreSentence scores a sentence based on various factors
func (s *Summarizer) scoreSentence(sentence, fullText string) float64 {
        // Normalize sentence by removing stopwords and converting to lowercase
        normalizedSentence := s.normalizeSentence(sentence)
        words := strings.Fields(normalizedSentence)
        
        // Skip very short sentences
        if len(words) < s.config.MinSentenceLength {
                return 0.0
        }
        
        // Calculate position score (sentences at beginning/end are often more important)
        positionScore := s.calculatePositionScore(sentence, fullText)
        
        // Calculate keyword frequency score
        keywordScore := s.calculateKeywordScore(words, fullText)
        
        // Calculate context relevance score
        contextScore := s.calculateContextScore(sentence, fullText)
        
        // Calculate entity presence score if enabled
        entityScore := 0.0
        if s.config.EnableEntityRecog {
                entityScore = s.calculateEntityScore(sentence)
        }
        
        // Weighted combination of scores
        totalScore := (s.config.SentenceImportance * positionScore) +
                (s.config.KeywordImportance * keywordScore) +
                (s.config.ContextImportance * contextScore) +
                (entityScore * 0.5) // Fixed weight for entity score
        
        return totalScore
}

// normalizeSentence normalizes a sentence for processing
func (s *Summarizer) normalizeSentence(sentence string) string {
        // Convert to lowercase
        normalized := strings.ToLower(sentence)
        
        // Remove punctuation
        normalized = strings.Map(func(r rune) rune {
                if strings.ContainsRune(".,;:!?\"'()[]{}", r) {
                        return -1
                }
                return r
        }, normalized)
        
        // Remove stopwords
        words := strings.Fields(normalized)
        filteredWords := make([]string, 0, len(words))
        
        for _, word := range words {
                if !s.stopwords[word] {
                        filteredWords = append(filteredWords, word)
                }
        }
        
        return strings.Join(filteredWords, " ")
}

// calculatePositionScore scores sentences based on their position
func (s *Summarizer) calculatePositionScore(sentence, fullText string) float64 {
        sentences := splitIntoSentences(fullText)
        position := 0
        
        for i, sent := range sentences {
                if sent == sentence {
                        position = i
                        break
                }
        }
        
        // Sentences at the beginning and end of the text are often more important
        numSentences := len(sentences)
        if numSentences <= 1 {
                return 1.0
        }
        
        // Higher score for sentences at the beginning, slightly lower for end, lowest for middle
        if position < int(float64(numSentences)*0.2) {
                return 1.0 - (float64(position) / float64(numSentences))
        } else if position > int(float64(numSentences)*0.8) {
                return 0.5 + (float64(position) / float64(numSentences)) * 0.5
        } else {
                return 0.5 * (1.0 - math.Abs(float64(position-numSentences/2)/float64(numSentences/2)))
        }
}

// calculateKeywordScore scores sentences based on presence of keywords
func (s *Summarizer) calculateKeywordScore(words []string, fullText string) float64 {
        // Extract the top keywords from the text
        keywords := s.extractKeywords(fullText, 10)
        keywordSet := make(map[string]bool)
        
        for _, keyword := range keywords {
                keywordSet[keyword] = true
        }
        
        // Count how many keywords are in the sentence
        keywordCount := 0
        for _, word := range words {
                if keywordSet[word] {
                        keywordCount++
                }
        }
        
        // Return the ratio of keywords present
        if len(keywords) == 0 {
                return 0.0
        }
        return float64(keywordCount) / float64(len(keywords))
}

// calculateContextScore scores sentences based on their context relevance
func (s *Summarizer) calculateContextScore(sentence, fullText string) float64 {
        sentences := splitIntoSentences(fullText)
        sentenceIndex := -1
        
        // Find the index of the current sentence
        for i, sent := range sentences {
                if sent == sentence {
                        sentenceIndex = i
                        break
                }
        }
        
        if sentenceIndex == -1 {
                return 0.0
        }
        
        // Get the context window around the sentence
        windowSize := s.config.ContextWindowSize
        startIdx := max(0, sentenceIndex-windowSize)
        endIdx := min(len(sentences)-1, sentenceIndex+windowSize)
        
        // Calculate similarity with surrounding sentences
        similarity := 0.0
        windowSentences := 0
        
        for i := startIdx; i <= endIdx; i++ {
                if i != sentenceIndex {
                        similarity += s.calculateSentenceSimilarity(sentence, sentences[i])
                        windowSentences++
                }
        }
        
        if windowSentences == 0 {
                return 0.0
        }
        
        return similarity / float64(windowSentences)
}

// calculateEntityScore scores sentences based on presence of named entities
func (s *Summarizer) calculateEntityScore(sentence string) float64 {
        // In a real implementation, this would use NER to identify entities
        // For this simulation, we'll use a simple heuristic based on capitalization
        words := strings.Fields(sentence)
        entityCount := 0
        
        for _, word := range words {
                if len(word) > 1 && isCapitalized(word) {
                        entityCount++
                }
        }
        
        return float64(entityCount) / float64(max(1, len(words)))
}

// extractKeywords extracts the top keywords from the text
func (s *Summarizer) extractKeywords(text string, numKeywords int) []string {
        // Normalize the text
        normalized := strings.ToLower(text)
        
        // Remove punctuation
        normalized = strings.Map(func(r rune) rune {
                if strings.ContainsRune(".,;:!?\"'()[]{}", r) {
                        return -1
                }
                return r
        }, normalized)
        
        // Split into words and remove stopwords
        words := strings.Fields(normalized)
        filteredWords := make([]string, 0, len(words))
        
        for _, word := range words {
                if !s.stopwords[word] && len(word) > 2 { // Ignore very short words
                        filteredWords = append(filteredWords, word)
                }
        }
        
        // Count word frequencies and calculate TF (Term Frequency)
        wordFreq := make(map[string]int)
        for _, word := range filteredWords {
                wordFreq[word]++
        }
        
        // Calculate TF-IDF-like score (simplified)
        // Higher weight for unique/distinctive words
        totalWords := float64(len(filteredWords))
        wordScores := make(map[string]float64)
        
        for word, count := range wordFreq {
                // Basic TF calculation
                tf := float64(count) / totalWords
                
                // Add weight for unique words (words that appear less frequently gain importance)
                // This helps differentiate between documents
                uniqueness := 1.0
                if count > 1 {
                        uniqueness = 1.0 / math.Log(float64(count)+1.0)
                }
                
                // Weight by word length (longer words often carry more meaning)
                lengthWeight := math.Log(float64(len(word)) + 1.0)
                
                // Combine the factors
                wordScores[word] = tf * uniqueness * lengthWeight
        }
        
        // Create sorted slice of word scores
        type wordScore struct {
                word  string
                score float64
        }
        
        scores := make([]wordScore, 0, len(wordScores))
        for word, score := range wordScores {
                scores = append(scores, wordScore{word, score})
        }
        
        // Sort by score (descending)
        sort.Slice(scores, func(i, j int) bool {
                return scores[i].score > scores[j].score
        })
        
        // Take the top keywords
        n := min(numKeywords, len(scores))
        keywords := make([]string, n)
        
        for i := 0; i < n; i++ {
                keywords[i] = scores[i].word
        }
        
        return keywords
}

// sortFrequencies sorts word frequencies in descending order
func sortFrequencies(frequencies []struct{ word string; count int }) {
        // Simple bubble sort for demonstration
        n := len(frequencies)
        for i := 0; i < n-1; i++ {
                for j := 0; j < n-i-1; j++ {
                        if frequencies[j].count < frequencies[j+1].count {
                                frequencies[j], frequencies[j+1] = frequencies[j+1], frequencies[j]
                        }
                }
        }
}

// calculateSentenceSimilarity calculates similarity between two sentences
func (s *Summarizer) calculateSentenceSimilarity(sent1, sent2 string) float64 {
        // Normalize sentences
        norm1 := s.normalizeSentence(sent1)
        norm2 := s.normalizeSentence(sent2)
        
        // Extract word sets
        words1 := strings.Fields(norm1)
        words2 := strings.Fields(norm2)
        
        // Create word sets
        set1 := make(map[string]bool)
        set2 := make(map[string]bool)
        
        for _, word := range words1 {
                set1[word] = true
        }
        
        for _, word := range words2 {
                set2[word] = true
        }
        
        // Calculate Jaccard similarity
        intersection := 0
        for word := range set1 {
                if set2[word] {
                        intersection++
                }
        }
        
        union := len(set1) + len(set2) - intersection
        
        if union == 0 {
                return 0.0
        }
        
        return float64(intersection) / float64(union)
}

// buildSimilarityGraph builds a graph of sentence similarities
func (s *Summarizer) buildSimilarityGraph(sentences []string) [][]float64 {
        n := len(sentences)
        graph := make([][]float64, n)
        
        for i := range graph {
                graph[i] = make([]float64, n)
        }
        
        for i := 0; i < n; i++ {
                for j := 0; j < n; j++ {
                        if i != j {
                                graph[i][j] = s.calculateSentenceSimilarity(sentences[i], sentences[j])
                        }
                }
        }
        
        return graph
}

// applyPageRank applies the PageRank algorithm to rank sentences
func (s *Summarizer) applyPageRank(graph [][]float64, dampingFactor float64, iterations int) map[string]float64 {
        n := len(graph)
        if n == 0 {
                return map[string]float64{}
        }
        
        // Initialize scores
        scores := make([]float64, n)
        for i := range scores {
                scores[i] = 1.0 / float64(n)
        }
        
        // Iterate to convergence
        for iter := 0; iter < iterations; iter++ {
                newScores := make([]float64, n)
                
                for i := 0; i < n; i++ {
                        newScores[i] = (1.0 - dampingFactor) / float64(n)
                        
                        for j := 0; j < n; j++ {
                                if i != j && graph[j][i] > 0 {
                                        // Calculate outbound link sum
                                        outSum := 0.0
                                        for k := 0; k < n; k++ {
                                                if j != k {
                                                        outSum += graph[j][k]
                                                }
                                        }
                                        
                                        if outSum > 0 {
                                                newScores[i] += dampingFactor * scores[j] * graph[j][i] / outSum
                                        }
                                }
                        }
                }
                
                // Update scores
                scores = newScores
        }
        
        return map[string]float64{}  // Placeholder return value
}

// extractTopSentences extracts the top-ranked sentences
func (s *Summarizer) extractTopSentences(sentences []string, scores map[string]float64) []string {
        // Determine how many sentences to include in the summary
        numSentencesToExtract := calculateNumSentencesToExtract(len(sentences), s.config.MaxSummaryLength)
        
        // Create a list of (sentence, score) pairs
        type scoredSentence struct {
                sentence string
                score    float64
        }
        
        scoredSentences := make([]scoredSentence, 0, len(sentences))
        for _, sentence := range sentences {
                scoredSentences = append(scoredSentences, scoredSentence{
                        sentence: sentence,
                        score:    scores[sentence],
                })
        }
        
        // Sort by score in descending order
        for i := 0; i < len(scoredSentences)-1; i++ {
                for j := 0; j < len(scoredSentences)-i-1; j++ {
                        if scoredSentences[j].score < scoredSentences[j+1].score {
                                scoredSentences[j], scoredSentences[j+1] = scoredSentences[j+1], scoredSentences[j]
                        }
                }
        }
        
        // Take the top sentences
        n := min(numSentencesToExtract, len(scoredSentences))
        topSentences := make([]string, n)
        
        for i := 0; i < n; i++ {
                topSentences[i] = scoredSentences[i].sentence
        }
        
        return topSentences
}

// combineSentencesInOrder combines top sentences in their original order
func (s *Summarizer) combineSentencesInOrder(topSentences []string, allSentences []string) string {
        // Create a set of top sentences for quick lookup
        topSentenceSet := make(map[string]bool)
        for _, sentence := range topSentences {
                topSentenceSet[sentence] = true
        }
        
        // Build summary by including top sentences in original order
        summary := make([]string, 0, len(topSentences))
        
        for _, sentence := range allSentences {
                if topSentenceSet[sentence] {
                        summary = append(summary, sentence)
                }
        }
        
        // Join sentences with proper spacing
        return strings.Join(summary, " ")
}

// extractKeyConcepts extracts key concepts from the text
func (s *Summarizer) extractKeyConcepts(text string) []string {
        // In a real implementation, this would use NLP techniques to extract key concepts
        // For this simulation, we'll use the keyword extraction function
        return s.extractKeywords(text, s.config.KeyphrasesPerDoc)
}

// generateSummaryFromConcepts generates a summary based on key concepts
func (s *Summarizer) generateSummaryFromConcepts(concepts []string, originalText string) string {
        // In a real implementation, this would generate new sentences based on concepts
        // For this simulation, we'll use a more sophisticated approach to weight sentences
        sentences := splitIntoSentences(originalText)
        
        // Score sentences based on multiple factors
        sentenceScores := make(map[string]float64)
        
        for _, sentence := range sentences {
                normalized := s.normalizeSentence(sentence)
                words := strings.Fields(normalized)
                
                // Skip very short sentences
                if len(words) < s.config.MinSentenceLength {
                        sentenceScores[sentence] = 0.0
                        continue
                }
                
                // 1. Concept coverage score - sentences that contain key concepts are important
                conceptScore := 0.0
                conceptMatches := 0
                
                for _, concept := range concepts {
                        if strings.Contains(normalized, concept) {
                                conceptMatches++
                                // Higher weight for exact matches vs partial matches
                                if containsWholeWord(normalized, concept) {
                                        conceptScore += 1.5  // Exact match bonus
                                } else {
                                        conceptScore += 0.7  // Partial match
                                }
                        }
                }
                
                // 2. Density score - sentences with higher concept density are better
                // (shorter sentences that contain concepts are more focused)
                densityScore := 0.0
                if len(words) > 0 && conceptMatches > 0 {
                        densityScore = float64(conceptMatches) / float64(len(words)) * 3.0
                }
                
                // 3. Uniqueness score - sentences that contain rare words are important
                uniquenessScore := 0.0
                for _, word := range words {
                        // Count word occurrences in the entire text to identify rare/unique words
                        occurCount := strings.Count(originalText, word)
                        if occurCount <= 3 && len(word) > 4 { // Significant but rare words
                                uniquenessScore += 0.5
                        }
                }
                
                // 4. Position score - sentences at the beginning/end matter more
                positionScore := s.calculatePositionScore(sentence, originalText) * 2.0
                
                // Combine all factors with weighted importance
                totalScore := (conceptScore * 2.5) + 
                              (densityScore * 1.5) + 
                              (uniquenessScore * 1.0) + 
                              (positionScore * 0.8)
                
                sentenceScores[sentence] = totalScore
        }
        
        // Extract top sentences
        topSentences := s.extractTopSentences(sentences, sentenceScores)
        
        // Combine sentences in their original order
        return s.combineSentencesInOrder(topSentences, sentences)
}

// containsWholeWord checks if text contains the whole word (not just as part of another word)
func containsWholeWord(text, word string) bool {
        // Add word boundaries to find the whole word
        wordWithBoundaries := "\\b" + word + "\\b"
        matched, _ := regexp.MatchString(wordWithBoundaries, text)
        return matched
}

// refineExtractiveWithAbstractive refines extractive summary with abstractive techniques
func (s *Summarizer) refineExtractiveWithAbstractive(extractiveSummary string, concepts []string) string {
        // In a real implementation, this would apply abstractive techniques to refine
        // For this simulation, we'll create a more varied and content-specific intro
        if len(concepts) == 0 {
                return extractiveSummary
        }
        
        // Use different intro templates based on the content and number of concepts
        // This makes summaries more varied across different texts
        var intro string
        
        // Get number of top concepts to include (1-3 depending on how many we have)
        numConcepts := min(3, len(concepts))
        topConcepts := concepts[:numConcepts]
        
        // Generate a hash-like value from the concepts to create variation
        // This ensures different text inputs get different intro styles
        contentHash := 0
        for _, concept := range topConcepts {
                for _, char := range concept {
                        contentHash += int(char)
                }
        }
        
        // Use the hash to select different intro styles
        switch contentHash % 5 {
        case 0:
                intro = fmt.Sprintf("This manuscript centers around %s. ", strings.Join(topConcepts, ", "))
        case 1:
                intro = fmt.Sprintf("The document primarily explores themes of %s. ", strings.Join(topConcepts, ", "))
        case 2:
                if numConcepts == 1 {
                        intro = fmt.Sprintf("The key focus of this text is %s. ", topConcepts[0])
                } else {
                        intro = fmt.Sprintf("Key topics covered include %s. ", strings.Join(topConcepts, ", "))
                }
        case 3:
                intro = fmt.Sprintf("This ancient text contains significant references to %s. ", strings.Join(topConcepts, ", "))
        default:
                intro = fmt.Sprintf("Analysis of this document reveals emphasis on %s. ", strings.Join(topConcepts, ", "))
        }
        
        return intro + extractiveSummary
}

// Helper functions

// splitIntoSentences splits text into sentences
func splitIntoSentences(text string) []string {
        // Simple sentence splitting - in a real implementation, this would be more sophisticated
        text = strings.ReplaceAll(text, "...", "###ELLIPSIS###")
        text = strings.ReplaceAll(text, "Mr.", "Mr###DOT###")
        text = strings.ReplaceAll(text, "Mrs.", "Mrs###DOT###")
        text = strings.ReplaceAll(text, "Dr.", "Dr###DOT###")
        text = strings.ReplaceAll(text, "Prof.", "Prof###DOT###")
        text = strings.ReplaceAll(text, "e.g.", "e.g###DOT###")
        text = strings.ReplaceAll(text, "i.e.", "i.e###DOT###")
        
        // Split on sentence-ending punctuation
        parts := strings.FieldsFunc(text, func(r rune) bool {
                return r == '.' || r == '!' || r == '?'
        })
        
        // Process each part into a proper sentence
        sentences := make([]string, 0, len(parts))
        for _, part := range parts {
                part = strings.TrimSpace(part)
                if part == "" {
                        continue
                }
                
                // Restore special cases
                part = strings.ReplaceAll(part, "###ELLIPSIS###", "...")
                part = strings.ReplaceAll(part, "###DOT###", ".")
                
                // Add back the terminal punctuation (assuming . for simplicity)
                sentences = append(sentences, part+".")
        }
        
        return sentences
}

// buildStopwordsMap builds a map of common stopwords
func buildStopwordsMap() map[string]bool {
        stopwordsList := []string{
                "a", "about", "above", "after", "again", "against", "all", "am", "an", "and", "any", "are", "as", "at",
                "be", "because", "been", "before", "being", "below", "between", "both", "but", "by",
                "could", "did", "do", "does", "doing", "down", "during",
                "each", "few", "for", "from", "further",
                "had", "has", "have", "having", "he", "her", "here", "hers", "herself", "him", "himself", "his", "how",
                "i", "if", "in", "into", "is", "it", "its", "itself",
                "me", "more", "most", "my", "myself",
                "no", "nor", "not", "now",
                "of", "off", "on", "once", "only", "or", "other", "our", "ours", "ourselves", "out", "over", "own",
                "same", "she", "should", "so", "some", "such",
                "than", "that", "the", "their", "theirs", "them", "themselves", "then", "there", "these", "they", "this", "those", "through", "to", "too",
                "under", "until", "up",
                "very",
                "was", "we", "were", "what", "when", "where", "which", "while", "who", "whom", "why", "will", "with", "would",
                "you", "your", "yours", "yourself", "yourselves",
        }
        
        stopwordsMap := make(map[string]bool)
        for _, word := range stopwordsList {
                stopwordsMap[word] = true
        }
        
        return stopwordsMap
}

// isCapitalized checks if a word is capitalized
func isCapitalized(word string) bool {
        if len(word) == 0 {
                return false
        }
        return word[0] >= 'A' && word[0] <= 'Z'
}

// calculateNumSentencesToExtract calculates how many sentences to include in the summary
func calculateNumSentencesToExtract(numSentences, maxSummaryLength int) int {
        // Heuristic: extract about 30% of sentences, or enough to fit maxSummaryLength
        targetPercentage := 0.3
        targetSentences := int(float64(numSentences) * targetPercentage)
        
        // Ensure at least 1 sentence and at most maxSummaryLength/10 (estimating ~10 words per sentence)
        return max(1, min(targetSentences, maxSummaryLength/10))
}

// min returns the minimum of two integers
func min(a, b int) int {
        if a < b {
                return a
        }
        return b
}

// max returns the maximum of two integers
func max(a, b int) int {
        if a > b {
                return a
        }
        return b
}
