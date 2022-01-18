package main


import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"context"
	"time"
	"path/filepath"
)

func main(){
	timeout := flag.Int("t", 10, "timeout time")
	status_log_path := flag.String("o", "./log/status.txt", "job status log file path")
	logfile := flag.String("log", "./log/log.txt", "job itself log file path")
	command := flag.String("cmd", "ls -la", "exec command")
	name := flag.String("name","hoge-job","in status file job name")
	flag.Parse()

	cmstr := *command
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout*100) * time.Millisecond)

    defer cancel()
    _, err := exec.CommandContext(ctx, "sh", "-c", cmstr).CombinedOutput()
	res := check_nil(err)
	message,err := make_log_message(*name,res,get_date)

	if err != nil{
		write_log(*logfile,"somthing fail in make_log_message")
	}

	if err := write_log(*status_log_path,message); err != nil{
		write_log(*logfile,"somthing fail in writing info to status log file")
	}
}

func check_nil(err error) bool {
	if err != nil{
		return true
	}else{
		return false
	}
}


func write_log(path string, message string) error{
	dirname := filepath.Dir(path)

	if err := os.MkdirAll(dirname, 0744); err != nil {
        fmt.Println(err)
    }
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
		fmt.Println(err)
        return err
    }
    defer file.Close()

	if _, err := file.WriteString(message); err != nil {
		return err
	}

	return err
}


func make_log_message(job_name string,status bool,f func() string) (message string, err error) {
	
	var format string

	if status {
		format = "[SUCCESS] %s %s is successful\n"
	}else{
		format = "[FAIL] %s %s is fail\n"
	}

	time := f()
	res := fmt.Sprintf(format, time, job_name)
	return res,err
}

func get_date()(date_str string){
	time_str := time.Now().Format("2006-01-02T15:04:05.000000000Z")
	return time_str
}