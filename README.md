### Job scheduler ###
[![GoDoc](http://godoc.org/github.com/kastermester/jobs?status.png)](http://godoc.org/github.com/kastermester/jobs)

This is a simple job scheduling library for go. It was greatly inspired by [robfig/cron](https://github.com/robfig/cron).

However, I needed a library that could do more than simply act as a built-in cron, a more general job scheduler, with support for concurrent running jobs.

Right now, all the building blocks for creating truly scheduled jobs are in place (implement the Schedule interface for custom scheduling). However none of that is in the library yet. The library is focused around getting the concurrency primitives right, ensuring that it works no matter how many goroutines are calling it.

A simple example may look like this:

```go
package main

import (
	"fmt"
	"jobs"
	"time"
)

func main() {
	r := jobs.NewRunner()

	r.RunFuncAt(time.Now().Add(time.Second*10), func() {
		fmt.Println("My job was run successfully!")
	})

	r.RunNamedFuncEvery("My job", time.Duration(time.Second), func() {
		fmt.Println("My job runs every second!")
	})

	r.Start()

	<-time.After(time.Second * 11)
}
```

This is a very early version, but every method of the API have been documented and as such, godoc should be able to provide more documentation to work from.
