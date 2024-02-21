package commands

import "github.com/bwmarrin/discordgo"

type SlashCommand interface {
	Identity() string
	Create() *discordgo.ApplicationCommand
	Handle(session *discordgo.Session, interaction *discordgo.InteractionCreate) error
}
