package scheduler

import (
	"github.com/go-co-op/gocron/v2"
	"log"
	"time"
)

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
		log.Fatalf("ジョブ登録失敗: %v", err)
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

// SendRemindMessage リマインドメッセージを送信する関数
func SendRemindMessage(team string, roleID string) {
	// リマインドメッセージを送信するロジックをここに追加
	log.Printf("リマインドメッセージを送信: チーム=%s, 役職ID=%s", team, roleID)
}
