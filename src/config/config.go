package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var (
	// C 全局配置文件，在Init调用前为nil
	C *Config
)

// Config 配置
type Config struct {
	Ws         ws         `yaml:"ws"`
	Postgresql postgresql `yaml:"postgresql"`
	Redis      redis      `yaml:"redis"`
	Jwt        jwt        `yaml:"jwt"`
	Log        log        `yaml:"log"`
	Mail       mail       `yaml:"mail"`
	Rabbitmq   rabbitmq   `yaml:"rabbitmq"`
}

type ws struct {
	Addr string `yaml:"addr"`
}

type postgresql struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}

type redis struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type jwt struct {
	Secret string `yaml:"secret"`
}

type log struct {
	Path string `yaml:"path"`
	File string `yaml:"file"`
}

type mail struct {
	Host      string `yaml:"host"`
	Addr      string `yaml:"addr"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Goroutine int    `yaml:"goroutine"`
	Template  string `yaml:"template"`
}

type rabbitmq struct {
	Host            string   `yaml:"host"`
	Port            string   `yaml:"port"`
	User            string   `yaml:"user"`
	Password        string   `yaml:"password"`
	ExchangeName    string   `yaml:"exchange_name"`
	WsRoutingKey    []string `yaml:"ws_routing_key"`
	EmailRoutingKey []string `yaml:"email_routing_key"`
}

func init() {
	configFile := "default.yml"

	if v, ok := os.LookupEnv("ENV"); ok {
		configFile = v + ".yml"
	}

	data, err := ioutil.ReadFile(fmt.Sprintf("./env/config/%s", configFile))

	if err != nil {
		panic(err)
		return
	}

	config := &Config{}

	err = yaml.Unmarshal(data, config)

	if err != nil {
		fmt.Println("Unmarshal config error!")
		panic(err)
		return
	}

	C = config
	fmt.Println("------- " + configFile + " loaded" + " -------")
}
