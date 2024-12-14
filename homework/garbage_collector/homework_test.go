package main

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func Trace(stacks [][]uintptr) []uintptr {
	visited := make(map[uintptr]bool)
	added := make(map[uintptr]bool)
	var pointers []uintptr

	// Собираем все указатели из стеков в порядке их появления
	for _, stack := range stacks {
		for _, ptr := range stack {
			if ptr == 0 {
				continue
			}
			// Добавляем указатель в pointers, если он ещё не добавлен
			if !added[ptr] {
				pointers = append(pointers, ptr)
				added[ptr] = true
			}
		}
	}

	// Обрабатываем указатели и добавляем достижимые объекты
	index := 0
	for index < len(pointers) {
		ptr := pointers[index]

		if visited[ptr] {
			index++
			continue
		}
		visited[ptr] = true
		fmt.Printf("\nОбрабатываем указатель: 0x%x\n", ptr)

		value := *(*uintptr)(unsafe.Pointer(ptr))
		fmt.Printf("Значение по адресу 0x%x: 0x%x\n", ptr, value)

		if value != 0 && !added[value] {
			index++ // Инкрементируем индекс перед вставкой
			// Вставляем новый указатель сразу после текущего
			pointers = append(pointers[:index], append([]uintptr{value}, pointers[index:]...)...)
			added[value] = true
			fmt.Printf("Добавлен новый указатель при обходе: 0x%x на позицию %d\n", value, index)
		} else {
			index++
		}
	}

	return pointers
}

func TestTrace(t *testing.T) {
	var heapObjects = []int{
		0x00, 0x00, 0x00, 0x00, 0x00,
	}

	var heapPointer1 *int = &heapObjects[1]
	var heapPointer2 *int = &heapObjects[2]
	var heapPointer3 *int = nil
	var heapPointer4 **int = &heapPointer3

	var stacks = [][]uintptr{
		{
			uintptr(unsafe.Pointer(&heapPointer1)), 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[0])),
			0x00, 0x00, 0x00, 0x00,
		},
		{
			uintptr(unsafe.Pointer(&heapPointer2)), 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[1])),
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[2])),
			uintptr(unsafe.Pointer(&heapPointer4)), 0x00, 0x00, 0x00,
		},
		{
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[3])),
		},
	}

	pointers := Trace(stacks)
	expectedPointers := []uintptr{
		uintptr(unsafe.Pointer(&heapPointer1)),
		uintptr(unsafe.Pointer(&heapObjects[0])),
		uintptr(unsafe.Pointer(&heapPointer2)),
		uintptr(unsafe.Pointer(&heapObjects[1])),
		uintptr(unsafe.Pointer(&heapObjects[2])),
		uintptr(unsafe.Pointer(&heapPointer4)),
		uintptr(unsafe.Pointer(&heapPointer3)),
		uintptr(unsafe.Pointer(&heapObjects[3])),
	}

	assert.True(t, reflect.DeepEqual(expectedPointers, pointers))
}
