package ui

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"github.com/liujiaxin/chartSystem/internal/config"
	"github.com/liujiaxin/chartSystem/internal/data"
	"github.com/liujiaxin/chartSystem/internal/render"
	"github.com/liujiaxin/chartSystem/pkg/fileio"
)

// MainWindow 表示应用程序的主窗口
type MainWindow struct {
	app         fyne.App
	window      fyne.Window
	dataModel   *data.DataModel
	renderer    *render.Renderer
	channelList *widget.List
	chartView   *ChartView
	config      *config.Config
}

// NewMainWindow 创建一个新的主窗口
func NewMainWindow(theFyneApp fyne.App, cfg *config.Config, dataModel *data.DataModel) (*MainWindow, error) {
	// 创建应用程序窗口，使用传入的 fyne.App 实例
	window := theFyneApp.NewWindow("Chart系统")

	// 创建渲染器
	renderer := render.NewRenderer(800, 600)

	var chartView *ChartView // 提前声明 chartView 变量

	// 创建通道列表
	channelList := widget.NewList(
		func() int { return len(cfg.Channels) },
		func() fyne.CanvasObject {
			// 创建一个容器，包含通道名称和一个显示开关
			nameLabel := widget.NewLabel("模板项")
			visibleCheck := widget.NewCheck("", nil)
			return container.NewBorder(nil, nil, nil, visibleCheck, nameLabel)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			// 获取通道配置
			channel := cfg.Channels[id]
			
			// 获取通道数据模型
			dataChannel, exists := dataModel.Channels[channel.ID]
			
			// 设置通道名称
			nameLabel := obj.(*fyne.Container).Objects[0].(*widget.Label)
			nameLabel.SetText(channel.Name)
			
			// 设置可见性开关
			visibleCheck := obj.(*fyne.Container).Objects[1].(*widget.Check)
			
			// 如果通道存在，使用其可见性状态
			if exists {
				visibleCheck.Checked = dataChannel.Visible
			} else {
				visibleCheck.Checked = channel.Visible
			}
			
			// 设置可见性变更回调
			visibleCheck.OnChanged = func(checked bool) {
				if dataChannel, ok := dataModel.Channels[channel.ID]; ok {
					dataChannel.Visible = checked
					// 直接刷新需要重绘的对象
					// visibleCheck.Refresh() // 可以保留，确保复选框状态更新
					if chartView != nil { // 确保 chartView 已被初始化
						chartView.Refresh() 
					}
				}
			}
		},
	)
	
	// 初始化 chartView (在 channelList 定义之后，因为它可能在回调中被引用)
	chartView = NewChartView(dataModel, renderer, cfg)

	// 设置列表选中回调
	channelList.OnSelected = func(id widget.ListItemID) {
		// 延迟一下取消选中状态，以便UI可以显示点击效果
		go func() {
			// 直接调用对象自身的Refresh方法
			channelList.UnselectAll()
			channelList.Refresh()
		}()
	}

	// 创建布局
	split := container.NewHSplit(
		container.NewBorder(
			widget.NewLabel("通道列表"),
			nil, nil, nil,
			container.NewScroll(channelList),
		),
		chartView,
	)
	split.SetOffset(0.2)

	// 创建菜单
	mainMenu := fyne.NewMainMenu(
		fyne.NewMenu("文件",
			fyne.NewMenuItem("打开", func() {
				// 创建文件对话框
				dialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
					if err != nil {
						dialog.ShowError(err, window)
						return
					}
					if reader == nil {
						return // 用户取消
					}

					// 获取文件路径
					filePath := reader.URI().Path()
					reader.Close()

					// 加载EDF文件
					err = loadEDFFile(filePath, dataModel)
					if err != nil {
						dialog.ShowError(err, window)
						return
					}

					// 刷新视图
					chartView.dataModel = dataModel
					channelList.Refresh()
					chartView.Refresh()
				}, window)
				dialog.SetFilter(storage.NewExtensionFileFilter([]string{".edf"}))
				dialog.Show()
			}),
			fyne.NewMenuItem("保存配置", func() {}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("退出", func() { theFyneApp.Quit() }),
		),
		fyne.NewMenu("视图",
			fyne.NewMenuItem("缩放重置", func() {
				// 重置图表视图的缩放和偏移
				chartView.scaleX = 1.0
				chartView.offsetX = 0.0
				chartView.RefreshView()
			}),
			fyne.NewMenuItem("显示网格", func() {
				// 切换网格显示状态
				renderer.GridVisible = !renderer.GridVisible
				chartView.RefreshView()
			}),
		),
		fyne.NewMenu("工具",
			fyne.NewMenuItem("设置", func() {}),
		),
		fyne.NewMenu("帮助",
			fyne.NewMenuItem("关于", func() {
				dialog.ShowInformation("关于", "Chart系统 v1.0\n用于显示和分析医疗数据的软件", window)
			}),
		),
	)
	window.SetMainMenu(mainMenu)

	// 设置窗口内容
	window.SetContent(split)
	window.Resize(fyne.NewSize(1024, 768))

	// 在窗口显示前，主动刷新一次 ChartView，确保初始数据能绘制
	if chartView != nil {
		chartView.Refresh()
	}

	return &MainWindow{
		app:         theFyneApp,
		window:      window,
		dataModel:   dataModel,
		renderer:    renderer,
		channelList: channelList,
		chartView:   chartView,
		config:      cfg,
	}, nil
}

