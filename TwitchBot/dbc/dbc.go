package dbc

import (
	"database/sql"
)
//Channel structure corresponding to state of a channel for the bot
type Channel struct {
	Id int
	Name string
	Active bool
}
// structure to contain database info
type DatabaseConnection struct{
	Database *sql.DB
}
//Get a slice of channels for a name.
// Should normally return a slice smaller than 2
func (data *DatabaseConnection) GetChannel(name string) *[]*Channel{
	pStatement, _ := data.Database.Prepare("SELECT Name, Active FROM channels WHERE Name = ?")
	defer pStatement.Close()

	q, _ := pStatement.Query(name)

	channels := make([]*Channel,0)

	for(q.Next()){
		c := Channel{}
		q.Scan(&c.Name, &c.Active)
		channels = append(channels, &c)
	}
	return &channels
}
// Get a slice of all channels in the database
func (data *DatabaseConnection) GetChannels() *[]*Channel{
	pStatement, _ := data.Database.Prepare("SELECT Name, Active FROM channels")
	defer pStatement.Close()

	q, _ := pStatement.Query()

	channels := make([]*Channel,0)

	for(q.Next()){
		c := Channel{}
		q.Scan(&c.Name, &c.Active)
		channels = append(channels, &c)
	}
	return &channels
}
// Get list of all active channels in database
func (data *DatabaseConnection) GetActiveChannels() *[]*Channel{
	pStatement, _ := data.Database.Prepare("SELECT Name, Active FROM channels WHERE Active = ?")
	defer pStatement.Close()

	q, _ := pStatement.Query(true)

	channels := make([]*Channel,0)

	for(q.Next()){
		c := Channel{}
		q.Scan(&c.Name, &c.Active)
		channels = append(channels, &c)
	}
	return &channels
}
//Insert channel in database
func (data *DatabaseConnection) InsertChannel(ch *Channel){
	pStatement, _ := data.Database.Prepare("INSERT INTO channels VALUES(?, ?)")
	defer pStatement.Close()

	pStatement.Exec(ch.Name, ch.Active)
}
// Update channel in database
func (data *DatabaseConnection) UpdateChannel(ch *Channel){
	pStatement, _ := data.Database.Prepare("UPDATE channels SET Active = ? WHERE Name = ?")
	defer pStatement.Close()
	pStatement.Exec(ch.Active, ch.Name)
}
//Create a connection for the database
//Url should be in the form user:pass@\dbname
func CreateDatabaseConnection(url string) *DatabaseConnection{
	sqldb, err := sql.Open("mysql", url)
	if err != nil{

	}
	dbc := DatabaseConnection{sqldb}

	return &dbc
}