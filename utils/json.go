package utils

import (
	"encoding/json"
	"fmt"
	"log"
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
		log.Fatalf("ファイル取得失敗: %v", err)
	}
	defer f.Close()

	output, err := json.MarshalIndent(data, "", "\t\t")
	if err != nil {
		log.Fatalf("JSONエンコード失敗: %v", err)
	}

	if _, err := f.Write(output); err != nil {
		log.Fatalf("JSON書き込み失敗: %v", err)
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
		log.Fatalf("ファイル取得失敗: %v", err)
	}
	defer f.Close()

	var data []JobData
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&data); err != nil {
		log.Fatalf("JSONデコード失敗: %v", err)
	}

	fmt.Println(data)

	return data
}
