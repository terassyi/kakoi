FROM golang:1.14

RUN apt update -y

COPY . /home/kakoi

RUN cd /home/kakoi && \
	/home/kakoi/setup.sh && \
	go mod download && \
	go build . && \
	mv kakoi /usr/local/bin/kakoi

CMD ['bash']
