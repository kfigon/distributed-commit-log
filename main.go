package main

import (
	"fmt"
	"log"
	"net/http"
	"commit-log/appendlog"
	"commit-log/web"
)

func main() {
	const port = 8080

	theLog := appendlog.NewAppendLog()
	log.Println(http.ListenAndServe(fmt.Sprintf(":%d", port), web.NewServer(theLog)))
}
