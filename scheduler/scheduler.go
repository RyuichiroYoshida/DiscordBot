package scheduler

import (
	"DiscordBot/utils"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron/v2"
	"log/slog"
	"os"
	"time"
)

var ChannelId = map[string]string{}

// ReminderTiming 登録したスケジュールから、この時間分前にリマインドを送信する (Hour)
const ReminderTiming = 3

// init は環境変数を読み込む
func init() {
	envLoader := &utils.DotenvLoader{}
	envLoader.LoadEnv("channel.env")

	ChannelId = map[string]string{
		"a": os.Getenv("TEAM_A"),
		"b": os.Getenv("TEAM_B"),
		"c": os.Getenv("TEAM_C"),
		"d": os.Getenv("TEAM_D"),
		"e": os.Getenv("TEAM_E"),
	}
}

// Scheduler スケジューラを管理するインターフェース
type Scheduler interface {
	RegisterJob(cronExpr string, jobFunc interface{}, params ...interface{})
	Start()
	Jobs() []gocron.Job
}

// GoCronScheduler gocronを使用したスケジューラの実装
type GoCronScheduler struct {
	scheduler gocron.Scheduler
	err       error
}

// NewGoCronScheduler GoCronSchedulerのインスタンスを作成
func NewGoCronScheduler(location *time.Location) *GoCronScheduler {
	s, err := gocron.NewScheduler(gocron.WithLocation(location))
	return &GoCronScheduler{
		scheduler: s,
		err:       err,
	}
}

// RegisterJob ジョブをスケジューラに登録
func (s *GoCronScheduler) RegisterJob(cronExpr string, jobFunc interface{}, params ...interface{}) {
	if _, err := s.scheduler.NewJob(
		gocron.CronJob(cronExpr, false),
		gocron.NewTask(jobFunc, params...),
	); err != nil {
		slog.Warn("ジョブ登録失敗: %v", err)
	}
}

// Start スケジューラを開始
func (s *GoCronScheduler) Start() {
	s.scheduler.Start()
}

// Jobs 登録されているジョブを返す
func (s *GoCronScheduler) Jobs() []gocron.Job {
	return s.scheduler.Jobs()
}

// RemoveJob ジョブをスケジューラから削除
func (s *GoCronScheduler) RemoveJob(jobID int) {
	err := s.scheduler.RemoveJob(s.Jobs()[jobID].ID())
	if err != nil {
		slog.Warn("ジョブ削除失敗: %v", err)
	}
}

// SendRemindMessage リマインドメッセージを送信する関数
func SendRemindMessage(team string, roleID string, dgs *discordgo.Session) {

	txt := fmt.Sprintf("<@&%s>\nMTGリマインドです～", roleID)
	_, err := dgs.ChannelMessageSend(ChannelId[team], txt)
	if err != nil {
		slog.Warn("メッセージ送信失敗: %v", err)
	}
}
