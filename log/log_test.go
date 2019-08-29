package log

import "testing"

func TestLog(t *testing.T) {
	ls, err := NewLogSystem("testlog.txt")
	if err != nil {
		t.Error(err)
		return
	}
	lg := ls.GetInstance("TESTTOPIC")
	lg2 := lg.GetSubInstance("TEST2")
	lg3 := lg.GetSubInstance("TEST3")
	lg.Log("hhhhh")
	lg.Log("hhhhh2333")
	lg2.Log("MMMMM")
	lg3.Log("0000")
	lg.Err("KKK")
	lg2.Err("KMM")
	lg3.Err("114514")
}
