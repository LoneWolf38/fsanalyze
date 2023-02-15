package analyze

import (
	"encoding/binary"
	"io"
	"os"

	pb "fsanalyze/protofiles/src/hadoop_hdfs_fsimage"
)

func decodeFileSummaryLength(fileLength int64, imageFile *os.File, r io.Reader) (int32, error) {
	var (
		fSumLenBytes   = make([]byte, FILE_SUM_BYTES)
		fSummaryLength int32
	)
	fileSummaryLengthStart := fileLength - FILE_SUM_BYTES
	_, err := imageFile.ReadAt(fSumLenBytes, fileSummaryLengthStart)
	if err != nil {
		if err != io.EOF {
			return fSummaryLength, nil
		}
	}
	if err = binary.Read(r, binary.BigEndian, &fSummaryLength); err != nil {
		return fSummaryLength, nil
	}
	return fSummaryLength, nil
}

func parseFileSummary(imageFile *os.File, fileLength int64, fSummaryLength int32, r io.Reader) (map[string]*pb.FileSummary_Section, error) {
	var (
		sectionMap  = make(map[string]*pb.FileSummary_Section)
		fileSummary = &pb.FileSummary{}
	)
	// last 4 bytes says how many bytes should be read from end to get the FileSummary message
	fSummaryLength64 := int64(fSummaryLength)
	readAt := fileLength - fSummaryLength64 - FILE_SUM_BYTES
	_, err := imageFile.Seek(readAt, io.SeekCurrent)
	if err != nil {
		return sectionMap, err
	}
	_, err = readDelimited(r, fileSummary)
	if err != nil {
		return sectionMap, err
	}
	for _, value := range fileSummary.GetSections() {
		sectionMap[value.GetName()] = value
	}
	return sectionMap, nil
}
