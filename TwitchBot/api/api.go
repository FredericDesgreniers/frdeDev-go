package api
import(
	"../irc"
	"net/http"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

var channelPath = regexp.MustCompile("^/(channel)/([a-zA-Z0-9_]+)/(join|leave)$")

var ircConnection *irc.IrcConnection

func SetIrcConnection(ircC *irc.IrcConnection){
	ircConnection = ircC
}

//Used to set the status of a channel
// Setting active to true will make the irc join the channel
func SetChannelStatus(w http.ResponseWriter, r *http.Request) {
	m := channelPath.FindStringSubmatch(r.URL.Path)
	if m == nil{
		http.NotFound(w, r)
	}
	name := strings.ToLower(m[2])
	command := m[3]

	switch command{
	case "join":
		ircConnection.JoinChannel(name)
		break
	case "leave":
		ircConnection.LeaveChannel(name)
		break
	}
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

	http.HandleFunc("/channel/", SetChannelStatus)

	http.HandleFunc("/channels/", GetChannelsStatus)

	error := http.ListenAndServe(":8081", nil)

	fmt.Println(error.Error())
}