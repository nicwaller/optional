package optional

import (
	"testing"
)

func TestWrapper_Nil(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		opt := OptionalPointer[any](nil)
		if opt.Nil() != true {
			t.Error("expected Nil() == true")
		}
	})
	t.Run("notNil", func(t *testing.T) {
		expected := "Hello, World"
		opt := OptionalPointer(&expected)
		if opt.Nil() {
			t.Error("expected Nil() == false")
		}
	})
}

func TestWrapper_Unwrap(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		opt := OptionalPointer[any](nil)
		opt.Unwrap(func(a *any) {
			t.Error("should not unwrap")
		})
	})
	t.Run("notNil", func(t *testing.T) {
		expected := "Hello, World"
		opt := OptionalPointer(&expected)
		didUnwrap := false
		opt.Unwrap(func(unwrapped *string) {
			if *unwrapped != "Hello, World" {
				t.Error("expected: Hello, World")
			}
			didUnwrap = true
		})
		if !didUnwrap {
			t.Error("it should have unwrapped")
		}
	})
}

func TestWrapper_IfElse(t *testing.T) {
	t.Run("nilFunc", func(t *testing.T) {
		opt := OptionalPointer[any](nil)
		opt.IfElse(nil, nil)
		// if we didn't panic, all good.
	})
	t.Run("nil", func(t *testing.T) {
		opt := OptionalPointer[any](nil)
		var elseCalled bool
		opt.IfElse(func(a *any) {
			t.Error("should not unwrap")
		}, func() {
			elseCalled = true
		})
		if !elseCalled {
			t.Errorf("else func should have been called")
		}
	})
	t.Run("notNil", func(t *testing.T) {
		expected := "Hello, World"
		opt := OptionalValue(expected)
		var actual string
		opt.IfElse(func(safePtr *string) {
			actual = *safePtr
		}, func() {
			t.Error("else should not be called: it had a real value")
		})
		if actual != expected {
			t.Errorf("expected %q but got %q", expected, actual)
		}
	})
}

func TestWrapper_Or(t *testing.T) {
	testVals := []any{
		false, true,
		0, 1,
		"", "hello",
		0.0, 0.1,
	}
	t.Run("nil", func(t *testing.T) {
		opt := OptionalPointer[any](nil)
		for _, v := range testVals {
			result := opt.Or(v)
			if *result != v {
				t.Error()
			}
		}
	})
	t.Run("notNil", func(t *testing.T) {
		for _, v := range testVals {
			opt := OptionalPointer(&v)
			result := opt.Or("wrong")
			if *result != v {
				t.Error("expected equality")
			}
		}
	})
}

func TestWrapper_Equal(t *testing.T) {
	testVals := []any{
		false, true,
		0, 1,
		"", "hello",
		0.0, 0.1,
	}
	t.Run("nil", func(t *testing.T) {
		opt := OptionalPointer[any](nil)
		if !opt.Equal(nil) {
			t.Error("expected nil equality")
		}
		for _, v := range testVals {
			if opt.Equal(v) {
				t.Error("expected inequality")
			}
		}
	})
	t.Run("notNil", func(t *testing.T) {
		for _, v := range testVals {
			opt := OptionalPointer(&v)
			if !opt.Equal(v) {
				t.Error("expected equality")
			}
		}
	})
}

func TestOptional_SetValue(t *testing.T) {
	t.Run("fromNil", func(t *testing.T) {
		var opt Optional[int]
		opt.SetValue(7)
		if *opt.Or(0) != 7 {
			t.Error("expected 7")
		}
	})
	t.Run("fromNotNil", func(t *testing.T) {
		original := 5
		opt := OptionalPointer(&original)
		opt.SetValue(7)
		if *opt.Or(0) != 7 {
			t.Error("expected 7")
		}
	})
}

func TestInSlice(t *testing.T) {
	vals := map[int]int{1: 1}
	foo := vals[0]
	_ = foo
	type Scenario struct {
		input []int
		index int
		nil   bool
		value int
	}
	scenarios := []Scenario{
		{
			input: nil, index: 0,
			nil: true, value: 0,
		},
		{
			input: []int{}, index: -1,
			nil: true, value: 0,
		},
		{
			input: []int{}, index: 0,
			nil: true, value: 0,
		},
		{
			input: []int{}, index: 1,
			nil: true, value: 0,
		},
		{
			input: []int{7}, index: -1,
			nil: true, value: 0,
		},
		{
			input: []int{7}, index: 0,
			nil: false, value: 7,
		},
		{
			input: []int{7}, index: 1,
			nil: true, value: 0,
		},
	}
	for _, expected := range scenarios {
		t.Log(expected)
		actual := OptionalFromSlice(expected.input, expected.index)
		if expected.nil != actual.Nil() {
			t.Fatalf("expected nil=%t", expected.nil)
		}
		if actual.Nil() {
			continue
		} else {
			if !actual.Equal(expected.value) {
				t.Errorf("failed equality test")
			}
		}
		actual.Unwrap(func(safePtr *int) {
			if *safePtr != expected.value {
				t.Errorf("expected %d but got %d", expected.value, *safePtr)
			}
		})
	}
}

