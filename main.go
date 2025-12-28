package main

import (
	"fmt"
	"scheduler/internal/config"
	"scheduler/internal/localdb"
)

func testCallback() {
	fmt.Println("Job done")
}

func runJobTest() {
	opts := config.GetOptions()
	dbCon, err := localdb.NewLocalDb(opts)
	if err != nil {
		fmt.Println("Got error when init dbcon", err)
		return
	}
	dbCon.InitializeSchema()
	sjt := localdb.NewSchedulerJobTable(dbCon)

	// sjob := localdb.NewReminderJob(
	// 	"Test JOB", 0, 0, 0, 1, 5,
	// )
	// sjt.CreateSchedulerJob(sjob)

	schedulerJob, err := sjt.GetAllJob()

	if err != nil {
		fmt.Println("Got error when get all", err)
		return
	}

	fmt.Println(schedulerJob)

	return

	// schedulerHandler := scheduler.NewSchedulerHandler()

	// job := scheduler.NewDailyJob(
	// 	"TestJob", 0, 49, testCallback,
	// )

	// schedulerHandler.AddJob(job)

	// go schedulerHandler.Run()

	// schedulerHandler.Stop(time.Second * 100)
}

func main() {
	runJobTest()

	// opts := config.GetOptions()

	// httpService := servicehandler.NewHttpService(opts)
	// serviceHandler := servicehandler.NewServiceHandler(httpService)

	// serviceHandler.RunService()
}
