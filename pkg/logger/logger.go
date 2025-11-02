package logger

import "log"

type Logger interface {
	Info(args ...any)
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
}

type std struct{ env string }

func New(env string) Logger { return &std{env: env} }

func (s *std) Info(a ...any)             { log.Println(a...) }
func (s *std) Infof(f string, a ...any)  { log.Printf(f, a...) }
func (s *std) Errorf(f string, a ...any) { log.Printf("[ERROR] "+f, a...) }
func (s *std) Fatalf(f string, a ...any) { log.Fatalf(f, a...) }
