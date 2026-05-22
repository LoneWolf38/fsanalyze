

build_proto:
	echo "Building protofiles"
	mkdir -p protofiles/src
	protoc --go_out=protofiles/src/ --go_opt=module=fsanalyze/protofiles/src protofiles/common/Security.proto
	protoc --go_out=protofiles/src/ --go_opt=module=fsanalyze/protofiles/src -Iprotofiles/common -Iprotofiles/ protofiles/hdfs/*.proto
	protoc --go_out=protofiles/src/ --go_opt=module=fsanalyze/protofiles/src -Iprotofiles/common -Iprotofiles/hdfs -Iprotofiles/ protofiles/hdfs/fsimage/fsimage.proto

build:
	echo "Building fsanalyze"
	env GOOS=linux GOARCH=amd64 go build -o bin/fsa .
	echo "uploading fsanalyze"
	scp bin/fsa root@sac01.acceldata.dvl:/data01/acceldata/fsa

build_expand:
	echo "Building fsa-expand (linux/amd64)"
	env GOOS=linux GOARCH=amd64 go build -o bin/fsa-expand ./cmd/expand/

build_expand_local:
	echo "Building fsa-expand (local)"
	go build -o bin/fsa-expand ./cmd/expand/