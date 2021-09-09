package sort

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const (
	tmpFolderPath = "large_sort_tmp_folder"

	leftBatchName  = "l_batch.txt"
	rightBatchName = "r_batch.txt"

	inputMirror  = "in_mirror.txt"
	outputMirror = "out_mirror.txt"
)

func initBuffers() (*os.File, *os.File, *os.File, *os.File, error) {
	files := make([]*os.File, 0, 4)
	for _, name := range []string{leftBatchName, rightBatchName, inputMirror, outputMirror} {
		path := filepath.Join(tmpFolderPath, name)
		file, err := os.OpenFile(path, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			for _, f := range files {
				if f != nil {
					_ = f.Close()
				}
			}
			return nil, nil, nil, nil, fmt.Errorf("cannot create file: %s: %w", path, err)
		}
		files = append(files, file)
	}
	return files[0], files[1], files[2], files[3], nil
}

func truncateFile(file *os.File) error {
	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("cannot truncate file: %s :%w", file.Name(), err)
	}
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("cannot seek file: %s :%w", file.Name(), err)
	}
	return nil
}

func truncateFiles(files ...*os.File) error {
	for _, file := range files {
		if err := truncateFile(file); err != nil {
			return err
		}
	}
	return nil
}

func copyLines(maxLines int, reader *bufio.Scanner, writer io.Writer) error {
	for i := 0; i < maxLines && reader.Scan(); i++ {
		if _, err := fmt.Fprintln(writer, reader.Text()); err != nil {
			return fmt.Errorf("cannot write line: %w", err)
		}
	}
	if reader.Err() != nil {
		return fmt.Errorf("cannot scan %d lines: %w", maxLines, reader.Err())
	}
	return nil
}

func mergeBatches(lBatch, rBatch *bufio.Scanner, out io.Writer) (err error) {
	var lStr, rStr string
	lOk, rOk := lBatch.Scan(), rBatch.Scan()
	if lOk {
		lStr = lBatch.Text()
	}
	if rOk {
		rStr = rBatch.Text()
	}
	for lOk && rOk {
		if lStr < rStr {
			if _, err = fmt.Fprintln(out, lStr); err != nil {
				return
			}
			lOk = lBatch.Scan()
			if lOk {
				lStr = lBatch.Text()
			}
		} else {
			if _, err = fmt.Fprintln(out, rStr); err != nil {
				return
			}
			rOk = rBatch.Scan()
			if rOk {
				rStr = rBatch.Text()
			}
		}
	}
	for lOk {
		if _, err = fmt.Fprintln(out, lStr); err != nil {
			return
		}
		lOk = lBatch.Scan()
		if lOk {
			lStr = lBatch.Text()
		}
	}
	for rOk {
		if _, err = fmt.Fprintln(out, rStr); err != nil {
			return
		}
		rOk = rBatch.Scan()
		if rOk {
			rStr = rBatch.Text()
		}
	}
	if lBatch.Err() != nil {
		return fmt.Errorf("cannot merge: left batch: %w", lBatch.Err())
	}
	if rBatch.Err() != nil {
		return fmt.Errorf("cannot merge: right batch: %w", rBatch.Err())
	}
	return
}

