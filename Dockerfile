FROM golang:1.20 as build

WORKDIR /go/src/app

COPY . .

RUN go mod download && go mod tidy
ARG OS $OS
ARG ARCH $ARCH
ENV PLACE="auth-$(OS)-(ARCH)"
RUN go build -o /go/bin/auth/auth ./main/

FROM gcr.io/distroless/cc-debian11
COPY --from=build /go/bin/auth /

EXPOSE 8083

ENTRYPOINT ["/auth"]