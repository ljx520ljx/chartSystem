package data

import (
	"strconv"
)

// DataPoint 表示一个数据点
type DataPoint struct {
	X float64
	Y float64
}

// Channel 表示一个数据通道
type Channel struct {
	ID            string
	Name          string
	Data          []DataPoint
	ProcessedData []DataPoint
	Visible       bool
	Color         string
	Scale         float64
	YAxisMin      float64
	YAxisMax      float64
}

// NewChannel 创建一个新的通道
func NewChannel(id string, name string) *Channel {
	return &Channel{
		ID:            id,
		Name:          name,
		Data:          make([]DataPoint, 0),
		ProcessedData: make([]DataPoint, 0),
		Visible:       true,
		Color:         "#FF0000",
		Scale:         1.0,
		YAxisMin:      -1.0,
		YAxisMax:      1.0,
	}
}

// AddDataPoint 添加一个数据点到通道
func (c *Channel) AddDataPoint(x, y float64) {
	c.Data = append(c.Data, DataPoint{X: x, Y: y})
}

// ClearData 清除通道中的所有数据
func (c *Channel) ClearData() {
	c.Data = make([]DataPoint, 0)
	c.ProcessedData = make([]DataPoint, 0)
}

// DataModel 表示应用程序的数据模型
type DataModel struct {
	Channels map[string]*Channel
}

// NewDataModel 创建一个新的数据模型
func NewDataModel() *DataModel {
	return &DataModel{
		Channels: make(map[string]*Channel),
	}
}

// AddChannel 添加一个通道到数据模型
func (m *DataModel) AddChannel(channel *Channel) {
	m.Channels[channel.ID] = channel
}

// GetChannel 通过ID获取通道
func (m *DataModel) GetChannel(id string) *Channel {
	return m.Channels[id]
}

// RemoveChannel 通过ID移除通道
func (m *DataModel) RemoveChannel(id string) {
	delete(m.Channels, id)
}

// IDToString 将索引转换为字符串ID
func IDToString(id int) string {
	return strconv.Itoa(id + 1)
}
