package timerTask

import (
	"fmt"
	"testing"
	"time"
)

func printTest1()  {
	fmt.Println("test1:", time.Now().Second())
}

func printTest2()  {
	fmt.Println("test2:", time.Now().Second())
}

func printTest3()  {
	fmt.Println("test3:", time.Now().Second())
}

func printTest4()  {
	fmt.Println("test4:", time.Now().Second())
}

func printTest5()  {
	fmt.Println("test5:", time.Now().Second())
}

// 测试开始和结束任务
func TestNewTimerTask(t *testing.T) {
	t1 := 2 * time.Second
	tasks := []func(){printTest1, printTest2, printTest3, printTest4, printTest5}
	timer := NewTimerTask(t1, tasks)
	timer.Waiter.Add(1)
	timer.start()
	//go func() {
	//	select {
	//	case <-time.After(5*time.Second):
	//		timer.stop()
	//		timer.Waiter.Done()
	//		return
	//	}
	//
	//}()
	timer.Waiter.Wait()
	/*
	=== RUN   TestNewTimerTask
	test5: 10
	test2: 10
	test1: 10
	test3: 10
	test4: 10
	test1: 11
	test3: 11
	test2: 11
	test4: 11
	test5: 11
	test5: 12
	test2: 12
	test1: 12
	test3: 12
	test4: 12
	test5: 13
	test2: 13
	test1: 13
	test3: 13
	test4: 13
	time="2021-12-16T14:52:14+08:00" level=info msg="Timed task stop"
	--- PASS: TestNewTimerTask (5.00s)
	PASS
	*/

}

// 测试重置定时时间任务
func TestReset(t *testing.T) {
	t1 := 1 * time.Second
	tasks := []func(){printTest1, printTest2, printTest3, printTest4, printTest5}
	timer := NewTimerTask(t1, tasks)
	timer.Waiter.Add(3)
	timer.start()
	go func() {
		select {
		case <-time.After(10*time.Second):
			timer.stop()
			timer.Waiter.Done()
			return
		}

	}()
	go func() {
		select {
		case <-time.After(5*time.Second):
			timer.Reset(2*time.Second)
			timer.Waiter.Done()
			return
		}

	}()
	timer.Waiter.Wait()
	/*
	=== RUN   TestReset
	test5: 55
	test1: 55
	test2: 55
	test3: 55
	test4: 55
	test5: 56
	test3: 56
	test4: 56
	test1: 56
	test2: 56
	test5: 57
	test3: 57
	test4: 57
	test1: 57
	test2: 57
	test5: 58
	test3: 58
	test2: 58
	test1: 58
	test4: 58
	time="2021-12-16T14:50:59+08:00" level=info msg="Timed task stop"
	test5: 1
	test3: 1
	test2: 1
	test4: 1
	test1: 1
	test5: 3
	test1: 3
	test2: 3
	test3: 3
	test4: 3
	--- PASS: TestReset (10.00s)
	PASS
	*/
}

// 测试定时开始定时任务
func TestTimerStart(t *testing.T) {
	t1 := 1 * time.Second
	ts := time.Now()
	tasks := []func(){printTest1, printTest2, printTest3, printTest4, printTest5}
	timer := NewTimerTask(t1, tasks, SetTimerStart(3*time.Second))
	timer.TimerStart()
	t.Logf("定时时间为：%v", time.Now().Sub(ts))
	timer.Waiter.Add(2)
	timer.start()
	go func() {
		select {
		case <-time.After(5*time.Second):
			timer.stop()
			timer.Waiter.Done()
			return
		}

	}()
	timer.Waiter.Wait()
	/*
	=== RUN   TestTimerStart
	    timer_task_test.go:146: 定时时间为：3.003182964s
	time="2021-12-16T15:05:58+08:00" level=info msg="The timer task has started"
	test2: 59
	test1: 59
	test3: 59
	test4: 59
	test5: 59
	test5: 0
	test3: 0
	test4: 0
	test1: 0
	test2: 0
	test5: 1
	test1: 1
	test2: 1
	test3: 1
	test4: 1
	test5: 2
	test1: 2
	test2: 2
	test3: 2
	test4: 2
	time="2021-12-16T15:06:03+08:00" level=info msg="Timed task stop"
	--- PASS: TestTimerStart (8.01s)
	PASS
	*/
}

// 测试定时结束定时任务
func TestTimerStop(t *testing.T) {
	t1 := 1 * time.Second
	tasks := []func(){printTest1, printTest2, printTest3, printTest4, printTest5}
	timer := NewTimerTask(t1, tasks, SetTimerStop(3*time.Second))
	timer.Waiter.Add(1)
	timer.start()
	time.Sleep(1*time.Second)
	ts := time.Now()
	timer.TimerStop()
	t.Logf("定时时间为：%v", time.Now().Sub(ts))
	timer.Waiter.Wait()
	/*
		=== RUN   TestTimerStop
	test5: 8
	test2: 8
	test1: 8
	test3: 8
	test4: 8
	test5: 9
	test2: 9
	test1: 9
	test3: 9
	test4: 9
	test5: 10
	test1: 10
	test2: 10
	test3: 10
	test4: 10
	time="2021-12-16T15:11:11+08:00" level=info msg="Timed task stop"
	    timer_task_test.go:199: 定时时间为：3.001107379s
	--- PASS: TestTimerStop (4.00s)
	PASS
	*/
}