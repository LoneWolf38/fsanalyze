package analyze

import (
	"fmt"
	"os"

	"github.com/VictoriaMetrics/fastcache"
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

	fmt.Println("FileLength: ", ByteCountIEC(uint64(fileLength)))

	fmt.Println("Decoding the fileSummary Length")
	fSummaryLength, err := decodeFileSummaryLength(fileLength, f)
	if err != nil {
		return err
	}
	fmt.Println("FileSummary Length: ", ByteCountIEC(uint64(fSummaryLength)))
	fmt.Println("parsing file Summary")
	sectionMap, err := parseFileSummary(f, fileLength, fSummaryLength)
	if err != nil {
		return err
	}
	fmt.Println("section Map: ", sectionMap)
	c, err := fastcache.LoadFromFile(CacheFile)
	inodeSectionInfo := sectionMap["INODE"]
	fmt.Println("InodeSection Length: ", ByteCountIEC(inodeSectionInfo.GetLength()))
	fmt.Println("Parsing Inode Section")
	_, entityCount, err := parseInodeSection(inodeSectionInfo, f, c)
	if err != nil {
		return err
	}
	fmt.Println("Total Number of Files: ", entityCount.Files)
	fmt.Println("Total Number of Directories: ", entityCount.Directories)
	fmt.Println("Total Number of Symlinks: ", entityCount.Symlinks)

	inodeDirectorySectionInfo := sectionMap["INODE_DIR"]
	fmt.Println("Parsing Inode_Dir Section")
	fmt.Println("InodeDirectory Section Length: ", ByteCountIEC(inodeDirectorySectionInfo.GetLength()))
	_, err = parseInodeDirectorySection(inodeDirectorySectionInfo, f)
	if err != nil {
		return err
	}
	return nil
}
