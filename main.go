package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// サーバーID
// TODO: 環境変数に移動
const guildID = "1251543101756018769"

// チャンネルIDを格納したmap
var idMap = make(map[string]string)

func main() {
	// TODO: 環境変数に移動
	idMap["a"] = "1295673918463414343"

	getEnv("bot.env")

	// 環境変数からボットのトークンを取得
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_BOT_TOKEN 環境変数が設定されていません。")
		return
	}

	// 新しいDiscordセッションを作成
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Discordセッションの作成に失敗:", err)
		return
	}

	// メッセージ作成時のイベントハンドラーを追加
	// dg.AddHandler(messageCreate)
	dg.AddHandler(onInteractionCreate)

	// 必要なインテントを設定（メッセージの読み取りのためにGUILD_MESSAGESを有効に）
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsGuildMembers | discordgo.IntentsAll | discordgo.PermissionSendMessages
	openBot(dg)

	registerCommands(dg)

	// dg.AddHandler(onReady)

	// 終了シグナルを待機
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// セッションを閉じる
	dg.Close()
}

// .envファイルを読み込み、環境変数に設定する
func getEnv(filename string) {
	err := godotenv.Load(filename)
	if err != nil {
		log.Fatal(".envファイルの読み込みに失敗しました:", err)
		return
	}
}

// ボットにログイン
func openBot(dg *discordgo.Session) {
	err := dg.Open()
	if err != nil {
		log.Fatal("Discordセッションのオープンに失敗:", err)
		return
	}

	fmt.Println("ボットが起動しました。Ctrl+Cで終了します。")
}

// メッセージ作成時のイベントハンドラー
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ボット自身のメッセージには反応しない
	if m.Author.ID == s.State.User.ID {
		return
	}

	// 特定のメッセージに反応
	if m.Content == "!ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
}

func onReady(s *discordgo.Session, r *discordgo.Ready) {
	fmt.Println("BOT準備完了")
}

// コマンドを登録
func registerCommands(s *discordgo.Session) {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "add-schedule",
			Description: "リマインドしたいスケジュールを追加",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "team",
					Description: "所属するチームを選択",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "day_of_week",
					Description: "リマインドする曜日",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "時",
					Description: "0~23",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "分",
					Description: "0~59",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "role",
					Description: "リマインドする役職",
					Required:    true,
				},
			},
		},
		{
			Name:        "select",
			Description: "selecting roles",
		},
		{
			Name:        "hello",
			Description: "Responds with a greeting",
		},
	}

	for _, cmd := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, cmd)
		if err != nil {
			log.Fatalf("Failed to register command %s: %v", cmd.Name, err)
		}
	}

	fmt.Println("All commands registered successfully")
}

func onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		switch i.ApplicationCommandData().Name {
		case "add-schedule":
			line := i.ApplicationCommandData()
			team := i.ApplicationCommandData().Options[0].StringValue()
			//week := i.ApplicationCommandData().Options[1].StringValue()
			//hour := i.ApplicationCommandData().Options[2].IntValue()
			//minute := i.ApplicationCommandData().Options[3].IntValue()
			role := i.ApplicationCommandData().Options[4].RoleValue(s, guildID)

			var response string
			for _, opt := range line.Options {
				response += fmt.Sprintf("%s: %v\n", opt.Name, opt.Value)
			}

			var test = new(discordgo.Role)
			test.ID = role.ID

			sendRemindMessage(team, s, test)

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: response,
				},
			})

			if err != nil {
				log.Fatal("コマンド実行失敗")
				return
			}
		}
	}
}

// リマインドメッセージを送信
func sendRemindMessage(t string, s *discordgo.Session, test *discordgo.Role) {

	txt := fmt.Sprintf("<@&%s>", test.ID)

	_, err := s.ChannelMessageSend(idMap[t], txt+"\nリマインドです")
	if err != nil {
		log.Fatal("メッセージ送信失敗")
	}
}
