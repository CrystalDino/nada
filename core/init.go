package core

//init the system,using for parsing flags
import (
	"log"
)

func init() {
	var err error
	//load config file
	log.SetFlags(log.LstdFlags)

	if err = LoadConfigFile("./config.json"); err != nil {
		log.Fatalln("init config error,", err)
		return
	}
	log.Println("load system config file done")
}
