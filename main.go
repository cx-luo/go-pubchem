package main

import (
	_ "github.com/go-sql-driver/mysql"
	"go-pubchem/pkg"
	"go-pubchem/router"
	"math/rand"
	"time"
)

// var parsedNames []string
// var unparsedNames []string
func randomInt(min, max int64) int64 {
	rand.Seed(time.Now().UnixNano())

	if min >= max {
		return min
	}

	rangeFloat := float64(max-min) + 1
	randomFloat := rand.Float64() * rangeFloat

	return int64(randomFloat) + min
}

func main() {
	r := router.NewRouter("./go-pubchem.log", "INFO")
	pkg.Logger.Info(r.Run("0.0.0.0:8100"))
}
