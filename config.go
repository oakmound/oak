package plastic

//Load configs

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/dlog"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	plasticPath = "src/bitbucket.org/oakmoundstudio/plasticpiston/plastic"
)

var (
	conf    plasticConfig
	tmpConf plasticConfig
	err     error
)

func LoadConf(fileName string) error {
	wd, _ := os.Getwd()
	dlog.Verb(conf)

	tmpConf, err = loadPlasticConfig(filepath.Join(wd, fileName))
	if err != nil {
		return err
	}
	return err
}

func loadDefaultConf() error {

	dlog.Error(filepath.Join(os.Getenv("GOPATH"), plasticPath, "default.config"))

	conf, err = loadPlasticConfig(filepath.Join(os.Getenv("GOPATH"), plasticPath, "default.config"))
	if err != nil {
		return err
	}

	if tmpConf.Assets.AssetPath != "" {
		conf.Assets.AssetPath = tmpConf.Assets.AssetPath
	}
	if tmpConf.Assets.ImagePath != "" {
		conf.Assets.ImagePath = tmpConf.Assets.ImagePath
	}
	if tmpConf.Assets.AudioPath != "" {
		conf.Assets.AudioPath = tmpConf.Assets.AudioPath
	}

	if tmpConf.Debug.Filter != "" {
		conf.Debug.Filter = tmpConf.Debug.Filter
	}
	if tmpConf.Debug.Level != "" {
		conf.Debug.Level = tmpConf.Debug.Level
	}

	if tmpConf.Screen.Width != 0 {
		conf.Screen.Width = tmpConf.Screen.Width
	}
	if tmpConf.Screen.Height != 0 {
		conf.Screen.Height = tmpConf.Screen.Height
	}

	dlog.Verb(conf)

	return err
}

func loadPlasticConfig(fileName string) (plasticConfig, error) {

	confFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return plasticConfig{}, err
	}
	var config plasticConfig
	json.Unmarshal(confFile, &config)

	return config, nil
}

type plasticConfig struct {
	Assets struct {
		AssetPath string `json:"assetPath"`
		AudioPath string `json:"audioPath"`
		ImagePath string `json:"imagePath"`
	} `json:"assets"`
	Debug struct {
		Filter string `json:"filter"`
		Level  string `json:"level"`
	} `json:"debug"`
	Screen struct {
		Height int `json:"height"`
		Width  int `json:"width"`
	} `json:"screen"`
}
