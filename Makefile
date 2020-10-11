.PHONY: protos

protos:
	protoc -I protos/ protos/dta.proto --go_out=plugins=grpc:protos

protos-gradle-parser:
	protoc -I protos/ protos/gradle.proto --go_out=plugins=grpc:protos