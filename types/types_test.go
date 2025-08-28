package types

import (
	"testing"
	"time"
)

func TestMarshalUnmarshal(t *testing.T)  {
	t1 := CustomTime(time.Now())
	t.Log(time.Time(t1))
	t1marshaled, err := t1.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	t2 := CustomTime{}
	err = t2.UnmarshalJSON(t1marshaled)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(time.Time(t2))
}
