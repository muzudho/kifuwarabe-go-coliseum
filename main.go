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

	cne "github.com/muzudho/gtp-engine-to-nngs/entities" // CoNnector
	cnui "github.com/muzudho/gtp-engine-to-nngs/ui"
	g "github.com/muzudho/kifuwarabe-go-coliseum/global"
	"github.com/muzudho/kifuwarabe-go-coliseum/ui"
	kwe "github.com/muzudho/kifuwarabe-gtp/entities"
	kwui "github.com/muzudho/kifuwarabe-gtp/ui"
	kwu "github.com/muzudho/kifuwarabe-gtp/usecases"
)

func main() {
	// Default working directory
	dwd, err := os.Getwd()
	if err != nil {
		// ここでは、ログはまだ設定できてない
		panic(fmt.Errorf("...Coliseum... err=%s", err))
	}
	fmt.Printf("...Coliseum... DefaultWorkingDirectory=%s\n", dwd)

	// コマンドライン引数
	wd := flag.String("workdir", dwd, "Working directory path.")
	wdw := flag.String("workdirw", dwd, "Working directory path of White phase.")
	wdb := flag.String("workdirb", dwd, "Working directory path of Black phase.")
	flag.Parse()
	fmt.Printf("...Coliseum... flag.Args()=%s\n", flag.Args())
	fmt.Printf("...Coliseum... WorkingDirectory=%s\n", *wd)
	fmt.Printf("...Coliseum... WorkingDirectoryW=%s\n", *wdw)
	fmt.Printf("...Coliseum... WorkingDirectoryB=%s\n", *wdb)
	coliseumConfPath := filepath.Join(*wdw, "input/coliseum.conf.toml")
	connectorConfPathW := filepath.Join(*wdw, "input/connector.conf.toml")
	engineConfPathW := filepath.Join(*wdw, "input/engine.conf.toml")
	connectorConfPathB := filepath.Join(*wdb, "input/connector.conf.toml")
	engineConfPathB := filepath.Join(*wdb, "input/engine.conf.toml")
	fmt.Printf("...Coliseum... connectorConfPathW=%s\n", connectorConfPathW)
	fmt.Printf("...Coliseum... engineConfPathW=%s\n", engineConfPathW)
	fmt.Printf("...Coliseum... connectorConfPathB=%s\n", connectorConfPathB)
	fmt.Printf("...Coliseum... engineConfPathB=%s\n", engineConfPathB)

	// ロガーの作成。
	g.G.Log = *kwu.NewLogger(
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
	g.G.Chat = *kwu.NewChatter(g.G.Log)
	g.G.StderrChat = *kwu.NewStderrChatter(g.G.Log)

	g.G.Chat.Trace("...Coliseum... Start\n")

	// 設定ファイル読込
	// coliseumConfig, err := ui.LoadColiseumConf(coliseumConfPath)
	_, err = ui.LoadColiseumConf(coliseumConfPath)
	if err != nil {
		panic(g.G.Log.Fatal(fmt.Sprintf("...Engine... coliseumConfPath=[%s] err=[%s]", coliseumConfPath, err)))
	}

	// 設定ファイル読込
	engineConfW, err := kwui.LoadEngineConf(engineConfPathW)
	if err != nil {
		panic(g.G.Chat.Fatal("...Coliseum... engineConfPathW=[%s] err=[%s]", engineConfPathW, err))
	}

	connectorConfW, err := cnui.LoadConnectorConf(connectorConfPathW)
	if err != nil {
		panic(g.G.Chat.Fatal("...Coliseum... connectorConfPathW=[%s] err=[%s]", connectorConfPathW, err))
	}

	engineConfB, err := kwui.LoadEngineConf(engineConfPathB)
	if err != nil {
		panic(g.G.Chat.Fatal("...Coliseum... engineConfPathB=[%s] err=[%s]", engineConfPathB, err))
	}

	connectorConfB, err := cnui.LoadConnectorConf(connectorConfPathB)
	if err != nil {
		panic(g.G.Chat.Fatal("...Coliseum... connectorConfPathB=[%s] err=[%s]", connectorConfPathB, err))
	}

	// 思考エンジンを起動
	wg := sync.WaitGroup{}
	wg.Add(2)
	g.G.Chat.Trace("...Coliseum... Start cmdW\n")
	go startEngine(engineConfW, connectorConfW, &wg)
	g.G.Chat.Trace("...Coliseum... Sleep 4 seconds\n")
	time.Sleep(time.Second * 4)
	g.G.Chat.Trace("...Coliseum... Start cmdB\n")
	go startEngine(engineConfB, connectorConfB, &wg)
	g.G.Chat.Trace("...Coliseum... WaitGropu wait\n")
	wg.Wait()

	g.G.Chat.Trace("...Coliseum... End\n")
}

// コネクターを起動
func startEngine(engineConf *kwe.EngineConf, connectorConf *cne.ConnectorConf, wg *sync.WaitGroup) {
	defer wg.Done()

	parameters := strings.Split(connectorConf.User.EngineCommandOption, " ")
	parametersString := strings.Join(parameters, " ")
	parametersString = strings.TrimRight(parametersString, " ")
	g.G.Chat.Trace("...Coliseum... (^q^) EngineCommand=[%s] ArgumentList=[%s]\n", connectorConf.User.EngineCommand, parametersString)
	cmd := exec.Command(connectorConf.User.EngineCommand, parameters...)
	err := cmd.Start()
	if err != nil {
		panic(g.G.Chat.Fatal(fmt.Sprintf("...Coliseum... cmd.Run() --> [%s]", err)))
	}

	g.G.Chat.Trace("...Coliseum... Cmd Wait\n")
	cmd.Wait()
	g.G.Chat.Trace("...Coliseum... Cmd Waited\n")
}
