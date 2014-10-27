/*
Package jobs provides the tools to have simple job scheduling in your Go programs.

Usage

Based on the excellent cron package at github.com/robfig/cron, the entire package is based on the JobRunner type.
You get a pointer to a JobRunner by calling jobs.NewRunner(), then attach jobs to it, start it - and then the runner
will schedule your jobs as you have specified. Below is a simple example:

	r := jobs.NewRunner()
	r.RunFuncAt(time.Now().Add(time.Duration(time.Second)), func(){
		fmt.Println("My job was run a second after it was added")
	})
	r.Start()

	// You can also add jobs while the runner is running:
	r.RunNamedFuncEvery("Run every 1 second", time.Duration(time.Second), func(){
		fmt.Println("This runs every second")
	})

	// Run for a while
	<-time.After(time.Second*3)

	// If you want to stop the runner temporarily, call stop
	r.Stop()

	// Once you want to free up resources, destroy the runner using Destroy
	r.Destroy()

Instances acquired using NewRunner will only ever run one job at a time, and job taking a long time to run will prevent others from being properly scheduled.
If you require running jobs concurrently, this package supports supplying how many goroutines you want created to run jobs in.
Use the constructor function NewRunnerWithConcurrentExecutors and specify how many workers you want:

	// Creates a new runner with 10 worker goroutines
	r := jobs.NewRunnerWithConcurrentExecutors(10)
*/
package jobs
