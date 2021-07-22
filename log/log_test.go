package log

import (
	"errors"
	"fmt"
	"testing"

	"github.com/micro/go-micro/v2/util/log"
)

func TestLogger_Log(t *testing.T) {
	Init(KafkaTopic("test_topic"), KafkaBrokers([]string{"127.0.0.1:9092"}))
	Debug("debug message")
	Info("info message")
	Infof("this is an error(%v), type is (%s)", errors.New("test err"), "www")
	Infow("this is a test message", "name", "ericjwang")
	_ = Sync()
}

func Test_Log(t *testing.T) {
	log.Log("test")
}

func TestLogger_noInit(t *testing.T) {
	Info("dit not Init")
}

func TestLogger_turnOffStdout(t *testing.T) {
	Init(StdoutOff(true))
	Info("should not be recorded")
}

func TestSync(t *testing.T) {
	Info("test Sync")
	fmt.Printf("err is (%v) \n", Sync())
}
