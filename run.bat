@echo off
windres -o main-res.syso main.rc && go build -o bin/summation-trainer.exe -ldflags="-s -w -H windowsgui"
start bin/summation-trainer.exe