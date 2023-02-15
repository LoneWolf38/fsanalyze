

build_proto:
	echo "Building protofiles"
	mkdir -p protofiles/src
	protoc --go_out=protofiles/src/ --go_opt=module=fsanalyze/protofiles/src protofiles/common/Security.proto
	protoc --go_out=protofiles/src/ --go_opt=module=fsanalyze/protofiles/src -Iprotofiles/common -Iprotofiles/ protofiles/hdfs/*.proto
	protoc --go_out=protofiles/src/ --go_opt=module=fsanalyze/protofiles/src -Iprotofiles/common -Iprotofiles/hdfs -Iprotofiles/ protofiles/hdfs/fsimage/fsimage.proto