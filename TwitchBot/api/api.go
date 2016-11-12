package api
import(
	"../irc"
	"net/http"
	"encoding/json"
	"fmt"
)

var ircConnection *irc.IrcConnection

func SetIrcConnection(ircC *irc.IrcConnection){
	ircConnection = ircC
}

//Used to set the status of a channel
// Setting active to true will make the irc join the channel
func SetChannelStatus(w http.ResponseWriter, r *http.Request) {

}
// Get a list of channels that have been used
func GetChannelsStatus(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ircConnection.Channels)
}

// Initialize api web stuff
func Ini(){
	createHandlers()
}

// Create teh handlers for listening to api requestds
func createHandlers(){
	fmt.Println("API - Creating handlers")

	http.HandleFunc("/channel/set/", SetChannelStatus)

	http.HandleFunc("/channels/", GetChannelsStatus)

	error := http.ListenAndServe(":8081", nil)

	fmt.Println(error.Error())
}