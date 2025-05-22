package signal

import (
	"math"
	"math/cmplx"

	"github.com/liujiaxin/chartSystem/internal/data"
	"gonum.org/v1/gonum/dsp/fourier"
	"gonum.org/v1/gonum/dsp/window"
)

// 用于消除未使用函数的警告
// 当系统集成此包时可以移除
var _ = func() interface{} {
	// 只是为了消除未使用函数的警告，不实际执行
	p := NewProcessor(1000.0)
	return p
}()

// Processor 表示信号处理器
type Processor struct {
	SampleRate float64
}

// NewProcessor 创建一个新的信号处理器
func NewProcessor(sampleRate float64) *Processor {
	return &Processor{
		SampleRate: sampleRate,
	}
}

// ApplyDifferential 应用微分处理
func (p *Processor) ApplyDifferential(channel *data.Channel) {
	if len(channel.Data) < 2 {
		return
	}

	channel.ProcessedData = make([]data.DataPoint, len(channel.Data))

	// 第一个点的微分设为0
	channel.ProcessedData[0] = data.DataPoint{
		X: channel.Data[0].X,
		Y: 0,
	}

	// 计算其他点的微分
	for i := 1; i < len(channel.Data); i++ {
		dt := channel.Data[i].X - channel.Data[i-1].X
		if dt == 0 {
			dt = 1.0 / p.SampleRate
		}

		dy := channel.Data[i].Y - channel.Data[i-1].Y
		derivativeY := dy / dt

		channel.ProcessedData[i] = data.DataPoint{
			X: channel.Data[i].X,
			Y: derivativeY,
		}
	}
}

// ApplyLowPassFilter 应用低通滤波
func (p *Processor) ApplyLowPassFilter(channel *data.Channel, cutoffFreq float64) {
	if len(channel.Data) < 3 {
		return
	}

	// 计算RC时间常数
	RC := 1.0 / (2.0 * math.Pi * cutoffFreq)
	
	// 计算alpha系数
	dt := 1.0 / p.SampleRate
	alpha := dt / (RC + dt)
	
	// 分配结果数组
	channel.ProcessedData = make([]data.DataPoint, len(channel.Data))
	
	// 第一个点不变
	channel.ProcessedData[0] = data.DataPoint{
		X: channel.Data[0].X,
		Y: channel.Data[0].Y,
	}
	
	// 应用一阶低通滤波器
	for i := 1; i < len(channel.Data); i++ {
		y := channel.ProcessedData[i-1].Y + alpha * (channel.Data[i].Y - channel.ProcessedData[i-1].Y)
		
		channel.ProcessedData[i] = data.DataPoint{
			X: channel.Data[i].X,
			Y: y,
		}
	}
}

// ApplyHighPassFilter 应用高通滤波
func (p *Processor) ApplyHighPassFilter(channel *data.Channel, cutoffFreq float64) {
	if len(channel.Data) < 3 {
		return
	}

	// 计算RC时间常数
	RC := 1.0 / (2.0 * math.Pi * cutoffFreq)
	
	// 计算alpha系数
	dt := 1.0 / p.SampleRate
	alpha := RC / (RC + dt)
	
	// 分配结果数组
	channel.ProcessedData = make([]data.DataPoint, len(channel.Data))
	
	// 第一个点不变
	channel.ProcessedData[0] = data.DataPoint{
		X: channel.Data[0].X,
		Y: channel.Data[0].Y,
	}
	
	// 应用一阶高通滤波器
	for i := 1; i < len(channel.Data); i++ {
		y := alpha * (channel.ProcessedData[i-1].Y + channel.Data[i].Y - channel.Data[i-1].Y)
		
		channel.ProcessedData[i] = data.DataPoint{
			X: channel.Data[i].X,
			Y: y,
		}
	}
}

// ApplyBandPassFilter 应用带通滤波
func (p *Processor) ApplyBandPassFilter(channel *data.Channel, lowCutoff, highCutoff float64) {
	// 创建临时通道进行中间处理
	tempChannel := &data.Channel{
		ID:         "temp",
		Name:       "temp",
		Data:       channel.Data,
		ProcessedData: nil,
		Visible:    true,
		Color:      "#FF0000",
		Scale:      1.0,
		YAxisMin:   -1.0,
		YAxisMax:   1.0,
	}
	
	// 先应用高通滤波器
	p.ApplyHighPassFilter(tempChannel, lowCutoff)
	
	// 更新临时通道的数据
	tempChannel.Data = tempChannel.ProcessedData
	
	// 再应用低通滤波器
	p.ApplyLowPassFilter(tempChannel, highCutoff)
	
	// 将处理结果复制到原通道
	channel.ProcessedData = tempChannel.ProcessedData
}

