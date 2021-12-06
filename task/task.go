package task

import (
	"basic/color"
	"fmt"
	"github.com/robfig/cron"
	"log"
)

type (
	// Task 一个路由的结构
	Task struct {
		Spec string
		Name string
	}

	taskMap map[string]*cron.Cron
)

var tasks taskMap

func init() {
	tasks = taskMap{}
}

// Register 注册任务
func (t Task) Register(cmd func()) {
	if t.Spec == "" {
		log.Panic("spec为空")
		return
	}
	if _, ok := tasks[t.Name]; ok {
		log.Panicf("'%s' 任务已经存在", t.Name)
		return
	}
	//创建任务
	c := cron.New()
	err := c.AddFunc(t.Spec, cmd)
	if err != nil {
		log.Println(err)
		log.Panicf("'%s' 任务创建失败", t.Name)
	}
	c.Start()
	tasks[t.Name] = c
}

// Cancel 取消
func (t Task) Cancel() {
	if c, ok := tasks[t.Name]; ok {
		c.Stop()
		return
	}
	log.Println(fmt.Sprintf("[%s]任务不存在", t.Name))
}

// Run 当主协程能自己维持的化，block不用开启
func Run(block bool) {
	color.Success(fmt.Sprintf("[task] start success，tasks total: %d", len(tasks)))
	if block {
		select {}
	}
}
