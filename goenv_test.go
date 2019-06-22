package goenv

import (
	"reflect"
	"testing"
)

func TestPopulate(t *testing.T) {
	type Worker struct {
		Name *string `goenv:"myname"`
	}
	type App struct {
		Host        **string
		Timeout     *int `goenv:"-"`
		Workers     []*Worker
		Coordinates [][]*int
	}
	backup := Environ
	Environ = func() []string {
		return []string{
			"HOST=0.0.0.0:3000",
			"WORKERS_0_MYNAME=worker1",
			"WORKERS_2_MYNAME=worker2",
			"COORDINATES_0_0=1",
			"COORDINATES_0_1=2",
			"COORDINATES_2_0=3",
			"COORDINATES_2_1=4",
		}
	}
	defer func() {
		Environ = backup
	}()
	app := App{}
	err := Populate(&app)
	if err != nil {
		t.Errorf("err: %v", err)
	}
	host := "0.0.0.0:3000"
	ptrToHost := &host
	worker1 := "worker1"
	worker2 := "worker2"
	one := 1
	two := 2
	three := 3
	four := 4
	expected := App{
		Host: &ptrToHost,
		Workers: []*Worker{
			&Worker{
				Name: &worker1,
			},
			nil,
			&Worker{
				Name: &worker2,
			},
		},
		Coordinates: [][]*int{
			[]*int{&one, &two},
			nil,
			[]*int{&three, &four},
		},
	}
	if !reflect.DeepEqual(app, expected) {
		t.Errorf("%+v != %+v", app, expected)
	}
}
