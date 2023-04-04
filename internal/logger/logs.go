package logger

import (
	"log"
	"os"
)

var (
	outfile, _ = os.OpenFile("../internal/logs/socialNetwork.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	LogFile    = log.New(outfile, "", 0)
)

func ForError(err error) {
	LogFile.Println(err)
}
