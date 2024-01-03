package config

import (
	"flag"
	"fmt"
	"github.com/eidentitet/rest-go/cache"
	"github.com/kkyr/fig"
	cache2 "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

const (
	internalCacheKey string = "internalConfigCacheKey"
)

func GetFlags() string {
	var (
		memory              = cache.Memory()
		configFileName      *string
		defaultFileName     = "config"
		workingDirectory, _ = os.Getwd()
		configurationPath   = fmt.Sprintf("%s", workingDirectory)
	)

	fileName := fmt.Sprintf("%s/config/%s.yaml", configurationPath, defaultFileName)
	configFileName = flag.String("config-file", fileName, "Path of the configuration file")

	flag.Parse()
	memory.Set("flagsCache", *configFileName, cache2.NoExpiration)

	return *configFileName
}

func GetServiceConfiguration(config interface{}, configFilePath string) {
	var (
		err error
	)

	err = fig.Load(config,
		fig.UseEnv("app"),
		fig.File(filepath.Base(configFilePath)),
		fig.Dirs(filepath.Dir(configFilePath), "../../config"),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func GetAppConfig() *AppConfiguration {
	var (
		config              AppConfiguration
		memory              = cache.Memory()
		isFound             bool
		cachedConfiguration interface{}
	)

	// Check if the configuration is cached
	cachedConfiguration, isFound = memory.Get(internalCacheKey)
	if isFound {
		return cachedConfiguration.(*AppConfiguration)
	}

	configFilePath, isFound := memory.Get("flagsCache")
	if !isFound {
		log.Fatal("Flags cache not found")
	}
	// Get configuration
	GetServiceConfiguration(&config, configFilePath.(string))

	memory.Set(internalCacheKey, &config, cache2.NoExpiration)

	return &config
}
