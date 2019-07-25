package scaffold

import (
	"bufio"
	"io"
	"os"
	"regexp"
)

func readFile(filePath string, position string, replacer string) ([]byte, bool, error) {
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, false, err
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	needHandle := false
	output := make([]byte, 0)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				return output, needHandle, nil
			}
			return nil, needHandle, err
		}
		if ok, _ := regexp.Match(position, line); ok {
			reg := regexp.MustCompile(position)
			newByte := reg.ReplaceAll(line, []byte(replacer))
			output = append(output, newByte...)
			output = append(output, []byte("\n")...)
			if !needHandle {
				needHandle = true
			}
		} else {
			output = append(output, line...)
			output = append(output, []byte("\n")...)
		}
	}
	return output, needHandle, nil
}

func writeToFile(filePath string, outPut []byte) error {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(f)
	_, err = writer.Write(outPut)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}
