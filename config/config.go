package config

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/dblk/tinshop/repository"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type debug struct {
	Nfs        bool
	NoSecurity bool
}

type config struct {
	rootShop         string
	ShopHost         string             `mapstructure:"host"`
	ShopProtocol     string             `mapstructure:"protocol"`
	ShopPort         int                `mapstructure:"port"`
	Debug            debug              `mapstructure:"debug"`
	AllSources       repository.Sources `mapstructure:"sources"`
	Name             string             `mapstructure:"name"`
	shopTemplateData repository.ShopTemplate
}

var serverConfig config

// LoadConfig handles viper under the hood
func LoadConfig() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	viper.SetDefault("sources.directories", "./games")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("Config not found!")
		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed, update new configuration...")
		serverConfig = loadAndCompute()
	})
	viper.WatchConfig()

	serverConfig = loadAndCompute()
}

// GetConfig returns the current configuration
func GetConfig() repository.Config {
	return &serverConfig
}

func loadAndCompute() config {
	err := viper.Unmarshal(&serverConfig)

	if err != nil {
		log.Fatalln(err)
	}
	computeDefaultValues(&serverConfig)

	return serverConfig
}

func computeDefaultValues(config repository.Config) repository.Config {
	// ----------------------------------------------------------
	// Compute rootShop url
	// ----------------------------------------------------------
	var rootShop string
	if config.Protocol() == "" {
		rootShop = "http"
	} else {
		rootShop = config.Protocol()
	}
	rootShop += "://"
	if config.Host() == "" {
		// Retrieve current IP
		host, _ := os.Hostname()
		addrs, _ := net.LookupIP(host)
		var myIP = ""
		for _, addr := range addrs {
			if ipv4 := addr.To4(); ipv4 != nil {
				if myIP == "" {
					myIP = ipv4.String()
				}
			}
		}
		rootShop += myIP
	} else {
		rootShop += config.Host()
	}
	if config.Port() == 0 {
		rootShop += ":3000"
	} else {
		rootShop += ":" + strconv.Itoa(config.Port())
	}
	config.SetRootShop(rootShop)

	config.SetShopTemplateData(repository.ShopTemplate{
		ShopTitle: config.ShopTitle(),
	})

	return config
}

func (cfg *config) SetRootShop(root string) {
	cfg.rootShop = root
}

func (cfg *config) RootShop() string {
	return cfg.rootShop
}
func (cfg *config) Protocol() string {
	return cfg.ShopProtocol
}
func (cfg *config) Host() string {
	return cfg.ShopHost
}
func (cfg *config) Port() int {
	return cfg.ShopPort
}
func (cfg *config) DebugNfs() bool {
	return cfg.Debug.Nfs
}
func (cfg *config) DebugNoSecurity() bool {
	return cfg.Debug.NoSecurity
}
func (cfg *config) Directories() []string {
	return cfg.AllSources.Directories
}
func (cfg *config) NfsShares() []string {
	return cfg.AllSources.Nfs
}
func (cfg *config) Sources() repository.Sources {
	return cfg.AllSources
}
func (cfg *config) ShopTemplateData() repository.ShopTemplate {
	return cfg.shopTemplateData
}
func (cfg *config) SetShopTemplateData(data repository.ShopTemplate) {
	cfg.shopTemplateData = data
}
func (cfg *config) ShopTitle() string {
	return cfg.Name
}
