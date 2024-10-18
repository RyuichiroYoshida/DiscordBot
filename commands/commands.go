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
	response := fmt.Sprintf("リマインドスケジュールを追加しました\nチーム: %s\n曜日: %s\n時間: %d時%d分\n役職: %s", team, weekParseData[week], hour, minute, role.Name)
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: response},
	}); err != nil {
		log.Fatalf("コマンド実行失敗: %v", err)
	}
	log.Println(response)
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
