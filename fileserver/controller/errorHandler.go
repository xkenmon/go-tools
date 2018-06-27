package controller

import (
	"net/http"
	"log"
	"fmt"
)

func WriteErr(w http.ResponseWriter, err interface{}) {
	if err != nil {
		log.Println(err)
		w.Write([]byte(fmt.Sprintln(err)))
	}
}

func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
