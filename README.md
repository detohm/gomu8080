# Gomu8080
An emulator for Intel 8080 processor written in Go.

This project passed all the CPU test suites below,
- cpudiag.bin
- tst8080.com
- 8080pre.com
- 8080exm.com

![Gomu8080](https://github.com/detohm/gomu8080/blob/main/docs/gomu8080.jpg?raw=true)

## How to run
You can activate debug mode using flag debug below.
Test mode:
```shell
go run example/main.go -path=[path to rom file] -debug=true
```
Space Invader mode:
```shell
go run example/main.go -path=[path to rom directory] -debug=false -spaceinvader=true
```

## Important Notes
Even though this project is passed all CPU diagnostic tests above, the Space Invader mode doesn't work as expected. There are some glitches in the animation logic. Therefore, PRs are welcome :)

+[Space Invaders](https://github.com/detohm/gomu8080/blob/main/docs/space-invaders.jpg?raw=true)