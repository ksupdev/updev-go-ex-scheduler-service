# updev-go-ex-scheduler-service
## setup project
1. setup module ``go mod init github.com/ksupdev/updev-go-ex-scheduler-service``

## implement project
1. ``context.go`` interface is defind function for manage context (Log,GetParam,Response,ReadInput,ReadInputs)
2. ``micoservice.go`` has function for support this service types
    - ``exitChannel`` เป็น Chanel type อยู่ใน Microservice struct จะถูกใช้ในการ 
2. ``context_scheduler.go`` is implementation of context.go


## Noted

- main จะทำการเรียกใช้งาน Microservice Schedule โดยใน service จะมีการกำหนด timer ซึ่งเราจะสามารถกำหนด Duration เพื่อกำหนดความถี่ในการทำงานได้
- โดย ms.Schedule จะมีการ Return channel exitChan เพื่อสำหรับให้ main func สามารถทำการหยุดในส่วนของ timer ใน go routine ได้
- main func จะทำการกำหนด exitScheduler (instance of exitChan) = ture ทันทีเมื่อ method main กำลังจะหยุดทำงานหรือก็คือในกรณีนี้ทันทีที่เรา กด Kill process นั้นเอง เนื่องจากเราได้มีการ Implement ส่วนในการรอรับคำสั่งในการ  kill process ไว้ ``signal.Notify(osQuit, syscall.SIGTERM, syscall.SIGINT)``
- และเมื่อ exitScheduler = true ระบบจะหยุดตัว timer และ job ที่กำลัง run ทันที

```golang
[filename:main.go]
func main() {
	ms := NewMicroservice()

	timer := 1 + time.Second
	exitScheduler := ms.Schedule(timer, func(ctx IContext) error {
    ....
	})

	defer func() { exitScheduler <- true }()
    ...
}


[microservice.go]
func (ms *Microservice) Schedule(timer time.Duration.... ) chan bool {
    go func() {
        t := time.NewTicker(timer)
        ....
        go func() {
            <-exitChan // Block process until exitChan has value
            t.Stop()
        }

        for {
            ...
            Execute process ()
            ...
        }

    }
}

```