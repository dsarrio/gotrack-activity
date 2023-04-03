package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
)

func normalizeEntry(path string, title string) string {

	if strings.Contains(path, "Visual Studio Code.app") {
		path = "VSCode"
		titleParts := strings.Split(title, "—")
		title = titleParts[len(titleParts)-1]
	}

	if strings.Contains(path, "Google Chrome.app") {
		path = "Browser"
		titleParts := strings.Split(title, "—")
		title = titleParts[0]
	}

	if strings.Contains(path, "zoom.us.app") {
		path = "Zoom"
		if strings.Contains(title, "Zoom Meeting") {
			title = "Meeting"
		} else {
			title = "Settings"
		}
	}

	if strings.Contains(path, "Terminal.app") {
		path = "Terminal"
		titleParts := strings.Split(title, "—")
		title = titleParts[0]
	}

	if m, _ := regexp.MatchString("/Applications/([^/]+).app/", path); m {
		pathParts := strings.Split(path, "/")
		path = strings.TrimSuffix(pathParts[2], ".app")
	}

	return path + "[" + strings.Trim(title, " \t—-_") + "]"
}

func appendLog(key string, elapsed int) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	outputDir := path.Join(dirname, "gotrack")
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile(path.Join(outputDir, "activity.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("%s,%d\n", key, elapsed))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	history := make(map[string]int)
	defer println(history)

	var refTime time.Time = time.Now()
	var refTitle string
	var refPath string

	for true {
		robotgo.Sleep(1)

		activeWindowPID := robotgo.GetPID()
		activeWindowTitle := robotgo.GetTitle()
		activeWindowPath, err := robotgo.FindPath(activeWindowPID)
		if err != nil {
			//println(err)
			continue
		}

		// log.Printf("Active window [path: %s]  [title: %s]\n", activeWindowPath, activeWindowTitle)

		// active window changed
		if refPath != activeWindowPath || refTitle != activeWindowTitle {
			currTime := time.Now()

			// we were actively tracking previous window
			if len(refPath) > 0 {
				elapsed := int(currTime.Sub(refTime).Seconds())
				if elapsed > 3 {
					historyKey := normalizeEntry(refPath, refTitle)
					history[historyKey] += elapsed
					appendLog(historyKey, elapsed)
					log.Printf("Recorded %v seconds to %s\n", elapsed, historyKey)
				}
			}

			// track new window
			refTime = time.Now()
			refTitle = activeWindowTitle
			refPath = activeWindowPath
		}
	}
}
