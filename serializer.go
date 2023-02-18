package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func Serialize(t any, filename string) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func Deserialize(dataType interface{}, filename string) error {
	file, err := os.Open(filename)
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
