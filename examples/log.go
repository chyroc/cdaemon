package main

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/takama/daemon"

	"github.com/chyroc/cdaemon"
)

func main() {
	{
		f, err := os.OpenFile("/tmp/log.log2", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}
		logrus.SetOutput(f)
	}

	d, err := cdaemon.New("example", "desc", daemon.UserAgent)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	d.AddRunner(func() {
		for {
			logrus.Infof(time.Now().Format(time.RFC3339))
			time.Sleep(time.Second)
		}
	})
	status, err := d.Run()
	if err != nil {
		fmt.Printf("%s\nErr: %s\n", status, err)
		os.Exit(1)
	}
	fmt.Printf(status)
}
