package gen

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"time"
)

var alphabet = "0123456789_ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RndString(maxSize int) (string, error) {
	if maxSize < 0 {
		return "", fmt.Errorf("cannot generate string with negative len")
	}
	size := rand.Intn(maxSize)
	bs := make([]byte, size)
	for i := 0; i < size; i++ {
		bs[i] = alphabet[rand.Int63()%int64(len(alphabet))]
	}
	return string(bs), nil
}

func RndStrings(cnt int, maxLen int, out io.Writer) error {
	bufOut := bufio.NewWriterSize(out, maxLen)
	for i := 0; i < cnt; i++ {
		s, err := RndString(maxLen)
		if err != nil {
			return err
		}

		if _, err = fmt.Fprintln(bufOut, s); err != nil {
			return err
		}
	}
	return bufOut.Flush()
}

func RepeatedReversedAlphabeticStrings(cnt int, out io.Writer) error {
	alpha := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bufOut := bufio.NewWriter(out)
	var s string
	for i := 0; i < cnt; i++ {
		s = alpha[len(alpha)-i%len(alpha)-1:len(alpha)-i%len(alpha)] + s
		if _, err := fmt.Fprintln(bufOut, s); err != nil {
			return err
		}
	}
	return bufOut.Flush()
}
