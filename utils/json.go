package utils

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
)

// JSONWriter JSONデータを書き込むインターフェース
type JSONWriter interface {
	WriteJSON(filename string, data interface{})
}

// FileJSONWriter ファイルにJSONデータを書き込む構造体
type FileJSONWriter struct{}

// WriteJSON JSONデータをファイルに書き込む
func (w *FileJSONWriter) WriteJSON(filename string, data interface{}) {
	f, err := os.Create(filename)
	if err != nil {
		slog.Warn("ファイル作成失敗: %v", err)
	}
	defer f.Close()

	output, err := json.MarshalIndent(data, "", "\t\t")
	if err != nil {
		slog.Warn("JSON変換失敗: %v", err)
	}

	if _, err := f.Write(output); err != nil {
		slog.Warn("ファイル書き込み失敗: %v", err)
	}
}

// JSONReader JSONデータを読み込むインターフェース
type JSONReader interface {
	ReadJSON(filename string) []JobData
}

// FileJSONReader ファイルからJSONデータを読み込む構造体
type FileJSONReader struct{}

// ReadJSON JSONデータをファイルから読み込む
func (r *FileJSONReader) ReadJSON(filename string) []JobData {
	f, err := os.Open(filename)
	if err != nil {
		slog.Error("ファイルオープン失敗: %v", err)
	}
	defer f.Close()

	var data []JobData
	decoder := json.NewDecoder(f)
	if _, e := decoder.Token(); e == io.EOF {
		slog.Warn("JSONデータが空です")
		// jsonが空の場合終了
		return data
	}

	if err := decoder.Decode(&data); err != nil {
		slog.Error("JSONデータの読み込みに失敗: %v", err)
	}

	slog.Info("JSONデータを読み込みました", data)

	return data
}
