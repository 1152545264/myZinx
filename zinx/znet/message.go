package znet

type Message struct {
	ID      uint32 //消息ID
	DataLen uint32 //消息的长度
	Data    []byte //消息的内容
}

// NewMessagePackage 创建一个Message消息包
func NewMessagePackage(id uint32, data []byte) *Message {
	return &Message{
		ID:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

// GetMsgID 获取消息的ID
func (m *Message) GetMsgID() uint32 {
	return m.ID
}

// GetMsgLen 获取消息的长度
func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}

// GetData 获取消息的内容
func (m *Message) GetData() []byte {
	return m.Data
}

// SetMsgID 设置消息的ID
func (m *Message) SetMsgID(id uint32) {
	m.ID = id
}

// SetMsgData 设置消息的内容
func (m *Message) SetMsgData(data []byte) {
	m.Data = data
}

// SetDataLen 设置消息的长度
func (m *Message) SetDataLen(dataLen uint32) {
	m.DataLen = dataLen
}