// ApplyMovingAverage 应用移动平均滤波
func (p *Processor) ApplyMovingAverage(channel *data.Channel, windowSize int) {
	if len(channel.Data) < windowSize {
		return
	}

	// 分配结果数组
	channel.ProcessedData = make([]data.DataPoint, len(channel.Data))
	
	// 前windowSize-1个点保持原值
	for i := 0; i < windowSize-1; i++ {
		channel.ProcessedData[i] = data.DataPoint{
			X: channel.Data[i].X,
			Y: channel.Data[i].Y,
		}
	}
	
	// 应用移动平均
	for i := windowSize - 1; i < len(channel.Data); i++ {
		sum := 0.0
		for j := 0; j < windowSize; j++ {
			sum += channel.Data[i-j].Y
		}
		
		y := sum / float64(windowSize)
		
		channel.ProcessedData[i] = data.DataPoint{
			X: channel.Data[i].X,
			Y: y,
		}
	}
}

// ApplyFFT 应用快速傅里叶变换
func (p *Processor) ApplyFFT(channel *data.Channel) []complex128 {
	// 获取数据点数量
	n := len(channel.Data)
	if n < 2 {
		return nil
	}
	
	// 提取Y值
	yValues := make([]float64, n)
	for i := 0; i < n; i++ {
		yValues[i] = channel.Data[i].Y
	}
	
	// 寻找最接近的2的幂
	powerOfTwo := 1
	for powerOfTwo < n {
		powerOfTwo *= 2
	}
	
	// 如果需要，填充数据
	if n < powerOfTwo {
		paddedData := make([]float64, powerOfTwo)
		copy(paddedData, yValues)
		yValues = paddedData
		n = powerOfTwo
	}
	
	// 应用窗函数（汉宁窗）
	windowedData := make([]float64, len(yValues))
	copy(windowedData, yValues)
	
	// 修正: 手动实现汉宁窗函数 w(n) = 0.5 * (1 - cos(2π*n/(N-1)))
	for i := range windowedData {
		if len(windowedData) > 1 {
			hannCoef := 0.5 * (1 - math.Cos(2*math.Pi*float64(i)/float64(len(windowedData)-1)))
			windowedData[i] = windowedData[i] * hannCoef
		}
	}
	
	// 创建FFT实例
	fft := fourier.NewFFT(n)
	
	// 将实数转换为复数
	complexData := make([]complex128, n)
	for i, v := range windowedData {
		complexData[i] = complex(v, 0)
	}
	
	// 执行FFT
	// 修正: 使用正确的方法计算FFT结果
	result := fft.Coefficients(nil, complexData)
	
	// 计算处理后的数据
	channel.ProcessedData = make([]data.DataPoint, n/2)
	
	// 计算频率分辨率
	freqResolution := p.SampleRate / float64(n)
	
	for i := 0; i < n/2; i++ {
		// 计算幅度谱
		magnitude := cmplx.Abs(result[i])
		
		// 将幅度归一化
		magnitude /= float64(n)
		
		// 如果不是直流分量，乘以2（因为实信号的频谱是对称的）
		if i > 0 {
			magnitude *= 2
		}
		
		// 计算频率
		freq := float64(i) * freqResolution
		
		channel.ProcessedData[i] = data.DataPoint{
			X: freq,
			Y: magnitude,
		}
	}
	
	return result
}

// DetectPeaks 检测峰值
func (p *Processor) DetectPeaks(channel *data.Channel, threshold float64) []int {
	if len(channel.Data) < 3 {
		return nil
	}

	// 存储峰值索引
	peaks := make([]int, 0)
	
	// 检测峰值
	for i := 1; i < len(channel.Data)-1; i++ {
		// 当前点大于阈值且大于相邻点
		if channel.Data[i].Y > threshold &&
		   channel.Data[i].Y > channel.Data[i-1].Y &&
		   channel.Data[i].Y >= channel.Data[i+1].Y {
			peaks = append(peaks, i)
		}
	}
	
	return peaks
}

// CalculateHeartRate 计算心率
func (p *Processor) CalculateHeartRate(channel *data.Channel) float64 {
	// 检测R波峰值
	peaks := p.DetectPeaks(channel, 0.5) // 阈值可能需要调整
	
	if len(peaks) < 2 {
		return 0
	}
	
	// 计算RR间隔的平均值
	totalTime := 0.0
	for i := 1; i < len(peaks); i++ {
		interval := channel.Data[peaks[i]].X - channel.Data[peaks[i-1]].X
		totalTime += interval
	}
	
	// 计算平均RR间隔（秒）
	avgRRInterval := totalTime / float64(len(peaks)-1)
	
	// 计算心率（次/分）
	heartRate := 60.0 / avgRRInterval
	
	return heartRate
}
