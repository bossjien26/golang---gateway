package main

import (
	"encoding/json"

	"github.com/Kong/go-pdk"
)

type Config struct {
	Apikey string
}

func New() interface{} {
	return &Config{}
}

func (conf Config) Access(kong *pdk.PDK) {
	// key, err := kong.Request.GetQueryArg("key")
	key, err := kong.Request.GetHeader("key")
	// println(kong.Request.GetHeader("key"))

	apiKey := conf.Apikey
	if err != nil {
		kong.Log.Err(err.Error())
	}
	response := make([]string, 1)
	response[0] = "You have no correct key"
	j, err := json.Marshal(response)
	x := make(map[string][]string)
	x["Content-Type"] = append(x["Content-Type"], "application/json")
	if apiKey != key {
		kong.Response.Exit(403, string(j), x)
	}
}
