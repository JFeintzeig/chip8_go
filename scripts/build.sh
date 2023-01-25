#!/bin/bash
if [ ! -f app ]; then
  rm app
fi

go build cmd/app/app.go
go build cmd/disassemble/disassemble.go
