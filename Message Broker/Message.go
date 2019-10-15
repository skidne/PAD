package messagebroker

import (
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
)

// Message struct/kinda class
type Message struct {
	topic   string
	content string
}

func (message *Message) read(reader io.Reader) (content []byte, err error) {
	separator := []byte("\n")
	content, err = ioutil.ReadAll(reader)
	if err == nil {
		messageSplit := bytes.SplitN(content, separator, 2)
		*message = Message{
			topic:   string(messageSplit[0]),
			content: string(messageSplit[1]),
		}
	}
	return
}

func (message *Message) write(writer io.Writer) (n int, err error) {
	buff := []byte(message.topic + "\n" + message.content)
	length := uint32(len(buff))
	lengthBuff := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBuff, length)

	n, err = writer.Write(append(lengthBuff, buff...))
	return
}
