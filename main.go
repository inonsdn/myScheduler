package main

import (
	"fmt"
	"scheduler/internal/config"
	"scheduler/internal/localdb"
	servicehandler "scheduler/internal/serviceHandler"
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

func runService() {
	opts := config.GetOptions(config.SetHost(), config.SetPort())

	// localDb, err := localdb.NewLocalDb(opts)
	// if err != nil {
	// 	fmt.Println("Found error when connect to local db", err)
	// 	return
	// }
	// if err := localDb.InitializeSchema(); err != nil {
	// 	fmt.Println("Found error when init to local db", err)
	// 	return
	// }
	msgProcessor := servicehandler.NewMessageProcessor(nil)
	lineService := servicehandler.NewLineService(opts, msgProcessor)

	service := servicehandler.NewHttpService(opts.GetServerOptions(), lineService.InitLineRoute())

	handler := servicehandler.NewServiceHandler(service)

	handler.RunService()
}

func main() {

	runService()

	// runJobTest()

	// opts := config.GetOptions()

	// httpService := servicehandler.NewHttpService(opts)
	// serviceHandler := servicehandler.NewServiceHandler(httpService)

	// serviceHandler.RunService()
}
