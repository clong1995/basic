package task

import (
	"basic/color"
	"fmt"
	"github.com/robfig/cron"
	"log"
)

type (
	Server struct {
		Block bool //当主协程能自己维持，block不用开启
	}
)

type (
	// Task 一个路由的结构
	Task struct {
		Spec      string
		Name      string
		Immediate bool //要立即执行的任务
	}

	taskMap       map[string]*cron.Cron
	taskImmediate []func()
)

var tasks taskMap
var tasksImmediate taskImmediate

func init() {
	tasks = taskMap{}
	tasksImmediate = taskImmediate{}
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
	if t.Immediate {
		tasksImmediate = append(tasksImmediate, cmd)
	}
}

// Cancel 取消
func (t Task) Cancel() {
	if c, ok := tasks[t.Name]; ok {
		c.Stop()
		return
	}
	log.Println(fmt.Sprintf("[%s]任务不存在", t.Name))
}

func (s Server) Run() {
	color.Success(fmt.Sprintf("[task] start success，tasks total: %d", len(tasks)))
	if len(tasksImmediate) > 0 {
		for _, task := range tasksImmediate {
			task()
		}
	}
	if s.Block {
		select {}
	}
}

//参考cron
//https://segmentfault.com/a/1190000039647260
//https://blog.csdn.net/qq_39135287/article/details/95664533?utm_medium=distribute.pc_relevant.none-task-blog-BlogCommendFromMachineLearnPai2-2.control&depth_1-utm_source=distribute.pc_relevant.none-task-blog-BlogCommendFromMachineLearnPai2-2.control
