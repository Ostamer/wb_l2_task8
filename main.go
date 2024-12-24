package main

import (
	"fmt"
	"time"
)

// Функция для завершения всех каналов при завершении одного из них
func orChannels(channels ...<-chan interface{}) <-chan interface{} {
	// Если кол-во каналов ноль завершаем выполнение или если 1 возвращаем только его
	switch len(channels) {
	case 0:
		closedChan := make(chan interface{})
		close(closedChan)
		return closedChan
	case 1:
		return channels[0]
	}
	// Создание группы для всех каналов
	orChannelsConfirm := make(chan interface{})
	// Закрытие группы для всех каналов если первый канал завершился
	// И рекурсивный запуск текущей функии для всех оставшихся каналов
	go func() {
		defer close(orChannelsConfirm)
		select {
		case <-channels[0]:
		case <-orChannels(channels[1:]...):
		}
	}()

	return orChannelsConfirm
}

// Основной код
func main() {
	// Функция для безопасного закрытия канала после истечения времени
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}
	// Начала отсчета времени
	start := time.Now()
	// Запуск фукнкции со всеми каналами
	<-orChannels(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	// Вывод
	fmt.Printf("fone after %v\n", time.Since(start))
}
