package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"verifydog/common"
	"verifydog/discord/commands"
)

type Bot struct {
	Session     *discordgo.Session
	Environment *common.Environment

	CreatedCommands []*discordgo.ApplicationCommand
	SlashCommands   []commands.SlashCommand
}

func NewBot(environment *common.Environment) *Bot {
	slashCommands := []commands.SlashCommand{
		&commands.SetupVerification{
			AuthorizeURL: environment.AuthorizeURL,
		},
	}

	return &Bot{
		Environment:   environment,
		SlashCommands: slashCommands,
	}
}

func (bot *Bot) Startup() error {
	session, err := discordgo.New("Bot " + bot.Environment.BotToken)
	if err != nil {
		return err
	}

	intents := discordgo.IntentGuilds | discordgo.IntentsGuildMessages | discordgo.IntentDirectMessages
	session.Identify.Intents = intents

	session.Identify.Presence.Game.Name = "VerifyDog 👀"
	session.Identify.Presence.Game.Type = discordgo.ActivityTypeGame

	session.AddHandler(bot.OnReady)
	session.AddHandler(bot.OnInteractionCreate)

	bot.Session = session
	return session.Open()
}

func (bot *Bot) CloseHandle() error {
	for _, guild := range bot.Session.State.Guilds {
		for _, v := range bot.CreatedCommands {
			err := bot.Session.ApplicationCommandDelete(bot.Session.State.User.ID, guild.ID, v.ID)
			if err != nil {
				logrus.Error("failed to delete command: ", err)
			}
		}
	}
	return bot.Session.Close()
}

func (bot *Bot) IsInvalid(guildID string, user *discordgo.User) bool {
	return guildID != bot.Environment.GuildID || user.Bot || user.System
}
