package main

import (
	"regexp"

	"github.com/ajvb/kala/client"
	"github.com/ajvb/kala/job"
)

const format = "^[a-z0-9]{8}-[a-z0-9]{4}-[1-5][a-z0-9]{3}-[a-z0-9]{4}-[a-z0-9]{12}$"

// KalaClientEbook 定义kalaClient
type KalaClientEbook struct {
	Client *client.KalaClient
}

//NewKalaClient 根据kconf中的配置参数来启动kala client
func NewKalaClient() *KalaClientEbook {
	c := client.New(kconf.URLBase)

	return &KalaClientEbook{
		Client: c,
	}
}

//NewKalaClientWithURLBase 根据输入的 kala_url_base来启动kala client
func NewKalaClientWithURLBase(KalaURLBase string) *KalaClientEbook {
	c := client.New(KalaURLBase)

	return &KalaClientEbook{
		Client: c,
	}
}

// BuildKalaJobInfo 根据提供的 JobName,JobCMD,JobSchedule来构建任务结构体; JobSchedule例子: "R0/2020-02-13T15:25:16/"
func BuildKalaJobInfo(JobName string, JobCMD string, JobSchedule string) *job.Job {
	body := &job.Job{
		Schedule: JobSchedule,
		Name:     JobName,
		Command:  JobCMD,
	}
	return body
}

//VerifyJobID 检测JobID是否为正确的uuidV4格式
func (kc *KalaClientEbook) VerifyJobID(JobID string) bool {
	re := regexp.MustCompile(format)
	return re.MatchString(JobID)
}

//CreateJob return id of job
func (kc *KalaClientEbook) CreateJob(body *job.Job) (string, error) {
	/*
			c := New("http://127.0.0.1:8000")
		body := &job.Job{
			Schedule: "R0/2020-02-13T15:25:16/", //在 2020-02-13 15:25:16分启动任务，只运行一次
			Name:	  "test_job",
			Command:  "bash -c 'date'",
		}
		id, err := c.CreateJob(body)
	*/
	id, err := kc.Client.CreateJob(body)
	return id, err
}

//StartJob 根据输入的 id编号来启动 任务
func (kc *KalaClientEbook) StartJob(id string) (bool, error) {
	//id := "93b65499-b211-49ce-57e0-19e735cc5abd"
	ok, err := kc.Client.StartJob(id)
	return ok, err
}

//DeleteAllJobs 删除服务器上面所有的任务
func (kc *KalaClientEbook) DeleteAllJobs() (bool, error) {
	ok, err := kc.Client.DeleteAllJobs()
	return ok, err
}

//DeleteJob 删除id所指定的任务
func (kc *KalaClientEbook) DeleteJob(id string) (bool, error) {
	// id := "93b65499-b211-49ce-57e0-19e735cc5abd"
	ok, err := kc.Client.DeleteJob(id)
	return ok, err
}

//DisableJob 禁言id所指定的任务
func (kc *KalaClientEbook) DisableJob(id string) (bool, error) {
	// id := "93b65499-b211-49ce-57e0-19e735cc5abd"
	ok, err := kc.Client.DisableJob(id)
	return ok, err
}

//EnableJob 允许id执行所指定
func (kc *KalaClientEbook) EnableJob(id string) (bool, error) {
	//id := "93b65499-b211-49ce-57e0-19e735cc5abd"
	ok, err := kc.Client.EnableJob(id)
	return ok, err
}

// GetAllJobs 列举出服务器上面所有的 任务
func (kc *KalaClientEbook) GetAllJobs() (map[string]*job.Job, error) {
	jobs, err := kc.Client.GetAllJobs()
	return jobs, err
}

//GetJob 根据id获取指定任务信息
func (kc *KalaClientEbook) GetJob(id string) (*job.Job, error) {
	//id := "93b65499-b211-49ce-57e0-19e735cc5abd"
	job, err := kc.Client.GetJob(id)
	return job, err
}

//GetJobStats 获取id所指定任务的状态信息
func (kc *KalaClientEbook) GetJobStats(id string) ([]*job.JobStat, error) {
	//id := "93b65499-b211-49ce-57e0-19e735cc5abd"
	stats, err := kc.Client.GetJobStats(id)
	return stats, err
}

//GetKalaStats 获取kala服务器的状态信息
func (kc *KalaClientEbook) GetKalaStats() (*job.KalaStats, error) {
	stats, err := kc.Client.GetKalaStats()
	return stats, err
}
