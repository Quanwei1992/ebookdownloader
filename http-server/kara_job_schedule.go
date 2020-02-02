package main

import (
	"github.com/ajvb/kala/client"
	"github.com/ajvb/kala/job"
)

type KalaClientEbook struct {
	Client *client.KalaClient
}

/*
func NewKalaClient() *KalaClientEbook {
	c := client.New(kconf.URL_BASE)

	return &KalaClientEbook{
		Client: c,
	}
}
*/
func NewKalaClientWithURLBase(KALA_URL_BASE string) *KalaClientEbook {
	c := client.New(KALA_URL_BASE)

	return &KalaClientEbook{
		Client: c,
	}
}

//return id of job
func (kc *KalaClientEbook) CreateJob(body *job.Job) (string, error) {
	/*
			c := New("http://127.0.0.1:8000")
		body := &job.Job{
			Schedule: "R2/2015-06-04T19:25:16.828696-07:00/PT10S",
			Name:	  "test_job",
			Command:  "bash -c 'date'",
		}
		id, err := c.CreateJob(body)
	*/
	id, err := kc.Client.CreateJob(body)
	return id, err
}

func (kc *KalaClientEbook) StartJob(id string) (bool, error) {
	//id := "93b65499-b211-49ce-57e0-19e735cc5abd"
	ok, err := kc.Client.StartJob(id)
	return ok, err
}

func (kc *KalaClientEbook) DeleteAllJobs() (bool, error) {
	ok, err := kc.Client.DeleteAllJobs()
	return ok, err
}

func (kc *KalaClientEbook) DeleteJob(id string) (bool, error) {
	// id := "93b65499-b211-49ce-57e0-19e735cc5abd"
	ok, err := kc.Client.DeleteJob(id)
	return ok, err
}

func (kc *KalaClientEbook) DisableJob(id string) (bool, error) {
	// id := "93b65499-b211-49ce-57e0-19e735cc5abd"
	ok, err := kc.Client.DisableJob(id)
	return ok, err
}

func (kc *KalaClientEbook) EnableJob(id string) (bool, error) {
	//id := "93b65499-b211-49ce-57e0-19e735cc5abd"
	ok, err := kc.Client.EnableJob(id)
	return ok, err
}

func (kc *KalaClientEbook) GetAllJobs() (map[string]*job.Job, error) {
	jobs, err := kc.Client.GetAllJobs()
	return jobs, err
}

func (kc *KalaClientEbook) GetJob(id string) (*job.Job, error) {
	//id := "93b65499-b211-49ce-57e0-19e735cc5abd"
	job, err := kc.Client.GetJob(id)
	return job, err
}

func (kc *KalaClientEbook) GetJobStats(id string) ([]*job.JobStat, error) {
	//id := "93b65499-b211-49ce-57e0-19e735cc5abd"
	stats, err := kc.Client.GetJobStats(id)
	return stats, err
}

func (kc *KalaClientEbook) GetKalaStats() (*job.KalaStats, error) {
	stats, err := kc.Client.GetKalaStats()
	return stats, err
}
