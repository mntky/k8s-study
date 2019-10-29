package main

import (
	"fmt"
	"errors"
)

func valiadic(a string, b  ...string){
	fmt.Println(a)
	fmt.Println(b)
}

func check(i int) []error{
	if i==1 {
		return nil
	}

	errs := []error{}
	errs = append(errs, errors.New("errodayo1"))
	errs = append(errs, errors.New("errodayo2"))
	errs = append(errs, errors.New("errodayo3"))
	return errs
}

func main() {
	fmt.Println("--func variadic--")
	valiadic("a","b","c","d","e","b","c","d")

	fmt.Printf("\n--func check--\n")
	var test []error
	test = append(test, check(0)...)
	test = append(test, check(1)...)
	test = append(test, check(0)...)
	fmt.Println(test)
	fmt.Println(len(test))
}
