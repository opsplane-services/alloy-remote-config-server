package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func Start() {
	godotenv.Load()
	httpPort := 8080
	httpPortEnv := os.Getenv("HTTP_PORT")
	if len(httpPortEnv) > 0 {
		httpPort, _ = strconv.Atoi(httpPortEnv)
	}
	grpcPort := 8888
	grpcPortEnv := os.Getenv("GRPC_PORT")
	if len(grpcPortEnv) > 0 {
		grpcPort, _ = strconv.Atoi(grpcPortEnv)
	}
	configFolder := "conf"
	configFolderEnv := os.Getenv("CONFIG_FOLDER")
	if len(configFolderEnv) > 0 {
		configFolder = configFolderEnv
	}
	err := LoadTemplates(configFolder)
	if err != nil {
		log.Println(fmt.Sprintf("Error loading templates: %v", err))
	}
	globalStorage, err = InitStorage()
	if err != nil {
		log.Println(fmt.Sprintf("Error loading storage: %v", err))
	}
	go StartConnectGrpcServer(grpcPort)
	StartRestServer(httpPort)
}
