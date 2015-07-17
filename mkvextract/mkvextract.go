package mkvextract

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("mkvextract")

const mkvinfoTitlePrefix = "| + Title: "
const mkvinfoDurationPrefix = "| + Duration:"

type MkvInfo struct {
	FileName string
	Title    string
	Chapters []*Chapter
}
type Chapter struct {
	Name  string
	Start time.Duration
	End   time.Duration
}

type TitleNotFound struct{}

func (TitleNotFound) Error() string {
	return "Title not found"
}

func ExtractMetadata(fileName string) (info MkvInfo, err error) {
	info = MkvInfo{FileName: fileName, Chapters: make([]*Chapter, 0, 2)}
	var titleDuration time.Duration
	if info.Title, titleDuration, err = extractTitle(fileName); err != nil {
		return
	}
	err = extractChapter(fileName, &info.Chapters)
	for i, chapter := range info.Chapters {
		if i == len(info.Chapters)-1 {
			chapter.End = titleDuration
		} else {
			chapter.End = info.Chapters[i+1].Start
		}
		log.Debug("%s Found chapter: %s", fileName, chapter)
	}
	return
}

func extractTitle(fileName string) (title string, duration time.Duration, err error) {
	info := exec.Command("mkvinfo", "--output-charset", "UTF-8", fileName)
	//Force default Language to find title field
	info.Env = append([]string{"LANG=C"}, os.Environ()...)
	cmdReader, err := info.StdoutPipe()
	infoReader := bufio.NewReader(cmdReader)
	if err != nil {
		return
	}
	err = info.Start()
	//Search for mkvinfoTitlePrefix in Output
	found := false
	for !found {
		line, err := infoReader.ReadString('\n')
		if err != nil {
			break
		} else if strings.HasPrefix(line, mkvinfoTitlePrefix) {
			//remove prefix
			title = strings.TrimPrefix(line, mkvinfoTitlePrefix)
			title = strings.TrimSuffix(title, "\n")
			found = true
		} else if strings.HasPrefix(line, mkvinfoDurationPrefix) {
			//remove prefix
			durationStr := strings.Split(line, "(")[1]
			durationStr = strings.TrimSuffix(durationStr, ")\n")
			duration, err = parseDuration(durationStr)
			if err != nil {
				break
			}
		}
	}
	err = info.Wait()
	if err != nil {
		return
	}
	if !found {
		err = TitleNotFound{}
	}
	return
}

func extractChapter(fileName string, chapters *[]*Chapter) (err error) {
	info := exec.Command("mkvextract", "chapters", "-s", "--output-charset", " UTF-8", fileName)
	cmdReader, err := info.StdoutPipe()
	infoReader := bufio.NewReader(cmdReader)
	if err != nil {
		return
	}
	err = info.Start()

	//a chapter consits of two lines, frist timecode then a name
	for {
		var lineTimeCode, lineTitle string
		var duration time.Duration
		lineTimeCode, err = infoReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
		lineTitle, err = infoReader.ReadString('\n')
		if err != nil {
			return
		}
		duration, err = parseDuration(extractField(lineTimeCode))
		if err != nil {
			return
		}
		*chapters = append(*chapters, &Chapter{Name: extractField(lineTitle), Start: duration})

	}
	err = info.Wait()
	if err != nil {
		return
	}
	return
}

func extractField(input string) string {
	return strings.TrimSuffix(strings.Join(strings.Split(input, "=")[1:], ""), "\n")
}

func parseDuration(input string) (length time.Duration, err error) {
	result := strings.Split(input, ".")
	ms, err := strconv.Atoi(result[1])
	if err != nil {
		return
	}
	result = strings.Split(result[0], ":")

	h, err := strconv.Atoi(result[0])
	if err != nil {
		return
	}
	m, err := strconv.Atoi(result[1])
	if err != nil {
		return
	}

	s, err := strconv.Atoi(result[2])
	if err != nil {
		return
	}
	length = time.Duration(h)*time.Hour +
		time.Duration(m)*time.Minute +
		time.Duration(s)*time.Second +
		time.Duration(ms)*time.Millisecond

	return
}
