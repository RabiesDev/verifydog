package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
	"verifydog/common"
	"verifydog/discord/commands"
)

func (bot *Bot) OnReady(session *discordgo.Session, ready *discordgo.Ready) {
	logrus.Info("Logged in to ", session.State.User.Username)

	for _, guild := range session.State.Guilds {
		for _, v := range bot.SlashCommands {
			command, err := session.ApplicationCommandCreate(session.State.User.ID, guild.ID, v.Create())
			if err != nil {
				logrus.WithField("error", err).Error("failed to create command")
				continue
			}

			bot.CreatedCommands = append(bot.CreatedCommands, command)
		}
	}
}

func (bot *Bot) OnGuildMemberAdd(session *discordgo.Session, guildMemberAdd *discordgo.GuildMemberAdd) {
	if bot.IsInvalid(guildMemberAdd.GuildID, guildMemberAdd.Member.User) {
		return
	}

	authenticated, err := common.Find("snowflake = ?", guildMemberAdd.Member.User.ID)
	if err != nil || time.Since(authenticated.UpdatedAt) > (time.Hour*24)*7 {
		return
	}

	bot.OnAuthenticate(*authenticated)
}

func (bot *Bot) OnInteractionCreate(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if bot.IsInvalid(interaction.GuildID, interaction.Member.User) || interaction.Type != discordgo.InteractionApplicationCommand {
		return
	}

	interactionData := interaction.ApplicationCommandData()
	command, ok := lo.Find(bot.SlashCommands, func(command commands.SlashCommand) bool {
		return strings.EqualFold(interactionData.Name, command.Identity())
	})

	if !ok {
		logrus.Warnln("Invalid handler execution: ", interactionData.Name)
		return
	}

	if err := command.Handle(session, interaction); err != nil {
		logrus.Error(err)
	}
}

func (bot *Bot) OnAuthenticate(authenticated common.Authenticated) {
	if err := bot.Session.GuildMemberRoleAdd(bot.Environment.GuildID, authenticated.Snowflake, bot.Environment.AuthRoleID); err != nil {
		logrus.Error(err)
	}
}
