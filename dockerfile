FROM golang:1.19-alpine3.16

##buat folder APP
RUN mkdir /bengcall

##set direktori utama
WORKDIR /bengcall

##copy seluruh file
ADD . .

##buat executeable
RUN go build -o main .

##jalankan executeable
CMD ["./main"]