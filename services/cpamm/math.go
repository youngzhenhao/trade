package cpamm

import (
	"errors"
	"strconv"
)

// THIS FILE IS ACTUAL NOT BEEN USED

func add(x uint, y uint) (z uint, err error) {
	z = x + y
	if !(z >= x) {
		err = errors.New("math-add-overflow(" + strconv.FormatUint(uint64(x), 10) + " add " + strconv.FormatUint(uint64(y), 10) + ")")
		return 0, err
	}
	return z, nil
}

func sub(x uint, y uint) (z uint, err error) {
	z = x - y
	if !(z <= x) {
		err = errors.New("math-sub-underflow(" + strconv.FormatUint(uint64(x), 10) + " sub " + strconv.FormatUint(uint64(y), 10) + ")")
		return 0, err
	}
	return z, nil
}

func mul(x uint, y uint) (z uint, err error) {
	z = x * y
	if !(y == 0 || z/y == x) {
		err = errors.New("math-mul-overflow(" + strconv.FormatUint(uint64(x), 10) + " mul " + strconv.FormatUint(uint64(y), 10) + ")")
		return 0, err
	}
	return z, nil
}

// Function 'min' collides with the 'builtin' function
func _min(x uint, y uint) (z uint) {
	if x < y {
		z = x
		return z
	}
	z = y
	return z
}

// https://en.wikipedia.org/wiki/Methods_of_computing_square_roots#Heron's_method
func sqrt(y uint) (z uint) {
	if y > 3 {
		z = y
		var x = y/2 + 1
		for x < z {
			z = x
			x = (y/x + x) / 2
		}
	} else if y != 0 {
		z = 1
	}
	return z
}
