package render

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"sort"

	"github.com/liujiaxin/chartSystem/internal/config"
	"github.com/liujiaxin/chartSystem/internal/data"
	"github.com/liujiaxin/chartSystem/pkg/util"
)

// Renderer 负责将数据渲染为图像
type Renderer struct {
	Width       int
	Height      int
	OffsetX     float64
	ScaleX      float64
	GridVisible bool
	GridColor   color.RGBA
	BackColor   color.RGBA
}

// NewRenderer 创建一个新的渲染器
func NewRenderer(width, height int) *Renderer {
	return &Renderer{
		Width:       width,
		Height:      height,
		OffsetX:     0,
		ScaleX:      1.0,
		GridVisible: true,
		GridColor:   color.RGBA{200, 200, 200, 255},
		BackColor:   color.RGBA{255, 255, 255, 255},
	}
}

// RenderChannel 将通道数据渲染为图像
func (r *Renderer) RenderChannel(channel *data.Channel, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, r.Width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{r.BackColor}, image.Point{}, draw.Src)

	// 如果通道不可见，或者没有数据，则绘制背景（和网格，如果需要）并返回
	if !channel.Visible || len(channel.Data) == 0 {
		if r.GridVisible { // 即便通道不显示数据，如果网格是全局可见的，也应绘制网格背景
			r.drawGrid(img)
		}
		return img // 返回带有背景和可选网格的图像
	}

	// 如果通道可见且有数据，则先绘制网格（如果需要）
	if r.GridVisible {
		r.drawGrid(img)
	}

	// 然后绘制波形
	r.drawWaveform(img, channel)

	return img
}

// RenderAllChannels 渲染所有通道
func (r *Renderer) RenderAllChannels(cfg *config.Config, model *data.DataModel, channelHeight int) *image.RGBA {
	var orderedChannelIDs []string
	if cfg != nil && len(cfg.Channels) > 0 {
		for _, confChannel := range cfg.Channels {
			orderedChannelIDs = append(orderedChannelIDs, confChannel.ID)
		}
	} else if len(model.Channels) > 0 {
		// Fallback: 如果没有config，尝试从datamodel的key排序 (不稳定，但比直接迭代map好一点)
		for id := range model.Channels {
			orderedChannelIDs = append(orderedChannelIDs, id)
		}
		sort.Strings(orderedChannelIDs) // 简单排序
	} else {
		// 完全没有通道信息，绘制一个空白图像
		h := r.Height
		if h <= 0 { h = 300 }
		img := image.NewRGBA(image.Rect(0, 0, r.Width, h))
		draw.Draw(img, img.Bounds(), &image.Uniform{r.BackColor}, image.Point{}, draw.Src)
		if r.GridVisible { r.drawGrid(img) }
		return img
	}

	numAllChannels := len(orderedChannelIDs)
	if numAllChannels == 0 || channelHeight <= 0 { // 再次检查，以防 fallback 也为空
		h := r.Height; if h <= 0 { h = 300 }
		img := image.NewRGBA(image.Rect(0, 0, r.Width, h)); draw.Draw(img, img.Bounds(), &image.Uniform{r.BackColor}, image.Point{}, draw.Src)
		if r.GridVisible { r.drawGrid(img) }; return img
	}

	totalHeight := channelHeight * numAllChannels
	img := image.NewRGBA(image.Rect(0, 0, r.Width, totalHeight))
	draw.Draw(img, img.Bounds(), &image.Uniform{r.BackColor}, image.Point{}, draw.Src)

	currentY := 0
	for _, channelID := range orderedChannelIDs {
		channelData, exists := model.Channels[channelID]
		if !exists {
			// 如果配置中定义的通道在数据模型中找不到，创建一个临时的空通道用于占位
			// (这通常意味着数据文件未加载或不包含此通道)
			tempName := channelID // 默认使用ID作为名字
			if cfg != nil {
				for _, confCh := range cfg.Channels {
					if confCh.ID == channelID {
						tempName = confCh.Name
						break
					}
				}
			}
			channelData = data.NewChannel(channelID, tempName) 
			channelData.Visible = false // 确保这个占位通道默认是不可见的，或根据配置
		}
		
		// RenderChannel 内部会根据 channelData.Visible 来决定是否绘制波形
		channelImg := r.RenderChannel(channelData, channelHeight)
		draw.Draw(img, image.Rect(0, currentY, r.Width, currentY+channelHeight),
		        channelImg, image.Point{0, 0}, draw.Over)
		currentY += channelHeight
	}
	return img
}

