package main

import (
	"errors"
	"fmt"
	"testing"
)

type Record struct {
	Name string
	Age  int
}

type recorder interface {
	Record(name string) (Record, error)
}

func Greeter(name string, fetch recorder) (string, error) {
	rec, err := fetch.Record(name)
	if err != nil {
		return "", err
	}
	if rec.Name != name {
		return "", fmt.Errorf("server returned record for %s, not %s", rec.Name, name)
	}
	if rec.Age < 18 {
		return "Greetings young one", nil
	}
	return fmt.Sprintf("Greetings %s", name), nil
}

type fakeRecorder struct {
	rec Record
	err bool
}

func (f fakeRecorder) Record(name string) (Record, error) {
	if f.err {
		return Record{}, errors.New("error")
	}
	return f.rec, nil
}

func TestGreeter(t *testing.T) {
	tests := []struct {
		desc      string
		name      string
		recorder  recorder
		want      string
		expectErr bool
	}{
		{
			desc:      "Error: recorder had some server error",
			name:      "John",
			recorder:  fakeRecorder{err: true},
			expectErr: true,
		},
		{
			desc: "Error: server returned wrong name",
			name: "John",
			recorder: fakeRecorder{
				rec: Record{Name: "Bob", Age: 20},
			},
			expectErr: true,
		},
		{
			desc: "Success",
			name: "John",
			recorder: fakeRecorder{
				rec: Record{Name: "John", Age: 20},
			},
			want: "Greetings John",
		},
		{
		//this object won't pass a test
			desc: "Success",
			name: "Mella",
			recorder: fakeRecorder{
				rec: Record{Name: "Mella", Age: 19},
			},
			want: "Greetings young one",
		},
	}

	// Executes each test.
	for _, test := range tests {
		got, err := Greeter(test.name, test.recorder)
		switch {
		// We did not get an error, but expected one
		case err == nil && test.expectErr:
			t.Errorf("TestGreet(%s): got err == nil, want err != nil", test.desc)
			continue
		// We got an error but did not expect one
		case err != nil && !test.expectErr:
			t.Errorf("TestGreet(%s): got err == %s, want err == nil", test.desc, err)
			continue
		// We got an error we expected, so just go to the next test
		case err != nil:
			continue
		}

		// We did not get the result we expected
		if got != test.want {
			t.Errorf("TestGreet(%s): got result %q, want %q", test.desc, got, test.want)
		}
	}
}
