package snet

import (
	"errors"
	"fmt"
	"server/iface"
	"sync"
)

/*
* 连接管理模块的实现层
 */

type ConnManager struct {
	connections map[uint32]iface.IConnection //管理的连接的集合
	connLock    sync.RWMutex                 //保护连接集合的锁
}

// 创建当前连接的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]iface.IConnection),
	}
}

// 添加连接
func (connMgr *ConnManager) Add(conn iface.IConnection) {
	//保护共享资源map，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//将conn加入到ConnManager中
	connMgr.connections[conn.GetConnID()] = conn

	fmt.Println("connection add to ConnManager successfully: conn num = ", connMgr.Len())
}

// 删除连接
func (connMgr *ConnManager) Remove(conn iface.IConnection) {
	//保护共享资源map，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//将conn从ConnManager中删除
	delete(connMgr.connections, conn.GetConnID())

	fmt.Println("connection Remove from ConnManager successfully: conn num = ", connMgr.Len())
}

// 根据connID获取连接
func (connMgr *ConnManager) Get(connID uint32) (iface.IConnection, error) {
	//保护共享资源map，加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not FOUND")
	}
}

// 得到当前连接总数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

// 清楚并终止所有连接
func (connMgr *ConnManager) ClearConn() {
	//保护共享资源map，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除conn并停止conn的工作
	for connID, conn := range connMgr.connections {
		conn.Stop()

		delete(connMgr.connections, connID)
	}

	fmt.Println("Clear All connection succ! conn num = ", connMgr.Len())
}
