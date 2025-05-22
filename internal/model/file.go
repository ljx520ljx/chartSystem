package model

import (
	"time"
)

// File 文件模型
type File struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Name         string         `json:"name" gorm:"size:255;not null"`
	Description  string         `json:"description" gorm:"size:1000"`
	FilePath     string         `json:"file_path" gorm:"size:500;not null"`
	FileSize     int64          `json:"file_size" gorm:"not null"`
	ContentType  string         `json:"content_type" gorm:"size:100;not null"`
	UserID       uint           `json:"user_id" gorm:"not null"`
	DataChannels []DataChannel  `json:"data_channels,omitempty" gorm:"foreignKey:FileID"`
	Processing   *FileProcessing `json:"processing,omitempty" gorm:"foreignKey:FileID"`
	Markers      []Marker       `json:"markers,omitempty" gorm:"foreignKey:FileID"`
	Analyses     []Analysis     `json:"analyses,omitempty" gorm:"foreignKey:FileID"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

// DataChannel 数据通道模型
type DataChannel struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	FileID      uint      `json:"file_id" gorm:"not null"`
	Name        string    `json:"name" gorm:"size:100;not null"`
	Description string    `json:"description" gorm:"size:500"`
	Unit        string    `json:"unit" gorm:"size:50"`
	SampleRate  float64   `json:"sample_rate" gorm:"not null"`
	DataFormat  string    `json:"data_format" gorm:"size:50;not null"`
	DataOffset  int64     `json:"data_offset" gorm:"not null"`
	DataLength  int64     `json:"data_length" gorm:"not null"`
	MinValue    float64   `json:"min_value"`
	MaxValue    float64   `json:"max_value"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// FileProcessing 文件处理状态模型
type FileProcessing struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	FileID    uint      `json:"file_id" gorm:"not null;uniqueIndex"`
	Status    string    `json:"status" gorm:"size:50;not null;default:'pending'"`
	Progress  float64   `json:"progress" gorm:"default:0"`
	Message   string    `json:"message" gorm:"size:500"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Marker 标记点模型
type Marker struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	FileID      uint      `json:"file_id" gorm:"not null"`
	ChannelID   uint      `json:"channel_id" gorm:"not null"`
	Position    float64   `json:"position" gorm:"not null"` // 时间位置（秒）
	Type        string    `json:"type" gorm:"size:50;not null"`
	Label       string    `json:"label" gorm:"size:200"`
	Description string    `json:"description" gorm:"size:500"`
	Color       string    `json:"color" gorm:"size:20;default:'#FF0000'"`
	CreatedBy   uint      `json:"created_by" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Analysis 分析结果模型
type Analysis struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	FileID      uint      `json:"file_id" gorm:"not null"`
	UserID      uint      `json:"user_id" gorm:"not null"`
	Name        string    `json:"name" gorm:"size:100;not null"`
	Type        string    `json:"type" gorm:"size:50;not null"`
	Parameters  string    `json:"parameters" gorm:"type:text"`
	Results     string    `json:"results" gorm:"type:longtext"`
	StartTime   float64   `json:"start_time"`
	EndTime     float64   `json:"end_time"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
