package main

import (
	"testing"
)

func TestMD5(t *testing.T) {
	t.Run("TestOneWordMD5", func(t *testing.T) {
		wordlist := []string{"test_md5"}

		wanted := wordlist[0]
		hash := "9050bddcf415f2d0518804e551c1be98"
		isFound, result := recoverHashFromSlice(wordlist, 1, prettyMD5, hash)
		if !isFound {
			t.Errorf("Not found MD5 hash %s in wordlist. Plaintest is %s", hash, wanted)
		} else if result != wanted {
			t.Errorf("Wrong plaintext for MD5 hash. Got %s, wanted %s, hash - %s", result, wanted, hash)
		}
	})

	t.Run("TestHashNotFoundMD5", func(t *testing.T) {
		wordlist := []string{"a", "b", "c", "d"}

		hash := "e1671797c52e15f763380b45e841ec32" // md5("e")
		isFound, result := recoverHashFromSlice(wordlist, 1, prettyMD5, hash)
		if isFound || result != "" {
			t.Errorf(
				"Found hash which is not in wordlist. Requested hash - %s, result - %s",
				hash,
				result,
			)
		}
	})
}

func TestSHA1(t *testing.T) {
	t.Run("TestOneWordSHA1", func(t *testing.T) {
		wordlist := []string{"test_sha1"}

		wanted := wordlist[0]
		hash := "9db4507552981975bccac89a41dab2cc821bff2e"
		isFound, result := recoverHashFromSlice(wordlist, 1, prettySHA1, hash)
		if !isFound {
			t.Errorf("Not found SHA1 hash %s in wordlist. Plaintest is %s", hash, wanted)
		} else if result != wanted {
			t.Errorf("Wrong plaintext for SHA1 hash. Got %s, wanted %s, hash - %s", result, wanted, hash)
		}
	})

	t.Run("TestHashNotFoundSHA1", func(t *testing.T) {
		wordlist := []string{"a", "b", "c", "d"}

		hash := "58e6b3a414a1e090dfc6029add0f3555ccba127f" // sha1("e")
		isFound, result := recoverHashFromSlice(wordlist, 1, prettyMD5, hash)
		if isFound || result != "" {
			t.Errorf(
				"Found hash which is not in wordlist. Requested hash - %s, result - %s",
				hash,
				result,
			)
		}
	})
}

func TestOneWordSHA256(t *testing.T) {
	wordlist := []string{"test_sha256"}

	wanted := wordlist[0]
	hash := "fda177bb1336270b24e4df0fd0c1dd0596c44699204f57c83ce70a0f19173be4"
	isFound, result := recoverHashFromSlice(wordlist, 1, prettySHA256, hash)
	if !isFound {
		t.Errorf("Not found SHA256 hash %s in wordlist. Plaintest is %s", hash, wanted)
	} else if result != wanted {
		t.Errorf("Wrong plaintext for SHA256 hash. Got %s, wanted %s, hash - %s", result, wanted, hash)
	}
}

func TestOneWordSHA512(t *testing.T) {
	wordlist := []string{"test_sha512"}

	wanted := wordlist[0]
	hash := "e335ec8aa0e729469a06c50fe8f93621b544970ebdb99ab6351368f3541f63fc37ed92bb2fee40549de8ebfeb167386859391866541d9578684ec06ea7a70cea"
	isFound, result := recoverHashFromSlice(wordlist, 1, prettySHA512, hash)
	if !isFound {
		t.Errorf("Not found SHA512 hash %s in wordlist. Plaintest is %s", hash, wanted)
	} else if result != wanted {
		t.Errorf("Wrong plaintext for SHA512 hash. Got %s, wanted %s, hash - %s", result, wanted, hash)
	}
}

func TestEmptyWordlist(t *testing.T) {
	var wordlist []string

	hash := "dac0d8a5cf48040d1bb724ea18a4f103" // md5(hashgoat)
	isFound, result := recoverHashFromSlice(wordlist, 1, prettyMD5, hash)
	if isFound || result != "" {
		t.Errorf("Found hash in empty wordlist. Requested hash - %s, result - %s", hash, result)
	}
}

func TestEmptyWordlistWithOneWord(t *testing.T) {
	wordlist := []string{"a"}

	hash := "dac0d8a5cf48040d1bb724ea18a4f103" // md5(hashgoat)
	isFound, result := recoverHashFromSlice(wordlist, 1, prettyMD5, hash)
	if isFound || result != "" {
		t.Errorf("Found hash in wrong wordlist. Requested hash - %s, result - %s", hash, result)
	}
}
