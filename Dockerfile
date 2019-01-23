# iron/go:dev is the alpine image with the go tools added
FROM iron/go:dev
WORKDIR /app
# Set an env var that matches your github repo name, replace treeder/dockergo here with your repo name
ENV SRC_DIR=/go/src/github.com/kimbbakar/rest-api/api-1
# Add the source code:
ADD . $SRC_DIR
RUN go get github.com/gorilla/mux 
RUN go get github.com/globalsign/mgo
RUN go get gopkg.in/mgo.v2/bson

# Build it:
RUN cd $SRC_DIR; go build -o api-1; cp api-1 /app/
#CMD ["./api-1", "-f=mongo"]
ENTRYPOINT ["./api-1"]
 