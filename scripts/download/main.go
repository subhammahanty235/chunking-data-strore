package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type ChunkMetadata struct {
	ChunkName string `json:"chunk_name"`
	Disk      string `json:"disk"`
	Order     int    `json:"order"`
}

type FileMeta struct {
	FileName string          `json:"file_name"`
	Chunks   []ChunkMetadata `json:"chunks"`
}

func main() {
	input := os.Args[1]
	bucketName, fileName := parseInput(input)
	err := downloadFile(bucketName, fileName)
	if err != nil {
		fmt.Println("download failed", err)
		os.Exit(1)
	}

	fmt.Println("Download successful:)")

}

func parseInput(input string) (string, string) {
	parts := filepath.SplitList(input)

	if len(parts) == 1 {
		s := input
		idx := -1
		for i := 0; i < len(s); i++ {
			if s[i] == '/' {
				idx = i
			}
		}
		fmt.Println(s)
		fmt.Println(idx)
		fmt.Println(s[:idx])
		return s[:idx], s[idx+1:]
	}
	return "", ""
}

func downloadFile(bucketName, fileName string) error {
	fmt.Println(bucketName)
	metadatafilename := fmt.Sprintf("%s.meta.json", fileName)
	metaPath := filepath.Join("buckets", bucketName, metadatafilename)
	metaData, err := os.ReadFile(metaPath)
	if err != nil {
		return err
	}
	var filemeta FileMeta
	err = json.Unmarshal(metaData, &filemeta)
	if err != nil {
		return err
	}

	if filemeta.FileName != fileName {
		return fmt.Errorf("File not found")
	}

	sort.Slice(filemeta.Chunks, func(i, j int) bool {
		return filemeta.Chunks[i].Order < filemeta.Chunks[j].Order
	})

	err = os.MkdirAll("downloaded_files", 0755)
	if err != nil {
		return err
	}

	outputpath := filepath.Join("downloaded_files", fileName)
	outputFile, err := os.Create(outputpath)

	defer outputFile.Close()

	for _, chunk := range filemeta.Chunks {
		chunkpath := filepath.Join(
			"storage_nodes",
			chunk.Disk,
			chunk.ChunkName,
		)

		data, err := os.ReadFile(chunkpath)
		if err != nil {
			return err
		}

		_, err = outputFile.Write(data)
		if err != nil {
			return err
		}
	}
	return nil
}
