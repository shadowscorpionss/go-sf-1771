package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

//Напишите код, в котором несколько горутин увеличивают значение
//целочисленного счётчика и синхронизируют свою работу через канал.
//Нужно предусмотреть возможность настройки количества используемых горутин
//и конечного значения счётчика, до которого его следует увеличивать.

// Попробуйте реализовать счётчик с элементами ООП (в виде структуры и методов структуры).
// Попробуйте реализовать динамическую проверку достижения счётчиком нужного значения.
// В качестве ответа приложите архивный файл с кодом программы из Задания 17.7.1.
func heyHeyCounter(cap, endCount int) (chan<- int, <-chan int) {
	c, d := make(chan int, cap), make(chan int)
	go func() {
		defer close(c)
		counter := 0
		for counter < endCount {
			counter += <-c
			fmt.Println("Считаем: ", counter)
		}
		d <- counter

	}()
	return c, d
}

func heyHeySender(cnt int, send chan<- int) {

	for i := 0; i < cnt; i++ {
		send <- 1
		//искусственная случайная задержка
		p := rand.Intn(100)
		time.Sleep(time.Duration(p) * time.Millisecond)
	}
}

func main() {

	var routineCount, endCount int
	fmt.Println("Введите количество до которого будем считать")
	_, err := fmt.Scanln(&endCount)
	if err != nil {
		panic(err)
	}
	fmt.Println("Введите количество рутин в которых будем считать")
	_, err = fmt.Scanln(&routineCount)
	if err != nil {
		panic(err)
	}

	hc := NewHeyHeyCounterClass()
	outer := hc.Start(routineCount, endCount)

	fmt.Println("Итог: ", <-outer)

	time.Sleep(5 * time.Second)

	fmt.Println("Осталось рутин:", runtime.NumGoroutine())
}
