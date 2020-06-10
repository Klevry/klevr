package main


import (
	"strings"
	"time"
	"fmt"
	"strconv"
)


func main() {
	var a [3]string
	a[0] = "last_check=1591755764&ip=192.168.10.12&clientType=baremetal&masterConnection=fail&disconnected=0"
	a[1] = "last_check=1591755364&ip=192.168.10.13&clientType=baremetal&masterConnection=ok&disconnected=0"
	a[2] = "last_check=1591355764&ip=192.168.10.11&clientType=baremetal&masterConnection=ok&disconnected=0"

	var target_txt, time_arry []string
	var time_string string
	today := time.Now()
	today_unix := today.Unix()
	today_fmt := today.Format("2006-01-02 15:04:05")

	chk_day := today.AddDate(0, 0, -1)

	fmt.Println("NowU:", today_unix)
	fmt.Println("Date_N:",today_fmt)
	fmt.Println("Date_A:", chk_day)

	for i := 0; i < len(a) ; i++ {
		target_txt = strings.Split(string(a[i]), "&")
		time_arry = strings.Split(target_txt[0], "=")

		time_string = string(time_arry[1])
		time_parsing, err := strconv.ParseInt(time_string, 10, 64)
		if err != nil {
			panic(err)
		}

		tm := time.Unix(time_parsing, 0)
		fmt.Println("Unix:",time_string)

		if time.Since(tm).Hours() > 24 {
			println("day over!!")
		}else{
			println("It's okay")
		}

		fmt.Println("Date_B:",tm)
	}

}
