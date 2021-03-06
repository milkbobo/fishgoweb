package main

import (
	. "github.com/milkbobo/fishgoweb/assert"
	. "github.com/milkbobo/fishgoweb/language"
	. "github.com/milkbobo/fishgoweb/language/querygen/testdata"
	"testing"
)

func TestQuerySelect(t *testing.T) {
	data := []User{
		User{Name: "Man_a"},
		User{Name: "Woman_b"},
		User{Name: "Man_c"},
	}

	AssertEqual(t, QuerySelect(data, func(a User) Sex {
		if len(a.Name) >= 3 && a.Name[0:3] == "Man" {
			return Sex{IsMale: true}
		} else {
			return Sex{IsMale: false}
		}
	}), []Sex{
		Sex{IsMale: true},
		Sex{IsMale: false},
		Sex{IsMale: true},
	})
	AssertEqual(t, QuerySelect([]int{3, 4, 5, -1}, func(a int) User {
		return User{UserId: a}
	}), []User{
		User{UserId: 3},
		User{UserId: 4},
		User{UserId: 5},
		User{UserId: -1},
	})
}

func BenchmarkQuerySelectHand(b *testing.B) {
	data := make([]User, 1000, 1000)

	b.ResetTimer()
	for i := 0; i != b.N; i++ {
		newData := make([]Sex, len(data), len(data))
		for i, single := range data {
			newData[i] = func(a User) Sex {
				if len(a.Name) >= 3 && a.Name[0:3] == "Man" {
					return Sex{IsMale: true}
				} else {
					return Sex{IsMale: false}
				}
			}(single)
		}
	}
}

func BenchmarkQuerySelectMacro(b *testing.B) {
	data := make([]User, 1000, 1000)

	b.ResetTimer()
	for i := 0; i != b.N; i++ {
		QuerySelect(data, func(a User) Sex {
			if len(a.Name) >= 3 && a.Name[0:3] == "Man" {
				return Sex{IsMale: true}
			} else {
				return Sex{IsMale: false}
			}
		})
	}
}

func BenchmarkQuerySelectReflect(b *testing.B) {
	data := make([]User, 1000, 1000)

	b.ResetTimer()
	for i := 0; i != b.N; i++ {
		QuerySelect(data, func(a User) bool {
			if len(a.Name) >= 3 && a.Name[0:3] == "Man" {
				return true
			} else {
				return false
			}
		})
	}
}
