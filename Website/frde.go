package main

import (
	"net/http"
	"html/template"
	"encoding/json"
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



func getJson(url string, target interface{}) error{
	r, err := http.Get(url)
	if err != nil{
		return err;
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func adminHandler(w http.ResponseWriter, r *http.Request){
	var channels []channel
	getJson("http://localhost:8081/channels/", &channels)
	templates.ExecuteTemplate(w, "admin.html", channels)

}



func main(){
	http.HandleFunc("/admin/", adminHandler)

	http.ListenAndServe(":8080", nil)
}