package config

import (
	"github.com/BurntSushi/toml"
)

var (
	global *Config
)

// LoadGlobal 加载全局配置
func LoadGlobal(fpath string) error {
	c, err := Parse(fpath)
	if err != nil {
		return err
	}
	global = c
	return nil
}

// Global 获取全局配置
func Global() *Config {
	if global == nil {
		return &Config{}
	}
	return global
}

// Parse 解析配置文件
func Parse(fpath string) (*Config, error) {
	var c Config
	_, err := toml.DecodeFile(fpath, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// Config 配置参数
type Config struct {
	RunMode        string         `toml:"run_mode"`
	Log            Log            `toml:"log"`
	Email          Email          `toml:"email"`
	ServerChan     ServerChan     `toml:"server_chan"`
	GuilinlifeConf GuilinlifeConf `toml:"guilinlife"`
}

// IsDebugMode 是否是debug模式
func (c *Config) IsDebugMode() bool {
	return c.RunMode == "debug"
}

// Log 日志配置参数
type Log struct {
	Level      int    `toml:"level"`
	Format     string `toml:"format"`
	Output     string `toml:"output"`
	OutputFile string `toml:"output_file"`
}

// Email email配置参数
type Email struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	UserName string `toml:"username"`
	Password string `toml:"password"`
	Receiver string `toml:"receiver"`
}

// ServerChan server酱配置参数
type ServerChan struct {
	Url string `toml:"url"`
}

// GuilinlifeConf 桂林人论坛配置参数
type GuilinlifeConf struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
}
