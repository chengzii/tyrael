package kval

import(
	"sort"
	//"fmt"
	//"strconv"
)

var kval map[string]string

type kvarr []item
                 
type item struct {
    Key string
    Val string
}

func init(){
	kval=make(map[string]string,65535)
}

/*
func main(){
	Add("v1","7")
	Add("v3","4")
	Add("v33","3")
	Add("v23","1")
	Add("v13","2")
	fmt.Println(Getall())
}
*/

func (arr kvarr) Len() int {
	return len(arr)
}
                 
func (arr kvarr) Less(i, j int) bool {
	//fmt.Println(arr[i].Key, arr[j].Key,arr[i].Key < arr[j].Key)
	//return arr[i].Val < arr[j].Val
	return arr[i].Key < arr[j].Key
}
                 
func (arr kvarr) Swap(i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}
func NewMap(input map[string]string) kvarr{
	ms := make(kvarr, 0, len(input))
                 
    	for k, v := range input {
        	ms = append(ms, item{k, v})
    	}
                 
    	return ms
}

func Add(key string, val string) {
	kval[key]=val
}
func Del(key string) {
	if isexist(key){
		delete(kval,key)
	}
}
func Mod(key string, val string) (res bool){
	if isexist(key){
		kval[key]=val
		res=true
	}else{
		res=false
	}
	return	
}
func Getbykey(key string) (val string,ok bool){
	if isexist(key){
		val=kval[key]
		ok=true
	}else{
		ok=false
	}
	return	
}
func Getall() (res kvarr){
	res=NewMap(kval)
	sort.Sort(res)
	return	
}
func isexist(key string) (res bool){
	if len(kval)==0{
		res=false
	}else{
		_,ok:=kval[key]
		if ok{
			res=true
		}else{
			res=false	
		}
	}
	return
}
