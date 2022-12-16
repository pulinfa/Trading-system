package utils

/*
	存储一切有关服务器的全局参数，供其它模块使用
	一些参数是可以通过server.json由用户自定义进行配置
*/
import (
	"encoding/json"
	"io/ioutil"
	"server/iface"
)

type GlobalObj struct {
	TcpServer iface.IServer //当前全局的Server对象
	Host      string        //当前服务器主机监听的IP
	TcpPort   int           //当前服务器主机监听的端口号
	Name      string        //当前服务器的名称

	Version          string //当前的版本号
	MaxConn          int    //当前服务器主机允许的最大连接数量
	MaxPackageSize   uint32 //当前框架数据包的最大值
	WorkerPoolSize   uint32 //当前业务工作Worker池的Goroutine数量
	MaxWorkerTaskLen uint32 //允许用户最多开启多少个Worker（限定条件）
}

// 定义全局的对外Globalobj
var GlobalObject *GlobalObj

// 从server.json中加载自定义的参数
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/server.json")
	if err != nil {
		panic(err)
	}

	//将json文件数据解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// 提供一个初始化方法，初始化当前的GlobalObj
func init() {
	GlobalObject = &GlobalObj{
		Name:             "ServerApp",
		Version:          "v10.0",
		TcpPort:          8999,
		Host:             "0.0.0.0",
		MaxConn:          1000,
		WorkerPoolSize:   10,
		MaxPackageSize:   4096,
		MaxWorkerTaskLen: 1024,
	}

	//应该尝试从
	GlobalObject.Reload()
}
