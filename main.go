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
	isSync := flag.Bool("sync", false, "Read file, then hash each line (slower, uses more memory)")

	flag.Usage = printExample // Override standard usage
	flag.Parse()

	unknownHash := flag.Arg(0)

	if *wordlistPtr == "" {
		fmt.Println("Wordlist filename (-w) is not specified!")
		printExample()
		return
	}

	if *threadsCntPtr <= 0 {
		fmt.Println("Threads count must be not less than 0")
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
		fmt.Println(
			"Hash algorithm specified badly! Supported algorithms: md5, sha1, sha256, sha512",
		)
		printExample()
		return
	}

	if unknownHash == "" {
		fmt.Println("You need to specify unknown hash!")
		printExample()
		return
	}

	fmt.Printf("hashgoat - trying to recover %s\n", unknownHash)

	var isFound bool
	var result string
	if *isSync {
		isFound, result = runSync(*wordlistPtr, *threadsCntPtr, hashAlgo, unknownHash)
	} else {
		isFound, result = runAsync(*wordlistPtr, hashAlgo, unknownHash)
	}

	if isFound {
		fmt.Printf("Result: %s\n", result)
	} else {
		fmt.Println("Hash not found! Try another wordlist (-w) or hash algorithm (-a)")
	}
}

func printExample() {
	fmt.Println("Example:")
	fmt.Println("hashgoat -w wordlist.txt -a md5 -sync dac0d8a5cf48040d1bb724ea18a4f103")
	fmt.Println(
		"hashgoat -w wordlist.txt -t 1 -a sha256 4e6dc79b64c40a1d2867c7e26e7856404db2a97c1d5854c3b3ae5c6098a61c62",
	)
	fmt.Println()
	fmt.Println("More info: https://github.com/diduk001/hashgoat")
}

func runSync(wordlistFilename string, threadsCnt int, hashFunction func(string) string, unknownHash string) (bool, string) {
	lines, err := readLinesToList(wordlistFilename)
	if err != nil {
		fmt.Println("Error occurred while reading wordlist")
		panic(err)
	}
	isFound, result := recoverHashFromSlice(lines, threadsCnt, hashFunction, unknownHash)
	return isFound, result
}

func runAsync(wordlistFilename string, hashFunction func(string) string, unknownHash string) (bool, string) {
	lines := make(chan string)
	linesDone := make(chan struct{})
	hashResultChan := make(chan string)

	go readLinesToChan(wordlistFilename, lines, linesDone)
	go recoverHashFromChan(lines, linesDone, hashFunction, unknownHash, hashResultChan)

	select {
	case result := <-hashResultChan:
		if result == "" {
			return false, ""
		} else {
			return true, result
		}
	}
}

func recoverHashFromSlice(
	lines []string,
	threadsCnt int,
	hashFunction func(string) string,
	unknownHash string,
) (bool, string) {
	numLines := len(lines)

	if numLines == 0 {
		return false, ""
	} else if numLines == 1 {
		line := lines[0]
		hash := hashFunction(line)
		if hash == unknownHash {
			return true, line
		}
		return false, ""
	}

	hashPairs := make(chan pair)     // contains pairs {plain string, hashed string}
	foundString := make(chan string) // contains found plain string (if found any)
	isDone := make(chan struct{})    // closes when wait group is ready
	var wg sync.WaitGroup

	chunkSize := numLines / threadsCnt // Chunk is a slice of wordlist. Each thread operates with a chunk
	if chunkSize < 1 {
		chunkSize = 1
	}

	for i := 0; i < threadsCnt; i++ {
		wg.Add(1)

		sliceStart := i * chunkSize
		sliceEnd := (i + 1) * chunkSize
		if i == threadsCnt-1 { // Add the end of wordlist on last iteration
			sliceEnd = numLines
		}
		slice := lines[sliceStart:sliceEnd]
		go hashSlice(slice, &wg, hashPairs, isDone, hashFunction) // Start hashing goroutine
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
		return true, result
	case <-isDone: // if done before foundString contains something
		return false, ""
	}
}

func recoverHashFromChan(
	linesChan <-chan string,
	linesDoneChan <-chan struct{},
	hashFunction func(string) string,
	unknownHash string,
	hashResultChan chan string,
) {
	pairs := make(chan pair)

	go func() {
		for {
			select {
			case <-linesDoneChan:
				return
			case line := <-linesChan:
				hashedLine := hashFunction(line)
				pairs <- pair{line, hashedLine}
			}
		}
	}()

	go func() {
		for {
			select {
			case curPair := <-pairs:
				if curPair.hash == unknownHash {
					hashResultChan <- curPair.plain
					close(hashResultChan)
					return
				}
			}
		}
	}()
}

// Put pairs {hashed, plain} for passed hashFunction into pairs channel until all lines from slice is hashed or done channel is closed
func hashSlice(
	wordlist []string,
	wg *sync.WaitGroup,
	pairs chan<- pair,
	done <-chan struct{},
	hashFunction func(string) string,
) {
	defer wg.Done()
	if len(wordlist) == 0 {
		return
	}

	for _, line := range wordlist {
		hash := hashFunction(line)
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

func readLinesToList(filename string) ([]string, error) {
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

func readLinesToChan(filename string, lines chan<- string, linesDone chan<- struct{}) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error during opening file")
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines <- line
	}

	close(linesDone)
}
