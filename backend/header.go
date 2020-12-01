package main

import "encoding/json"

type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

func (p Header) ToJson() string {
	data, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return string(data)
}