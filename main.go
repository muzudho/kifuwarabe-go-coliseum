package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	cne "github.com/muzudho/gtp-engine-to-nngs/entities" // CoNnector
	cnui "github.com/muzudho/gtp-engine-to-nngs/ui"
	u "github.com/muzudho/kifuwarabe-go-coliseum/usecases"
	kwe "github.com/muzudho/kifuwarabe-gtp/entities"
	kwui "github.com/muzudho/kifuwarabe-gtp/ui"
	kwu "github.com/muzudho/kifuwarabe-gtp/usecases"
)

func main() {
	// Working directory
	wdir, err := os.Getwd()
	if err != nil {
		// ここでは、ログはまだ設定できてない
		panic(fmt.Sprintf("...Coliseum... wdir=%s", wdir))
	}
	fmt.Printf("...Coliseum... wdir=%s\n", wdir)

	// コマンドライン引数
	workdir := flag.String("workdir", wdir, "Working directory path.")
	workdirw := flag.String("workdirw", wdir, "Working directory path of White phase.")
	workdirb := flag.String("workdirb", wdir, "Working directory path of Black phase.")
	flag.Parse()
	fmt.Printf("...Coliseum... flag.Args()=%s\n", flag.Args())
	fmt.Printf("...Coliseum... workdir=%s\n", *workdir)
	fmt.Printf("...Coliseum... workdirw=%s\n", *workdirw)
	fmt.Printf("...Coliseum... workdirb=%s\n", *workdirb)
	entryConfPathW := filepath.Join(*workdirw, "input/entry.conf.toml")
	engineConfPathW := filepath.Join(*workdirw, "input/engine.conf.toml")
	entryConfPathB := filepath.Join(*workdirb, "input/entry.conf.toml")
	engineConfPathB := filepath.Join(*workdirb, "input/engine.conf.toml")
	fmt.Printf("...Coliseum... entryConfPathW=%s\n", entryConfPathW)
	fmt.Printf("...Coliseum... engineConfPathW=%s\n", engineConfPathW)
	fmt.Printf("...Coliseum... entryConfPathB=%s\n", entryConfPathB)
	fmt.Printf("...Coliseum... engineConfPathB=%s\n", engineConfPathB)

	// ロガーの作成。
	u.G.Log = *kwu.NewLogger(
		filepath.Join(*workdir, "output/trace.log"),
		filepath.Join(*workdir, "output/debug.log"),
		filepath.Join(*workdir, "output/info.log"),
		filepath.Join(*workdir, "output/notice.log"),
		filepath.Join(*workdir, "output/warn.log"),
		filepath.Join(*workdir, "output/error.log"),
		filepath.Join(*workdir, "output/fatal.log"),
		filepath.Join(*workdir, "output/print.log"))

	// 既存のログ・ファイルを削除
	u.G.Log.RemoveAllOldLogs()
	u.G.Log.OpenAllLogs()
	u.G.Log.Trace("...Coliseum Remove all old logs\n")

	// チャッターの作成。 標準出力とロガーを一緒にしただけです。
	u.G.Chat = *kwu.NewChatter(kwu.G.Log)
	u.G.StderrChat = *kwu.NewStderrChatter(kwu.G.Log)

	// fmt.Println("...GE2NNGS... 設定ファイルを読み込んだろ☆（＾～＾）")
	engineConfW, err := kwui.LoadEngineConf(engineConfPathW)
	if err != nil {
		panic(u.G.Chat.Fatal("engineConfPathW=[%s] err=[%s]", engineConfPathW, err))
	}

	entryConfW, err := cnui.LoadEntryConf(entryConfPathW)
	if err != nil {
		panic(u.G.Chat.Fatal("entryConfPathW=[%s] err=[%s]", entryConfPathW, err))
	}

	engineConfB, err := kwui.LoadEngineConf(engineConfPathB)
	if err != nil {
		panic(u.G.Chat.Fatal("engineConfPathB=[%s] err=[%s]", engineConfPathB, err))
	}

	entryConfB, err := cnui.LoadEntryConf(entryConfPathB)
	if err != nil {
		panic(u.G.Chat.Fatal("entryConfPathB=[%s] err=[%s]", entryConfPathB, err))
	}

	// 思考エンジンを起動
	go startEngine(engineConfW, entryConfW, workdirw)
	startEngine(engineConfB, entryConfB, workdirb)

	u.G.Chat.Trace("(^q^) コロシアムを終了するぜ")
	u.G.Log.Trace("...Coliseum... End\n")
	u.G.Log.CloseAllLogs()
}

// コネクターを起動
func startEngine(engineConf *kwe.EngineConf, entryConf *cne.EntryConf, workdir *string) {
	parameters := strings.Split("--workdir "+*workdir+" "+entryConf.User.EngineCommandOption, " ")
	u.G.Chat.Trace("(^q^) GTP対応の思考エンジンを起動するぜ☆ [%s] [%s]", entryConf.User.EngineCommand, strings.Join(parameters, " "))
	cmd := exec.Command(entryConf.User.EngineCommand, parameters...)
	err := cmd.Start()
	if err != nil {
		panic(u.G.Chat.Fatal(err.Error()))
	}
	// cmd.Wait()
}
