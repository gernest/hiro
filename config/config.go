package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

// ErrNotSet is an error returned when there is no configuration file found.
var ErrNotSet = errors.New("configuration file not set")

// Config configuration settings.
type Config struct {
	DBDriver  string `yaml:"driver,omitempty"`
	DBConn    string `yaml:"db_conn,omitempty"`
	Port      int    `yaml:"port,omitempty"`
	Secret    string `yaml:"secret,omitempty"`
	Host      string `yaml:"host,omitempty"`
	ImageHost string `yaml:"image_host,omitempty"`
}

// Minio minio client settings.
type Minio struct {
	Endpoint     string `yaml:"endpoint"`
	AccessKey    string `yaml:"access_key"`
	AccessSecret string `yaml:"access_secret"`
}

// NSQ nsq client options.
type NSQ struct {
	LookupD string `yaml:"lookupd"`
	NSQD    string `yaml:"nsqd"`
}

func load(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	err = yaml.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// FromCtx returns configuration with setting from cli context.
func FromCtx(ctx *cli.Context) (*Config, error) {
	cfg := ctx.String("config")
	if cfg != "" {
		return load(cfg)
	}
	if cfg != "" {
		return load(cfg)
	}
	h, err := GetHomeDir()
	if err != nil {
		return nil, err
	}
	cfgFile := filepath.Join(h, "config.yaml")
	_, err = os.Stat(cfgFile)
	if os.IsNotExist(err) {
		return fromCli(ctx), nil
	}
	fc, err := load(cfgFile)
	if err != nil {
		return nil, err
	}
	fc.Update(fromCli(ctx))
	return fc, nil
}

func fromCli(ctx *cli.Context) *Config {
	cfg := &Config{
		DBDriver:  ctx.String("driver"),
		DBConn:    ctx.String("db-conn"),
		Port:      ctx.Int("port"),
		Secret:    ctx.String("secret"),
		Host:      ctx.String("host"),
		ImageHost: ctx.String("image-host"),
	}
	return cfg
}

// GetHomeDir returns bq home directory. The bq home is $HOME/.config/bq
func GetHomeDir() (string, error) {
	h, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	h = filepath.Join(h, ".config", "bq")
	_, err = os.Stat(h)
	if os.IsNotExist(err) {
		err = os.MkdirAll(h, 0777)
		if err != nil {
			return "", err
		}
	}
	return h, nil
}

// Update updates c with n.
func (c *Config) Update(n *Config) {
	if n.DBDriver != "" {
		c.DBDriver = n.DBDriver
	}
	if n.DBConn != "" {
		c.DBConn = n.DBConn
	}
	if n.Port != 0 {
		c.Port = n.Port
	}
	if n.Secret != "" {
		c.Secret = n.Secret
	}
	if n.Host != "" {
		c.Host = n.Host
	}
	if n.ImageHost != "" {
		c.ImageHost = n.ImageHost
	}
}
