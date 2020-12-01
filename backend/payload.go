package main

import "encoding/json"

type Payload map[string]interface{}

func (p Payload) ToJson() string {
	data, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func (p *Payload) LoadJSON(data []byte) error {
	err := json.Unmarshal(data, p)
	return err
}

func (p Payload) Get(key string) string {
	value, found := p[key]
	if !found {
		return ""
	}
	str, ok := value.(string)
	if !ok {
		return ""
	}
	return str
}