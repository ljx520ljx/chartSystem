package config

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

// Config 表示应用程序配置
type Config struct {
	XMLName  xml.Name  `xml:"ChartConfig"`
	Channels []Channel `xml:"Channels>Channel"`
	Display  Display   `xml:"Display"`
}

// Channel 表示通道配置
type Channel struct {
	ID       string  `xml:"id,attr"`
	Name     string  `xml:"Name"`
	Color    string  `xml:"Color"`
	Scale    float64 `xml:"Scale"`
	Visible  bool    `xml:"Visible"`
	YAxisMin float64 `xml:"YAxisMin"`
	YAxisMax float64 `xml:"YAxisMax"`
}

// Display 表示显示配置
type Display struct {
	GridVisible bool    `xml:"GridVisible"`
	RefreshRate int     `xml:"RefreshRate"`
	TimeScale   float64 `xml:"TimeScale"`
}

// LoadConfig 从指定路径加载配置文件
func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := xml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig 保存配置到指定路径
func SaveConfig(config *Config, path string) error {
	data, err := xml.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0644)
}
