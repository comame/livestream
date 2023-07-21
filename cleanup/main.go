package main

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/comame/readenv-go"
)

type Env struct {
	HlsPath string `env:"HLS_PATH"`
}

const DurationCleanup = 15 * time.Minute

func main() {
	var env Env
	readenv.Read(&env)

	log.Println("cleanup")

	hlsPath, err := filepath.Abs(env.HlsPath)
	if err != nil {
		panic(err)
	}

	keys, err := listStreamKeys(hlsPath)
	if err != nil {
		panic(err)
	}

	for _, key := range keys {
		ended, err := isEndedStreamKey(hlsPath, key)
		if err != nil {
			panic(err)
		}

		if ended {
			log.Printf("rm-stream %s\n", key)
			rmStreamKey(hlsPath, key)
			continue
		}

		files, err := listExpiredTsFile(hlsPath, key)
		if err != nil {
			panic(err)
		}
		log.Printf("rm-file %s\n", strings.Join(files, ", "))
		rmTsFile(hlsPath, key, files)
	}
}

func isEndedStreamKey(hlsPath string, streamKey string) (bool, error) {
	p := path.Join(hlsPath, streamKey, "index.m3u8")
	stat, err := os.Lstat(p)
	if err != nil {
		return false, err
	}

	return stat.ModTime().Before(time.Now().Add(-DurationCleanup)), nil
}

func listExpiredTsFile(hlsPath string, streamKey string) ([]string, error) {
	folderPath := path.Join(hlsPath, streamKey)
	ls, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, name := range ls {
		p := path.Join(folderPath, name.Name())
		stat, err := os.Lstat(p)
		if err != nil {
			return nil, err
		}

		if !strings.HasSuffix(stat.Name(), ".ts") {
			continue
		}

		if stat.ModTime().Before(time.Now().Add(-DurationCleanup)) {
			result = append(result, stat.Name())
		}
	}

	return result, nil
}

func listStreamKeys(hlsPath string) ([]string, error) {
	dirs, err := os.ReadDir(hlsPath)
	if err != nil {
		return nil, err
	}

	var ls []string
	for _, dir := range dirs {
		if !dir.Type().IsDir() {
			continue
		}
		ls = append(ls, dir.Name())
	}

	return ls, nil
}

func rmStreamKey(hlsPath string, streamKey string) {
	folderPath := path.Join(hlsPath, streamKey)
	os.RemoveAll(folderPath)
}

func rmTsFile(hlsPath, streamKey string, files []string) {
	for _, file := range files {
		p := path.Join(hlsPath, streamKey, file)
		os.Remove(p)
	}
}
