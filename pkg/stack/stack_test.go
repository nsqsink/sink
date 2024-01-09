package stack

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Stack
	}{
		{
			name: "Test New",
			want: &Stack{nil, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStack_Push_Peek_Size(t *testing.T) {
	s := New()

	type args struct {
		value interface{}
	}

	type want struct {
		size     int
		topValue interface{}
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Test 0 items",
			args: args{
				value: "",
			},
			want: want{
				size:     0,
				topValue: nil,
			},
		},
		{
			name: "Test 1 items",
			args: args{
				value: "items",
			},
			want: want{
				size:     1,
				topValue: "items",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.Push(tt.args.value)

			if got := s.Size(); !reflect.DeepEqual(got, tt.want.size) {
				t.Errorf("Size() = %v, want %v", got, tt.want.size)
			}

			if got := s.Peek(); !reflect.DeepEqual(got, tt.want.topValue) {
				t.Errorf("Peek() = %v, want %v", got, tt.want.size)
			}
		})
	}
}

func TestStack_Pop(t *testing.T) {
	s1 := New()
	s1.Push("abc")

	s2 := New()
	s2.Push("abc")
	s2.Push("def")

	type want struct {
		topValue         interface{}
		topValueAfterPop interface{}
		sizeAfterPop     int
	}

	tests := []struct {
		name string
		s    *Stack
		want want
	}{
		{
			name: "Test emtpy",
			s:    New(),
			want: want{
				sizeAfterPop: 0,
			},
		},
		{
			name: "Test 1 data",
			s:    s1,
			want: want{
				topValue:     "abc",
				sizeAfterPop: 0,
			},
		},
		{
			name: "Test 2 data",
			s:    s2,
			want: want{
				topValue:         "def",
				sizeAfterPop:     1,
				topValueAfterPop: "abc",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotValue := tt.s.Pop(); !reflect.DeepEqual(gotValue, tt.want.topValue) {
				t.Errorf("Stack.Pop() = %v, want %v", gotValue, tt.want.topValue)
			}

			if gotValue := tt.s.Size(); !reflect.DeepEqual(gotValue, tt.want.sizeAfterPop) {
				t.Errorf("Stack.Size() = %v, want %v", gotValue, tt.want.sizeAfterPop)
			}

			if gotValue := tt.s.Peek(); !reflect.DeepEqual(gotValue, tt.want.topValueAfterPop) {
				t.Errorf("Stack.Peek() = %v, want %v", gotValue, tt.want.topValueAfterPop)
			}
		})
	}
}
