package irc

import (
	"net"
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"../dbc"
)
import (
	_ "github.com/go-sql-driver/mysql"
)



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
	reader     *bufio.Reader
	Info       *IrcConnectionInfo

	Database   *dbc.DatabaseConnection

	Cm *CommandManager
}

func (irc *IrcConnection) JoinActiveChannels(){
	channels :=irc.Database.GetActiveChannels()
	for _, c := range (*channels){
		if c != nil{
			irc.SendMessage("JOIN #"+c.Name)
		}
	}

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
	if strings.HasPrefix(status, "PING"){
		irc.SendMessage("PONG"+status[4:])
		return irc.GetMessage()
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

	channels := ircConnection.Database.GetChannel(channel)


	if len(*channels) == 0{
		ircConnection.Database.InsertChannel(&dbc.Channel{0, channel, true})
		fmt.Printf("Added new channel to database: #%s\n", channel)
	}else{
		ch := (*channels)[0]

		if(ch.Active){
			fmt.Printf("Channel #%s was already joined\n", channel)
			return
		}else{
			(*channels)[0].Active = true
			ircConnection.Database.UpdateChannel((*channels)[0])
		}
	}


	fmt.Printf("Channel #%s has been joined\n", channel)
	ircConnection.SendMessage("JOIN %s", "#"+channel)



}
//Leave an irc channel
func (ircConnection *IrcConnection) LeaveChannel(channel string){

	channels := ircConnection.Database.GetChannel(channel)


	if len(*channels) == 0{
		ircConnection.Database.InsertChannel(&dbc.Channel{0, channel, false})
		fmt.Printf("Added new channel to database: #%s\n", channel)
		return
	}else{
		ch := (*channels)[0]

		if(!ch.Active){
			fmt.Printf("Channel #%s is alrady non-active\n", channel)
			return
		}else{
			(*channels)[0].Active = false
			ircConnection.Database.UpdateChannel((*channels)[0])
		}
	}


	fmt.Printf("Channel #%s has been left\n", channel)
	ircConnection.SendMessage("PART %s", "#"+channel)



}
// Close the irc connection.
func (ircConnection *IrcConnection) CloseConnection(){
	ircConnection.connection.Close()
}
//Create an irc connection
// should be called first
func CreateIrcConnection(info *IrcConnectionInfo) (*IrcConnection, error){
	conn, err := net.Dial("tcp", info.ServerName+":"+strconv.Itoa(info.ServerPort))
	if err != nil{
		return &IrcConnection{}, err
	}

	ircC := &IrcConnection{conn, bufio.NewReader(conn), info,dbc.CreateDatabaseConnection("root@/frde"), &CommandManager{}}


	return ircC, nil
}