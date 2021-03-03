package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	l "github.com/muzudho/go-logger"
	"github.com/muzudho/kifuwarabe-go-coliseum/entities"
	g "github.com/muzudho/kifuwarabe-go-coliseum/global"
	"github.com/muzudho/kifuwarabe-go-coliseum/ui"
)

func main() {
	// Default working directory
	dwd, err := os.Getwd()
	if err != nil {
		// ここでは、ログはまだ設定できてない
		panic(fmt.Errorf("...Coliseum... err=%s", err))
	}
	fmt.Printf("...Coliseum... DefaultWorkingDirectory=%s\n", dwd)

	// コマンドライン引数登録
	wd := flag.String("workdir", dwd, "Working directory path.")
	// 解析
	flag.Parse()

	fmt.Printf("...Coliseum... flag.Args()=%s\n", flag.Args())
	fmt.Printf("...Coliseum... WorkingDirectory=%s\n", *wd)
	coliseumConfPath := filepath.Join(*wd, "input/coliseum.conf.toml")
	fmt.Printf("...Coliseum... coliseumConfPath=%s\n", coliseumConfPath)

	// ロガーの作成。
	g.G.Log = *l.NewLogger(
		filepath.Join(*wd, "output/trace.log"),
		filepath.Join(*wd, "output/debug.log"),
		filepath.Join(*wd, "output/info.log"),
		filepath.Join(*wd, "output/notice.log"),
		filepath.Join(*wd, "output/warn.log"),
		filepath.Join(*wd, "output/error.log"),
		filepath.Join(*wd, "output/fatal.log"),
		filepath.Join(*wd, "output/print.log"))

	// 既存のログ・ファイルを削除
	g.G.Log.RemoveAllOldLogs()

	// ログ・ファイルの開閉
	err = g.G.Log.OpenAllLogs()
	if err != nil {
		// ログ・ファイルを開くのに失敗したのだから、ログ・ファイルへは書き込めません
		panic(fmt.Sprintf("...Coliseum... %s", err))
	}
	defer g.G.Log.CloseAllLogs()

	g.G.Log.Trace("...Coliseum Remove all old logs\n")

	// チャッターの作成。 標準出力とロガーを一緒にしただけです。
	g.G.Chat = *l.NewChatter(g.G.Log)
	g.G.StderrChat = *l.NewStderrChatter(g.G.Log)

	g.G.Chat.Trace("...Coliseum... Start\n")

	// 設定ファイル読込
	coliseumConfig, err := ui.LoadColiseumConf(coliseumConfPath)
	if err != nil {
		panic(g.G.Log.Fatal(fmt.Sprintf("...Engine... coliseumConfPath=[%s] err=[%s]", coliseumConfPath, err)))
	}
	wdw := coliseumConfig.White.Workspace
	wdb := coliseumConfig.Black.Workspace
	fmt.Printf("...Coliseum... WorkingDirectoryW=%s\n", wdw)
	fmt.Printf("...Coliseum... WorkingDirectoryB=%s\n", wdb)
	connectorConfPathW := filepath.Join(wdw, "input/connector.conf.toml")
	engineConfPathW := filepath.Join(wdw, "input/engine.conf.toml")
	connectorConfPathB := filepath.Join(wdb, "input/connector.conf.toml")
	engineConfPathB := filepath.Join(wdb, "input/engine.conf.toml")
	fmt.Printf("...Coliseum... connectorConfPathW=%s\n", connectorConfPathW)
	fmt.Printf("...Coliseum... engineConfPathW=%s\n", engineConfPathW)
	fmt.Printf("...Coliseum... connectorConfPathB=%s\n", connectorConfPathB)
	fmt.Printf("...Coliseum... engineConfPathB=%s\n", engineConfPathB)

	// コネクターを起動
	wg := sync.WaitGroup{}
	wg.Add(2)
	g.G.Chat.Trace("...Coliseum... Start cmdW\n")
	go startConnector(coliseumConfig.White, &wg)
	g.G.Chat.Trace("...Coliseum... Sleep 4 seconds\n")
	time.Sleep(time.Second * 4)
	g.G.Chat.Trace("...Coliseum... Start cmdB\n")
	go startConnector(coliseumConfig.Black, &wg)

	g.G.Chat.Trace("...Coliseum... WaitGropu.wait begin\n")
	wg.Wait()
	g.G.Chat.Trace("...Coliseum... WaitGropu.wait end\n")

	g.G.Chat.Trace("...Coliseum... End\n")
}

// コネクターを起動
func startConnector(colorConf entities.Color, wg *sync.WaitGroup) {
	defer wg.Done()

	p1 := fmt.Sprintf("--workdir %s", colorConf.Workspace)
	parameters := []string{p1}
	parametersString := strings.Join(parameters, " ")
	g.G.Chat.Trace("...Coliseum... (^q^) Exe=[%s] ArgumentList=[%s]\n", colorConf.Connector, parametersString)

	cmd := exec.Command(colorConf.Connector, parameters...)
	err := cmd.Start()
	if err != nil {
		panic(g.G.Chat.Fatal(fmt.Sprintf("...Coliseum... cmd.Run() --> [%s]", err)))
	}

	g.G.Chat.Trace("...Coliseum... Cmd Wait\n")
	cmd.Wait()
	g.G.Chat.Trace("...Coliseum... Cmd Waited\n")
}
