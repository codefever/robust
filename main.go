package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"

	"github.com/codefever/robust/subprocess"
)

var (
	endpoints    = flag.String("endpoints", "http://127.0.0.1:2379", "Endpoints for etcd cluster.")
	dialTimeout  = flag.Int("dialTimeout", 2, "Timeout seconds for connecting.")
	electionName = flag.String("electionName", "/robust", "Path name in etcd.")
	electionTTL  = flag.Int("electionTTL", 5, "TTL for elections")
	command      = flag.String("command", "", "Command to be run under elections.")
	selfName     = flag.String("name", "", "Name, which would be registered into etcd.")
)

// Return false if it's stopped.
func runCampaign(cli *clientv3.Client, id string, quitc <-chan struct{}) bool {
	// create elections
	sess, err := concurrency.NewSession(cli, concurrency.WithTTL(*electionTTL))
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()
	election := concurrency.NewElection(sess, *electionName)

	// run for new campaigns
	err = election.Campaign(context.Background(), id)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Become leader!")

	_, errc, cancel := subprocess.RunCommand(*command)
	select {
	case err := <-errc:
		if err != nil {
			log.Print("Command exits unexpectedly: ", err)
		}
	case <-quitc:
		log.Print("Quit!")
		cancel()
		<-errc
		return false
	case <-sess.Done():
		log.Print("Session expired.")
		cancel()
		<-errc
	}
	return true
}

func main() {
	flag.Parse()

	if *command == "" {
		log.Fatal("Not specified command.")
	}
	log.Printf("CMD: %q", *command)

	endpointList := strings.Split(*endpoints, ",")
	log.Printf("Endpoints: %v", endpointList)

	// create etcd clients
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpointList,
		DialTimeout: time.Duration(*dialTimeout) * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	signalc := make(chan os.Signal, 1)
	signal.Notify(signalc, syscall.SIGINT, syscall.SIGTERM)
	quitc := make(chan struct{})
	go func() {
		s := <-signalc
		log.Print("Got signal: ", s)
		close(quitc)
	}()

	var round uint64
	for {
		log.Printf("Start campaign for Round[%v].", round)
		if !runCampaign(cli, *selfName, quitc) {
			break
		}
		round++
		time.Sleep(time.Second * 1)
	}
}
