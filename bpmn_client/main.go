package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/go-resty/resty"
)

const (
	engineUrl = "http://localhost:8080/engine-rest"
)

type ExternalTask struct {
	Id string
}

var myWorkerId string
var chDone chan int

func startRemoteProcess(key string, n int) {
	client := resty.New()
	for i := 0; i < n; i++ {
		//resp, err :=
		client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Accept", "application/json").
			SetBody("{}").
			Post(engineUrl + "/process-definition/key/" + key + "/start")
		//fmt.Printf("%v, %v\n", resp, err)
		time.Sleep(10 * time.Millisecond)
	}
	chDone <- 0
}

func fetchTasks(client *resty.Client, topic string) ([]ExternalTask, error) {
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(`{
		"workerId": "` + myWorkerId + topic + `",
		"maxTasks": 1,
		"usePriority": true,
		
		"topics": [
			{
				"topicName": "` + topic + `",
				"lockDuration": 2000
			}
		]
		}`).
		Post(engineUrl + "/external-task/fetchAndLock/")
	if err != nil {
		fmt.Printf("%v\n", resp)
		return nil, err
	}

	var tasks []ExternalTask
	err = json.Unmarshal(resp.Body(), &tasks)
	return tasks, err
}

func completeWorker(client *resty.Client, topic string, ch chan ExternalTask, n int) {
	count := 0
	var startTime time.Time
	started := false
	for {
		task := <-ch
		if !started {
			startTime = time.Now()
			started = true
		}
		//_, err :=
		client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Accept", "application/json").
			SetBody(`{
				"workerId": "` + myWorkerId + topic + `"
			}`).
			Post(engineUrl + "/external-task/" + task.Id + "/complete")
		//fmt.Printf("Complete: %s err=%v\n", task.Id, err)
		count++
		if count%100 == 0 {
			fmt.Printf("%d %s tasks in %s\n", count, topic, time.Since(startTime).String())
		}
		if count >= n {
			break
		}
	}
	duration := time.Since(startTime)
	fmt.Printf("\nElapsed time=%s\n", duration.String())
	chDone <- 1
}

func fetchWorker(client *resty.Client, topic string, ch chan ExternalTask) {
	const startBackOff = 100 * time.Millisecond
	const maxBackOff = 400 * time.Millisecond
	backOff := startBackOff
	for {
		tasks, err := fetchTasks(client, topic)
		if err != nil {
			fmt.Printf("Fetch error: %v\n", err)
			break
		}
		for _, task := range tasks {
			ch <- task
		}
		if len(tasks) > 0 {
			backOff = startBackOff
		} else {
			time.Sleep(backOff)
			backOff *= 2
			if backOff > maxBackOff {
				backOff = maxBackOff
			}
		}
	}
}

func topicWorker(topic string, n int) {
	client := resty.New()
	ch := make(chan ExternalTask)
	go fetchWorker(client, topic, ch)
	go completeWorker(client, topic, ch, n)
}

func main() {
	n := 100
	args := os.Args[1:]
	if len(args) > 0 {
		n, _ = strconv.Atoi(args[0])
	}
	myWorkerId = fmt.Sprintf("worker-%d", rand.Int63())
	chDone = make(chan int)

	// start process instances
	go startRemoteProcess("Proc_Order_Test", n)
	// complete external tasks
	go topicWorker("ReserveGoods", n)
	// complete external tasks
	go topicWorker("Charge", n)

	// wait for workers
	_ = <-chDone
	_ = <-chDone
	_ = <-chDone
}
