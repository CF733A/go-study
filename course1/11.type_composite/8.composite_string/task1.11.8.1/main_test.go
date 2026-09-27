package main

import "testing"

func TestCountBytes(t *testing.T){
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Empty string", "", 0},
		{"ASCII only", "hello", 5},
		{"Cyrillic", "Привет", 12},       
		{"Chinese chars", "世界", 6},      
		{"Mixed string", "Hello, 世界!", 14}, 
		{"With emoji", "Go👍", 6},         
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := countBytes(tt.input); result != tt.expected {
				t.Errorf("countBytes(%q) got:%v want:%v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCountSymbols(t *testing.T){
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Empty string", "", 0},
		{"ASCII only", "hello", 5},
		{"Cyrillic", "Привет", 6},      
		{"Chinese chars", "世界", 2},    
		{"Mixed string", "Hello, 世界!", 10}, 
		{"With emoji", "Go👍", 3},       
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := countSymbols(tt.input); result != tt.expected {
				t.Errorf("countSymbols(%q) got:%v want:%v", tt.input, result, tt.expected)
			}
		})
	}
}