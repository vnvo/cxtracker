package cxgenerator

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

// Constants defining data generation
const (
	nUsers    = 1000 // Number of unique users
	nServices = 100  // Number of microservices
	dimension = 200  // Number of features per user vector
	datapath  = "user_behavior_vectors.csv"
)

type ServiceType string

const (
	APIServer       ServiceType = "API Server"
	APIWithDatabase ServiceType = "API with Database"
	KafkaProducer   ServiceType = "Kafka Producer"
	KafkaConsumer   ServiceType = "Kafka Consumer"
)

// Distribution of service types (higher portion for API servers and API with DBs)
var serviceDistribution = []struct {
	Type    ServiceType
	Portion float64
}{
	{APIServer, 0.4},
	{APIWithDatabase, 0.4},
	{KafkaProducer, 0.1},
	{KafkaConsumer, 0.1},
}

// Metrics for each service type
var metricTemplates = map[ServiceType][]string{
	APIServer:       {"http_resp", "http_rate", "http_err"},
	APIWithDatabase: {"db_lat", "db_rate", "db_err"},
	KafkaProducer:   {"msg_prod", "pub_lat", "pub_fail"},
	KafkaConsumer:   {"msg_cons", "proc_lat", "cons_err"},
}

// Generate random vector with service-specific patterns and realistic values
func generateServiceMetrics(serviceType ServiceType) []float64 {
	metrics := metricTemplates[serviceType]
	vector := make([]float64, len(metrics))
	for i := 0; i < len(metrics); i++ {
		switch serviceType {
		case APIServer:
			if metrics[i] == "http_resp" {
				vector[i] = rand.Float64()*100 + 50
			} else if metrics[i] == "http_rate" {
				vector[i] = rand.Float64()*50 + 10
			} else if metrics[i] == "http_err" {
				vector[i] = rand.Float64() * 5
			}
		case APIWithDatabase:
			if metrics[i] == "db_lat" {
				vector[i] = rand.Float64()*150 + 100
			} else if metrics[i] == "db_rate" {
				vector[i] = rand.Float64()*30 + 5
			} else if metrics[i] == "db_err" {
				vector[i] = rand.Float64() * 10
			}
		case KafkaProducer:
			if metrics[i] == "msg_prod" {
				vector[i] = rand.Float64()*1000 + 500
			} else if metrics[i] == "pub_lat" {
				vector[i] = rand.Float64()*100 + 50
			} else if metrics[i] == "pub_fail" {
				vector[i] = rand.Float64() * 1
			}
		case KafkaConsumer:
			if metrics[i] == "msg_cons" {
				vector[i] = rand.Float64()*1200 + 600
			} else if metrics[i] == "proc_lat" {
				vector[i] = rand.Float64()*80 + 20
			} else if metrics[i] == "cons_err" {
				vector[i] = rand.Float64() * 2
			}
		}
	}
	return vector
}

// Randomly pick a service type based on the defined distribution
func pickServiceType() ServiceType {
	r := rand.Float64()
	sum := 0.0
	for _, entry := range serviceDistribution {
		sum += entry.Portion
		if r < sum {
			return entry.Type
		}
	}
	return APIServer // Default fallback
}

// Generate user metrics for all services
func generateUserMetrics() ([]float64, []string) {
	userMetrics := []float64{}
	headers := []string{}
	for i := 0; i < nServices; i++ {
		serviceType := pickServiceType()
		serviceMetrics := generateServiceMetrics(serviceType)
		metrics := metricTemplates[serviceType]
		for j, value := range serviceMetrics {
			userMetrics = append(userMetrics, value)
			headers = append(headers, fmt.Sprintf("s%d_%s", i+1, metrics[j]))
		}
	}
	return userMetrics, headers
}

// Generate and save data to a file
func saveDataToFile() error {
	file, err := os.Create(datapath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Generate a single user's metrics to establish headers
	_, headers := generateUserMetrics()

	// Write header
	file.WriteString("user_id")
	for _, header := range headers {
		file.WriteString("," + header)
	}
	file.WriteString("\n")

	// Generate user data and write to file
	for userID := 1; userID <= nUsers; userID++ {
		userMetrics, _ := generateUserMetrics()
		line := fmt.Sprintf("user_%d", userID)
		for _, value := range userMetrics {
			line += fmt.Sprintf(",%.4f", value)
		}
		file.WriteString(line + "\n")
	}

	return nil
}

func Generate() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	fmt.Println("Generating user behavior vectors with service metrics...")
	err := saveDataToFile()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Data generation complete. File saved at %s\n", datapath)
}
