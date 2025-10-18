package config

import (
	"encoding/json"
	"fmt"
	"github.com/Gitrupesh20/real-time-notification-system/pkg"
	"log"
	"time"
)

type Config struct {
	Port             string        `json:"port"`
	Mode             int           `json:"mode"` // 1p, 0d
	HandShakeTimeout time.Duration `json:"hand_shake_timeout_ms"`
	ReadBufferSize   int           `json:"read_buffer_size"`
	WriteBufferSize  int           `json:"write_buffer_size"`
	AllowedOrigins   []string      `json:"allowed_origins"`
	AllowNoOrigin    bool          `json:"allow_no_origin"`
	PingInterval     time.Duration `json:"ping_interval_ms"`
	AllowWsWithSSL   bool          `json:"allow_ws_w_ssl"`
	MqAddr           string        `json:"mq_addr"`
	MqQueueName      string        `json:"mq_queue_name"`
	NoOfWorker       int           `json:"no_of_worker"`
}

const (
	configPath = "./config"
	configName = "config.json"
)

func LoadConfig() Config {
	var config Config

	if err := pkg.LoadFile(configPath, configName, &config, json.Unmarshal); err != nil {
		log.Fatal("error while loading config error: ", err)
	} else if err = config.validateConfig(); err != nil {
		log.Fatal("error while validating config: ", err)
	}
	config.NoOfWorker = 10
	return config
}

// validateConfig ensure that config should contain necessary filed with value
func (c *Config) validateConfig() error {
	if c.Port == "" {
		log.Println("error: port is empty, setting default port")
		c.Port = "8080"
	} else if c.Mode > 2 {
		return fmt.Errorf("error: mode is unkown")
	} else if c.HandShakeTimeout == 0 {
		return fmt.Errorf("error: hand_shake_timeout is unkown")
	} else if c.ReadBufferSize == 0 {
		return fmt.Errorf("error: read_buffer_size is unkown")
	} else if c.WriteBufferSize == 0 {
		return fmt.Errorf("error: write_buffer_size is unkown")
	} else if c.AllowedOrigins == nil && len(c.AllowedOrigins) == 0 {
		return fmt.Errorf("error: allowed_origins is unkown")
	} else if c.PingInterval == 0 {
		return fmt.Errorf("error: ping_interval is unkown")
	}

	return nil
}
func (c *Config) IsAllowOrigin(origin string) bool {
	for _, o := range c.AllowedOrigins {
		if o == origin {
			return true
		}
	}
	return false
}
