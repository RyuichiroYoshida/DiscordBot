package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// サーバーID
// TODO: 環境変数に移動
const guildID = "1251543101756018769"

// チャンネルIDを格納したmap
var idMap = make(map[string]string)

var jst, _ = time.LoadLocation("Asia/Tokyo")

var ns, _ = gocron.NewScheduler(gocron.WithLocation(jst))

var dgs *discordgo.Session

func main() {
	scheduler()

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
	dgs = dg

	// メッセージ作成時のイベントハンドラーを追加
	dg.AddHandler(onInteractionCreate)

	// 必要なインテントを設定（メッセージの読み取りのためにGUILD_MESSAGESを有効に）
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsGuildMembers | discordgo.IntentsAll | discordgo.PermissionSendMessages
	openBot(dg)

	registerCommands(dg)

	// 終了シグナルを待機
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// セッションを閉じる
	dg.Close()
}

// リマインドメッセージを送信するタスク
func task() {
	_, err := dgs.ChannelMessageSend("1295673918463414343", "\nリマインドです")
	if err != nil {
		log.Fatal("メッセージ送信失敗")
	}
}

// スケジューラーを起動
func scheduler() {
	ns.Start()
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

// スケジュールを表示
func showSchedules() {
	ns.Jobs()
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
					Description: "所属するチームを選択 (a~e)",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "曜日",
					Description: "リマインドする曜日 (英語3文字で指定 sun~sat)",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "時",
					Description: "リマインドする時間 (0~23)",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "分",
					Description: "リマインドする分 (0~59)",
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
			Name:        "show-schedules",
			Description: "登録されているスケジュールを表示",
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

// コマンド実行時処理
func onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		switch i.ApplicationCommandData().Name {
		case "add-schedule":
			// コマンド引数で入力された要素を読み取る
			line := i.ApplicationCommandData()
			team := i.ApplicationCommandData().Options[0].StringValue()
			week := i.ApplicationCommandData().Options[1].StringValue()
			hour := i.ApplicationCommandData().Options[2].IntValue()
			minute := i.ApplicationCommandData().Options[3].IntValue()
			role := i.ApplicationCommandData().Options[4].RoleValue(s, guildID)

			// コマンド引数で入力された全ての要素の名前と値を改行区切りで繋げる
			var response string
			for _, opt := range line.Options {
				response += fmt.Sprintf("%s: %v\n", opt.Name, opt.Value)
			}

			var test = new(discordgo.Role)
			test.ID = role.ID

			// コマンド引数で入力された文字を小文字に変換
			team = strings.ToLower(team)
			week = strings.ToLower(week)

			// コマンド引数で入力されたスケジュールをcron形式に整形
			cronText := fmt.Sprintf("%d %d * * %s", minute, hour, week)

			// ジョブ登録
			_, er := ns.NewJob(
				gocron.CronJob(cronText, false),
				gocron.NewTask(task),
			)
			if er != nil {
				log.Fatal("ジョブ登録失敗")
				return
			}

			// コマンド実行時に入力内容をリマインドする
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

		case "show-schedules":
			var response string
			for _, s := range ns.Jobs() {
				nj, _ := s.NextRun()
				response += fmt.Sprintf("%s\n", nj)
			}

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: response,
				},
			})

			if err != nil {
				log.Fatal("スケジュール表示失敗")
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