func Sort(inputPath string, outputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("cannot open input file: %w", err)
	}
	defer func() { _ = inputFile.Close() }()

	// buffers files tmp folder
	if err = os.Mkdir(tmpFolderPath, 0777); err != nil {
		return fmt.Errorf("cannot create tmp dir: %w", err)
	}
	defer func() { _ = os.RemoveAll(tmpFolderPath) }()

	// buffers files for sorting
	lBatchFile, rBatchFile, inputMirrorFile, outputMirrorFile, err := initBuffers()
	if err != nil {
		return fmt.Errorf("cannot initialize sort buffers: %w", err)
	}
	defer func() {
		_ = lBatchFile.Close()
		_ = rBatchFile.Close()
		_ = inputMirrorFile.Close()
		_ = outputMirrorFile.Close()
	}()

	linesCount := 0

	// copy input file to input mirror
	inputScanner := bufio.NewScanner(inputFile)
	inputMirrorWriter := bufio.NewWriter(inputMirrorFile)
	for inputScanner.Scan() {
		if _, err = fmt.Fprintln(inputMirrorWriter, inputScanner.Text()); err != nil {
			return fmt.Errorf("cannot copy input file to mirror: %w", err)
		}
		linesCount++
	}
	if inputScanner.Err() != nil {
		return fmt.Errorf("cannot copy input file to input mirror: %w", inputScanner.Err())
	}
	if err = inputMirrorWriter.Flush(); err != nil {
		return fmt.Errorf("cannot flush mirror: %w", err)
	}

	for batchSize := 1; batchSize < linesCount; batchSize *= 2 {

		// truncate used output mirror
		if err = truncateFile(outputMirrorFile); err != nil {
			return fmt.Errorf("cannot truncate mirror: %w", err)
		}

		// input mirror was updated before
		if _, err = inputMirrorFile.Seek(0, 0); err != nil {
			return fmt.Errorf("cannot seek file: %s :%w", inputMirrorFile.Name(), err)
		}

		inputMirrorScanner := bufio.NewScanner(inputMirrorFile)
		outputMirrorWriter := bufio.NewWriter(outputMirrorFile)

		for line := 0; line < linesCount; line += batchSize {
			// truncate used batches buffers
			if err = truncateFiles(lBatchFile, rBatchFile); err != nil {
				return fmt.Errorf("cannot truncate batches: %w", err)
			}

			// copy left batch from input mirror
			lBatchW, rBatchW := bufio.NewWriter(lBatchFile), bufio.NewWriter(rBatchFile)
			if err = copyLines(batchSize, inputMirrorScanner, lBatchW); err != nil {
				return fmt.Errorf("cannot copy left batch: %w", err)
			}
			if err = lBatchW.Flush(); err != nil {
				return fmt.Errorf("cannot flush left batch: %w", err)
			}

			// copy right batch from input mirror
			if err = copyLines(batchSize, inputMirrorScanner, rBatchW); err != nil {
				return fmt.Errorf("cannot copy right batch: %w", err)
			}
			if err = rBatchW.Flush(); err != nil {
				return fmt.Errorf("cannot flush right batch: %w", err)
			}

			// we need read from batches later
			if _, err = lBatchFile.Seek(0, 0); err != nil {
				return fmt.Errorf("cannot seek file: %s :%w", lBatchFile.Name(), err)
			}
			if _, err = rBatchFile.Seek(0, 0); err != nil {
				return fmt.Errorf("cannot seek file: %s :%w", rBatchFile.Name(), err)
			}

			// merge batches to output mirror
			lBatch, rBatch := bufio.NewScanner(lBatchFile), bufio.NewScanner(rBatchFile)
			if err = mergeBatches(lBatch, rBatch, outputMirrorWriter); err != nil {
				return fmt.Errorf("cannot merge batches: %w", err)
			}
			if err = outputMirrorWriter.Flush(); err != nil {
				return fmt.Errorf("cannot flush mirror: %w", err)
			}
		}

		// swap mirrors
		inputMirrorFile, outputMirrorFile = outputMirrorFile, inputMirrorFile
	}

	// input mirror was updated before
	if _, err = inputMirrorFile.Seek(0, 0); err != nil {
		return fmt.Errorf("cannot seek mirror: %w", err)
	}

	// if output path is input path => we just open this file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("cannot open output file: %w", err)
	}
	defer func() { _ = outputFile.Close() }()

	// copy sorted data to output file
	if _, err = io.Copy(bufio.NewWriter(outputFile), bufio.NewReader(inputMirrorFile)); err != nil {
		return fmt.Errorf("cannot copy sorted data from mirror: %w", err)
	}

	return nil
}
