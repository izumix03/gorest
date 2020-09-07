package gorest

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
)

type multipartSetting struct {
	key      string
	fileName string
	reader   io.Reader
}

func (cli *client) MultipartData(key string, reader io.Reader) Multipart {
	return cli.MultipartAsFormFile(key, "", reader)
}

func (cli *client) MultipartAsFormFile(key string, fileName string, reader io.Reader) Multipart {
	cli.multipartSettings = append(cli.multipartSettings, multipartSetting{
		key:      key,
		fileName: fileName,
		reader:   reader,
	})
	return cli
}

func (cli *client) setupMultipartRequest() (io.Reader, error) {
	var body bytes.Buffer
	multipartWriter := multipart.NewWriter(&body)
	var err error

	for _, v := range cli.multipartSettings {
		if err = func() error {
			reader := v.reader
			if x, ok := reader.(io.Closer); ok {
				defer func() {
					// ignore closing error
					_ = x.Close()
				}()
			}

			var writer io.Writer
			if file, ok := reader.(*os.File); ok {
				if writer, err = multipartWriter.CreateFormFile(v.key, file.Name()); err != nil {
					return err
				}
			} else if v.fileName != "" {
				if writer, err = multipartWriter.CreateFormFile(v.key, v.fileName); err != nil {
					return err
				}
			} else {
				if writer, err = multipartWriter.CreateFormField(v.key); err != nil {
					return err
				}
			}

			if _, err = io.Copy(writer, reader); err != nil {
				return err
			}

			if err := multipartWriter.Close(); err != nil {
				return err
			}
			return nil
		}(); err != nil {
			return nil, err
		}
	}

	cli.contentType = contentType(multipartWriter.FormDataContentType())
	return &body, nil
}
