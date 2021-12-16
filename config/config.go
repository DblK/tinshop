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

type config struct {
	rootShop         string
	debugNfs         bool
	directories      []string
	nfsShares        []string
	shopTitle        string
	shopTemplateData repository.ShopTemplate
}

func LoadConfig() repository.Config {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory
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
		fmt.Println("Config file changed:", e.Name)
		// TODO: Reload config on change
	})
	viper.WatchConfig()

	serverConfig := &config{}

	// ----------------------------------------------------------
	// General config
	// ----------------------------------------------------------
	host := viper.Get("host")
	protocol := viper.Get("protocol")
	port := viper.Get("port")

	var rootShop string
	if protocol == nil {
		rootShop = "http"
	} else {
		rootShop = protocol.(string)
	}
	rootShop += "://"
	if host == nil {
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
		rootShop += host.(string)
	}
	if port == nil {
		rootShop += ":3000"
	} else {
		rootShop += ":" + strconv.Itoa(port.(int))
	}
	serverConfig.rootShop = rootShop

	// ----------------------------------------------------------
	// Debug
	// ----------------------------------------------------------
	serverConfig.debugNfs = viper.GetBool("debug.nfs")

	// ----------------------------------------------------------
	// Sources
	// ----------------------------------------------------------
	// Directories
	cfgDirectories := viper.GetStringSlice("sources.directories")
	if cfgDirectories == nil {
		// Default search
		serverConfig.directories = make([]string, 0)
		serverConfig.directories = append(serverConfig.directories, "./games")
	} else {
		serverConfig.directories = cfgDirectories
	}

	// NFS
	cfgNfs := viper.GetStringSlice("sources.nfs")
	if cfgNfs != nil {
		serverConfig.nfsShares = cfgNfs
	}

	// ----------------------------------------------------------
	// Shop Template
	// ----------------------------------------------------------
	serverConfig.shopTitle = viper.GetString("name")
	if serverConfig.shopTitle == "" {
		serverConfig.shopTitle = "TinShop"
	}
	serverConfig.shopTemplateData = repository.ShopTemplate{
		ShopTitle: serverConfig.shopTitle,
	}

	return serverConfig
}

func (cfg *config) RootShop() string {
	return cfg.rootShop
}
func (cfg *config) DebugNfs() bool {
	return cfg.debugNfs
}
func (cfg *config) Directories() []string {
	return cfg.directories
}
func (cfg *config) NfsShares() []string {
	return cfg.nfsShares
}
func (cfg *config) ShopTemplateData() repository.ShopTemplate {
	return cfg.shopTemplateData
}
func (cfg *config) ShopTitle() string {
	return cfg.shopTitle
}
