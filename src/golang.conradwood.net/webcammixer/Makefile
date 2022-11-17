.PHONY:	client server
all:	client server
client:
	cd client && go install ${LDFLAGS} webcammixer-client.go `ls -1 *.go|grep -v webcammixer-client.go`
server:
	cd server && go install ${LDFLAGS} webcammixer-server.go `ls -1 *.go|grep -v webcammixer-server.go`
