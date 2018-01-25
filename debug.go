package main

import(
	"log"
	"net/http"
	_ "net/http/pprof"
)

func pprof(){
	log.Println(http.ListenAndServe("0.0.0.0:7777", nil))
}