package optional

import (
	"testing"
)

func TestWrapper_Nil(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		opt := NewOptional[any](nil)
		if opt.Nil() != true {
			t.Error("expected Nil() == true")
		}
	})
	t.Run("notNil", func(t *testing.T) {
		expected := "Hello, World"
		opt := NewOptional(&expected)
		if opt.Nil() {
			t.Error("expected Nil() == false")
		}
	})
}

func TestWrapper_Unwrap(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		opt := NewOptional[any](nil)
		opt.Unwrap(func(a *any) {
			t.Error("should not unwrap")
		})
	})
	t.Run("notNil", func(t *testing.T) {
		expected := "Hello, World"
		opt := NewOptional(&expected)
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

func TestWrapper_Or(t *testing.T) {
	testVals := []any{
		false, true,
		0, 1,
		"", "hello",
		0.0, 0.1,
	}
	t.Run("nil", func(t *testing.T) {
		opt := NewOptional[any](nil)
		for _, v := range testVals {
			result := opt.Or(v)
			if *result != v {
				t.Error()
			}
		}
	})
	t.Run("notNil", func(t *testing.T) {
		for _, v := range testVals {
			opt := NewOptional(&v)
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
		opt := NewOptional[any](nil)
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
			opt := NewOptional(&v)
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
		opt := NewOptional(&original)
		opt.SetValue(7)
		if *opt.Or(0) != 7 {
			t.Error("expected 7")
		}
	})
}
