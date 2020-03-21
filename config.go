package telebotex

import (
	"io/ioutil"
	"os"

	"github.com/json-iterator/go"
)

func newConfig(configFile string) (map[string]jsoniter.RawMessage, error) {
	f, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	config := make(map[string]jsoniter.RawMessage)
	err = jsoniter.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func UnmarshalFromConfig(config map[string]jsoniter.RawMessage, key string, v interface{}) error {
	rawMessage, ok := config[key]
	if !ok {
		return configLoadErr
	}
	if err := jsoniter.Unmarshal(rawMessage, v); err != nil {
		return configLoadErr
	}
	return nil
}
