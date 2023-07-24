package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Item サンプルのデータ構造
type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var items []Item

func init() {
	// サンプルデータの初期化
	items = append(items, Item{ID: 1, Name: "Item 1"})
	items = append(items, Item{ID: 2, Name: "Item 2"})
}

func main() {
	// ルーティングを設定
	http.HandleFunc("/health", healthCheckHandler)
	http.HandleFunc("/items", logRequest(getItemsHandler))
	http.HandleFunc("/items/", logRequest(getItemHandler))

	// サーバーを開始
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// ハンドラーを実行
		next(w, r)

		// リクエストログを出力
		log.Printf(
			"[%s] %s %s %s",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			time.Since(start),
		)
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// ヘルスチェック用のレスポンス
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func getItemsHandler(w http.ResponseWriter, r *http.Request) {
	// データをJSON形式でレスポンス
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func getItemHandler(w http.ResponseWriter, r *http.Request) {
	// パスパラメータからIDを取得
	idStr := r.URL.Path[len("/items/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	// IDに該当するアイテムを検索
	var item *Item
	for _, i := range items {
		if i.ID == id {
			item = &i
			break
		}
	}

	// アイテムが見つからない場合
	if item == nil {
		http.NotFound(w, r)
		return
	}

	// アイテムをJSON形式でレスポンス
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}
