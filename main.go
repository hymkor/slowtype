package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

var (
	flagSecond    = flag.Uint("ms", 100, "milli seconds")
	flagBytes     = flag.Uint("b", 0, "bytes")
	flagKiloBytes = flag.Uint("kb", 0, "kilo bytes")
	flagMegaBytes = flag.Uint("mb", 0, "mega bytes")
	flagHang      = flag.Bool("hang", false, "Sleep after output all")
)

func Cat(ticker <-chan time.Time, r io.Reader) error {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		text := sc.Text()
		fmt.Println(text)
		<-ticker
	}
	return sc.Err()
}

type binCat int64

func (size binCat) Cat(ticker <-chan time.Time, r io.Reader) error {
	for {
		_, err := io.CopyN(os.Stdout, r, int64(size))
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		<-ticker
	}
}

func mains(args []string) error {
	ticker := time.NewTicker(time.Duration(*flagSecond) * time.Millisecond)
	defer ticker.Stop()

	cat := Cat
	if *flagMegaBytes > 0 {
		cat = binCat(int64(*flagMegaBytes) * 1024 * 1024).Cat
	} else if *flagKiloBytes > 0 {
		cat = binCat(int64(*flagKiloBytes) * 1024).Cat
	} else if *flagBytes > 0 {
		cat = binCat(int64(*flagBytes)).Cat
	}
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
	if *flagHang {
		time.Sleep(time.Hour)
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
