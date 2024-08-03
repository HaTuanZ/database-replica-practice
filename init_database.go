package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type ComposeFile struct {
	Version  string             `yaml:"version"`
	Services map[string]Service `yaml:"services"`
}

type Service struct {
	Image       string            `yaml:"image"`
	Environment map[string]string `yaml:"environment"`
	Ports       []string          `yaml:"ports"`
	Volumes     []string          `yaml:"volumes"`
}

func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Print("Error loading environment")
	}

	masterReplica, err := strconv.Atoi(os.Getenv("MASTER_REPLICA"))
	if err != nil {
		fmt.Print("Invalid master replica value")
		return
	}

	slaveReplicas, err := strconv.Atoi(os.Getenv("SLAVE_REPLICA"))
	if err != nil {
		fmt.Println("Invalid SLAVE_REPLICA value")
		return
	}

	compose := ComposeFile{
		Version:  "3.4",
		Services: make(map[string]Service),
	}

	for i := 1; i <= masterReplica; i++ {
		serviceName := fmt.Sprintf("db-master-%d", i)
		compose.Services[serviceName] = Service{
			Image: "mysql:8.4",
			Environment: map[string]string{
				"MYSQL_ROOT_PASSWORD": generatePassword(),
				"MYSQL_DATABASE":      "example_db",
			},
			Ports:   []string{fmt.Sprintf("%d:3306", 3306+1)},
			Volumes: []string{fmt.Sprintf("./mysql-data-master-%d:/var/lib/mysql", i)},
		}
	}

	for i := 1; i <= slaveReplicas; i++ {
		serviceName := fmt.Sprintf("db-slave-%d", i)
		compose.Services[serviceName] = Service{
			Image: "mysql:8.4",
			Environment: map[string]string{
				"MYSQL_ROOT_PASSWORD": generatePassword(),
				"MYSQL_DATABASE":      "example_db",
			},
			Ports:   []string{fmt.Sprintf("%d:3306", 3306+i+masterReplica)},
			Volumes: []string{fmt.Sprintf("./mysql-data-slave-%d:/var/lib/mysql", 1)},
		}
	}

	data, err := yaml.Marshal(&compose)

	if err != nil {
		panic((err))
	}

	os.WriteFile("docker-compose.yml", data, 0644)

	fmt.Printf("docker-compose.yml generated successfully")
}

func generatePassword() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}|;:,.<>?/"
	passwordLength := 12
	password := make([]byte, passwordLength)
	for i := range password {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			fmt.Println("Error generating password:", err)
		}
		password[i] = charset[num.Int64()]
	}

	return string(password)
}
