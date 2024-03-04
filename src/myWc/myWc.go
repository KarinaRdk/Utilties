package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
)

type flags struct {
	l, m, w bool
}

//  NB! according to README.md we may consider spaces as the only word delimiters

func main() {
	if len(os.Args) < 3 {
		fmt.Println("You need to provide a flag(s) and a path")
		return
	}

	var f flags

	flag.BoolVar(&f.l, "l", false, "for counting lines")
	flag.BoolVar(&f.m, "m", false, "for counting symbols")
	flag.BoolVar(&f.w, "w", false, "for counting words")

	flag.Parse()

	if flag.NFlag() > 1 {
		fmt.Println("Error: there could only be  one flag")
		return
	}
	var waitGroup sync.WaitGroup
	for _, file := range os.Args[2:] {
		//  need a waitGroup to ensure that the main returns only after all goroutines has been finished
		waitGroup.Add(1)
		go ProcesFile(file, f, &waitGroup)
	}
	waitGroup.Wait()
}

func ProcesFile(file string, f flags, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	r, err := os.Open(file)

	if err != nil {
		fmt.Println("Couldn't open ", file)
		return
	}
	defer r.Close()

	if f.l {
		lines, _ := Counter(r, '\n')
		fmt.Println(lines, file)
	}
	if f.w || (!f.m && !f.l) {
		words, _ := Counter(r, ' ')
		fmt.Println(words, file)
	}
	if f.m {
		symbols := SymbolsCounter(r)
		fmt.Println(symbols, file)
	}
}

func SymbolsCounter(r io.Reader) int {
	var count int
	scan := bufio.NewScanner(r)
	//  Need i variable to count '/n'
	i := 0
	for ; scan.Scan(); i++ {
		str := scan.Text()
		count = count + len([]rune(str))
	}
	count += i
	return count
}

func Counter(r io.Reader, separator byte) (int, error) {
	var count int
	buf := make([]byte, bufio.MaxScanTokenSize)
	for {
		//Read() populates the given byte slice with data and returns the number of bytes populated and an error value.
		bufferSize, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return 0, err
		}
		var buffPosition int
		for {
			//  find bytes that have a particular symbol - separator
			i := bytes.IndexByte(buf[buffPosition:], separator)
			if i == -1 || bufferSize == buffPosition {
				break
			}
			// to continue the search from the byte next to the one detected previously:
			buffPosition += i + 1
			count++
		}
		if err == io.EOF {
			break
		}
	}
	return count, nil
}
