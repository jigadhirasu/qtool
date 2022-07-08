package env

import (
	"math/rand"
	"os"
	"time"
)

func Develop(ff ...func()) {

	rand.Seed(time.Now().UnixNano())

	os.Setenv("MARIADB_HOST", "35.201.200.61:3306")
	os.Setenv("MARIADB_USER", "develop")
	os.Setenv("MARIADB_PASS", "develop@asg")

	os.Setenv("PREFIX_DBNAME", "c")
	os.Setenv("REDIS_ADDR", "35.201.200.61:6379")

	os.Setenv("Namespace", "develop")
	os.Setenv("HOSTNAME", "localhost")
	os.Setenv("KUBEMQ_ADDR", "35.201.200.61")
	os.Setenv("KUBEMQ_PORT", "50000")

	os.Setenv("PORT", "8080")
	for _, f := range ff {
		f()
	}
}
