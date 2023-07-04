package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"os"
	"sync"
)

func main() {
	wordlistPtr := flag.String("w", "", "Path to wordlist")
	threadsCntPtr := flag.Int("t", 10, "Number of threadsCnt")
	hashAlgoPtr := flag.String("a", "", "Hashing algorithm (md5/sha1/sha256/sha512)")
	flag.Usage = printExample // Override standard usage
	flag.Parse()

	unknownHash := flag.Arg(0)

	if *wordlistPtr == "" {
		fmt.Println("Wordlist filename (-w) is not specified!")
		printExample()
		return
	}

	var hashAlgo func(string) string // Hashing algorithm
	switch *hashAlgoPtr {
	case "":
		fmt.Println("Hash not specified!")
		printExample()
		return
	case "md5":
		hashAlgo = prettyMD5
	case "sha1":
		hashAlgo = prettySHA1
	case "sha256":
		hashAlgo = prettySHA256
	case "sha512":
		hashAlgo = prettySHA512
	default:
		fmt.Println("Hash algorithm specified badly! Supported algorithms: md5, sha1, sha256, sha512")
		printExample()
		return
	}

	if unknownHash == "" {
		fmt.Println("You need to specify unknown hash!")
		printExample()
		return
	}

	lines, err := readLines(*wordlistPtr)
	if err != nil {
		fmt.Println("Error occurred while reading wordlist")
		panic(err)
		return
	}

	fmt.Printf("hashgoat - trying to recover %s\n", unknownHash)
	numLines := len(lines)

	hashPairs := make(chan pair)     // contains pairs {plain string, hashed string}
	foundString := make(chan string) // contains found plain string (if found any)
	isDone := make(chan string)      // closes when wait group is ready
	var wg sync.WaitGroup

	threadsCnt := *threadsCntPtr
	chunkSize := numLines / threadsCnt // Chunk is a slice of wordlist. Each thread operates with a chunk
	for i := 0; i < threadsCnt; i++ {
		wg.Add(1)

		sliceStart := i * chunkSize
		sliceEnd := (i + 1) * chunkSize
		if i == threadsCnt-1 { // Add the end of wordlist on last iteration
			sliceEnd = numLines
		}
		slice := lines[sliceStart:sliceEnd]
		go hashSlice(slice, &wg, hashPairs, isDone, hashAlgo) // Start hashing goroutine
	}

	// Checks if there is a hash equal to user's input in hashPairs
	go func() {
		for p := range hashPairs {
			if p.hash == unknownHash {
				foundString <- p.plain
				close(foundString)
				return
			}
		}
	}()

	go func() {
		wg.Wait()
		close(isDone)
	}()

	select {
	case result := <-foundString: // if foundString contains something
		fmt.Printf("Result: %s\n", result)
	case <-isDone: // if done before foundString contains something
		fmt.Println("Hash not found! Try another wordlist (-w) or hash algorithm (-a)")
	}
}

func printExample() {
	fmt.Println("Example:")
	fmt.Println("hashgoat -w wordlist.txt -a md5 dac0d8a5cf48040d1bb724ea18a4f103")
	fmt.Println("hashgoat -w wordlist.txt -t 1 -a sha256 4e6dc79b64c40a1d2867c7e26e7856404db2a97c1d5854c3b3ae5c6098a61c62")
	fmt.Println()
	fmt.Println("More info: https://github.com/diduk001/hashgoat")
}

func hashSlice(wordlist []string, wg *sync.WaitGroup, pairs chan<- pair, done <-chan string, hashFunc func(string) string) {
	defer wg.Done()
	for _, line := range wordlist {
		hash := hashFunc(line)
		select {
		case pairs <- pair{line, hash}:
		case <-done:
			return
		}
	}
}

type pair struct {
	plain string
	hash  string
}

func prettyMD5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func prettySHA1(s string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(s)))
}

func prettySHA256(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}
func prettySHA512(s string) string {
	return fmt.Sprintf("%x", sha512.Sum512([]byte(s)))
}

func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}
