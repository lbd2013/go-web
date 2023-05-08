// Package console 命令行辅助方法
package console

import (
	"fmt"
	"goweb/pkg/app"
	"os"
	"time"

	"github.com/mgutz/ansi"
)

// Log 打印一条消息，白色输出
func Log(msg string) {
	currTime := time.Now().Format("2006-01-02 15:04:05")
	msg = fmt.Sprintf("%s\t %s", currTime, msg)
	colorOut(msg, "white")
}

// Success 打印一条成功消息，绿色输出
func Success(msg string) {
	currTime := time.Now().Format("2006-01-02 15:04:05")
	msg = fmt.Sprintf("%s\t %s", currTime, msg)
	colorOut(msg, "green")
}

// Error 打印一条报错消息，红色输出
func Error(msg string) {
	currTime := time.Now().Format("2006-01-02 15:04:05")
	msg = fmt.Sprintf("%s\t %s", currTime, msg)
	colorOut(msg, "red")
}

// Warning 打印一条提示消息，黄色输出
func Warning(msg string) {
	currTime := time.Now().Format("2006-01-02 15:04:05")
	msg = fmt.Sprintf("%s\t %s", currTime, msg)
	colorOut(msg, "yellow")
}

// Exit 打印一条报错消息，并退出 os.Exit(1)
func Exit(msg string) {
	currTime := time.Now().Format("2006-01-02 15:04:05")
	msg = fmt.Sprintf("%s\t %s", currTime, msg)
	Error(msg)
	os.Exit(1)
}

// ExitIf 语法糖，自带 err != nil 判断
func ExitIf(err error) {
	if err != nil {
		Exit(err.Error())
	}
}

// colorOut 内部使用，设置高亮颜色
func colorOut(message, color string) {
	if app.IsLocalNeedColor() || app.IsTestingNeedColor() || app.IsProductionNeedColor() {
		fmt.Fprintln(os.Stdout, ansi.Color(message, color))
	} else {
		fmt.Fprintln(os.Stdout, message)
	}
}
