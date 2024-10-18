package discord

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

// SessionManager Discordセッションを管理するインターフェース
type SessionManager interface {
	InitializeSession(token string) *discordgo.Session
}

// DiscordSessionManager Discordセッションを管理する構造体
type DiscordSessionManager struct{}

// InitializeSession Discordセッションを初期化する
func (d *DiscordSessionManager) InitializeSession(token string) *discordgo.Session {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Discordセッションの作成に失敗: %v", err)
	}
	return dg
}
