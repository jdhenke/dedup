package main

import (
	"crypto/sha512"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	force bool
)

func init() {
	flag.BoolVar(&force, "force", false, "force removal of files")
	flag.Parse()
}

func main() {
	if !force {
		log.Fatal("TODO")
	}
	bySize := make(map[int64][]string)
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if !info.Mode().IsRegular() {
			return nil
		}
		size := info.Size()
		var sameSize []string
		var ok bool
		if sameSize, ok = bySize[size]; !ok {
			sameSize = make([]string, 0, 1)
		}
		sameSize = append(sameSize, path)
		bySize[size] = sameSize
		return nil
	})
	for _, paths := range bySize {
		byHash := make(map[string][]string)
		for _, path := range paths {
			h := hash(path)
			var sameHash []string
			var ok bool
			if sameHash, ok = byHash[h]; !ok {
				sameHash = make([]string, 0, len(paths))
			}
			sameHash = append(sameHash, path)
			byHash[h] = sameHash
		}
		for _, paths := range byHash {
			shortest := paths[0]
			for _, path := range paths {
				if len(path) < len(shortest) {
					shortest = path
				}
			}
			for _, path := range paths {
				if path == shortest {
					continue
				}
				fmt.Printf("Removing %v (same as %v)\n", path, shortest)
				os.Remove(path)
			}
		}

	}
}

func hash(path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	h := sha512.New()
	_, err = io.Copy(h, f)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
