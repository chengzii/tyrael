package main

import (
	"strconv"
	"goconf/conf"
	"html/template"
	"info"
	"net/http"
	"redisdb"
	"time"
	"os"
	"runtime"
	"syscall"
	"encoding/json"
	//"fmt"
)

var (
	indexTemp = template.Must(template.ParseFiles("../tpl/index.html"))
	showTemp  = template.Must(template.ParseFiles("../tpl/show.html"))
)

func main() {
//	daemon(1, 1)
	go func(){
		test()
	}()
	info.Logsave("Start create a httpserver!")
	starthttpserver()
}

func starthttpserver() {
	confarr, err := getconf()
	if err != nil {
		info.Logsave("Get http config error!")
		return
	}
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/show", showHandler)
	http.HandleFunc("/getinfo", getinfoHandler)
	if err := http.ListenAndServe(":"+confarr["httpport"], nil); err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	info.Logsave("REQUES FROM: "+r.RemoteAddr+" URI:"+r.RequestURI)
	if err := indexTemp.Execute(w, nil); err != nil {
		info.Logsave("indexHandler.indexTemp.Execute: " + err.Error())
		return
	}
}

func showHandler(w http.ResponseWriter, r *http.Request) {
	info.Logsave("REQUES FROM: "+r.RemoteAddr+" URI:"+r.RequestURI)
	if err := showTemp.Execute(w, nil); err != nil {
		info.Logsave("showHandler.showTemp.Execute: " + err.Error())
		return
	}
}
func getinfoHandler(w http.ResponseWriter, r *http.Request) {
	info.Logsave("REQUES FROM: "+r.RemoteAddr+" URI:"+r.RequestURI)
	rearr := make([]interface{}, 0)
	title := make([]interface{}, 2)
	title[0] = "Label"
	title[1] = "Value"
	rearr = append(rearr, title)
	infoarr,err:=info.Info()
	//fmt.Println(infoarr)
	if err!=nil{
		return
	}
	for k,v:=range infoarr{
		arr := make([]interface{}, 2)
		fv,_:=strconv.ParseFloat(v, 2)
		arr[0] = k
		arr[1] = fv
		rearr = append(rearr,arr)	
	}
	result, _ := json.Marshal(rearr)
	w.Write(result)
}
func getconf() (map[string]string, error) {
	arr := make(map[string]string)
	conffile := "../conf/public.conf"
	c, err := conf.ReadConfigFile(conffile)
	if err != nil {
		info.Logsave(conffile + err.Error())
		return arr, err
	}
	httpport, err := c.GetString("http", "port") // returns false
	intervaltime, err := c.GetString("system", "intervaltime") // returns false
	totaltime, err := c.GetString("system", "totaltime") // returns false
	issave, err := c.GetString("system", "issave") // returns false

	if httpport != "" {
		arr["httpport"] = httpport
	}
	if intervaltime != "" {
		arr["intervaltime"] = intervaltime
	}
	if totaltime != "" {
		arr["totaltime"] = totaltime
	}
	if issave != "" {
		arr["issave"] = issave
	}
	return arr, err
}

func daemon(nochdir, noclose int) int {
    var ret, ret2 uintptr
    var err syscall.Errno
 
    darwin := runtime.GOOS == "darwin"
 
    // already a daemon
    if syscall.Getppid() == 1 {
        return 0
    }
 
    // fork off the parent process
    ret, ret2, err = syscall.RawSyscall(syscall.SYS_FORK, 0, 0, 0)
    if err != 0 {
        return -1
    }
 
    // failure
    if ret2 < 0 {
        os.Exit(-1)
    }
 
    // handle exception for darwin
    if darwin && ret2 == 1 {
        ret = 0
    }
 
    // if we got a good PID, then we call exit the parent process.
    if ret > 0 {
        os.Exit(0)
    }
 
    /* Change the file mode mask */
    _ = syscall.Umask(0)
 
    // create a new SID for the child process
    s_ret, s_errno := syscall.Setsid()
    if s_errno != nil {
   	info.Logsave("Error: syscall.Setsid errno: "+s_errno.Error())
    }
    if s_ret < 0 {
        return -1
    }
 
    if nochdir == 0 {
        os.Chdir("/")
    }
 
    if noclose == 0 {
        f, e := os.OpenFile("/dev/null", os.O_RDWR, 0)
        if e == nil {
            fd := f.Fd()
            syscall.Dup2(int(fd), int(os.Stdin.Fd()))
            syscall.Dup2(int(fd), int(os.Stdout.Fd()))
            syscall.Dup2(int(fd), int(os.Stderr.Fd()))
        }
    }
 
    return 0
}

func test() {
	confarr, err := getconf()
	intervaltime,err := strconv.Atoi(confarr["intervaltime"])
	totaltime,err := strconv.Atoi(confarr["totaltime"])
	issave,err := strconv.Atoi(confarr["issave"])
	flag := true
	if err != nil {
		flag = false
		info.Logsave("Get intervaltime config error!")
	}
	var nowtime int
	for{	
		nowtime = int(time.Now().Unix())
		if !flag{
			break
		}
		time.Sleep(time.Duration(intervaltime)*time.Second)
		rearr,err:=info.Info()
		if err!=nil{
			continue
		}
		for k,v:=range rearr{
			key:=k+":"+strconv.Itoa(nowtime)
			info.Logsave("KEY:"+key+" => VALUE:"+v)
			if issave==1 {
				redisdb.Rset(key,[]byte(v))
				redisdb.Rsetval(key,int64(totaltime))
			}
		}
	}
}
