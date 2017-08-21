package controller

//data manage

import "log"

func init() {
	if Server == nil {
		log.Fatalln("data:init web server error")
		return
	}

}
