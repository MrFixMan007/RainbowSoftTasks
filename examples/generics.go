package main

import (
	"fmt"
	"strconv"
)

type myint int

func (i myint) String() string { //Класс myint реализовал приведение к string
	return strconv.Itoa(int(i))
}

// Stringer это ограничение типа которое требует, чтобы аргумент типа
// имел метод String и позволяет дженерик функции вызвать String().
// Метод String должен возвращать строку.
type Stringer interface {
	String() string
}

// Plusser это ограничение типа, которое требует, чтобы аргумент типа
// имел метод Plus. Ожидается что Plus метод добавить получаемую
// строку с внутренней строке и вернет получившийся результат
type Plusser interface {
	Plus(string) string
}

// ConcatTo принимает слайс елементов, которые имеют метод String
// и слайс элементов с методом Plus. Слайсы должны быть одинакового размера.
// Функция конвертирует каждый элемент слайса s в строку и передает его в
// метод Plus соответствующего элемента из слайса p и возвращает строк,
// полученных в результате
func ConcatTo[S Stringer, P Plusser](s []S, p []P) []string {
	r := make([]string, len(s))
	for i, v := range s {
		r[i] = p[i].Plus(v.String())
	}
	return r
}

func main() {
	var s = make([]int, 0)
	s = append(s, 5, 0, 7)
	Print(s)

	x := []myint{myint(1), myint(2), myint(3)}
	Stringify(x)

	// Инициализация мапы с интовыми значениями
	ints := map[string]int64{
		"first":  34,
		"second": 12,
	}

	// Инициализация мапы с float значениями
	floats := map[string]float64{
		"first":  35.98,
		"second": 26.99,
	}

	SumIntsOrFloats(ints)
	SumIntsOrFloats(floats)

	tesla := Car{}
	boing := Aircraft{}
	drive(tesla)
	drive(boing)
	move(tesla)
	move(boing)
}

func F[T any](p T) {}

func Print[T any](s []T) {
	for _, v := range s {
		fmt.Println(v)
	}
}

func Stringify[T Stringer](s []T) (ret []string) {
	for _, v := range s {
		ret = append(ret, v.String()) // Это не сработает
	}
	return ret
}

// SumIntsOrFloats суммирует значения мамы m.
// Поддерживает int64 и float64 типы элементов в мапе.
func SumIntsOrFloats[K comparable, V int64 | float64](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

type Vehicle interface {
	move()
}

func drive(vehicle Vehicle) {
	vehicle.move()
}

type Car struct{}
type Aircraft struct{}

func (c Car) move() {
	fmt.Println("Автомобиль едет")
}
func (a Aircraft) move() {
	fmt.Println("Самолет летит")
}
func move[T Vehicle](v T) {
	v.move()
}
