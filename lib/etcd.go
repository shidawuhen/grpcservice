package lib

import (
	"github.com/coreos/etcd/clientv3"
	"time"
	"fmt"
	"context"
)
const (
	GROUP = "b2c"
	TEAM =  "i18n"
)

var (
	config  clientv3.Config
	err     error
	client  *clientv3.Client
	kv      clientv3.KV
	putResp *clientv3.PutResponse
)

func init(){
	//配置
	config = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second * 5,
	}
	//连接 创建一个客户端
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}
}

func EtcdPut(port string) {
	if client == nil {
		return
	}
	//获取ip
	ip, err := ExternalIP()
	if err != nil {
		fmt.Println(err)
		return
	}
	address := ip.String() + port
	fmt.Println(address)
	//用于读写etcd的键值对
	kv = clientv3.NewKV(client)
	putResp, err = kv.Put(context.TODO(), "/"+GROUP+ "/" + TEAM + "/" + address, address, clientv3.WithPrevKV())
	if err != nil {
		fmt.Println(err)
	} else {
		//获取版本信息
		fmt.Println("Revision:", putResp.Header.Revision)
		if putResp.PrevKv != nil {
			fmt.Println("key:", string(putResp.PrevKv.Key))
			fmt.Println("Value:", string(putResp.PrevKv.Value))
			fmt.Println("Version:", string(putResp.PrevKv.Version))
		}
	}
}

func EtcdDelete(port string){
	fmt.Println("etcddelete")
	if client == nil {
		return
	}
	//获取ip
	ip, err := ExternalIP()
	if err != nil {
		fmt.Println(err)
		return
	}
	address := ip.String() + port
	fmt.Println(address)

	//用于读写etcd的键值对
	kv = clientv3.NewKV(client)

	delResp,err := kv.Delete(context.TODO(),"/"+GROUP+ "/" + TEAM + "/" + address,clientv3.WithPrevKV())
	if err != nil{
		fmt.Println(err)
		return
	}else{
		if len(delResp.PrevKvs) > 0 {
			for idx,kvpair := range delResp.PrevKvs{
				idx = idx
				fmt.Println("删除了",string(kvpair.Key),string(kvpair.Value))
			}
		}
	}
}
