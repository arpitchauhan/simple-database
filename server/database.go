package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

type ErrorCode uint32

const (
	OK            ErrorCode = 0
	KeyNotFound   ErrorCode = 1
	InternalError ErrorCode = 2
)

type database struct {
	filepath     string
	initialized  bool
	keyPositions map[string]int64
}

func (d *database) openForReading() (*os.File, ErrorCode) {
	f, err := os.OpenFile(d.filepath, os.O_RDONLY|os.O_CREATE, 0o644)
	if err != nil {
		log.Printf("Failed to open the database file: %v", err)
		return nil, InternalError
	}
	return f, OK
}

func (d *database) openForWriting() (*os.File, ErrorCode) {
	f, err := os.OpenFile(d.filepath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		log.Printf("Failed to open the database file: %v", err)
		return nil, InternalError
	}
	return f, OK
}

func (d *database) ensureInitialized() {
	if !d.initialized {
		log.Fatal("Database should have been initialized by now")
	}
}

func (d *database) initialize() ErrorCode {
	if d.initialized {
		log.Fatal("Database was already initialized")
	}

	code := d.initializeKeyPositions()

	if code != OK {
		return code
	}

	d.initialized = true

	return OK
}

func (d *database) initializeKeyPositions() ErrorCode {
	d.keyPositions = make(map[string]int64)

	f, code := d.openForReading()
	if code != OK {
		return code
	}

	csvReader := csv.NewReader(f)

	for {
		pos := csvReader.InputOffset()
		record, err := csvReader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Error while reading file: %v", err)
			return InternalError
		}

		key := record[0]
		d.updateKeyPosition(key, pos)
	}

	return OK
}

func (d *database) getKeyPosition(key string) (bool, int64) {
	keyPosition, keyFound := d.keyPositions[key]

	return keyFound, keyPosition
}

func (d *database) updateKeyPosition(key string, pos int64) {
	d.keyPositions[key] = pos
}

func (d *database) getKey(key string) (string, ErrorCode) {
	d.ensureInitialized()
	keyFound, keyPosition := d.getKeyPosition(key)

	if !keyFound {
		return "", KeyNotFound
	}

	f, code := d.openForReading()
	defer f.Close()

	if code != OK {
		return "", code
	}

	_, err := f.Seek(keyPosition, io.SeekStart)
	if err != nil {
		log.Printf("Failed to seek position of key: %v", err)
		return "", InternalError
	}

	csvReader := csv.NewReader(f)
	record, err := csvReader.Read()
	if err != nil {
		log.Printf("Error while reading: %v", err)
		return "", InternalError
	}

	readKey := record[0]

	if readKey != key {
		log.Printf("Key at stored position is not correct")
		return "", InternalError
	}

	return record[1], OK
}

func (d *database) setKey(key string, value string) ErrorCode {
	d.ensureInitialized()

	f, code := d.openForWriting()
	defer f.Close()

	if code != OK {
		return code
	}

	currentPosition, err := f.Seek(0, io.SeekEnd)
	if err != nil {
		log.Printf("Error while getting current position in file: %v", err)
		return InternalError
	}

	csvWriter := csv.NewWriter(f)

	if err := csvWriter.Write([]string{key, value}); err != nil {
		log.Printf("Error while writing to file: %v", err)
		return InternalError
	}

	csvWriter.Flush()

	if err := csvWriter.Error(); err != nil {
		log.Printf("Error after flushing: %v", err)
		return InternalError
	}

	d.updateKeyPosition(key, currentPosition)
	return OK
}
