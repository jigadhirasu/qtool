package mariadb

import (
	"net/url"
	"regexp"
	"testing"

	"github.com/jigadhirasu/qtool/log4"
	"github.com/kubemq-io/kubemq-go/pkg/uuid"
)

func TestRegexp(t *testing.T) {
	log4.Debug(regexp.MatchString("^LT\\([\\w-]+\\)", "LT("+uuid.New()+")"))

	log4.Debug(url.PathEscape("LT_123"))
	log4.Debug(url.QueryEscape("LT_123"))
}
