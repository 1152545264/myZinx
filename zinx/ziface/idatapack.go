package ziface

/*
封包拆包模块
直接面向tcp连接种的数据流，用于处理TCP粘包问题
*/
type IDataPack interface {
	//GetHeadLen 获取包头长度
	GetHeadLen() uint32
	//Pack 封包方法
	Pack(msg IMessage) ([]byte, error)
	//UnPack 拆包方法
	UnPack([]byte) (IMessage, error)
}
