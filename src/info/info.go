package info

import(
	"fmt"
	"goconf/conf"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
	"errors"
	"os/exec"
	"strconv"
)

func Info() (map[string]string,error){
	arr:=make(map[string]string)
	ostype := runtime.GOOS
	confarr, err := getconf()
	if err != nil {
		Logsave(err.Error())
		return arr,err
	}
	if strings.ToLower(ostype) == "windows" {
		arr,err = wininfo(confarr)
	} else if strings.ToLower(ostype) == "linux" {
		arr,err = linuxinfo(confarr)
	} else {
		Logsave("Your OS is not windows or linux!")
		err = errors.New("Your OS is not windows or linux!")
	}
	return arr,err
}
func getconf() (confs map[string]string, err error) {
	arr := make(map[string]string)
	conffile := "../conf/public.conf"
	c, err := conf.ReadConfigFile(conffile)
	if err != nil {
		Logsave(conffile + err.Error())
		return
	}
	process, err := c.GetString("system", "process") // returns false
	port, err := c.GetString("system", "port")       // returns false
	aspect, err := c.GetString("system", "aspect")   // returns false
	if process != "" {
		arr["process"] = process
	}
	if port != "" {
		arr["port"] = port
	}
	if aspect != "" {
		arr["aspect"] = aspect
	}
	return arr, err
}

func wininfo(confs map[string]string) (map[string]string , error){
	arr:=make(map[string]string)
	err:=errors.New("We can not support for windows now. Please waiting for the next version ...")
	Logsave("We can not support for windows now. Please waiting for the next version ...")
	return arr,err
}
func linuxinfo(confs map[string]string) (map[string]string , error){
	arr:=make(map[string]string)
	var narr []string
	var nport []string
	var nas []string
	var err error
	flag := true
	for k, v := range confs {
		Logsave(k +" : "+ v)
		if k=="process"{
			narr=strings.Split(v,",")	
		}else if k=="port"{
			nport=strings.Split(v,",")
		}else if k=="aspect"{
			nas=strings.Split(v,",")	
		}
	}
	if len(nport)==0{
		flag=false
	}
	
	if flag && len(narr)!=len(nport){
		err=errors.New("process个数和port格式个数不同！")
		Logsave("FATAL: process个数和port格式个数不同！")
		return arr,err
	}	
	for k1,v1:= range narr{
		var pid string
		var ps string
		if !flag{
			pid=getpid(v1,"-1")
		}else{
			pid=getpid(v1,nport[k1])
			v1=v1+"_"+nport[k1]
		}
		if pid==""{
			break
		}
		ps=getps(pid)
		if ps == ""{
			break
		}		
		for _,v2:=range nas{
			if v2=="cpu"{
				arr[v1+":"+"cpu"]=getcpu(ps)
			}else if v2=="mem"{
				arr[v1+":"+"mem"]=getmem(ps)
			}
		}
	}
	return arr,err
}
func getpid(process string,port string) (pid string){
	Logsave("PROCESS NAME: "+ process+" && PORT: " +port)
	intp,err:=strconv.Atoi(port)
	if err!=nil{
		Logsave("there is incorrect port in config")
		return 
	}
	if intp == -1{
		res,_ := exec.Command("/bin/sh", "-c", `ps -ef | grep -v "grep" | grep `+process +`|awk '{print $2}'`).Output()
		pid=strings.TrimSpace(string(res))
	}else if intp > 0{
		as:=[]byte(process)
		if len(as)>9{
			process=string(as[:9])
		}
		res,_ := exec.Command("/bin/sh", "-c", `lsof -i -P | grep  "LISTEN" | grep `+process+ `| grep `+port + `|awk '{print $2}'`).Output()	
		pid=strings.TrimSpace(string(res))
	}else{
		return 
	}
	return
}
func getps(pid string) (string){
	//res,_ := exec.Command("/bin/sh", "-c", `ps -aux | grep -v "grep" | grep  ` + pid).Output()
	res,_ := exec.Command("/bin/sh", "-c", `top -b -n 1 -p `+pid+` | sed '$d' | sed -n '$p'`).Output()
	Logsave("PID : "+pid+" => PSINFO : "+string(res))
	return	string(res)
} 
func getcpu(ps string) (msg string){
	psarr:=strings.Fields(ps)
	if len(psarr)<12 {
		return
	}
	msg = psarr[8]
	Logsave("CPU: "+msg)
	return 
}
func getmem(ps string) (msg string){
	psarr:=strings.Fields(ps)
	if len(psarr)<12 {
		return
	}
	msg = psarr[9]
	Logsave("MEM: "+msg)
	return 	
}
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func Logsave(msg string) {
	//log.SetFlags(log.Lshortfile | log.LstdFlags)
	timeflag := time.Now().Format("2006-01-02")
	logfile := "../log/" + timeflag + ".log"
	if !checkFileIsExist(logfile) {
		fc, _ := os.Create(logfile)
		defer fc.Close()
	}
	f, err := os.OpenFile(logfile, os.O_RDWR|os.O_APPEND, 0660)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
	}
	//logger := log.New(f, "", log.Ldate|log.Ltime|log.Llongfile)
	logger := log.New(f, "", log.Ldate|log.Ltime)
	logger.Print(msg)
}
