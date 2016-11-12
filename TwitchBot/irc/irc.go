package irc

import (
	"net"
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"database/sql"
)
import _ "github.com/go-sql-driver/mysql"



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
	reader     *bufio.Reader
	Info       *IrcConnectionInfo

	Database   *sql.DB
}

func (irc *IrcConnection) JoinActiveChannels(){
	var name string

	qeury, _ := irc.Database.Query("SELECT name FROM channels WHERE Active = 1")
	defer qeury.Close()
	for(qeury.Next()){
		qeury.Scan(&name)
		irc.SendMessage("JOIN #"+name)
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

	query, err := ircConnection.Database.Prepare("SELECT Active FROM channels WHERE Name = ?")
	if err !=nil{
		fmt.Println(err.Error())
		return
	}
	defer query.Close()
	var active bool

	err = query.QueryRow(channel).Scan(&active)
	if err != nil{
		iQeury, _ := ircConnection.Database.Prepare("INSERT INTO channels VALUES(? , ?)")
		defer iQeury.Close()
		iQeury.Exec(channel, true)

		fmt.Printf("Added new channel to database: #%s\n", channel)
	}else{
		if(active){
			fmt.Printf("Channel #%s was already joined\n", channel)
			return
		}else{
			update, _ := ircConnection.Database.Prepare("UPDATE channels SET Active = ? WHERE Name = ?")
			update.Exec(true, channel)
		}
	}


	fmt.Printf("Channel #%s has been joined\n", channel)
	ircConnection.SendMessage("JOIN %s", "#"+channel)



}
//Leave an irc channel
func (ircConnection *IrcConnection) LeaveChannel(channel string){

	query, _ := ircConnection.Database.Prepare("SELECT Active FROM channels WHERE Name = ?")
	defer query.Close()
	var active bool

	query.QueryRow(channel).Scan(&active)

	if(!active){
		fmt.Printf("Irc has never joined #%s\n", channel)
		return
	}

	update, _ := ircConnection.Database.Prepare("UPDATE channels SET Active = ? WHERE Name = ?")
	update.Exec(false, channel)

	ircConnection.SendMessage("PART %s", "#"+channel)

	fmt.Printf("Left channel %s\n", channel)

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
	sqldb, err := sql.Open("mysql", "root@/frde")
	if err != nil{
		panic(err.Error())
	}
	ircC := &IrcConnection{conn, bufio.NewReader(conn), info,sqldb}


	return ircC, nil
}