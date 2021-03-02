package ui

import (
	"io/ioutil"

	e "github.com/muzudho/kifuwarabe-gtp/entities"
	"github.com/pelletier/go-toml"
)

// LoadColiseumConf - ゲーム設定ファイルを読み込みます。
func LoadColiseumConf(path string) (*e.EngineConf, error) {

	// ファイル読込
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Toml解析
	binary := []byte(string(fileData))
	config := &e.EngineConf{}
	toml.Unmarshal(binary, config)

	return config, nil
}
