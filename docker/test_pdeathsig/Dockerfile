FROM golang:latest

COPY . /robust/

RUN /bin/sh -c "cd /robust/subprocess/wrapper && GOPROXY=https://goproxy.io go build -v"
RUN ln -s /robust/docker/test_pdeathsig/_test.sh /
WORKDIR /

CMD ["/bin/bash"]
