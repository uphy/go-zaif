package example

import (
	"io/ioutil"
	"os"

	"github.com/uphy/go-zaif/zaif"
	"gopkg.in/yaml.v2"
)

type (
	Config struct {
		Key    string `yaml:"key"`
		Secret string `yaml:"secret"`
	}
)

func newPrivateAPI() *zaif.PrivateAPI {
	var config Config
	file, err := os.Open("config.yml")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		panic(err)
	}
	if len(config.Key) == 0 || len(config.Secret) == 0 {
		panic("no api key found")
	}
	return zaif.NewPrivateApi(config.Key, config.Secret)
}
