package practise

import "time"

func ScheduleTest() {
	after5Second()

}

func every5Second() {
	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for range ticker.C {
			println("test")
		}
	}()
	time.Sleep(time.Minute)
}

func after5Second()  {
	timer1:=time.NewTimer(time.Second*5)
	<-timer1.C
	println("test")
}
