#!/bin/bash
if [ ! -f chip8 ]; then
  rm chip8
fi

go build cmd/chip8/chip8.go
