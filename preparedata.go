package main

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"
)

func readData() {
	fmt.Println("I'm alive")
	fmt.Println(time.Now())
	printMemUsage()

	{
		file, err := os.Open("/tmp/data/options.txt")
		if err != nil {
			fmt.Println("File /tmp/data/options.txt not found")
			fmt.Println(err)
		} else {
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				scanner.Text()

				i, e := strconv.Atoi(scanner.Text())
				if e != nil {
					fmt.Println("time_not_parse")
					fmt.Println(e)
				} else {
					nowtime = int32(i)
				}

				break
			}
		}

		//nowtime
	}
	err := unzip("/home/scherbina/Documents/highloadcup.ru/go_v0/data/hard/data/data.zip", "/tmp/mydata/")
	//err := unzip("/tmp/data/data.zip", "/tmp/mydata/")
	if err != nil {
		log.Fatal(err)
	}

	loadAccountsFromFiles("/tmp/mydata")
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func loadAccountsFromFiles(fld string) {

	//runtime.GOMAXPROCS(1)

	files, err := ioutil.ReadDir(fld)
	if err != nil {
		log.Fatal(err)
	}

	printMemUsage()

	for i, f := range files {
		b, err := ioutil.ReadFile(fld + "/" + f.Name())

		af := AccountFile{}
		err = json.Unmarshal(b, &af)

		if err != nil {
			log.Fatal(err)
		}

		for _, aa := range af.Accounts {
			a := AccountCreate(aa)
			accounts.Append(&a)
		}
		//printMemUsage()
		//fmt.Println(f.Name())
		debug.FreeOSMemory()
		fmt.Print(i, " ")
	}

	{
		fmt.Println("Do Indexes")
		go func() {
			fmt.Println("Do like Indexes & group")
			for i, a := range accounts.data {
				if a != nil {
					a.FillOtherLikeBack()

					groupindex.Append(a)
				}
				if i%30000 == 0 {
					runtime.GC()
					fmt.Print("l-", i/10000, " ")
				}
			}
			fmt.Println("Done like Indexes & group")
			printMemUsage()
			debug.FreeOSMemory()
		}()

		for i, a := range accounts.data {
			if i%30000 == 0 {
				//printMemUsage()
				debug.FreeOSMemory()
				fmt.Print("i-", i/10000, " ")
				//fmt.Println("Doing Indexes", kv)
			}
			if a != nil {
				index.Append(a)
			}

		}

		fmt.Println("Done Indexes")
	}
	runtime.GC()

	//fmt.Println("i_li_map", i_li_map, " i_li", i_li)

	printMemUsage()
	go func() {
		var i int64
		for {
			i++
			time.Sleep(time.Millisecond * 1800)
			if i%10 == 0 {
				printMemUsage()
				debug.FreeOSMemory()
			} else {
				//runtime.GC()
				debug.FreeOSMemory()
			}
			//printMemUsage()

		}
	}()

	//runtime.GOMAXPROCS(runtime.NumCPU())
}

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tHeapAlloc = %v MiB", bToMb(m.HeapAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v", m.NumGC)
	fmt.Printf("\tStackSys = %v MiB \t %v\t %v\n", bToMb(m.StackSys), time.Now(), time.Now().Unix())
}
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
