package utils

import (
	"github.com/joho/godotenv"
	"log/slog"
)

// EnvLoader 環境変数を読み込むインターフェース
type EnvLoader interface {
	LoadEnv(filename string)
}

// DotenvLoader .envファイルを読み込む構造体
type DotenvLoader struct{}

// LoadEnv .envファイルを読み込み、環境変数に設定する
func (d *DotenvLoader) LoadEnv(filename string) {
	err := godotenv.Load(filename)
	if err != nil {
		slog.Error(".envファイルの読み込みに失敗しました:", err)
	}
}
