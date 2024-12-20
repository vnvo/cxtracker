package cxsimilarity

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	datapath  = "user_behavior_vectors.csv" // Path to the generated data
	simThresh = 0.8                         // Similarity threshold for top matches
)

// LoadData loads the dataset into a slice of user vectors
func LoadData(filepath string) ([][]float64, []string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	headers, err := reader.Read()
	if err != nil {
		return nil, nil, err
	}

	fmt.Printf("File contains %d Service+Metrics\n", len(headers)-1)

	var data [][]float64
	var userIDs []string

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		userIDs = append(userIDs, record[0])
		var vector []float64
		for _, value := range record[1:] {
			if strings.TrimSpace(value) == "-1" {
				vector = append(vector, -1.0)
			} else {
				floatVal, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return nil, nil, err
				}
				vector = append(vector, floatVal)
			}
		}
		data = append(data, vector)
	}

	return data, userIDs, nil
}

// CosineSimilarity computes the cosine similarity between two vectors
func CosineSimilarity(vec1, vec2 []float64) float64 {
	dotProduct := 0.0
	normA := 0.0
	normB := 0.0

	for i := 0; i < len(vec1); i++ {
		if vec1[i] != -1 && vec2[i] != -1 { // Ignore missing data (-1)
			dotProduct += vec1[i] * vec2[i]
			normA += vec1[i] * vec1[i]
			normB += vec2[i] * vec2[i]
		}
	}

	if normA == 0 || normB == 0 {
		return 0 // Avoid division by zero
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

// FindSimilarUsers finds the most similar user and other top matches above a threshold
func FindSimilarUsers(data [][]float64, userIDs []string, targetIndex int, threshold float64) (string, float64, []string) {
	bestSimilarity := -1.0
	bestUserID := ""
	topMatches := []string{}

	targetVector := data[targetIndex]
	for i, vector := range data {
		if i == targetIndex {
			continue // Skip comparing the user to themselves
		}
		similarity := CosineSimilarity(targetVector, vector)
		if similarity > bestSimilarity {
			bestSimilarity = similarity
			bestUserID = userIDs[i]
		}
		if similarity >= threshold {
			topMatches = append(topMatches, fmt.Sprintf("%s (%.4f)", userIDs[i], similarity))
		}
	}

	return bestUserID, bestSimilarity, topMatches
}

func LoadAndCheck() {
	fmt.Println("Loading data...")
	data, userIDs, err := LoadData(datapath)
	if err != nil {
		fmt.Printf("Error loading data: %v\n", err)
		return
	}

	fmt.Printf("Data loaded successfully. Sample Size = %d\n", len(userIDs))

	// Example: Find the most similar user and top matches to the first user
	targetIndex := 0 // Index of the target user
	bestUserID, bestSimilarity, topMatches := FindSimilarUsers(data, userIDs, targetIndex, simThresh)

	fmt.Printf("The most similar user to %s is %s with a similarity score of %.4f\n",
		userIDs[targetIndex], bestUserID, bestSimilarity)

	fmt.Printf("Top matches above threshold(%.2f):\n", simThresh)
	for _, match := range topMatches {
		fmt.Println(match)
	}
}
