package algo

import "fmt"

// r: "default" rating = 3
// w: initial belief "worth" = 10
const (
	defaultR = 3
	defaultW = 10
)

// balanceUserRatings algorithm inspired by Bayesian probability to balance
// the number of ratings versus the ratings themselves.
func balanceUserRatings(rating float64, count int) string {
	top := (defaultW * defaultR) + (float64(count) * rating)
	bot := defaultW + count
	res := top / float64(bot)
	return fmt.Sprintf("%.2f", res)
}