// Show 显示主窗口并运行应用事件循环
func (w *MainWindow) Show() error {
	w.window.Show()
	w.app.Run()
	return nil
}

// ChartView 表示图表视图
type ChartView struct {
	widget.BaseWidget
	dataModel *data.DataModel
	renderer  *render.Renderer
	config    *config.Config
	image     *canvas.Image
	offsetX   float64
	scaleX    float64
}

// 实现Tappable接口
func (c *ChartView) Tapped(pe *fyne.PointEvent) {
	// 单击事件 - 可以在这里添加功能
}

// 实现SecondaryTappable接口
func (c *ChartView) TappedSecondary(pe *fyne.PointEvent) {
	// 右键单击事件 - 可以添加上下文菜单
}

// 实现Scrollable接口
func (c *ChartView) Scrolled(se *fyne.ScrollEvent) {
	// 鼠标滚轮事件 - 控制缩放
	if se.Scrolled.DY > 0 {
		// 放大
		c.scaleX *= 1.1
	} else {
		// 缩小
		c.scaleX *= 0.9
	}
	
	// 重绘
	c.Refresh()
}

// 实现Draggable接口
func (c *ChartView) Dragged(de *fyne.DragEvent) {
	// X轴拖拽 - 控制平移
	c.offsetX -= float64(de.Dragged.DX) / c.scaleX
	
	// 重绘
	c.Refresh()
}

func (c *ChartView) DragEnd() {
	// 拖拽结束
}

func (c *ChartView) DragBegin() {
	// 拖拽开始
}

// NewChartView 创建一个新的图表视图
func NewChartView(dataModel *data.DataModel, renderer *render.Renderer, cfg *config.Config) *ChartView {
	c := &ChartView{
		dataModel: dataModel,
		renderer:  renderer,
		config:    cfg,
		offsetX:   0,
		scaleX:    1.0,
	}
	c.ExtendBaseWidget(c)
	return c
}

// CreateRenderer 实现Widget接口，创建渲染器
func (c *ChartView) CreateRenderer() fyne.WidgetRenderer {
	// 创建一个空白图像
	img := image.NewRGBA(image.Rect(0, 0, c.renderer.Width, c.renderer.Height))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	
	rasterImg := canvas.NewImageFromImage(img)
	rasterImg.FillMode = canvas.ImageFillContain
	rasterImg.ScaleMode = canvas.ImageScaleSmooth
	
	c.image = rasterImg
	
	// 返回渲染器
	return &chartViewRenderer{
		chartView: c,
		image:     rasterImg,
		objects:   []fyne.CanvasObject{rasterImg},
	}
}

// MinSize 返回最小尺寸
func (c *ChartView) MinSize() fyne.Size {
	return fyne.NewSize(400, 300)
}

