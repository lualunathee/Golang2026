package main

import "testing"

func TestAddTableDriven(t *testing.T) {
	tests := []struct {
		name   string
		a, b   int
		want   int
	}{
		{"both positive", 2, 3, 5},
		{"positive + zero", 5, 0, 5},
		{"negative + positive", -1, 4, 3},
		{"both negative", -2, -3, -5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Add(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestSubtractTableDriven(t *testing.T) {
	tests := []struct {
		name   string
		a, b   int
		want   int
	}{
		{"both positive", 5, 2, 3},
		{"positive minus zero", 5, 0, 5},
		{"negative minus positive", -2, 3, -5},
		{"both negative", -5, -2, -3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Subtract(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Subtract(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		got, err := Divide(10, 2)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if got != 5 {
			t.Errorf("Divide(10, 2) = %d; want 5", got)
		}
	})

	t.Run("division by zero", func(t *testing.T) {
		_, err := Divide(10, 0)
		if err == nil {
			t.Error("Expected error for division by zero, got nil")
		}
	})
}