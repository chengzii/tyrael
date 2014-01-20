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
	"kval"
	"strings"
)

var (
	indexTemp = template.Must(template.ParseFiles("../tpl/index.html"))
	showTemp  = template.Must(template.ParseFiles("../tpl/show.html"))
	lineTemp  = template.Must(template.ParseFiles("../tpl/line.html"))
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
	http.Handle("/img/", http.FileServer(http.Dir("../public")))
	http.Handle("/css/", http.FileServer(http.Dir("../public")))
	http.Handle("/js/", http.FileServer(http.Dir("../public")))
	
	http.HandleFunc("/index", indexHandler)
	http.HandleFunc("/show", showHandler)
	http.HandleFunc("/line", lineHandler)
	http.HandleFunc("/getinfo", getinfoHandler)
	http.HandleFunc("/getline", getlineHandler)
	http.HandleFunc("/", notfoundHandler)
	if err := http.ListenAndServe(":"+confarr["httpport"], nil); err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
func notfoundHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/index", http.StatusFound)
	}
	t, err := template.ParseFiles("../tpl/404.html")
	if (err != nil) {
		info.Logsave(err.Error())
	}
	t.Execute(w, nil)
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
func lineHandler(w http.ResponseWriter, r *http.Request) {
	info.Logsave("REQUES FROM: "+r.RemoteAddr+" URI:"+r.RequestURI)
	if err := lineTemp.Execute(w, nil); err != nil {
		info.Logsave("lineHandler.lineTemp.Execute: " + err.Error())
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
type One struct{
	date	string
	size	float64
}
type Ones struct{
	list	[]One
}
func getlineHandler(w http.ResponseWriter, r *http.Request) {
	info.Logsave("REQUES FROM: "+r.RemoteAddr+" URI:"+r.RequestURI)
	kvarr:=kval.Getall()

	childarr:=make([]interface{},0)
	var rearr map[string]interface{}
	flag:=""
	types:=make([]string,2)
	types[0]="date"
	types[1]="使用率"
	for i:=0;i<len(kvarr);i++{
		v:=kvarr[i]		
		ks := strings.Split(v.Key,":")
		ktime,_:=strconv.Atoi(ks[2])
		kvalue,_:=strconv.ParseFloat(v.Val,0)
		date := time.Unix(int64(ktime), 0).Format("2006-01-02 15:04:05")
		newone:=make([]interface{},2)
		newone[0]=date
		newone[1]=kvalue
		//var newone map[string]interface{}=map[string]interface{}{"0":date,"1":kvalue}
		if i==len(kvarr)-1{
			if ks[0]+":"+ks[1] != flag{
				//rearr = map[string]interface{}{flag:childarr}
				if len(rearr)==0{
					rearr = map[string]interface{}{flag:childarr}
				}else{	
					rearr[flag]=childarr
				}
				flag = ks[0]+":"+ks[1]
				childarr = childarr[:0]	
				childarr = append(childarr,types)
			}
			childarr = append(childarr,newone)
			//rearr = map[string]interface{}{flag:childarr}
			if len(rearr)==0{
				rearr = map[string]interface{}{flag:childarr}
			}else{	
				rearr[flag]=childarr
			}
		}else if i==0{
			flag = ks[0]+":"+ks[1]
			childarr = append(childarr,types)
			childarr = append(childarr,newone)
		}else if ks[0]+":"+ks[1] != flag && i>0{
			//rearr = map[string]interface{}{flag:childarr}
			if len(rearr)==0{
				rearr = map[string]interface{}{flag:childarr}
			}else{	
				rearr[flag]=childarr
			}
			flag = ks[0]+":"+ks[1]
			if len(childarr)>0{
				childarr = childarr[:0]
				childarr = append(childarr,types)
			}
			childarr = append(childarr,newone)
		}else{
			childarr = append(childarr,newone)
		}	
		//fmt.Println("rearr",rearr)
	}
	result,_:= json.Marshal(rearr)
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
			}else{
				kval.Add(key,v)	
				//fmt.Println(kval.Getall())	
			}
		}
	}
}
