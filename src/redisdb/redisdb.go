package redisdb

import (
	"fmt"
	"github.com/hoisie/redis"
	"goconf/conf"
	"info"
	"strconv"
	"sync"
	//"time"
)

var (
	client *redis.Client
	mutex  sync.Mutex
)

func init() {

	mutex.Lock()
	defer mutex.Unlock()

	if client != nil {
		return
	}
	confarr, err := getconf()
	db, err := strconv.Atoi(confarr["db"])
	maxpoolsize, err := strconv.Atoi(confarr["maxpoolsize"])
	issave, err := strconv.Atoi(confarr["issave"])
	if err != nil {
		info.Logsave("Get db config error!")
		return
	}
	client = &redis.Client{
		Addr:        confarr["ip"] + ":" + confarr["port"],
		Db:          db, // default db is 0
		Password:    confarr["password"],
		MaxPoolSize: maxpoolsize,
	}
	if issave!=1{
		info.Logsave("Save config is not open")
		return
	}
	if err := client.Auth("chengzi"); err != nil {
		fmt.Println("chengzi: ", err.Error())
		return
	}
}
func Rset(key string, value []byte) error {
	return client.Set(key, value)
}
func Rsetval(key string, interval int64) (bool, error) {
	return client.Expire(key, interval)
}
func Rget(key string) (value []byte, err error) {
	return client.Get(key)
}
func Rdel(key string) (bool, error) {
	return client.Del(key)
}
func getconf() (map[string]string, error) {
	arr := make(map[string]string)
	conffile := "../conf/public.conf"
	c, err := conf.ReadConfigFile(conffile)
	if err != nil {
		info.Logsave(conffile + err.Error())
		return arr, err
	}
	ip, err := c.GetString("db", "ip")                   // returns false
	port, err := c.GetString("db", "port")               // returns false
	db, err := c.GetString("db", "db")                   // returns false
	password, err := c.GetString("db", "password")       // returns false
	maxpoolsize, err := c.GetString("db", "maxpoolsize") // returns false
	issave, err := c.GetString("system", "issave") // returns false
	if ip != "" {
		arr["ip"] = ip
	}
	if port != "" {
		arr["port"] = port
	}
	if db != "" {
		arr["db"] = db
	}
	if password != "" {
		arr["password"] = password
	} else {
		arr["password"] = ""
	}
	if maxpoolsize != "" {
		arr["maxpoolsize"] = maxpoolsize
	}
	if issave != "" {
		arr["issave"] = issave
	}
	return arr, err
}
