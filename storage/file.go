package storage

import (
	"gitlab.com/likekimo/goproxyspider/pkg/logger"
	"go.uber.org/zap"

	"os"
)

//fileStorage структура для хранения в файле
type fileStorage struct {
	filePath    string
	filePathDef string
	logger      *logger.Logger
}

func FileStorage(logger *logger.Logger, path string) Storage {
	svc := &fileStorage{
		filePath:    path,
		logger:      logger,
		filePathDef: "./proxy.txt",
	}

	return svc
}

func (fs fileStorage) CheckStorage() error {
	return nil
}

func (fs fileStorage) WriteStorage(source string) error {
	var filename string
	if fs.filePath == "" {
		filename = fs.filePathDef
	} else {
		filename = fs.filePath
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fs.logger.Error("Write proxy error", zap.Error(err))
		file, err = os.Create(filename)
	}

	if err != nil {
		fs.logger.Error("Write proxy error", zap.Error(err))
	} else {
		wString := source + "\n"
		file.WriteString(wString)
		// fs.logger.Log("msg", "Write proxy")
	}
	file.Close()
	return err
}
func (fs fileStorage) Stat() error {
	fs.logger.Info("Count file stat")
	return nil
}
