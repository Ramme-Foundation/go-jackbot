package controller

type SlackCommand = string

const (
	SlackCommandHelp    SlackCommand = "help"
	SlackCommandJackpot SlackCommand = "jackpot"
	SlackCommandChoose  SlackCommand = "choose"
	SlackCommandRandom  SlackCommand = "random"
)
