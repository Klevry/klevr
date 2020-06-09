package main

import (
        "strings"
        "flag"
        _"regexp"
        "encoding/json"
        "github.com/gorilla/mux"
	"github.com/ralfyang/klevr/communicator"
) 


var api_server string
var klevr_server = "http://192.168.10.11:8080"
var api_key_string string

func Get_apiserver_info() string{
	api_key_string = communicator.Get_http(klevr_server+"/apikey", _ )
	return api_key_string
}

func Get_apiserver_info() string{
	api_server = communicator.Get_http(klevr_server+"/apiserver", api_key_string )
	return api_server
}


func main(){
	println("apiserver :", api_server)
	println("apikey :", api_key_string)

}


