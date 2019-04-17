package main

import (
	"gov/backend/common/config"
	"github.com/go-errors/errors"
	"gov/public"
	"fmt"
	"log"
	"os"
)

func main() {
	defer logRecover() // log panic errors
	config.Get()		 // init config, define Set variable
	public.Public()	 // launch app
}

func logRecover() {
	if err := recover(); err != nil {
		// set log file
		file, er := os.OpenFile("logfile", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
		if er != nil {
			log.Fatalf("error opening file: %v", er)
		}
		defer file.Close()
		// ------------

		// set logger
		log.SetOutput(file)

		// append error into the log
		log.Println(errors.Wrap(err, 2).ErrorStack())
		// show error in the output
		fmt.Println(errors.Wrap(err, 2).ErrorStack())
	}
}
