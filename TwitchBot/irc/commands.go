package irc

import "regexp"
// DoCommand is called when a command is amtched
type DoCommand func(string, []string)

//Command structure
type Command struct {
	Name string
	RegexStr *regexp.Regexp
	Run DoCommand
}

//Manages different commands
type CommandManager struct {
	RegisteredCommands []Command
}

// Register a command in the command manager
func (cm *CommandManager) RegisterCommand(c Command){
	cm.RegisteredCommands = append(cm.RegisteredCommands, c)
}
