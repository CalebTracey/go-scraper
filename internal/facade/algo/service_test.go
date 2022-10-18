package algo

import "testing"

func Test_balanceUserRatings(t *testing.T) {
	type args struct {
		rating float64
		count  int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Happy Path",
			args: args{
				rating: 5.0,
				count:  34,
			},
			want: "4.55",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := balanceUserRatings(tt.args.rating, tt.args.count); got != tt.want {
				t.Errorf("balanceUserRatings() = %v, want %v", got, tt.want)
			}
		})
	}
}
