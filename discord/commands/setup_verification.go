package commands

import (
	"github.com/bwmarrin/discordgo"
)

type SetupVerification struct {
	AuthorizeURL string
}

func (command *SetupVerification) Identity() string {
	return "setup-verification"
}

func (command *SetupVerification) Create() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        command.Identity(),
		Description: "Create Verification Embed",
	}
}

func (command *SetupVerification) Handle(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{{
				Title:       "️**認証が必要です**",
				Description: "下のボタンを押してロールを手に入れましょう",
			}},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label: "認証する",
						URL:   command.AuthorizeURL,
						Emoji: discordgo.ComponentEmoji{Name: "✅"},
						Style: discordgo.LinkButton,
					},
				}},
			},
		},
	})
}
