package analyze

import (
	"bufio"
	"fmt"
	"os"
)

const (
	FILE_SUM_BYTES = 4
	ROOT_INODE_ID  = 16385
)

func Parse(filePath string) error {

	fInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("cannot find the %s file", filePath)
	}
	if err != nil {
		return fmt.Errorf("cannot access the %s file, Reason: %s", filePath, err.Error())
	}
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("cannot open the %s file, Reason: %s", filePath, err.Error())
	}
	defer f.Close()
	fileLength := fInfo.Size()
	bReader := bufio.NewReader(f)

	fmt.Println("Decoding the fileSummary Length")
	fSummaryLength, err := decodeFileSummaryLength(fileLength, f, bReader)
	if err != nil {
		return err
	}
	fmt.Println("parsing file Summary")
	sectionMap, err := parseFileSummary(f, fileLength, fSummaryLength, bReader)
	if err != nil {
		return err
	}
	fmt.Println("section Map: ", sectionMap)

	return nil
}
