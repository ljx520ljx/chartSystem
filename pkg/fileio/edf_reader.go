package fileio

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/liujiaxin/chartSystem/internal/data"
)

// EDFHeader 表示EDF文件头
type EDFHeader struct {
	Version       string    // 8字节
	PatientID     string    // 80字节
	RecordingID   string    // 80字节
	StartTime     time.Time // 从日期和时间字段解析
	HeaderBytes   int       // 8字节
	Reserved      string    // 44字节
	DataRecords   int       // 8字节
	Duration      float64   // 8字节
	NumSignals    int       // 4字节
	SignalHeaders []SignalHeader
}

// SignalHeader 表示信号头
type SignalHeader struct {
	Label       string  // 16字节
	Transducer  string  // 80字节
	PhysicalDim string  // 8字节
	PhysicalMin float64 // 8字节
	PhysicalMax float64 // 8字节
	DigitalMin  float64 // 8字节
	DigitalMax  float64 // 8字节
	Prefiltering string  // 80字节
	Samples     int     // 8字节
	Reserved    string  // 32字节
}

// EDFReader 表示EDF文件读取器
type EDFReader struct {
	file   *os.File
	header EDFHeader
}

// OpenEDF 打开一个EDF文件
func OpenEDF(path string) (*EDFReader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader := &EDFReader{
		file: file,
	}

	if err := reader.readHeader(); err != nil {
		file.Close()
		return nil, err
	}

	return reader, nil
}

// Close 关闭文件
func (r *EDFReader) Close() error {
	return r.file.Close()
}

// 读取固定长度的字符串
func readString(file *os.File, length int) (string, error) {
	buf := make([]byte, length)
	_, err := file.Read(buf)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(buf)), nil
}

// 读取整数
func readInt(file *os.File, length int) (int, error) {
	str, err := readString(file, length)
	if err != nil {
		return 0, err
	}
	val, err := strconv.Atoi(strings.TrimSpace(str))
	if err != nil {
		return 0, err
	}
	return val, nil
}

