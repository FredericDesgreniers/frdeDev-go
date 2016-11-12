package irc

import (
	"net"
	"bufio"
	"fmt"
	"strconv"
)
// ChannelInfo store the status of a channel on the irc
// Name being the channel name and Active being if the channel is currently being listened too
type ChannelInfo struct{
	// Name of the channel
	Name string
	// If the channel is being listened too
	Active bool
}
// IrcConnectionInfo store information about an irc connection
type IrcConnectionInfo struct{
	// irc server name address
	ServerName string
	// irc server port
	ServerPort int
	// authentication username
	Username string
	// authentication password
	Password string
}
// IrcConnection keeps track of the connections and channels
type IrcConnection struct{
	connection net.Conn
	reader *bufio.Reader
	Info *IrcConnectionInfo

	Channels map[string]*ChannelInfo
}

//Send a message using the printF format
func (irc *IrcConnection) SendMessage(message string, args ...interface{}) (error){

	fmt.Printf("Sending: "+message+"\r\n", args...)
	_, err := fmt.Fprintf(irc.connection, message+" \r\n", args...)
	if err != nil{
		return err;
	}

	return nil
}
// Get the next message from the irc
func (irc *IrcConnection) GetMessage() (string, error){
	status, err := bufio.NewReader(irc.connection).ReadString('\n')
	if err != nil{
		return "", err
	}

	return status, nil
}
//Default enthentication
//Uses PASS and NICK
func (ircConnection *IrcConnection) Authenticate() {
	ircConnection.SendMessage("PASS %s",  ircConnection.Info.Password)
	ircConnection.SendMessage("NICK %s",  ircConnection.Info.Username)

}
// Join an irc channel
func (ircConnection *IrcConnection) JoinChannel(channel string){
	ircConnection.SendMessage("JOIN %s", "#"+channel)

	if channelInfo, ok := ircConnection.Channels[channel]; ok{
		channelInfo.Active = true
	}else{
		ircConnection.Channels[channel] = &ChannelInfo{channel, true}
	}


}
//Leave an irc channel
func (ircConnection *IrcConnection) LeaveChannel(channel string){
	ircConnection.SendMessage("PART %s", "#"+channel)
	if channelInfo, ok := ircConnection.Channels[channel]; ok{
		channelInfo.Active = false
	}else{
		ircConnection.Channels[channel] = &ChannelInfo{channel, false}
	}

}
//Create an irc connection
// should be called first
func CreateIrcConnection(info *IrcConnectionInfo) (*IrcConnection, error){

	conn, err := net.Dial("tcp", info.ServerName+":"+strconv.Itoa(info.ServerPort))
	if err != nil{
		return &IrcConnection{}, err
	}
	ircC := &IrcConnection{conn, bufio.NewReader(conn), info, make(map[string]*ChannelInfo)}

	return ircC, nil


}