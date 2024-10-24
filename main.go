package main

import (
	"DiscordBot/commands"
	"DiscordBot/discord"
	"DiscordBot/scheduler"
	"DiscordBot/utils"
	"github.com/bwmarrin/discordgo"
	"log/slog"
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
		slog.Warn("Discordセッションの開始に失敗: %v", err)
	}

	defer dgs.Close()
	// 一括コマンド削除関数
	//defer unregisterAllCommands(dgs)

	slog.Info("ボットが起動しました。Ctrl+Cで終了します。")

	createCommands()

	waitForExitSignal()
}

// initializeEnvは環境変数の読み込みとJSONファイルの読み込みを行う
func initializeEnv() {
	envLoader := &utils.DotenvLoader{}
	// 開発時と本番時で環境変数を読み込むファイルを変更
	envLoader.LoadEnv("bot.env")

	GuildID = os.Getenv("DISCORD_GUILD_ID")

	commands.JobDataSlice = append(reader.ReadJSON("jobData.json"))
}

// initializeScheduleはスケジューラを初期化する
func initializeSchedule() {
	for _, jobData := range commands.JobDataSlice {
		s.RegisterJob(jobData.Cron, scheduler.SendRemindMessage, jobData.Team, jobData.Role, dgs)
	}

	s.Start()
}

// clearCacheはアプリケーションコマンドを削除する
func unregisterAllCommands(ds *discordgo.Session) {
	cs, err := ds.ApplicationCommands(ds.State.User.ID, "")
	if err != nil {
		slog.Warn("アプリケーションコマンドの取得に失敗しました: %v", err)
	}

	for _, cmd := range cs {
		err := ds.ApplicationCommandDelete(ds.State.User.ID, "", cmd.ID)
		if err != nil {
			slog.Warn("コマンドの削除に失敗しました: %v", err)
		} else {
			slog.Info("コマンドを削除しました: %s", cmd.Name)
		}
	}
}

// createCommandsはDiscordのコマンドを登録する
func createCommands() {
	add := &commands.CreateAddScheduleCommand{}
	show := &commands.CreateShowSchedulesCommand{}
	remove := &commands.CreateRemoveScheduleCommand{}
	// デバッグ用
	// showEnv := &commands.CreateShowEnvCommand{}

	for _, cmd := range add.CreateCommand() {
		_, err := dgs.ApplicationCommandCreate(dgs.State.User.ID, "", cmd)
		if err != nil {
			slog.Warn("コマンド登録失敗: ", err)
		}
	}

	for _, cmd := range show.CreateCommand() {
		_, err := dgs.ApplicationCommandCreate(dgs.State.User.ID, "", cmd)
		if err != nil {
			slog.Warn("コマンド登録失敗: ", err)
		}
	}

	for _, cmd := range remove.CreateCommand() {
		_, err := dgs.ApplicationCommandCreate(dgs.State.User.ID, "", cmd)
		if err != nil {
			slog.Warn("コマンド登録失敗: ", err)
		}
	}
	// デバッグ用
	//for _, cmd := range showEnv.CreateCommand() {
	//	_, err := dgs.ApplicationCommandCreate(dgs.State.User.ID, "", cmd)
	//	if err != nil {
	//		slog.Warn("コマンド登録失敗: ", err)
	//	}
	//}
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
		//case "show-env":
		//	cmd = &commands.ShowEnvCommand{}
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
