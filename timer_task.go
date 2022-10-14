package timerTask

//定时单个或者对个任务：
//1. 定时定时任务
//2. 停止定时任务
//3. 重置定时时间，并执行定时任务
//4. 定时开启定时任务
//5. 定时结束定时任务


import (
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

// timer配置
type TimerConfig struct {
	timing       time.Duration      // 定时时间
	startTiming  time.Duration      // 定时开启定时任务时间
	stopTiming   time.Duration      // 定时停止定时任务时间
	running   	 bool               // 定时器运行的状态,运作中为true/没有运行为false
	tasks        []func()           // 定时任务
	stopChan     chan struct{}      // 停止定时器日任务
	runningMu    sync.Mutex
	Waiter       sync.WaitGroup
}

type Option func(*TimerConfig)

// 添加定时开启定时任务时间
func SetTimerStart(d time.Duration) Option {
	return func(t *TimerConfig) {
		t.startTiming = d
	}
}

// 添加定时停止定时任务时间
func SetTimerStop(d time.Duration) Option {
	return func(t *TimerConfig) {
		t.stopTiming = d
	}
}

// 添加定时时间
func (t *TimerConfig)addTiming(d time.Duration) {
	t.timing = d
}

// 任务停止
func (t *TimerConfig)stop()  {
	t.runningMu.Lock()
	defer t.runningMu.Unlock()
	if !t.running{
		log.Info("The timer task is not started and cannot be stopped")
		return
	}
	t.stopChan <- struct{}{}
	t.running = false

}

// 任务执行
func (t *TimerConfig)run()  {
	t.running = true
	for {
		timer := time.NewTimer(t.timing)
		select {
		case <- timer.C:
			for _, task := range t.tasks{
				go task()
			}
		case <- t.stopChan:
			timer.Stop()
			log.Info("Timed task stop")
			return
		}
	}
}

// 任务开始
func (t *TimerConfig)start()  {
	t.runningMu.Lock()
	defer t.runningMu.Unlock()
	t.Waiter.Add(1)
	defer t.Waiter.Done()
	if t.running{
		log.Info("The timer task has started")
		return
	}
	go t.run()

}

func (t *TimerConfig)Stop() {
	t.stop()
}

func (t *TimerConfig)Start() {
	  t.start()
}

// 定时开启定时任务
func (t *TimerConfig)TimerStart() {
	// 如果定时开启时间未开启或定时任务正在执行，直接返回
	if t.startTiming <= 0 || t.running{
		log.Info("The scheduled start time is not open or the scheduled task is being executed")
		return
	}
	t.Waiter.Add(1)
	go func() {
		select {
		case <- time.After(t.startTiming):
			t.start()
			t.Waiter.Done()
			return
		}
	}()
	t.Waiter.Wait()
}

// 定时停止定时任务
func (t *TimerConfig)TimerStop() {
	// 如果定时结束时间未开启或定时任务未执行，直接返回
	if t.stopTiming <= 0 || !t.running{
		log.Info("The timing end time is not turned on or the timing task is not executed")
		return
	}
	t.Waiter.Add(1)
	go func() {
		select {
		case <- time.After(t.stopTiming):
			t.stop()
			t.Waiter.Done()
			return
		}
	}()
	t.Waiter.Wait()
}

// 重置定时时间
func (t *TimerConfig)Reset(d time.Duration)  {
	if !t.running{
		log.Info("The timer task is not started and cannot be stopped")
		return
	}
	t.addTiming(d)
	t.stop()
	t.start()
}

// 得到运行状态
func (t *TimerConfig)GetRunStatus() bool {
	return t.running
}

func NewTimerTask(d time.Duration, tasks []func(), opts ... Option) *TimerConfig {
	timer := &TimerConfig{
		timing:d,
		startTiming:-1,
		stopTiming:-1,
		running:false,
		tasks: tasks,
		stopChan:make(chan struct{}),
	}
	for _, opt := range opts{
		opt(timer)
	}
	return timer
}
