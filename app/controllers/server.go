package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"app/app/models"
	"app/config"
)

// リクエストに対応したHTMLを生成
func generateHTML(w http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	// ファイルパスを格納したスライスを作成
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("app/views/templates/%s.html", file))
	}
	// HTMLテンプレートの生成 & 値の埋め込み
	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(w, "layout", data)
}

// ユーザー認証 (各ページ遷移時) リクエストからCookieを取得し、取得したUUIDがセッションに存在しない(ログインしていない)場合エラーを返し、存在する(ログインしている)場合UUIDに一致したセッションを返す
// session -> getSessionByCookie、 checkLogin (関数名を分かりやすいよう変更することを検討)
func session(w http.ResponseWriter, r *http.Request) (sess models.Session, err error) {
	cookie, err := r.Cookie("_cookie")
	if err == nil {
		sess = models.Session{UUID: cookie.Value}
		if ok, _ := sess.CheckSession(); !ok {
			err = fmt.Errorf("invalid session")
		}
	}
	return sess, err
}

// 正規表現によるURLのバリデーションを設定 (正規表現オブジェクトを作成)
var validPath = regexp.MustCompile("^/todos/(edit|update|delete)/([0-9]+)$")

// リクエストされたURLを解析 (ハンドラーを返却)
func parseURL(fn func(http.ResponseWriter, *http.Request, int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 対象とするURL -> 例) /todos/edit/1

		// validPathとマッチしたURLをスライスとして取得
		q := validPath.FindStringSubmatch(r.URL.Path)
		if q == nil {
			http.NotFound(w, r)
			return
		}
		// URLの最後のパスをint型として取得 (int型でない場合はエラー)
		qi, err := strconv.Atoi(q[2])
		if err != nil {
			http.NotFound(w, r)
		}
		// 引数で受け取った関数を実行 (メソッドチェーン)
		fn(w, r, qi)
	}
}

// HTTPサーバーの起動
func StartMainServer() error {
	// '/static/'から静的ファイルを参照できるようにパスを変更
	files := http.FileServer(http.Dir(config.Config.Static))
	http.Handle("/static/", http.StripPrefix("/static/", files))

	// ルーティング (URLに対応したハンドラーをを登録)
	http.HandleFunc("/", top)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/authenticate", authenticateUser)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/todos", index)
	http.HandleFunc("/todos/new", todoNew)
	http.HandleFunc("/todos/save", todoSave)
	http.HandleFunc("/todos/edit/", parseURL(todoEdit))     // ハンドラー関数をチェーン
	http.HandleFunc("/todos/update/", parseURL(todoUpdate)) // ハンドラー関数をチェーン
	http.HandleFunc("/todos/delete/", parseURL(todoDelete)) // ハンドラー関数をチェーン

	port := os.Getenv("PORT") // 環境変数の読み込み
	return http.ListenAndServe(":"+port, nil)
	// return http.ListenAndServe(":"+config.Config.Port, nil)
}
