package main

import (
	"fmt"
	"./irc"
	"./api"
	"strings"
	"regexp"
	"time"
	"os"
	"encoding/json"
)

type Config struct{
	Token string
}

var botIrcInfo = irc.IrcConnectionInfo{"irc.twitch.tv", 6667, "frde_bot", ""}


//REGEX to parse chat messages
var chatMsgRegex = regexp.MustCompile("(:[a-zA-Z1-9_!@.]+) PRIVMSG (#[a-zA-Z1-9_]+) (:[a-zA-Z1-9!_@ ]+)")

func main(){
	file, _ := os.Open("irc.json")
	decoder := json.NewDecoder(file)
	conf := Config{}
	err := decoder.Decode(&conf)
	if err != nil{
		panic(err.Error())
	}

	botIrcInfo.Password = conf.Token
	fmt.Println("using token "+conf.Token)
	err = runBot();
	if err != nil{
		fmt.Println(err.Error())
	}
}


func runBot() (err error){

	ircConnection, err := irc.CreateIrcConnection(&botIrcInfo)
	//Register commands:
	//time command
	ircConnection.Cm.RegisterCommand(irc.Command{"time", regexp.MustCompile("^(time)$"),func(channel string,args []string){
		ircConnection.SendMessage("PRIVMSG #"+channel+" :"+time.Now().String())
	}})

	defer ircConnection.CloseConnection()
	api.SetIrcConnection(ircConnection)
	go api.Ini()
	if err != nil{
		return err
	}

	ircConnection.Authenticate()
	ircConnection.JoinActiveChannels()
	for{
		line, err := ircConnection.GetMessage()

		if err != nil{
			return err
		}
		fmt.Println(line)
		//match with chat messages
		m := chatMsgRegex.FindStringSubmatch(line)
		if len(m)>3 {
			//extract info
			channel := m[2][1:] //removes '#' character
			message := m[3][1:] //removes ':' character
			//TODO Remove this and add a toggle for which channels can have messages sent to them
			if channel == "winter_squirrel" {
				//Make sure message is a command
				if (strings.HasPrefix(message, "!")) {
					//remove '!' character
					message = message[1:]
					//go through command and try find a match
					for _, command := range ircConnection.Cm.RegisteredCommands {
						match := command.RegexStr.FindStringSubmatch(message)

						if (len(match) > 1) {
							//Run command when match is found
							command.Run(channel, match)
						}
					}
				}
			}
		}


	}

	return nil
}


