package app

import (
	"log"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"github.com/liujiaxin/chartSystem/internal/config"
	"github.com/liujiaxin/chartSystem/internal/data"
	"github.com/liujiaxin/chartSystem/internal/ui"
	"github.com/liujiaxin/chartSystem/pkg/fileio"
)

// App 表示图表应用程序
type App struct {
	fyneApp    fyne.App
	Config     *config.Config
	DataModel  *data.DataModel
	MainWindow *ui.MainWindow
}

// NewApp 创建并初始化一个新的应用程序实例
func NewApp(theFyneApp fyne.App) (*App, error) {
	// 加载配置
	cfg, err := config.LoadConfig("config/config.xml")
	if err != nil {
		log.Printf("配置加载失败，使用默认配置: %v", err)
		// 如果配置加载失败，使用默认配置
		cfg = getDefaultConfig()
	}

	// 创建数据模型
	dataModel := data.NewDataModel()

	// 为每个配置的通道创建数据模型通道
	for _, channelCfg := range cfg.Channels {
		channel := data.NewChannel(channelCfg.ID, channelCfg.Name)
		channel.Color = channelCfg.Color
		channel.Scale = channelCfg.Scale
		channel.Visible = channelCfg.Visible
		channel.YAxisMin = channelCfg.YAxisMin
		channel.YAxisMax = channelCfg.YAxisMax

		dataModel.AddChannel(channel)
	}

	// generateSimulatedData(dataModel) // <<< 注释掉或移除此行，以实现初始无曲线状态

	// 创建主窗口，传递 fyne.App 实例
	mainWindow, err := ui.NewMainWindow(theFyneApp, cfg, dataModel)
	if err != nil {
		return nil, err
	}

	return &App{
		fyneApp:    theFyneApp,
		Config:     cfg,
		DataModel:  dataModel,
		MainWindow: mainWindow,
	}, nil
}

// Run 运行应用程序
func (a *App) Run() error {
	// 显示主窗口
	return a.MainWindow.Show()
}

// LoadEDFFile 加载EDF文件
func (a *App) LoadEDFFile(path string) error {
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
	a.DataModel = data.NewDataModel()

	// 加载每个信号到通道
	for i := 0; i < numSignals && i < 4; i++ {
		// 获取信号信息
		label, _, physMin, physMax := edfReader.GetChannelInfo(i) // 忽略未使用的physDim变量

		// 创建通道
		channel := data.NewChannel(strconv.Itoa(i), label)
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
		a.DataModel.AddChannel(channel)
	}

	return nil
}

// 生成模拟数据
func generateSimulatedData(model *data.DataModel) {
	// 生成心电数据
	if channel, ok := model.Channels["1"]; ok {
		fileio.CreateSimulatedEDFData(channel, "ecg", 10.0, 250.0)
	}

	// 生成血压数据
	if channel, ok := model.Channels["2"]; ok {
		fileio.CreateSimulatedEDFData(channel, "bp", 10.0, 250.0)
	}

	// 生成血氧数据
	if channel, ok := model.Channels["3"]; ok {
		fileio.CreateSimulatedEDFData(channel, "spo2", 10.0, 250.0)
	}

	// 生成呼吸数据
	if channel, ok := model.Channels["4"]; ok {
		fileio.CreateSimulatedEDFData(channel, "resp", 10.0, 250.0)
	}
}

// 获取默认配置
func getDefaultConfig() *config.Config {
	return &config.Config{
		Channels: []config.Channel{
			{
				ID:       "1",
				Name:     "心电",
				Color:    "#FF0000",
				Scale:    1.0,
				Visible:  true,
				YAxisMin: -1.0,
				YAxisMax: 1.0,
			},
			{
				ID:       "2",
				Name:     "血压",
				Color:    "#00FF00",
				Scale:    1.0,
				Visible:  true,
				YAxisMin: -1.0,
				YAxisMax: 1.0,
			},
			{
				ID:       "3",
				Name:     "血氧",
				Color:    "#0000FF",
				Scale:    1.0,
				Visible:  true,
				YAxisMin: -1.0,
				YAxisMax: 1.0,
			},
			{
				ID:       "4",
				Name:     "呼吸",
				Color:    "#FFFF00",
				Scale:    1.0,
				Visible:  true,
				YAxisMin: -1.0,
				YAxisMax: 1.0,
			},
		},
		Display: config.Display{
			GridVisible: true,
			RefreshRate: 30,
			TimeScale:   1.0,
		},
	}
}
