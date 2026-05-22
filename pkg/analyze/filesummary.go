package analyze

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	pb "fsanalyze/protofiles/src/hadoop_hdfs_fsimage"

	"google.golang.org/protobuf/proto"
)

func decodeFileSummaryLength(fileLength int64, imageFile *os.File) (int32, error) {
	var (
		fSumLenBytes   = make([]byte, FILE_SUM_BYTES)
		fSummaryLength int32
	)
	fileSummaryLengthStart := fileLength - FILE_SUM_BYTES
	bReader := bytes.NewReader(fSumLenBytes)
	_, err := imageFile.ReadAt(fSumLenBytes, fileSummaryLengthStart)
	if err != nil {
		if err != io.EOF {
			return fSummaryLength, err
		}
	}
	if err = binary.Read(bReader, binary.BigEndian, &fSummaryLength); err != nil {
		return fSummaryLength, err
	}
	return fSummaryLength, nil
}

func parseFileSummary(imageFile *os.File, fileLength int64, fSummaryLength int32) (map[string]*pb.FileSummary_Section, error) {
	var (
		sectionMap  = make(map[string]*pb.FileSummary_Section)
		fileSummary = &pb.FileSummary{}
	)
	// last 4 bytes says how many bytes should be read from end to get the FileSummary message
	fSummaryLength64 := int64(fSummaryLength)
	readAt := fileLength - fSummaryLength64 - FILE_SUM_BYTES

	fSummaryBytes := make([]byte, fSummaryLength)
	_, err := imageFile.ReadAt(fSummaryBytes, readAt)
	if err != nil {
		if err != io.EOF {
			return sectionMap, err
		}
	}

	_, c := binary.Uvarint(fSummaryBytes)
	if c <= 0 {
		return sectionMap, fmt.Errorf("buf too small(0) or overflows(-1): %d", c)
	}

	fSummaryBytes = fSummaryBytes[c:]
	if err = proto.Unmarshal(fSummaryBytes, fileSummary); err != nil {
		return sectionMap, err
	}

	for _, value := range fileSummary.GetSections() {
		sectionMap[value.GetName()] = value
	}
	return sectionMap, nil
}
