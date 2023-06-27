package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

//连接管理模块的具体实现

type ConnManager struct {
	connections map[uint32]ziface.IConnection //管理的连接信息
	connLock    sync.RWMutex                  //保护连接的读写锁

}

// Remove 删除连接
func (cm *ConnManager) Remove(conn ziface.IConnection) {

	//保护共享资源, 加写锁
	cm.connLock.Lock()
	//defer cm.connLock.Unlock()

	//删除连接信息
	delete(cm.connections, conn.GetConnID())
	cm.connLock.Unlock()

	fmt.Println("connectionID = ", conn.GetConnID(), "remove from ConnManager success: conn num = ", cm.Len())
}

// Len 获取当前连接总数
func (cm *ConnManager) Len() int {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()
	return len(cm.connections)
}

// Add 添加连接
func (cm *ConnManager) Add(conn ziface.IConnection) {
	//保护共享资源, 加写锁
	cm.connLock.Lock()
	cm.connections[conn.GetConnID()] = conn
	cm.connLock.Unlock()

	fmt.Println("connectionID = ", conn.GetConnID(), "connection add to ConnManager success: conn num = ", cm.Len())
}

// Get 根据connID获取连接
func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()
	if conn, ok := cm.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found ")
	}
}

// ClearConn 清除并终止所有连接
func (cm *ConnManager) ClearConn() {
	cm.connLock.Lock()
	//删除conn并停止conn的工作
	for connID, conn := range cm.connections {
		//停止
		conn.Stop()
		//删除
		delete(cm.connections, connID)
	}
	cm.connLock.Unlock()

	fmt.Println("Clear all connections!, now connection num=", cm.Len())
}

// NewConnManager 创建当前连接管理模块的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}
