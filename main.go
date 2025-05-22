package main

import (
	"log"

	"fyne.io/fyne/v2"
	fyneApp "fyne.io/fyne/v2/app" // Renamed to avoid conflict with your app package
	"fyne.io/fyne/v2/theme"
	"github.com/liujiaxin/chartSystem/internal/app"
)

// MyTheme 结构体嵌入了 fyne.Theme，允许我们覆盖默认字体
type MyTheme struct {
	fyne.Theme
}

// Font 返回指定文本样式的字体资源
func (m *MyTheme) Font(style fyne.TextStyle) fyne.Resource {
	fontResource, err := fyne.LoadResourceFromPath("assets/STKAITI.TTF")
	if err != nil {
		log.Printf("错误：无法加载字体，将使用默认字体: %v", err)
		// 如果字体加载失败，回退到 Fyne 的默认字体
		if m.Theme == nil {
			// Fallback if m.Theme is somehow nil, though it shouldn't be with proper initialization
			return theme.DefaultTheme().Font(style)
		}
		return m.Theme.Font(style)
	}
	return fontResource
}

func main() {
	// 创建 Fyne 应用实例
	myFyneApp := fyneApp.New()

	// 创建并设置自定义主题
	// 确保 MyTheme 正确初始化了其嵌入的 Theme 字段
	customTheme := &MyTheme{Theme: theme.DefaultTheme()}
	myFyneApp.Settings().SetTheme(customTheme)

	// 初始化应用程序，传递 fyne.App 实例
	// 您需要修改 internal/app/app.go 中的 NewApp 函数以接受此参数
	chartApp, err := app.NewApp(myFyneApp)
	if err != nil {
		log.Fatalf("应用程序初始化失败: %v", err)
	}

	// 运行应用程序
	// 假设您的 chartApp.Run() 方法内部会调用 myFyneApp.Run() 或者窗口的 ShowAndRun()
	if err := chartApp.Run(); err != nil {
		log.Fatalf("应用程序运行失败: %v", err)
	}
}
