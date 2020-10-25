package helpers

import (
	"github.com/gen2brain/beeep"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Config struct {
	Discord struct {
		Client    string `yaml:"client"`
		ImageName string `yaml:"imagename"`
	} `yaml:"discord"`
	LastFM struct {
		ApiKey    string `yaml:"api_key"`
		ApiKeySec string `yaml:"api_secret_key"`
		Username  string `yaml:"username"`
	} `yaml:"lastfm"`
	ServiceSettings struct {
		DebugMode bool `yaml:"debug_mode"`
		BlockNoti bool `yaml:"block_notifications"`
		NoTime    bool `yaml:"no_time_display"`
	} `yaml:"serviceSettings"`
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func MakeEmptyFile() {
	_, _ = os.Create("config.yaml")
	//config := &Config{}
	file, err := os.Open("config.yaml")
	if err != nil {
		panic(err)
	}

	emptyConf := Config{
		Discord: struct {
			Client    string `yaml:"client"`
			ImageName string `yaml:"imagename"`
		}{},
		LastFM: struct {
			ApiKey    string `yaml:"api_key"`
			ApiKeySec string `yaml:"api_secret_key"`
			Username  string `yaml:"username"`
		}{},
		ServiceSettings: struct {
			DebugMode bool `yaml:"debug_mode"`
			BlockNoti bool `yaml:"block_notifications"`
			NoTime    bool `yaml:"no_time_display"`
		}{},
	}

	d, err2 := yaml.Marshal(&emptyConf)

	if err2 != nil {
		panic(err2)
	}

	_ = ioutil.WriteFile("config.yaml", d, 0600)

	println("An Empty configuration file has been created - please insert your data in there and relaunch!")
	_ = beeep.Alert("Edit config!", "An Empty configuration file has been created - please insert your data in there and relaunch!", "")
	file.Close()
	os.Exit(0)
}

func LoadSettings() Config {
	f, err := os.Open("config.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}
