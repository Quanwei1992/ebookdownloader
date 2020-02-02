package main

import (
	"fmt"
	"github.com/ajvb/kala/job"
	//"github.com/ajvb/kala/client"
	"time"
	//"github.com/stretchr/testify/assert"
	"testing"
)

func test1(t *testing.T) {
	c := NewKalaClientWithURLBase("http://192.168.13.100:8091")
	// our job just run once,after 5 minute
	schedule := fmt.Sprintf("R0/%s/", time.Now().Add(time.Minute*5).Format(time.RFC3339))
	fmt.Println(schedule)
	body := &job.Job{
		Schedule: schedule,
		Name:     "test_job",
		Command:  "/home/pi/gowork/src/github.com/sndnvaps/ebookdownloader/ebookdownloader_cli --bookid=91_91911 --txt --mobi",
	}
	id, err := c.CreateJob(body)
	fmt.Println("Job Created: ", id)
	if err != nil {
		fmt.Println(err.Error())
	}
}
