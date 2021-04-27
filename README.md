# NumbersTrainer
This simple application for training numbers

# Windows
```cmd
windres -o main-res.syso main.rc && go build -o trainer.exe -ldflags="-s -w -H windowsgui"
```
