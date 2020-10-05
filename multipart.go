package gorest

import (
	"bytes"
	"fmt"
	"io"
	mineMultipart "mime/multipart"
	"net/textproto"
	"os"
)

type multipartSetting struct {
	key            string
	fileName       string
	forceMultipart bool
	reader         io.Reader
}

func (cli *client) MultipartData(key string, reader io.Reader, forceMultipart bool) Multipart {
	return cli.MultipartAsFormFile(key, "", reader, forceMultipart)
}

func (cli *client) MultipartAsFormFile(key string, fileName string, reader io.Reader, forceMultipart bool) Multipart {
	cli.multipartSettings = append(cli.multipartSettings, multipartSetting{
		key:            key,
		fileName:       fileName,
		reader:         reader,
		forceMultipart: forceMultipart,
	})
	return cli
}

func (cli *client) setupMultipartRequest() (io.Reader, error) {
	var body bytes.Buffer
	multipartWriter := mineMultipart.NewWriter(&body)
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
				if writer, err = cli.createFormFileAsMultipart(
					v.key,
					file.Name(),
					reader,
					v.forceMultipart,
					multipartWriter,
				); err != nil {
					return err
				}
			} else if v.fileName != "" {
				if writer, err = cli.createFormFileAsMultipart(
					v.key,
					v.fileName,
					reader,
					v.forceMultipart,
					multipartWriter,
				); err != nil {
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

			if err = multipartWriter.Close(); err != nil {
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

func (cli *client) createFormFileAsMultipart(
	fieldName,
	fileName string,
	reader io.Reader,
	forceMultipart bool,
	writer *mineMultipart.Writer,
) (io.Writer, error) {
	if !forceMultipart {
		return writer.CreateFormFile(fieldName, fileName)
	}

	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, fileName))
	header.Set("Content-Type", string(multipart))
	return writer.CreatePart(header)
}