// Refresh 刷新视图
func (c *ChartView) RefreshView() {
	c.Refresh()
}

// 图表视图渲染器
type chartViewRenderer struct {
	chartView *ChartView
	image     *canvas.Image
	objects   []fyne.CanvasObject
}

// MinSize 返回最小尺寸
func (r *chartViewRenderer) MinSize() fyne.Size {
	return r.chartView.MinSize()
}

// Layout 布局
func (r *chartViewRenderer) Layout(size fyne.Size) {
	r.image.Resize(size)
}

// Refresh 刷新
func (r *chartViewRenderer) Refresh() {
	r.updateImage()
	r.image.Refresh()
	canvas.Refresh(r.chartView)
}

// Objects 返回对象列表
func (r *chartViewRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// Destroy 销毁
func (r *chartViewRenderer) Destroy() {}

// 更新图像
func (r *chartViewRenderer) updateImage() {
	width := int(r.chartView.Size().Width)
	height := int(r.chartView.Size().Height)
	
	if width <= 0 || height <= 0 {
		width = 400 
		height = 300
	}
	
	r.chartView.renderer.Width = width
	// r.chartView.renderer.Height = height // Renderer 的 Height 是总高度，由其内部管理，或不设置
	
	numExpectedChannels := 0
	if r.chartView.config != nil && len(r.chartView.config.Channels) > 0 {
		numExpectedChannels = len(r.chartView.config.Channels)
	} else {
        numExpectedChannels = 4 // Default if no config
    }
    if numExpectedChannels == 0 { numExpectedChannels = 1 } 

	channelHeight := height / numExpectedChannels
	if channelHeight < 50 { channelHeight = 50 }
	
	// RenderAllChannels 应该能处理 dataModel 为空或无通道的情况
	// if len(r.chartView.dataModel.Channels) == 0 && (r.chartView.config == nil || len(r.chartView.config.Channels) == 0) {
	// 	 // 之前尝试在这里画空网格，但RenderAllChannels应该自己处理
	// 	 emptyImg := image.NewRGBA(image.Rect(0, 0, width, height))
	// 	 draw.Draw(emptyImg, emptyImg.Bounds(), &image.Uniform{r.chartView.renderer.BackColor}, image.Point{}, draw.Src)
	// 	 if r.chartView.renderer.GridVisible {
	// 		 // r.chartView.renderer.drawGrid(emptyImg) // renderer 没有导出的 drawGrid
	// 	 }
	// 	 r.image.Image = emptyImg
	// 	 return
	// }
	
r.chartView.renderer.SetViewport(r.chartView.offsetX, r.chartView.scaleX)
	img := r.chartView.renderer.RenderAllChannels(r.chartView.config, r.chartView.dataModel, channelHeight)
	r.image.Image = img
}

// loadEDFFile 加载EDF文件
func loadEDFFile(path string, model *data.DataModel) error {
	// 检查文件是否存在
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return err
	}

	// 打开EDF文件
	edfReader, err := fileio.OpenEDF(path)
	if err != nil {
		return err
	}
	defer edfReader.Close()

	// 获取信号数量
	numSignals := edfReader.GetNumSignals()

	// 清空数据模型
	*model = *data.NewDataModel()

	// 加载每个信号到通道
	for i := 0; i < numSignals && i < 4; i++ {
		// 获取信号信息
		label, _, physMin, physMax := edfReader.GetChannelInfo(i)

		// 创建通道
		channel := data.NewChannel(strconv.Itoa(i+1), label)
		channel.YAxisMin = physMin
		channel.YAxisMax = physMax

		// 设置颜色
		switch i {
		case 0:
			channel.Color = "#FF0000" // 红色
		case 1:
			channel.Color = "#00FF00" // 绿色
		case 2:
			channel.Color = "#0000FF" // 蓝色
		case 3:
			channel.Color = "#FFFF00" // 黄色
		}

		// 加载信号数据
		if err := edfReader.LoadSignalToChannel(i, channel); err != nil {
			log.Printf("加载信号%d失败: %v", i, err)
			continue
		}

		// 添加通道到数据模型
		model.AddChannel(channel)
	}

	return nil
}
