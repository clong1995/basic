// Package page 获取一张动态网页
package page

import (
	"basic/color"
	"context"
	"github.com/chromedp/chromedp"
	"log"
	"time"
)

var Page *server

type server struct {
}

type Server struct {
}

// Doc 获取页面
func (s server) Doc(url string, selector string) (string, error) {

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-extensions", false),
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent(`'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36'`),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(
		context.Background(),
		opts...,
	)
	defer cancel()

	// create context
	chromeCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	// 执行一个空task, 用提前创建Chrome实例
	err := chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)
	if err != nil {
		log.Println(err)
		return "", err
	}

	//创建一个上下文，超时时间为40s
	timeoutCtx, cancel := context.WithTimeout(chromeCtx, 40*time.Second)
	defer cancel()

	var htmlContent string

	err = chromedp.Run(timeoutCtx,
		chromedp.Navigate(url),
		chromedp.OuterHTML(selector, &htmlContent, chromedp.NodeVisible),
	)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return htmlContent, nil
}

func (s Server) Run() {
	//防止多次创建
	if Page != nil {
		return
	}
	//创建对象
	Page = new(server)
	color.Success("[reptile page] create success")
}
