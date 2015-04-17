package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"time"
)

func cp(dst, src string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	// no need to check errors on read only file, we already got everything
	// we need from the filesystem, so nothing can go wrong now.
	defer s.Close()
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}

func main() {
	dataDir := os.Args[1]
	backupDir := os.Args[2]

	for {
		info, err := ioutil.ReadDir(dataDir)
		if err != nil {
			fmt.Printf("readdir: %v\n", err)
		}

		for _, i := range info {
			err = cp(path.Join(backupDir, i.Name()), path.Join(dataDir, i.Name()))
			if err != nil {
				fmt.Printf("cp %s: %v\n", i.Name(), err)
				continue
			}
			fmt.Printf("cp %s: success\n", i.Name())
		}

		time.Sleep(time.Second)
	}
}
