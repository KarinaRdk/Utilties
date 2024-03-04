package main

//   ls -l  This command will display a detailed list of files and directories in the current directory, including any symlinks.
//   ln -s /path/to/original/file /path/to/symlink  command to create a symlink

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("\tYou need to provide at least one flag (-sl, -d, or -f) and a path")
		return
	}

	var sl, d, f bool
	flag.BoolVar(&sl, "sl", false, "")
	flag.BoolVar(&d, "d", false, "")
	flag.BoolVar(&f, "f", false, "")
	ext := flag.String("ext", "", "Only files with a certain extensio\n")
	flag.Parse()

	if flag.Arg(0) == "" {
		fmt.Println("\tYou need to provide a path")
		return
	}

	if err := filepath.Walk(flag.Arg(0), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if d {
			if info.IsDir() {
				fmt.Println(path)
			}
		}

		if f {
			if info.IsDir() {
				return nil
			}
			if *ext != "" {
				if *ext == strings.Trim(filepath.Ext(path), ".") {
					fmt.Println(path)
				}
			} else {
				fmt.Println(path)
			}
		}

		if sl {
			if info.Mode()&os.ModeSymlink != 0 {
				link, error := filepath.EvalSymlinks(path)
				if error != nil {
					fmt.Println(path, "->", "[broken]")
				} else {
					fmt.Println(path, "->", link)
				}
			}
		}

		return err

	}); err != nil {
		fmt.Println(err)
	}
}