func TestInMap(t *testing.T) {
	type Scenario struct {
		input map[string]string
		key   string
		nil   bool
		value string
	}
	scenarios := []Scenario{
		{
			input: nil,
			key:   "", value: "", nil: true,
		},
		{
			input: map[string]string{"test": "test"},
			key:   "", value: "", nil: true,
		},
		{
			input: map[string]string{"test": "test"},
			key:   "test", value: "test", nil: false,
		},
	}
	for _, expected := range scenarios {
		t.Log(expected)
		actual := OptionalFromMap(expected.input, expected.key)
		if expected.nil != actual.Nil() {
			t.Fatalf("expected nil=%t", expected.nil)
		}
		if actual.Nil() {
			continue
		} else {
			if !actual.Equal(expected.value) {
				t.Errorf("failed equality test")
			}
		}
		actual.Unwrap(func(safePtr *string) {
			if *safePtr != expected.value {
				t.Errorf("expected %s but got %s", expected.value, *safePtr)
			}
		})
	}
}

func TestDelve(t *testing.T) {
	type Scenario struct {
		collection any
		indices    []int
		value      int
		nil        bool
	}

	mtx1D := []int{0, 1}

	mtx2D := [][]int{
		{0, 1},
		{2, 3},
	}

	mtx3D := [][][]int{
		{
			{0, 1},
			{2, 3},
		},
		{
			{4, 5},
			{6, 7},
		},
	}

	scenarios := []Scenario{
		{collection: nil, indices: nil, value: 0, nil: true},
		{collection: nil, indices: []int{}, value: 0, nil: true},
		{collection: nil, indices: []int{0}, value: 0, nil: true},
		{collection: nil, indices: []int{0, 0}, value: 0, nil: true},

		{collection: mtx1D, indices: nil, value: 0, nil: true},
		{collection: mtx1D, indices: []int{}, value: 0, nil: true},
		{collection: mtx1D, indices: []int{-1}, value: 0, nil: true},
		{collection: mtx1D, indices: []int{0}, value: 0, nil: false},
		{collection: mtx1D, indices: []int{1}, value: 1, nil: false},
		{collection: mtx1D, indices: []int{2}, value: 0, nil: true},

		{collection: mtx2D, indices: nil, value: 0, nil: true},
		{collection: mtx2D, indices: []int{}, value: 0, nil: true},
		{collection: mtx2D, indices: []int{0}, value: 0, nil: true},
		{collection: mtx2D, indices: []int{0, 0}, value: 0, nil: false},
		{collection: mtx2D, indices: []int{0, 1}, value: 1, nil: false},
		{collection: mtx2D, indices: []int{1, 0}, value: 2, nil: false},
		{collection: mtx2D, indices: []int{1, 1}, value: 3, nil: false},

		{collection: mtx3D, indices: []int{0, 0, 0}, value: 0, nil: false},
		{collection: mtx3D, indices: []int{0, 0, 1}, value: 1, nil: false},
		{collection: mtx3D, indices: []int{0, 1, 0}, value: 2, nil: false},
		{collection: mtx3D, indices: []int{0, 1, 1}, value: 3, nil: false},
		{collection: mtx3D, indices: []int{1, 0, 0}, value: 4, nil: false},
		{collection: mtx3D, indices: []int{1, 0, 1}, value: 5, nil: false},
		{collection: mtx3D, indices: []int{1, 1, 0}, value: 6, nil: false},
		{collection: mtx3D, indices: []int{1, 1, 1}, value: 7, nil: false},
		{collection: mtx3D, indices: nil, value: 0, nil: true},
		{collection: mtx3D, indices: []int{}, value: 0, nil: true},
		{collection: mtx3D, indices: []int{0}, value: 0, nil: true},
		{collection: mtx3D, indices: []int{0, 0}, value: 0, nil: true},
		{collection: mtx3D, indices: []int{0, 0, 0, 0}, value: 0, nil: true},
		{collection: mtx3D, indices: []int{-1, -1, -1}, value: 0, nil: true},
		{collection: mtx3D, indices: []int{5, 5, 5}, value: 0, nil: true},
	}
	for _, scenario := range scenarios {
		actual := Delve[int](scenario.collection, scenario.indices...)
		if actual.Nil() != scenario.nil {
			t.Log(scenario)
			t.Errorf("failed nil test")
		}
		if actual.Nil() {
			continue
		}
		if !actual.Equal(scenario.value) {
			t.Log(scenario)
			t.Errorf("expected %d but got %v", scenario.value, actual)
		}
	}
}
