package main

import (
	"net/http"
	"html/template"
	"encoding/json"
	"time"
	"fmt"
	"io/ioutil"
	"log"
)

var WebsitePath = "Website/"

var pagesPath = WebsitePath+"pages/"
var templatesPath = WebsitePath+"templates/"

var templates = template.Must(template.ParseFiles(templatesPath+"admin.html"))

type channel struct{
	Id int
	Name string
	Active bool
}


//Get the json from a url and map it to an interface
func getJson(url string, target interface{}) error{
	//timout after 5 seconds
	timeout := time.Duration(5 * time.Second)
	client := http.Client{Timeout: timeout}

	r, err := client.Get(url)
	if err != nil{
		return err;
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func adminHandler(w http.ResponseWriter, r *http.Request){
	var channels []channel
	if err := getJson("http://localhost:8081/channels/", &channels); err != nil{

		return
	}
	templates.ExecuteTemplate(w, "admin.html", channels)
}

//This handler is a proxy between the website and the twitch bot\
// a user can specify a channel  (ex: winter_squirrel) and an action (ex: leave / join)
// And the request will be forwarded to the twitch bot
func channelHandler(w http.ResponseWriter, r *http.Request) {
	client := http.Client{Timeout: time.Duration(5 * time.Second)}

	//Get channel and action from body
	channel := r.FormValue("channel")
	action := r.FormValue("action")

	//TODO add other values to proxy url so complex commands can also be called

	//Ask the twitch bot to do the action
	response, err := client.Get(fmt.Sprintf("http://localhost:8081/channel/%s/%s", channel, action))
	if err != nil{
		log.Fatal(err.Error())
	}
	defer response.Body.Close()
	//Give back the twitch bot response
	body, err := ioutil.ReadAll(response.Body)
	w.Write([]byte(body))
	return
}


func main(){
	//TODO add authentication for admin subpath
	http.HandleFunc("/admin/panel/", adminHandler)
	http.HandleFunc("/admin/channel/", channelHandler)

	http.ListenAndServe(":8080", nil)
}