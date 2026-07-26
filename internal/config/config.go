package config

import (
	"encoding/json"
	"os"
)
const configFileName = "/.gatorconfig.json"
type Config struct {
	Db_url string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func (c Config) SetUser(user string) {
	c.Current_user_name = user
	write(c)
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	location := home + configFileName
	return location, nil
}

func Read() (Config, error) {
	var json_output Config
	
	path, err := getConfigFilePath()
	if err != nil {
	
		return json_output, err
	}
	
	json_input, err := os.ReadFile(path)
	if err != nil {
		
		return json_output, err
	}
	
	if err := json.Unmarshal(json_input, &json_output); err != nil {
		return json_output, err
	}
	
	return json_output, nil
}

func write(cfg Config) error {
	location, err := getConfigFilePath()
	if err != nil {
		return err
	}
	
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	
	err = os.WriteFile(location, jsonData, 666)
	if err != nil {
		return err
	}
	return nil
}