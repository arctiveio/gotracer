package gotracer

import (
	"log"
	"runtime"
	"time"
)

const layout = "Jan 2, 2006 at 3:04pm (MST)"

type Tracer struct {
	EmailHost     string
	EmailPort     string
	EmailUsername string
	EmailPassword string
	EmailSender   string
	EmailFrom     string
	ErrorTo       string
}

func (self Tracer) Notify(extra ...func() string) {
	err := recover()
	if err != nil {
		const size = 4096
		var exc_message string

		buf := make([]byte, size)
		buf = buf[:runtime.Stack(buf, false)]
		buffer := string(buf)

		switch err.(type) {
		case string:
			_err, ok := err.(string)
			if ok == true {
				exc_message = _err
			}
		case interface{}:
			_err, ok := err.(error)
			if ok == true {
				exc_message = _err.Error()
			}
		}

		extras := ""

		for i := range extra {
			extras += extra[i]()
			extras += " "
		}

		self.sendException(&ErrorStack{
			Subject:   exc_message,
			Traceback: buffer,
			Extra:     extras,
			Timestamp: time.Now().Format(layout),
		})
	}
}

func (self Tracer) sendException(stack *ErrorStack) {
	log.Println(stack.Traceback)

	connection := MakeConn(&self)
	connection.SenderName += " Exception"

	connection.SendEmail(Message{
		self.EmailFrom,
		[]string{self.ErrorTo},
		stack.Subject,
		ErrorTemplate(stack),
	})
}

type ErrorStack struct {
	Subject   string
	Extra     string
	Traceback string
	Timestamp string
}
