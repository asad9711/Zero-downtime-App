FROM golang:1.16-alpine as builder


#RUN apk add --update --upgrade --no-cache  git curl



#ENV GOPATH /opt
#ADD . $GOPATH/src
#COPY . /opt/src
ADD . /opt/src

WORKDIR /opt/src
RUN go mod download

RUN GOOS=linux go build -o listener .



# use scratch (base for a docker image)
FROM alpine:latest
#FROM ubuntu:18.04
# set working directory
WORKDIR /root
# copy the binary from builder
COPY --from=builder /opt/src/listener .

# run the binary
CMD ["./listener"]
