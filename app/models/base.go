package models

import (
	"app/config"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/lib/pq"
	// _ "github.com/lib/pq"
	// _ "github.com/mattn/go-sqlite3"
)

// 変数の宣言
var Db *sql.DB
var err error

/*
const (
	tableNameUser    = "users"
	tableNameTodo    = "todos"
	tableNameSession = "sessions"
)
*/

func init() {
	// Heroku用設定
	url := os.Getenv("DATABASE_URL")                        // 環境変数の読み込み
	connection, _ := pq.ParseURL(url)                       // DBコネクション(文字列)の取得
	connection += "sslmode=require"                         // 文字列を連結
	Db, err = sql.Open(config.Config.SQLDriver, connection) // データベースに接続
	if err != nil {
		log.Fatalln(err)
	}

	/*
		// データベースに接続
		Db, err = sql.Open(config.Config.SQLDriver, config.Config.DbName)
		if err != nil {
			log.Fatalln(err)
		}

		// 各テーブルの作成

		// Usersテーブル
		cmdU := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			uuid STRING NOT NULL UNIQUE,
			name STRING,
			email STRING,
			password STRING,
			created_at DATETIME
		)`, tableNameUser)
		Db.Exec(cmdU)

		// Todosテーブル
		cmdT := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT,
			user_id INTEGER,
			created_at DATETIME
		)`, tableNameTodo)
		Db.Exec(cmdT)

		// Sessionsテーブル
		cmdS := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			uuid STRING NOT NULL UNIQUE,
			email STRING,
			user_id INTEGER,
			created_at DATETIME
		)`, tableNameSession)
		Db.Exec(cmdS)
	*/
}

// UUIDを生成
func createUUID() (uuidobj uuid.UUID) {
	uuidobj, _ = uuid.NewUUID()
	return uuidobj
}

// パスワードのハッシュ化
func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return cryptext
}
