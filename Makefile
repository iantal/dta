.PHONY: protos

protos:
	protoc -I protos/ protos/dta.proto --go_out=plugins=grpc:protos