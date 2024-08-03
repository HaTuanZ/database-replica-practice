package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
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
		panic(err)
	}

	slaveReplicas, err := strconv.Atoi(os.Getenv("SLAVE_REPLICA"))
	if err != nil {
		panic(err)
	}

	compose := ComposeFile{
		Version:  "3.4",
		Services: make(map[string]Service),
	}

	for i := 1; i <= masterReplica; i++ {
		serviceName := fmt.Sprintf("db-master-%d", i)
		configFileName := createConfigFile(fmt.Sprintf("my-master-%d.conf", i), i)

		compose.Services[serviceName] = Service{
			Image: "mysql:8.4",
			Environment: map[string]string{
				"MYSQL_ROOT_PASSWORD": generatePassword(),
				"MYSQL_DATABASE":      "example_db",
			},
			Ports: []string{fmt.Sprintf("%d:3306", 3306+1)},
			Volumes: []string{
				fmt.Sprintf("./mysql-data-master-%d:/var/lib/mysql", i),
				fmt.Sprintf("./%s:/etc/mysql/conf.d/my.cnf", configFileName),
			},
		}
	}

	for i := 1; i <= slaveReplicas; i++ {
		serviceName := fmt.Sprintf("db-slave-%d", i)
		configFileName := createConfigFile(fmt.Sprintf("my-slave-%d.conf", i), i)

		compose.Services[serviceName] = Service{
			Image: "mysql:8.4",
			Environment: map[string]string{
				"MYSQL_ROOT_PASSWORD": generatePassword(),
				"MYSQL_DATABASE":      "example_db",
			},
			Ports: []string{fmt.Sprintf("%d:3306", 3306+i+masterReplica)},
			Volumes: []string{
				fmt.Sprintf("./mysql-data-slave-%d:/var/lib/mysql", 1),
				fmt.Sprintf("./%s:/etc/mysql/conf.d/my.cnf", configFileName),
			},
		}
	}

	data, err := yaml.Marshal(&compose)

	if err != nil {
		panic(err)
	}

	err = os.WriteFile("docker-compose.yml", data, 0644)

	if err != nil {
		panic(err)
	}

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

func createConfigFile(fileName string, serverId int) string {
	config := fmt.Sprintf(`[mysqld]
server-id=%d
log_bin=mysql-bin
binlog_format=ROW
	`, serverId)

	filePath := "conf/" + fileName

	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	err := os.WriteFile(filePath, []byte(config), 0644)

	if err != nil {
		panic(err)
	}

	return filePath
}
