package slack

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// cowart is the formatting string for printing the cow art.
const cowart = "Moo! %s"

//encore:api public raw path=/cowsay
func Cowsay(w http.ResponseWriter, req *http.Request) {
	text := req.FormValue("text")
	data, _ := json.Marshal(map[string]string{
		"response_type": "in_channel",
		"text":          fmt.Sprintf(cowart, text),
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
}
