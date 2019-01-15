FROM golang:1.11.4
WORKDIR /go/src/github.com/Cidan/sheep
ADD . .

# Setup apt deps
RUN \
apt-get update && \
apt-get install -y unzip

# Install protoc
RUN \
wget https://github.com/protocolbuffers/protobuf/releases/download/v3.6.1/protoc-3.6.1-linux-x86_64.zip && \
unzip protoc-3.6.1-linux-x86_64.zip && \
mv include/* /usr/include/ && \
mv bin/* /usr/bin/ && \
chmod +x /usr/bin/protoc

# Install deps
RUN \
go get -v github.com/Masterminds/glide && \
cd $GOPATH/src/github.com/Masterminds/glide && \
git checkout 245caced2b16358b1c5e267691b17e9ee9952127 && \
go install && cd - && \
glide install && \
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway && \
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger && \
go install ./vendor/github.com/golang/protobuf/protoc-gen-go/ # silly

RUN make && make test