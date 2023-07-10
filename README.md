# 🐐hashgoat🐐

*⚡Blazing-fast⚡ ~~hash brute-forcing~~ password recovery tool written in Golang*

## Installation

* [Install `go`](https://go.dev/doc/install) 
* Install with `go install`: `go install github.com/diduk001/hashgoat@latest`

  *or*
* Build from source: `git clone https://github.com/diduk001/hashgoat && cd hashgoat && go get && go build`

## Usage

`hashgoat -w path-to-wordlist -a hashing-algorithm [-t number-of-threads] unknown-hash`

Currently implemented hashing algorithms:
- MD5 (`md5`)
- SHA1 (`sha1`)
- SHA256 (`sha256`)
- SHA512 (`sha512`)

## Examples

```
hashgoat -w wordlist.txt -a md5 dac0d8a5cf48040d1bb724ea18a4f103
hashgoat -w wordlist.txt -t 1 -a sha256 4e6dc79b64c40a1d2867c7e26e7856404db2a97c1d5854c3b3ae5c6098a61c62
```

*(Hashed string is `hashgoat`)*

## TODO

✅ Add basic hash algorithms (MD5, SHA1, SHA256, SHA512)

⬜ Add unit tests (IN PROGRESS)

⬜ Add asynchronous I/O

⬜ Add progress bar

⬜ Add automatic hash detection with regular expressions

⬜ Add benchmark

⬜ Add searching by mask

⬜ Add flag to compatibility with `hashcat` and `john` options


## Why "hashgoat"?

Because it's like [hashcat](https://github.com/hashcat/hashcat), but written in ***GO***. And it doesn't use CUDA. 
And I doubt it's actually better. But it's fun to write.