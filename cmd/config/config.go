package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

var (
	logo = `    ___       ___       ___       ___       ___       ___       ___       ___       ___       ___       ___       ___       ___       ___   
   /\  \     /\  \     /\  \     /\__\     /\  \     /\  \     /\  \     /\  \     /\  \     /\__\     /\  \     /\__\     /\  \     /\  \  
  /::\  \   /::\  \   /::\  \   /:/  /    /::\  \   /::\  \   /::\  \   /::\  \   /::\  \   /:/  /    /::\  \   /:| _|_   /::\  \   /::\  \ 
 /:/\:\__\ /::\:\__\ /::\:\__\ /:/__/    /:/\:\__\ /::\:\__\ /:/\:\__\ /::\:\__\ /::\:\__\ /:/__/    /::\:\__\ /::|/\__\ /:/\:\__\ /::\:\__\
 \:\:\/__/ \:\:\/  / \:\:\/  / \:\  \    \:\/:/  / \/\::/  / \:\/:/  / \:\::/  / \/\::/  / \:\  \    \/\::/  / \/|::/  / \:\ \/__/ \:\:\/  /
  \::/  /   \:\/  /   \:\/  /   \:\__\    \::/  /    /:/  /   \::/  /   \::/  /    /:/  /   \:\__\     /:/  /    |:/  /   \:\__\    \:\/  / 
   \/__/     \/__/     \/__/     \/__/     \/__/     \/__/     \/__/     \/__/     \/__/     \/__/     \/__/     \/__/     \/__/     \/__/  
`
)

type Config struct {
	ListenSchema           string       `yaml:"listen_schema"`
	ListenAddr             string       `yaml:"listen_addr"`
	ListenPort             int          `yaml:"listen_port"`
	SSLCertificate         string       `yaml:"ssl_certificate"`
	SSLCertificateKey      string       `yaml:"ssl_certificate_key"`
	TCPHealthCheck         bool         `yaml:"tcp_health_check"`
	TCPHealthCheckInterval uint         `yaml:"tcp_health_check_interval"`
	MaxAllowed             uint         `yaml:"max_allowed"`
	Locations              []*Locations `yaml:"locations"`
}

type Locations struct {
	Prefix      string   `yaml:"prefix"`
	BalanceMode string   `yaml:"balance_mode"`
	Servers     []string `yaml:"servers"`
}

// NewReadConfig 用于读取配置文件
func NewReadConfig(file string) (*Config, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// Print 用于打印配置信息
func (c *Config) Print() {
	fmt.Printf("%s\nSchema: %s\nPort: %d\nInterval: %d\nHealth Check: %v\nLocation:\n", logo, c.ListenSchema, c.ListenPort, c.TCPHealthCheckInterval, c.TCPHealthCheck)
	for _, l := range c.Locations {
		fmt.Printf("\tRoute: %s\n\tServers: %s\n\tMode: %s\n\n", l.Prefix, l.Servers, l.BalanceMode)
	}
}

// Validation 用于验证配置信息
func (c *Config) Validation() error {
	if c.ListenSchema != "http" && c.ListenSchema != "https" {
		return fmt.Errorf("invalid schema")
	}

	if len(c.Locations) == 0 {
		return fmt.Errorf("no location found")
	}

	if c.ListenSchema == "https" && (len(c.SSLCertificate) == 0 || len(c.SSLCertificateKey) == 0) {
		return errors.New("the https proxy requires ssl_certificate_key and ssl_certificate")
	}

	if c.TCPHealthCheckInterval < 1 {
		return fmt.Errorf("invalid health check interval")
	}

	return nil
}