// 读取浮点数
func readFloat(file *os.File, length int) (float64, error) {
	str, err := readString(file, length)
	if err != nil {
		return 0, err
	}
	val, err := strconv.ParseFloat(strings.TrimSpace(str), 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// 读取文件头
func (r *EDFReader) readHeader() error {
	// 重置文件指针到开头
	_, err := r.file.Seek(0, 0)
	if err != nil {
		return err
	}

	// 读取版本
	r.header.Version, err = readString(r.file, 8)
	if err != nil {
		return err
	}

	// 读取病人ID
	r.header.PatientID, err = readString(r.file, 80)
	if err != nil {
		return err
	}

	// 读取记录ID
	r.header.RecordingID, err = readString(r.file, 80)
	if err != nil {
		return err
	}

	// 读取开始日期
	startDate, err := readString(r.file, 8)
	if err != nil {
		return err
	}

	// 读取开始时间
	startTime, err := readString(r.file, 8)
	if err != nil {
		return err
	}

	// 解析日期和时间
	r.header.StartTime, err = time.Parse("02.01.06 15.04.05", startDate+" "+startTime)
	if err != nil {
		// 如果解析失败，使用当前时间
		r.header.StartTime = time.Now()
	}

	// 读取头部大小
	r.header.HeaderBytes, err = readInt(r.file, 8)
	if err != nil {
		return err
	}

	// 读取保留字段
	r.header.Reserved, err = readString(r.file, 44)
	if err != nil {
		return err
	}

	// 读取数据记录数
	r.header.DataRecords, err = readInt(r.file, 8)
	if err != nil {
		return err
	}

	// 读取每个数据记录的持续时间
	r.header.Duration, err = readFloat(r.file, 8)
	if err != nil {
		return err
	}

	// 读取信号数量
	r.header.NumSignals, err = readInt(r.file, 4)
	if err != nil {
		return err
	}

	// 读取信号头
	r.header.SignalHeaders = make([]SignalHeader, r.header.NumSignals)

	// 按字段类型批量读取所有信号的头部信息
	for i := 0; i < r.header.NumSignals; i++ {
		label, err := readString(r.file, 16)
		if err != nil {
			return err
		}
		r.header.SignalHeaders[i].Label = label
	}

	for i := 0; i < r.header.NumSignals; i++ {
		transducer, err := readString(r.file, 80)
		if err != nil {
			return err
		}
		r.header.SignalHeaders[i].Transducer = transducer
	}

	for i := 0; i < r.header.NumSignals; i++ {
		physicalDim, err := readString(r.file, 8)
		if err != nil {
			return err
		}
		r.header.SignalHeaders[i].PhysicalDim = physicalDim
	}

	for i := 0; i < r.header.NumSignals; i++ {
		physicalMin, err := readFloat(r.file, 8)
		if err != nil {
			return err
		}
		r.header.SignalHeaders[i].PhysicalMin = physicalMin
	}

	for i := 0; i < r.header.NumSignals; i++ {
		physicalMax, err := readFloat(r.file, 8)
		if err != nil {
			return err
		}
		r.header.SignalHeaders[i].PhysicalMax = physicalMax
	}

	for i := 0; i < r.header.NumSignals; i++ {
		digitalMin, err := readFloat(r.file, 8)
		if err != nil {
			return err
		}
		r.header.SignalHeaders[i].DigitalMin = digitalMin
	}

	for i := 0; i < r.header.NumSignals; i++ {
		digitalMax, err := readFloat(r.file, 8)
		if err != nil {
			return err
		}
		r.header.SignalHeaders[i].DigitalMax = digitalMax
	}

	for i := 0; i < r.header.NumSignals; i++ {
		prefiltering, err := readString(r.file, 80)
		if err != nil {
			return err
		}
		r.header.SignalHeaders[i].Prefiltering = prefiltering
	}

	for i := 0; i < r.header.NumSignals; i++ {
		samples, err := readInt(r.file, 8)
		if err != nil {
			return err
		}
		r.header.SignalHeaders[i].Samples = samples
	}

	for i := 0; i < r.header.NumSignals; i++ {
		reserved, err := readString(r.file, 32)
		if err != nil {
			return err
		}
		r.header.SignalHeaders[i].Reserved = reserved
	}

	return nil
}

// ReadSignalData 读取信号数据
func (r *EDFReader) ReadSignalData(signalIndex int, startRecord, numRecords int) ([]int16, error) {
	if signalIndex < 0 || signalIndex >= r.header.NumSignals {
		return nil, fmt.Errorf("信号索引超出范围: %d", signalIndex)
	}

	if startRecord < 0 || startRecord >= r.header.DataRecords {
		return nil, fmt.Errorf("起始记录超出范围: %d", startRecord)
	}

	if numRecords <= 0 {
		return nil, fmt.Errorf("记录数量必须大于0: %d", numRecords)
	}

	// 确保不超出总记录数
	if startRecord+numRecords > r.header.DataRecords {
		numRecords = r.header.DataRecords - startRecord
	}

	// 计算一个数据记录中所有信号的样本总数
	totalSamplesPerRecord := 0
	for i := 0; i < r.header.NumSignals; i++ {
		totalSamplesPerRecord += r.header.SignalHeaders[i].Samples
	}

	// 计算数据记录大小（字节数）
	recordSize := totalSamplesPerRecord * 2 // 每个样本2字节

	// 计算目标信号在记录中的偏移
	signalOffset := 0
	for i := 0; i < signalIndex; i++ {
		signalOffset += r.header.SignalHeaders[i].Samples * 2
	}

	// 每个记录中当前信号的样本数
	samplesPerSignal := r.header.SignalHeaders[signalIndex].Samples

	// 分配结果数组
	result := make([]int16, samplesPerSignal*numRecords)

	// 读取每个记录中的信号数据
	for i := 0; i < numRecords; i++ {
		// 计算记录在文件中的位置
		recordPos := int64(r.header.HeaderBytes) + int64(startRecord+i)*int64(recordSize)
		
		// 跳转到记录中当前信号的起始位置
		_, err := r.file.Seek(recordPos+int64(signalOffset), 0)
		if err != nil {
			return nil, err
		}

		// 读取当前信号的所有样本
		for j := 0; j < samplesPerSignal; j++ {
			var sample int16
			err := binary.Read(r.file, binary.LittleEndian, &sample)
			if err != nil {
				return nil, err
			}
			result[i*samplesPerSignal+j] = sample
		}
	}

	return result, nil
}

// ConvertToPhysical 将数字值转换为物理值
func (r *EDFReader) ConvertToPhysical(signalIndex int, digitalValue int16) float64 {
	if signalIndex < 0 || signalIndex >= r.header.NumSignals {
		return 0
	}

	sh := r.header.SignalHeaders[signalIndex]
	
	// 计算转换因子
	scale := (sh.PhysicalMax - sh.PhysicalMin) / (sh.DigitalMax - sh.DigitalMin)
	
	// 转换数值
	physicalValue := sh.PhysicalMin + float64(digitalValue-int16(sh.DigitalMin))*scale
	
	return physicalValue
}

// LoadSignalToChannel 将信号数据加载到通道
func (r *EDFReader) LoadSignalToChannel(signalIndex int, channel *data.Channel) error {
	// 读取所有数据记录
	digitalData, err := r.ReadSignalData(signalIndex, 0, r.header.DataRecords)
	if err != nil {
		return err
	}

	// 清除通道中现有数据
	channel.ClearData()

	// 计算每个样本的时间间隔
	timeStep := r.header.Duration / float64(r.header.SignalHeaders[signalIndex].Samples)

	// 将数字值转换为物理值并添加到通道
	for i, digitalValue := range digitalData {
		// 计算时间
		t := float64(i) * timeStep
		
		// 转换为物理值
		y := r.ConvertToPhysical(signalIndex, digitalValue)
		
		// 添加到通道
		channel.AddDataPoint(t, y)
	}

	return nil
}

// GetChannelInfo 获取通道信息
func (r *EDFReader) GetChannelInfo(signalIndex int) (string, string, float64, float64) {
	if signalIndex < 0 || signalIndex >= r.header.NumSignals {
		return "", "", 0, 0
	}

	sh := r.header.SignalHeaders[signalIndex]
	return sh.Label, sh.PhysicalDim, sh.PhysicalMin, sh.PhysicalMax
}

// GetHeader 获取文件头
func (r *EDFReader) GetHeader() EDFHeader {
	return r.header
}

// GetNumSignals 获取信号数量
func (r *EDFReader) GetNumSignals() int {
	return r.header.NumSignals
}

// GetSignalSamplingRate 获取信号采样率
func (r *EDFReader) GetSignalSamplingRate(signalIndex int) float64 {
	if signalIndex < 0 || signalIndex >= r.header.NumSignals {
		return 0
	}
	
	return float64(r.header.SignalHeaders[signalIndex].Samples) / r.header.Duration
}

// CreateSimulatedEDFData 创建模拟的EDF数据
func CreateSimulatedEDFData(channel *data.Channel, dataType string, duration float64, samplingRate float64) {
	// 清除通道中现有数据
	channel.ClearData()
	
	// 计算总样本数
	totalSamples := int(duration * samplingRate)
	
	// 时间步长
	timeStep := 1.0 / samplingRate
	
	// 生成数据
	for i := 0; i < totalSamples; i++ {
		t := float64(i) * timeStep
		var y float64
		
		switch dataType {
		case "sine":
			// 生成正弦波
			y = math.Sin(2 * math.Pi * t)
		case "ecg":
			// 模拟心电图数据
			period := 1.0 // 心跳周期（秒）
			phase := t - math.Floor(t/period)*period // 0到period之间的相位
			
			if phase < 0.1 {
				// P波
				y = 0.25 * math.Sin(2*math.Pi*phase/0.2)
			} else if phase >= 0.1 && phase < 0.4 {
				// 平坦段
				y = 0
			} else if phase >= 0.4 && phase < 0.45 {
				// Q波
				y = -0.5 * (phase - 0.4) / 0.05
			} else if phase >= 0.45 && phase < 0.5 {
				// R波
				y = -0.5 + 2 * (phase - 0.45) / 0.05
			} else if phase >= 0.5 && phase < 0.55 {
				// S波
				y = 1.5 - 2 * (phase - 0.5) / 0.05
			} else if phase >= 0.55 && phase < 0.7 {
				// T波
				peak := (phase - 0.55) / 0.15
				y = -0.5 + 0.75 * math.Sin(math.Pi * peak)
			} else {
				// 平坦段
				y = 0
			}
		case "bp":
			// 模拟血压数据
			period := 1.0 // 心跳周期（秒）
			phase := t - math.Floor(t/period)*period
			
			// 收缩压和舒张压的波形
			if phase < 0.3 {
				// 快速上升（收缩）
				y = 80 + 40 * math.Sin(math.Pi/2 + math.Pi*phase/0.3)
			} else {
				// 缓慢下降（舒张）
				y = 80 + 40 * math.Sin(math.Pi/2 + math.Pi*0.3/0.3) * math.Exp(-(phase-0.3)/0.5)
			}
		case "resp":
			// 模拟呼吸数据
			period := 4.0 // 呼吸周期（秒）
			y = math.Sin(2 * math.Pi * t / period)
		case "spo2":
			// 模拟血氧数据
			y = 98 + math.Sin(2*math.Pi*t) * 1 // 98% 左右波动
		default:
			// 默认为噪声
			y = (rand() - 0.5) * 2
		}
		
		// 添加一些随机噪声
		y += (rand() - 0.5) * 0.1
		
		// 添加到通道
		channel.AddDataPoint(t, y)
	}
}

// 生成0到1之间的随机数
func rand() float64 {
	return float64(time.Now().UnixNano()%1000) / 1000.0
}
