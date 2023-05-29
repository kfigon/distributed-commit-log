package main

import (
	"commit-log/appendlog"
	"commit-log/web"
)

func main() {
	const port = 8080

	theLog := appendlog.NewAppendLog()
	web.Run(web.NewServer(theLog), port)
}
