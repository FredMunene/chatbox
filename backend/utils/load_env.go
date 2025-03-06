package utils

import (
	"bufio"
	"os"
	"strings"
)


var (
	GoogleClientID, GoogleClientSecret, GithubClientID, GithubClientSecret string
)

func LoadEnvVariables(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan(){
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		parts := strings.SplitN(line,"=",2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		os.Setenv(key,value)
	}

	GithubClientID = os.Getenv("GithubClientID")
	GithubClientSecret = os.Getenv("GithubClientID")
	GoogleClientID = os.Getenv("GoogleClientID")
	GoogleClientSecret = os.Getenv("GoogleClientSecret")


	return scanner.Err()
}
