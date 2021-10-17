package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type ParamsLocal struct {
	Port  string `yaml:"port"`
	Redis struct {
		Adr      string `yaml:"adr"`
		Password string `yaml:"password"`
		Db       int    `yaml:"db"`
	} `yaml:"redis"`
	GrpcPort string `yaml:"grpc_port"`
}

func NewConfig(url string) ParamsLocal {
	var c ParamsLocal
	dat, err := ioutil.ReadFile(url)
	if err != nil {
		log.Fatal(err)
	}
	er := yaml.Unmarshal(dat, &c)
	if er != nil {
		log.Fatalf("error: %v", er)
	}
	return c
}
