package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

var flagSecond = flag.Uint("ms", 100, "milli seconds")

func cat(ticker <-chan time.Time, r io.Reader) error {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		text := sc.Text()
		fmt.Println(text)
		<-ticker
	}
	return sc.Err()
}

func mains(args []string) error {
	ticker := time.NewTicker(time.Duration(*flagSecond) * time.Millisecond)
	defer ticker.Stop()

	if len(args) <= 0 {
		return cat(ticker.C, os.Stdin)
	}
	for _, arg1 := range args {
		fd, err := os.Open(arg1)
		if err != nil {
			return err
		}
		err1 := cat(ticker.C, fd)
		err2 := fd.Close()
		if err1 != nil {
			return err1
		}
		if err2 != nil {
			return err2
		}
	}
	return nil
}

func main() {
	flag.Parse()
	if err := mains(flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
