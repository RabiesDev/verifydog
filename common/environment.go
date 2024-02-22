package common

import (
	"github.com/jessevdk/go-flags"
	"github.com/joho/godotenv"
	"os"
)

type Environment struct {
	// argument flags
	DevMode bool `short:"d" long:"debug" description:"start in development mode"`

	// environments
	Driver       string
	Database     string
	AuthorizeURL string
	RedirectURL  string
	BotToken     string
	ClientID     string
	ClientSecret string
	GuildID      string
	AuthRoleID   string
}

func (environment *Environment) Parse() error {
	parser := flags.NewParser(environment, flags.Default)
	if _, err := parser.Parse(); err != nil {
		return err
	}

	dotEnvFile := ".env"
	if environment.DevMode {
		dotEnvFile = ".env.local"
	}

	if err := godotenv.Load(dotEnvFile); err != nil {
		return err
	}

	environment.Driver = os.Getenv("DRIVER")
	environment.Database = os.Getenv("DATABASE")

	environment.AuthorizeURL = os.Getenv("AUTHORIZE_URL")
	environment.RedirectURL = os.Getenv("REDIRECT_URL")

	environment.BotToken = os.Getenv("BOT_TOKEN")
	environment.ClientID = os.Getenv("CLIENT_ID")
	environment.ClientSecret = os.Getenv("CLIENT_SECRET")

	environment.GuildID = os.Getenv("GUILD_ID")
	environment.AuthRoleID = os.Getenv("AUTH_ROLE_ID")

	return nil
}
