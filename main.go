package main

import (
	"DiscordBot/commands"
	"DiscordBot/discord"
	"DiscordBot/scheduler"
	"DiscordBot/utils"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	GuildID string
	dgs     *discordgo.Session
	reader  utils.JSONReader = &utils.FileJSONReader{}
	s                        = commands.Ns
)

func main() {
	initializeEnv()
	initializeSchedule()

	sessionManager := &discord.DiscordSessionManager{}
	dgs = sessionManager.InitializeSession(os.Getenv("DISCORD_BOT_TOKEN"))

	dgs.AddHandler(onInteractionCreate)
	dgs.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsGuildMembers | discordgo.IntentsAll | discordgo.PermissionSendMessages

	if err := dgs.Open(); err != nil {
		log.Fatalf("Discordセッションのオープンに失敗: %v", err)
	}
	defer dgs.Close()

	log.Println("ボットが起動しました。Ctrl+Cで終了します。")

	createCommands()

	waitForExitSignal()
}

// initializeEnvは環境変数の読み込みとJSONファイルの読み込みを行う
func initializeEnv() {
	envLoader := &utils.DotenvLoader{}
	envLoader.LoadEnv("bot.env")
	envLoader.LoadEnv("channel.env")

	GuildID = os.Getenv("DISCORD_GUILD_ID")

	commands.JobDataSlice = append(reader.ReadJSON("jobData.json"))
}

// initializeScheduleはスケジューラを初期化する
func initializeSchedule() {
	for _, jobData := range commands.JobDataSlice {
		s.RegisterJob(jobData.Cron, scheduler.SendRemindMessage, jobData.Team, jobData.Role)
	}

	s.Start()
}

// createCommandsはDiscordのコマンドを登録する
func createCommands() {
	add := commands.CreateAddScheduleCommand{}
	show := &commands.CreateShowSchedulesCommand{}
	remove := &commands.CreateRemoveScheduleCommand{}

	for _, cmd := range add.CreateCommand() {
		_, err := dgs.ApplicationCommandCreate(dgs.State.User.ID, GuildID, cmd)
		if err != nil {
			log.Fatalf("コマンド登録失敗: %v", err)
		}
	}

	for _, cmd := range show.CreateCommand() {
		_, err := dgs.ApplicationCommandCreate(dgs.State.User.ID, GuildID, cmd)
		if err != nil {
			log.Fatalf("コマンド登録失敗: %v", err)
		}
	}

	for _, cmd := range remove.CreateCommand() {
		_, err := dgs.ApplicationCommandCreate(dgs.State.User.ID, GuildID, cmd)
		if err != nil {
			log.Fatalf("コマンド登録失敗: %v", err)
		}
	}
}

// onInteractionCreateはDiscordからのインタラクションイベントを処理する
func onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var cmd commands.Command

	switch i.ApplicationCommandData().Name {
	case "add-schedule":
		cmd = &commands.AddScheduleCommand{}
	case "show-schedules":
		cmd = &commands.ShowSchedulesCommand{}
	case "remove-schedule":
		cmd = &commands.RemoveScheduleCommand{}
	}

	if cmd != nil {
		cmd.Execute(s, i)
	}
}

// waitForExitSignalは終了シグナルを待機してボットを安全にシャットダウンする
func waitForExitSignal() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