// 绘制网格
func (r *Renderer) drawGrid(img *image.RGBA) {
	// 垂直网格线
	for x := 0; x < r.Width; x += 50 {
		for y := 0; y < img.Bounds().Max.Y; y++ {
			img.Set(x, y, r.GridColor)
		}
	}

	// 水平网格线
	for y := 0; y < img.Bounds().Max.Y; y += 50 {
		for x := 0; x < r.Width; x++ {
			img.Set(x, y, r.GridColor)
		}
	}
}

// 绘制波形
func (r *Renderer) drawWaveform(img *image.RGBA, channel *data.Channel) {
	// 将十六进制颜色转换为RGBA
	waveColor, err := util.ParseColor(channel.Color)
	if err != nil {
		// 如果解析失败，使用默认颜色（红色）
		waveColor = color.RGBA{255, 0, 0, 255}
	}

	// 获取通道数据
	channelData := channel.Data
	if len(channelData) == 0 {
		return
	}

	// 计算Y轴缩放比例
	height := img.Bounds().Max.Y
	yScale := float64(height) / (channel.YAxisMax - channel.YAxisMin)

	// 计算可见数据范围
	startIdx := 0
	endIdx := len(channelData) - 1

	// 如果数据点太多，进行抽样
	if endIdx > r.Width*2 {
		// 简单抽样：每N个点取一个
		samplingRate := int(math.Ceil(float64(endIdx) / float64(r.Width*2)))
		sampledData := make([]data.DataPoint, 0, r.Width*2)

		for i := startIdx; i <= endIdx; i += samplingRate {
			sampledData = append(sampledData, channelData[i])
		}

		channelData = sampledData
		endIdx = len(channelData) - 1
	}

	// 绘制波形线段
	if endIdx > 0 {
		for i := 0; i < endIdx; i++ {
			// 计算坐标
			x1 := int((channelData[i].X - r.OffsetX) * r.ScaleX)
			y1 := height - int((channelData[i].Y-channel.YAxisMin)*yScale)

			x2 := int((channelData[i+1].X - r.OffsetX) * r.ScaleX)
			y2 := height - int((channelData[i+1].Y-channel.YAxisMin)*yScale)

			// 确保坐标在有效范围内
			if x1 >= 0 && x1 < r.Width && y1 >= 0 && y1 < height &&
				x2 >= 0 && x2 < r.Width && y2 >= 0 && y2 < height {
				// 绘制线段
				drawLine(img, x1, y1, x2, y2, waveColor)
			}
		}
	}
}

// SetViewport 设置视口参数（滚动和缩放）
func (r *Renderer) SetViewport(offsetX, scaleX float64) {
	r.OffsetX = offsetX
	r.ScaleX = scaleX
}

// 绘制直线（Bresenham算法）
func drawLine(img *image.RGBA, x0, y0, x1, y1 int, clr color.RGBA) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)

	sx, sy := 1, 1
	if x0 >= x1 {
		sx = -1
	}
	if y0 >= y1 {
		sy = -1
	}

	err := dx - dy

	for {
		// 设置像素
		img.Set(x0, y0, clr)

		// 检查是否到达终点
		if x0 == x1 && y0 == y1 {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

// 取绝对值
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
