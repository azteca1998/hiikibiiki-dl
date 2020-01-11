package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"mime"
	"net/http"
	"os"
	"path"
)

const indexUrl = "http://hkbk.fm/?feed=podcast"

func main() {
	log.Println("Downloading podcasts index...")
	res, err := http.Get(indexUrl)
	if err != nil {
		panic(err)
	}

	resDec := xml.NewDecoder(res.Body)

	log.Println("Decoding podcasts index")
	var resData Rss
	err = resDec.Decode(&resData)
	if err != nil {
		panic(err)
	}

	// There should only be one channel, but you never know...
	for _, channel := range resData.Channels {
		err = downloadChannel(&channel)
		if err != nil {
			panic(err)
		}
	}
}

func getDefaultPerms(perms *os.FileMode) (err error) {
	stat, err := os.Stat(".")
	if err != nil {
		return
	}

	*perms = stat.Mode()

	return nil
}

func downloadChannel(channel *RssChannel) (err error) {
	log.Println("Downloading channel", channel.Title)

	log.Println()
	log.Println("  Title:", channel.Title)
	log.Println("  Description:", channel.Description)
	log.Println("  Link:", channel.Link)
	log.Println()

	var defaultPerms os.FileMode
	err = getDefaultPerms(&defaultPerms)
	if err != nil {
		return
	}

	channelPath := path.Join(channel.Title, "")
	err = os.MkdirAll(channelPath, defaultPerms)
	if err != nil {
		return
	}

	err = downloadFile(path.Join(channelPath, channel.Image.Title), channel.Image.Url, nil, nil)
	if err != nil {
		return
	}

	for idx, item := range channel.Items {
		err = downloadFile(path.Join(channelPath, item.Title), item.Enclosure.Url, &item.Enclosure.Type, &Pos{
			Idx: idx,
			Len: len(channel.Items),
		})
		if err != nil {
			return
		}
	}

	log.Println("Done! Enjoy :D")

	return nil
}

type Pos struct {
	Idx int
	Len int
}

func downloadFile(path string, url string, mimeType *string, pos *Pos) (err error) {
	prefix := " "
	if pos != nil {
		lenLen := uint(math.Ceil(math.Log10(float64(pos.Len))))
		prefix = fmt.Sprintf("  [%*d/%d]", lenLen, pos.Idx+1, pos.Len)
	}

	var checkFileSize int64 = -1
	if mimeType == nil {
		// Query for the extension.
		res, err := http.Head(url)
		if err != nil {
			return err
		}

		tmp := res.Header.Get("Content-Type")
		mimeType = &tmp

		checkFileSize = res.ContentLength
	}

	extList, err := mime.ExtensionsByType(*mimeType)
	if err != nil {
		return err
	}
	path += extList[0]

	// Do nothing if the file already exists and has the correct size.
	if stat, err := os.Stat(path); !os.IsNotExist(err) && (checkFileSize == -1 || stat.Size() == checkFileSize) {
		log.Println(prefix, "Skipping", path)
		return nil
	}

	// Download the file.
	log.Println(prefix, "Downloading", path)

	out, err := os.Create(path)
	if err != nil {
		return
	}

	res, err := http.Get(url)
	if err != nil {
		return
	}

	n, err := io.Copy(out, res.Body)
	if err != nil {
		return
	}
	if n != res.ContentLength {
		return errors.New("file length mismatch")
	}

	err = out.Close()
	if err != nil {
		return
	}

	return nil
}
