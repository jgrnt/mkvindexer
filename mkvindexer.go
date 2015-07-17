//go:generate stringer -type Result
package main

import (
	"fmt"
	"os"

	"github.com/jgrnt/mkvindexer/mkvextract"

	"sort"

	"github.com/codegangsta/cli"
	"github.com/op/go-logging"
)

type ResultMessage struct {
	File   string
	Result Result
	Error  error
}
type Result int

const (
	Indexed Result = iota
	AlreadyIndexed
	Error
)

func (r ResultMessage) String() string {
	return fmt.Sprintf("%10s:%s %v", r.Result, r.File, r.Error)
}

//Sort
type ByName []ResultMessage

func (r ByName) Len() int           { return len(r) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].File < a[j].File }

var log = logging.MustGetLogger("mkvindexer")

func proccessMkv(file string, channel chan ResultMessage) {
	info, err := mkvextract.ExtractMetadata(file)
	if err != nil {
		channel <- ResultMessage{file, Error, err}
		return
	}
	log.Debug("%s", info)
	channel <- ResultMessage{file, Indexed, nil}

}

func coordinateMkvs(fileNames []string) {
	resultChannel := make(chan ResultMessage, len(fileNames))
	results := make([]ResultMessage, len(fileNames))
	for _, fileName := range fileNames {
		go proccessMkv(fileName, resultChannel)
	}
	for i := range fileNames {
		results[i] = <-resultChannel
	}
	sort.Sort(ByName(results))
	for _, res := range results {
		switch res.Result {
		case Indexed:
			log.Info("%s", res)
		case AlreadyIndexed:
			log.Warning("%s", res)
		case Error:
			log.Error("%s", res)
		}

	}
	//Output afterwards
}

func main() {

	logging.SetFormatter(logging.MustStringFormatter(
		"%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{color:reset} %{message}",
	))
	logging.SetLevel(logging.DEBUG, "")

	app := cli.NewApp()
	app.Name = "MkvIndexer"
	app.Action = func(c *cli.Context) {
		coordinateMkvs(c.Args())
	}
	app.Run(os.Args)
}
