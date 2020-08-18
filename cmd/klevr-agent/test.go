package main

import (
	"encoding/json"
	"fmt"
	"github.com/Klevry/klevr/pkg/agent"
	"io/ioutil"
	"os"
)

var test []agent.Secondary

func add(ip string, active bool){
	sec := agent.Secondary{}

	sec.IP = ip
	sec.IsAlive = active

	test = append(test, sec)
}
func main() {
	secon := agent.Cluster{}

	secon.Primary.IP = "192.168.1.2"

	add("test",true)

	secon.Secondary = test

	doc, _ := json.MarshalIndent(secon, "", "  ")

	err := ioutil.WriteFile(".test.json", doc, os.FileMode(0644))

	if err != nil{
		fmt.Println(err)
	}

	b, err2 := ioutil.ReadFile(".test.json")
	if err2 != nil{
		fmt.Println(err)
	}

	var data agent.Cluster
	json.Unmarshal(b, &data)

	fmt.Println(data.Secondary[0])
	fmt.Println(len(test))
}