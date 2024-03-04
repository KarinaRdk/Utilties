# Utilties

In this project I replicated several command line utilities:

# MyFind
Accepts a path and a set of command-line options to locate different types of entries, namely directories, regular files and symbolic links.
You can specify wich type of objects you are interested in by providing flags -d, -f, -sl or -ext. The latter works only when -f is specified and allows you search files with particular extention only.

Example of finding all files/directories/symlinks recursively in directory /foo
    
    ~$ ./myFind /foo
    /foo/bar
    /foo/bar/baz
    /foo/bar/baz/deep/directory
    /foo/bar/test.txt
    /foo/bar/buzz -> /foo/bar/baz
    /foo/bar/broken_sl -> [broken]

# MyWc
Gathers basic statistics about files. Works with three mutually exclusive flags f -l for counting lines, -m for counting characters and -w for counting words. 

Example of Finding only *.go files ignoring all the rest

    ~$ ./myFind -f -ext 'go' /go
    /go/src/github.com/mycoolproject/main.go
    /go/src/github.com/mycoolproject/magic.go

  
# MyXargs
Can be used to build and execute commands from standard input. It converts input from standard input into arguments to a command.
For instance  

    ~$ echo -e "/a\n/b\n/c" | ./myXargs ls -la

would render the same result as running  

    ~$ ls -la /a /b /c
