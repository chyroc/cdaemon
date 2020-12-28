package cdaemon

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/takama/daemon"
)

func New(name, desc string, kind daemon.Kind) (*Daemon, error) {
	r := new(Daemon)
	r.Name = name
	r.Desc = desc
	var err error
	r.daemon, err = daemon.New(name, desc, kind)
	if err != nil {
		return nil, fmt.Errorf("init daemon failed: %s", err)
	}

	return r, nil
}

type Daemon struct {
	Name   string
	Desc   string
	daemon daemon.Daemon
	runner func()
}

func (r *Daemon) AddRunner(f func()) {
	r.runner = f
}

func (r *Daemon) Run() (string, error) {
	usage := fmt.Sprintf("Usage: %s install | remove | start | stop | status", r.Name)
	// If received any kind of command, do it
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return r.daemon.Install()
		case "remove":
			return r.daemon.Remove()
		case "start":
			return r.daemon.Start()
		case "stop":
			return r.daemon.Stop()
		case "status":
			return r.daemon.Status()
		default:
			return usage, nil
		}
	}
	fmt.Println(usage)

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Create a new cron manager
	go r.runner()

	// Waiting for interrupt by system signal
	killSignal := <-interrupt
	fmt.Println("Got signal:", killSignal)
	return "Service exited", nil
}
