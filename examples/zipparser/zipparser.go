package main

import (
	"archive/zip"
	"io"

	"github.com/NextronSystems/jsonlog/thorlog/v3"
	"github.com/NextronSystems/thor-plugin"
)

func Init(config thor.Configuration, logger thor.Logger, actions thor.RegisterActions) {
	actions.AddYaraRule(thor.TypeMeta, `
rule DetectZipFiles: ZIPFILE {
    meta:
        score = 0
    condition: filetype == "ZIP"
}`)
	actions.AddRuleHook("ZIPFILE", func(scanner thor.Scanner, object thor.MatchingObject) {
		file, isFile := object.Object.(*thorlog.File)
		if !isFile {
			return
		}
		scanner.Debug("Scanning ZIP file", "path", file.Path)
		zipReader, err := zip.NewReader(object.Content, object.Content.Size())
		if err != nil {
			scanner.Error("Could not parse zip file", "path", file.Path, "error", err)
			return
		}
		for _, file := range zipReader.File {
			scanFile(config, scanner, object, file)
		}
	})
	logger.Info("ZipParser plugin loaded!")
}

func scanFile(config thor.Configuration, scanner thor.Scanner, object thor.MatchingObject, file *zip.File) {
	fileReader, err := file.Open()
	if err != nil {
		scanner.Error("Could not open file in zip file", "file", file.Name, "error", err)
		return
	}
	defer fileReader.Close()
	if file.UncompressedSize64 > config.MaxFileSize {
		scanner.Error("File in zip file is too large for analysis", "file", file.Name, "size", file.UncompressedSize64)
		return
	}
	data, err := io.ReadAll(fileReader)
	if err != nil {
		scanner.Error("Could not read file in zip file", "file", file.Name, "error", err)
		return
	}
	// This is possibly recursive: if the data is (another) ZIP, our hook will be called again
	scanner.ScanFile(file.Name, data, "ZIP")
}
