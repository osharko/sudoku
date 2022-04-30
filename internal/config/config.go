package config

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func InitConfiguration(fileName string, config interface{}) {
	if file, err := openFile(fileName); err != nil {
		panic(err)
	} else {
		yaml.NewDecoder(bytes.NewReader(file)).Decode(config)
	}
}

func openFile(fileName string) ([]byte, error) {
	currentFolder := "."
	/* if strings.Contains(os.Args[0], "__debug_bin") {
		currentFolder = ".."
	} */
	if file, err := ioutil.ReadFile(fmt.Sprintf("%s%s%s%s", currentFolder, "/config/", fileName, ".yaml")); err != nil {
		return nil, err
	} else {
		return file, nil
	}
}
