package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec" //  the os/exec package runs external commands
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	a := flag.String("a", "", "in case we need to put all the arcs in a directory")
	flag.Parse()

	if flag.Arg(0) == "" {
		fmt.Printf("\tYou need to provide a directory and a file\n")
		return
	}

	index := 1
	inCaseOfA := ""
	if *a != "" {
		//  Create a folder if there isn't one:
		if err := os.MkdirAll(os.Args[2], os.ModePerm); err != nil {
			fmt.Println(err)
			return
		}
		index = 3
		inCaseOfA = *a
	}

	times, err := getModifTime(index)
	if err != nil {
		if *a != "" {
			os.Remove(os.Args[2])
		}
		return
	}

	var waitGroup sync.WaitGroup

	for i, file := range os.Args[index:] {
		waitGroup.Add(1)
		err = isFile(file, &waitGroup)
		if err != nil {
			continue //  In case there are valid argumnts given after an invalid one
		}
		err = arc(file, times[i], inCaseOfA, &waitGroup)
		if err != nil {
			return
		}
	}
	waitGroup.Wait()
}

func getModifTime(index int) ([]string, error) {
	//  Get the modification time with the command "stat -f %m":
	cmd := exec.Command("stat", "-f")
	cmd.Args = append(cmd.Args, "%m")
	cmd.Args = append(cmd.Args, os.Args[index:]...)

	//  Stderr sets streaming STDERR if enabled, else nil
	var stderr bytes.Buffer
	var time bytes.Buffer
	cmd.Stderr = &stderr
	//  Stdout sets streaming STDOUT if enabled, else nil
	cmd.Stdout = &time
	//  The Run function starts the specified command and waits for it to complete
	err := cmd.Run()
	if err != nil {
		fmt.Println("\tERROR: ", stderr.String())
		return nil, err
	}
	readBuf, _ := ioutil.ReadAll(&time)
	//  Get a slice with modif time for all the files:
	return strings.Split(string(readBuf), "\n"), err
}

// In case user put a folder instead of a file
func isFile(file string, waitGroup *sync.WaitGroup) error {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		waitGroup.Done()

		return err
	}
	// defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		fmt.Println(err)
		waitGroup.Done()
		return err
	}
	if stat.IsDir() {
		err = fmt.Errorf("\tERROR: %s is a directory, not a file\n \tDid you meen to use flag -a? Do it explicitly", file)
		fmt.Println(err)
		waitGroup.Done()
	}
	return err
}

func arc(file string, time string, inCaseOfA string, waitGroup *sync.WaitGroup) error {
	defer waitGroup.Done()
	arcName := name(file, inCaseOfA, time)
	err := createArc(arcName, file)
	if err != nil {
		fmt.Println(err)
		return err
	}
	os.Remove(file)
	return err
}

func createArc(name string, filePath string) error {
	tarfile, err := os.Create(name)
	if err != nil {
		fmt.Println(err)
	}
	gzipWriter := gzip.NewWriter(tarfile)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	err = addFileToTarWriter(filePath, tarWriter)
	if err != nil {
		fmt.Println(err)
		err := os.Remove(name)
		if err != nil {
			fmt.Println(err)
		}
	}
	return err
}

func name(path string, inCaseOfA string, time string) (name string) {
	nameSplit := strings.Split(path, ".")
	var a rune = '/'
	if inCaseOfA != "" && []rune(inCaseOfA)[len(inCaseOfA)-1] != a {

		inCaseOfA += "/"
	}
	name = fmt.Sprint(inCaseOfA, nameSplit[0], "_", time, ".tag.gz")
	return name
}

func addFileToTarWriter(filePath string, tarWriter *tar.Writer) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	header := &tar.Header{
		Name:    filepath.Base(filePath),
		Size:    stat.Size(),
		Mode:    int64(stat.Mode()),
		ModTime: stat.ModTime(),
	}
	err = tarWriter.WriteHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(tarWriter, file)
	if err != nil {
		return err
	}
	return nil
}
