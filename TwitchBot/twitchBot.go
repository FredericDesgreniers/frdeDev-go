package main

import (
	"fmt"
	"./irc"
	"./api"
)



var botIrcInfo = irc.IrcConnectionInfo{"irc.twitch.tv", 6667, "frde_bot", "oauth:qmsz9bc54rqc5429r05oomjsvhbzkm"}


func main(){


	err := runBot();
	if err != nil{
		fmt.Println(err.Error())
	}
}


func runBot() (err error){

	ircConnection, err := irc.CreateIrcConnection(&botIrcInfo)
	defer ircConnection.CloseConnection()
	api.SetIrcConnection(ircConnection)
	go api.Ini()
	if err != nil{
		return err
	}

	ircConnection.Authenticate()

	ircConnection.JoinChannel("winter_squirrel")

	for{

		line, err := ircConnection.GetMessage()

		if err != nil{
			return err
		}

		fmt.Println(line)
	}

	return nil
}


