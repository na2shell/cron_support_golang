package main


import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"context"
	"time"
)

func main(){
	timeout := flag.Int("t", 10, "timeout time")
	status := flag.String("o", "./", "job status log file path")
	logfile := flag.String("log", "./", "job itself log file path")
	command := flag.String("cmd", "ls -la", "exec command")
	name := flag.String("name","hoge-job","in status file job name")
	flag.Parse()
    // fmt.Println(*timeout)
    fmt.Println(*status)
    fmt.Println(*logfile)
    // fmt.Println(*command)

	cmstr := *command
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout*100) * time.Millisecond)

	var res bool

    defer cancel()
    if output, err := exec.CommandContext(ctx, "sh", "-c", cmstr).CombinedOutput(); err != nil {
        fmt.Printf("Exceed the timeout 100ms. %v\n", err)
        fmt.Printf("Combine Output: %s\n", output)
		fmt.Println("err happaen")
		res = false
    }else{
		fmt.Println("suc")
		res = true
	}
	err := write_status("./logfile",res,*name,make_log_message)

	if err != nil {
		fmt.Println("can't write to status file -> logfile")
	}
}

func write_status(
	path string,
	status bool,
	job_name string,
	f func(string,bool) (string, error)) error{

	file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

	line,err := f(job_name,status)

	if err != nil {
		return err
	}

	if _, err := file.WriteString(line); err != nil {
		return err
	}

	return err
}

func make_log_message(job_name string,status bool) (message string, err error) {

	var line string
	// time := time.Now()

	if status {
		line = "[suc] time "+job_name + " is successful"
	}else{
		line = "[fail] time "+job_name + " is failed"
	}

	return line,err

}