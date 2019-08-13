package slave_side

import "testing"

func TestPyExecuter(t *testing.T) {
	exe := new(PyExecuter)
	code :=
		`def func():
     for i in input:
         print(i)
func()
output=["10","10","22"]`
	exe.Run(3, []string{"AAA", "BBB", "CCC", "DDD", "EEE", "FFF"}, code, nil)
}
