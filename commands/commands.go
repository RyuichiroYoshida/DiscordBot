package commands

import (
	"DiscordBot/scheduler"
	"DiscordBot/utils"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
	"time"
)

var (
	JobDataSlice []utils.JobData
	Ns           = scheduler.NewGoCronScheduler(time.FixedZone("Asia/Tokyo", 9*60*60))
	jsonWriter   = &utils.FileJSONWriter{}
)

// Command コマンドを処理するインターフェース
type Command interface {
	Execute(s *discordgo.Session, i *discordgo.InteractionCreate)
}

// CommandFactory コマンドを生成するインターフェース
type CommandFactory interface {
	CreateCommand() []*discordgo.ApplicationCommand
}

// CreateAddScheduleCommand スケジュールを追加するコマンドを生成するファクトリ
type CreateAddScheduleCommand struct{}

func (c *CreateAddScheduleCommand) CreateCommand() []*discordgo.ApplicationCommand {

	dc := []*discordgo.ApplicationCommand{
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
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "day_of_week",
					Description: "リマインドする曜日 (0: 日曜日, 1: 月曜日, 2: 火曜日, 3: 水曜日, 4: 木曜日, 5: 金曜日, 6: 土曜日)",
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
	}

	return dc
}

// AddScheduleCommand スケジュールを追加するコマンド
type AddScheduleCommand struct{}

func (c *AddScheduleCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	team := strings.ToLower(i.ApplicationCommandData().Options[0].StringValue())
	week := i.ApplicationCommandData().Options[1].IntValue()
	hour := i.ApplicationCommandData().Options[2].IntValue()
	minute := i.ApplicationCommandData().Options[3].IntValue()
	role := i.ApplicationCommandData().Options[4].RoleValue(s, i.GuildID)

	cronText := fmt.Sprintf("%d %d * * %d", minute, hour, week)
	jobData := utils.JobData{Team: team, Cron: cronText, Role: role.ID}
	JobDataSlice = append(JobDataSlice, jobData)
	jsonWriter.WriteJSON("jobData.json", JobDataSlice)
	Ns.RegisterJob(cronText, scheduler.SendRemindMessage, team, role.ID)

	weekParseData := [7]string{"日曜日", "月曜日", "火曜日", "水曜日", "木曜日", "金曜日", "土曜日"}
	response := fmt.Sprintf("リマインドスケジュールを追加しました (Number: %d)\nチーム: %s\n曜日: %s\n時間: %d時%d分\n役職: %s", len(JobDataSlice), team, weekParseData[week], hour, minute, role.Name)
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: response},
	}); err != nil {
		log.Fatalf("コマンド実行失敗: %v", err)
	}
	log.Println(response)
}

// CreateShowSchedulesCommand スケジュールを表示するコマンドを生成するファクトリ
type CreateShowSchedulesCommand struct{}

func (c *CreateShowSchedulesCommand) CreateCommand() []*discordgo.ApplicationCommand {

	dc := []*discordgo.ApplicationCommand{
		{
			Name:        "show-schedules",
			Description: "登録されているスケジュールを表示",
		},
	}

	return dc
}

// ShowSchedulesCommand スケジュールを表示するコマンド
type ShowSchedulesCommand struct{}

func (c *ShowSchedulesCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var response string
	for _, job := range Ns.Jobs() {
		if nextRun, err := job.NextRun(); err == nil {
			response += fmt.Sprintf("%s\n", nextRun)
		}
	}
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: response},
	}); err != nil {
		log.Fatalf("スケジュール表示失敗: %v", err)
	}
}

// CreateRemoveScheduleCommand スケジュールを削除するコマンドを生成するファクトリ
type CreateRemoveScheduleCommand struct{}

func (c *CreateRemoveScheduleCommand) CreateCommand() []*discordgo.ApplicationCommand {
	dc := []*discordgo.ApplicationCommand{
		{
			Name:        "remove-schedule",
			Description: "登録されているスケジュールを削除",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "job_number",
					Description: "削除したいスケジュールの番号",
					Required:    true,
				},
			},
		},
	}

	return dc
}

// RemoveScheduleCommand スケジュールを削除するコマンド
type RemoveScheduleCommand struct{}

func (c *RemoveScheduleCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// ナンバーは1から始まるため、1を引く
	jobNumber := i.ApplicationCommandData().Options[0].IntValue() - 1
	// 指定された番号のスケジュールが存在しない場合はエラーメッセージを返す
	if jobNumber < 0 || jobNumber >= (int64(len(JobDataSlice))) {
		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "指定された番号のスケジュールは存在しません"},
		}); err != nil {
			log.Fatalf("スケジュール削除失敗: %v", err)
		}
		return
	}

	Ns.RemoveJob(int(jobNumber))
	JobDataSlice = append(JobDataSlice[:jobNumber], JobDataSlice[jobNumber+1:]...)
	jsonWriter.WriteJSON("jobData.json", JobDataSlice)

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: "スケジュールを削除しました"},
	}); err != nil {
		log.Fatalf("スケジュール削除失敗: %v", err)
	}
}
