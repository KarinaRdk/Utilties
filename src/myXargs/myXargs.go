//  to test properly add to the folder bynary files myFind and myWc

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec" //  he os/exec package runs external commands
)

func main() {
	fmt.Println("os.Args", len(os.Args))

	if len(os.Args) < 2 {
		fmt.Println("USAGE: command that generates a list of argumens | ./myXargs command")
		return
	}

	var list []string
	output := bufio.NewScanner(os.Stdin)
	for output.Scan() {
		i := output.Text()
		list = append(list, i)
	}

	//  The Command returns the Cmd struct to execute the specified program with the given arguments.
	//  The first parameter is the program to be run; the other arguments are parameters to the program.
	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	//  Add to the Args element of the structure output from the first command given by the user
	cmd.Args = append(cmd.Args, list...)
	//  Stderr sets streaming STDERR if enabled, else nil
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	//  Stdout sets streaming STDOUT if enabled, else nil
	cmd.Stdout = os.Stdout
	//  The Run function starts the specified command and waits for it to complete
	err := cmd.Run()
	if err != nil {
		fmt.Println(stderr.String())
	}

}
