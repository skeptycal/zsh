package zsh

import "testing"

func TestRemoveDuplicates(t *testing.T) {
	type args struct {
		s string
		c byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
        // TODO: Add test cases.
        {"test", args{"test",32},"test"},
        {"test   data", args{"test   data",32},"test data"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveDuplicates(tt.args.s, tt.args.c); got != tt.want {
				t.Errorf("RemoveDuplicates() = %v, want %v", got, tt.want)
			}
		})
	}
}
