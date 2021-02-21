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
	fmt.Printf("...GE2NNGS... flag.Args()=%s\n", flag.Args())
	fmt.Printf("...GE2NNGS... workdir=%s\n", *workdir)
	fmt.Printf("...GE2NNGS... workdirw=%s\n", *workdirw)
	fmt.Printf("...GE2NNGS... workdirb=%s\n", *workdirb)
	entryConfPathW := filepath.Join(*workdirw, "input/entry.conf.toml")
	engineConfPathW := filepath.Join(*workdirw, "input/engine.conf.toml")
	entryConfPathB := filepath.Join(*workdirb, "input/entry.conf.toml")
	engineConfPathB := filepath.Join(*workdirb, "input/engine.conf.toml")
	fmt.Printf("...GE2NNGS... entryConfPathW=%s\n", entryConfPathW)
	fmt.Printf("...GE2NNGS... engineConfPathW=%s\n", engineConfPathW)
	fmt.Printf("...GE2NNGS... entryConfPathB=%s\n", entryConfPathB)
	fmt.Printf("...GE2NNGS... engineConfPathB=%s\n", engineConfPathB)

	// ロガーの作成。
	kwu.G.Log = *kwu.NewLogger(
		filepath.Join(*workdir, "output/trace.log"),
		filepath.Join(*workdir, "output/debug.log"),
		filepath.Join(*workdir, "output/info.log"),
		filepath.Join(*workdir, "output/notice.log"),
		filepath.Join(*workdir, "output/warn.log"),
		filepath.Join(*workdir, "output/error.log"),
		filepath.Join(*workdir, "output/fatal.log"),
		filepath.Join(*workdir, "output/print.log"))

	// 既存のログ・ファイルを削除。エンジンが起動時に行う

	// チャッターの作成。 標準出力とロガーを一緒にしただけです。
	kwu.G.Chat = *kwu.NewChatter(kwu.G.Log)
	kwu.G.StderrChat = *kwu.NewStderrChatter(kwu.G.Log)

	// fmt.Println("...GE2NNGS... 設定ファイルを読み込んだろ☆（＾～＾）")
	engineConfW := kwui.LoadEngineConf(engineConfPathW)
	entryConfW := cnui.LoadEntryConf(entryConfPathW)
	engineConfB := kwui.LoadEngineConf(engineConfPathB)
	entryConfB := cnui.LoadEntryConf(entryConfPathB)

	// 思考エンジンを起動
	go startEngine(engineConfW, entryConfW, workdirw)
	startEngine(engineConfB, entryConfB, workdirb)

	kwu.G.Chat.Trace("(^q^) コロシアムを終了するぜ")
}

// コネクターを起動
func startEngine(engineConf kwe.EngineConf, entryConf cne.EntryConf, workdir *string) {
	parameters := strings.Split("--workdir "+*workdir+" "+entryConf.User.EngineCommandOption, " ")
	kwu.G.Chat.Trace("(^q^) GTP対応の思考エンジンを起動するぜ☆ [%s] [%s]", entryConf.User.EngineCommand, strings.Join(parameters, " "))
	cmd := exec.Command(entryConf.User.EngineCommand, parameters...)
	err := cmd.Start()
	if err != nil {
		panic(kwu.G.Chat.Fatal(err.Error()))
	}
	// cmd.Wait()
}
