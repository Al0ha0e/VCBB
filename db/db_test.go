package db

import (
	"fmt"
	"testing"
)

type testStruct struct {
	A string
	B int
}

func TestMapDB(t *testing.T) {
	var tdb MapDB
	tt := testStruct{
		A: "shh",
		B: 19,
	}
	tdb.Set("test", tt)
	vv, _ := tdb.Get("test")
	fmt.Println(vv)
	v2 := testStruct{
		A: "k",
		B: 100,
	}
	v, get := tdb.GetOrSet("test2", v2)
	fmt.Println(v, get)
	v, get = tdb.GetOrSet("test2", 1)
	fmt.Println(v, get)
}
