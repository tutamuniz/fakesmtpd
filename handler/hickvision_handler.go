package handler

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/tutamuniz/fakesmtpd/helper/chat"
)

// HickVision is a handler for hickvision
type HickVision struct {
	Logger  *log.Logger
	Chat    chat.Chat
	DataDir string
}

func (hv *HickVision) DoHelo(args string) bool {
	return true
}

func (hv *HickVision) DoMailFrom(args string) bool {
	return true
}

func (hv *HickVision) DoRcptTo(args string) bool {
	return true
}

func (hv *HickVision) DoData(d []byte) bool {
	files, err := hv.processData(d)
	if err != nil {
		hv.Logger.Println("RET DATA", err)
	}

	for _, file := range files {
		hv.Chat.SendMessage("", file)
	}

	line := strings.Join(files, ",")
	hv.Logger.Printf("Files: %s\n", line)
	return true
}

func (hv HickVision) String() string {
	return "HickVision Handler"
}

func (hv *HickVision) processData(data []byte) ([]string, error) {
	timestamp := time.Now().Format("2006-01-02-15-04-05")

	email := fmt.Sprintf("%s/email-%s.txt", hv.DataDir, timestamp)

	err := os.WriteFile(email, data, 0o644)
	if err != nil {
		hv.Logger.Printf("error writing file %s (%s)", email, err)
	}

	d := bytes.NewReader(data)
	files := []string{}
	rm, err := mail.ReadMessage(d)
	if err != nil {
		return files, err
	}

	_, params, err := mime.ParseMediaType(rm.Header.Get("Content-Type"))
	if err != nil {
		return files, err
	}

	mtr := multipart.NewReader(rm.Body, params["boundary"])

	for {
		p, err := mtr.NextPart()

		if err == io.EOF {
			break
		}

		header := p.Header.Get("Content-type")

		if strings.Contains(header, "image/jpeg") {

			filename := strings.ToLower(p.FileName())
			fullpath := fmt.Sprintf("%s/imagem-%s-%s", hv.DataDir, timestamp, filename)

			hv.Logger.Println("->", fullpath)

			files = append(files, fullpath)

			br := bytes.Buffer{}

			buffer := make([]byte, 2048)

			for {

				n, err := p.Read(buffer)

				if err == io.EOF && n == 0 {
					break
				}

				_, err = br.Write(buffer[:n])

				if err != nil {
					hv.Logger.Println("ERROR DATA:", err)
					break
				}

			}

			b, err := base64.StdEncoding.DecodeString(br.String())
			// improve this
			if err != nil {
				_ = os.WriteFile(fmt.Sprintf("%s/dump-%s-%s.log", hv.DataDir, timestamp, filename), br.Bytes(), 0o644)
				hv.Logger.Println("ERROR DUMP:", err)
				break
			}

			im, _ := os.Create(fullpath)
			defer im.Close()

			imr := bufio.NewWriter(im)
			imr.Write(b)
			imr.Flush()

		}
	}
	return files, nil
}
