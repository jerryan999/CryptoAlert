# CryptoAlert

This is a demo program to watch realtime crypto price and when certain price alert is met, send the email to user.

## Start WebServer
```
    go run cmd/server/main.go
```

## Start Email Workers
```
    go run cmd/workers/main.go email
```

## Start Watch Workers
```
    go run cmd/workers/main.go watch
```
