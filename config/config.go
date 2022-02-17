// SPDX-License-Identifier: Apache-2.0

package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/zpatrick/go-config"
)

type Config struct {
	MinIOEndpoint string
	MinIOBucket   string
}

var GlobalConfig *Config

func Init() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalln("Error loading config: " + err.Error())
	}
	configPath := filepath.Join(configDir, "sogno-file-service/config.json")
	fmt.Println("Loading config from: " + configPath)
	configFile := config.NewJSONFile(configPath)
	c := config.NewConfig([]config.Provider{configFile})

	minioEndpoint, err := c.String("minio_endpoint")
	if err != nil {
		log.Fatalln("Error loading config: " + err.Error())
	}
	minioBucket, err := c.String("minio_bucket")
	if err != nil {
		log.Fatalln("Error loading config: " + err.Error())
	}
	GlobalConfig = &Config{
		MinIOEndpoint: minioEndpoint,
		MinIOBucket:   minioBucket,
	}
}
