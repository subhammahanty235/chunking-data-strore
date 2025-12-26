package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	chunksize = 1024 * 1024
)

var disks = []string{
	"storage_nodes/node1",
	"storage_nodes/node2",
	"storage_nodes/node3",
	"storage_nodes/node4",
	"storage_nodes/node5",
	"storage_nodes/node6",
}

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
	bucketName := os.Args[1]
	filename := os.Args[2]

	err := uploadFile(bucketName, filename)
	if err != nil {
		fmt.Println("Upload failed", err)
		os.Exit(1)
	}

	fmt.Println("Upload successful:)")

}

func uploadFile(bucketName, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	fileName := filepath.Base(filePath)
	var chunks []ChunkMetadata
	buffer := make([]byte, chunksize)

	part := 0
	diskIndex := 0

	for {
		n, err := file.Read(buffer)
		if n > 0 {
			chunkName := fmt.Sprintf("%s.part-%d", fileName, part)
			targetDisk := disks[diskIndex%len(disks)]

			chunkPath := filepath.Join(targetDisk, chunkName)

			err := os.WriteFile(chunkPath, buffer[:n], 0644)
			if err != nil {
				return err
			}

			chunks = append(chunks, ChunkMetadata{
				ChunkName: chunkName,
				Disk:      filepath.Base(targetDisk),
				Order:     part,
			})

			part++
			diskIndex++
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return nil
		}
	}
	return writeMetaData(bucketName, fileName, chunks)
}

func writeMetaData(bucketName string, fileName string, chunks []ChunkMetadata) error {
	bucketPath := filepath.Join("buckets", bucketName)
	err := os.MkdirAll(bucketPath, 0755)
	if err != nil {
		return err
	}
	metadatafilename := fmt.Sprintf("%s.meta.json", fileName)
	metapath := filepath.Join(bucketPath, metadatafilename)
	fileMeta := FileMeta{
		FileName: fileName,
		Chunks:   chunks,
	}

	data, err := json.MarshalIndent(fileMeta, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(metapath, data, 0644)
}
