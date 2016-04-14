package main

import (
	"fmt"
	"net/http"
	//"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("hello world")
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil) //local: http://192.168.2.100:8080/
}
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello home")
}