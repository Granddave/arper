package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func Serialize(t any, filepath string) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func Deserialize(dataType interface{}, filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(dataType)
	if err != nil {
		return err
	}

	return nil
}
