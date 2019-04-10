routine
===

job routine for golang program with the following strategies:

- retry
- repeat
- crontab
- guarantee

more details at the quick start.

## Quick Start

````go

import (
	"context"
	"errors"
	"log"
    "time"
    
    "github.com/x-mod/routine"
)

func main(){
	err := routine.Main(context.TODO(), routine.ExecutorFunc(func(ctx context.Context) error {
		log.Println("main executing begin ...")

		ch1 := routine.Go(ctx, routine.Retry(3, routine.ExecutorFunc(func(arg1 context.Context) error {
			log.Println("Go1 retry begin ...", routine.FromRetry(arg1))
			time.Sleep(1 * time.Second)
			log.Println("Go1 retry end")
			return errors.New("Go1 error")
		})))
		log.Println("Go1 result: ", <-ch1)

		ch2 := routine.Go(ctx, routine.Repeat(2, time.Second, routine.ExecutorFunc(func(arg1 context.Context) error {
			log.Println("Go2 repeat begin ...", routine.FromRepeat(arg1))
			time.Sleep(2 * time.Second)
			log.Println("Go2 repeat end")
			return nil
		})))
		log.Println("Go2 result: ", <-ch2)

		routine.Go(ctx, routine.Repeat(2, time.Second, routine.Guarantee(routine.ExecutorFunc(func(arg1 context.Context) error {
			log.Println("Go4 repeat guarantee begin ...")
			log.Println("Go4 repeat guarantee end")
			return errors.New("Go4 failed")
		}))))

		routine.Go(ctx, Crontab("* * * * *", routine.ExecutorFunc(func(arg1 context.Context) error {
			log.Println("Go3 crontab begin ...", routine.FromCrontab(arg1))
			log.Println("Go3 crontab end")
			return nil
		})))

		ch5 := routine.Go(ctx, Repeat(3, time.Second, routine.Command("echo", "hello", "routine")))
		log.Println("Go5 result: ", <-ch5)

		ch6 := routine.Go(ctx, routine.Timeout(3*time.Second, Command("sleep", "6")))
		log.Println("Go6 timeout result: ", <-ch6)

		ch7 := routine.Go(ctx, routine.Deadline(time.Now().Add(time.Second), Command("sleep", "6")))
		log.Println("Go7 deadline result: ", <-ch7)

		log.Println("main executing end")
		return nil
	}), routine.DefaultCancelInterruptors...)
	log.Println("main exit: ", err)
}

````
