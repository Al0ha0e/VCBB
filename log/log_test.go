package log

import "testing"

func TestLog(t *testing.T) {
	ls, err := NewLogSystem("testlog.txt")
	if err != nil {
		t.Error(err)
		return
	}
	lg := ls.GetInstance("TESTTOPIC")
	err = lg.Log("hhhhh")
	if err != nil {
		t.Error(err)
	}
	err = lg.Log("hhhhh2333")
	if err != nil {
		t.Error(err)
	}
}
