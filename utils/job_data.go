package utils

// JobData は、チーム、cronスケジュール、および役職を持つジョブの構造体
type JobData struct {
	Team string `json:"team"`
	Cron string `json:"cron"`
	Role string `json:"role"`
}
