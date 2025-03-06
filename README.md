# shell-car-remote
## Disclaimer
All reverse engineering was taken from [Shell Motorsport Bluetooth Remote Control](https://gist.github.com/scrool/e79d6a4cb50c26499746f4fe473b3768). Awesome cheatsheet for this project.

# Controller mapping
### DualShock mapping
+ **Right analog**: directions right/left
+ **PS button**: Headlights 😎
+ **Left trigger**: move forward (press mor than 40% for turbo 🚀)
+ **Right trigger**: move backward (press mor than 40% for turbo 🚀)

# ⚠️ Developers zone start here⚠️
## implement other input?
in `input/input.go` you can find
```go 
type Input interface {
	Listen() *chan Command
}
```
to implement other input device you need to implement your input device interface.
this implementation will send into a channel all inputs to `chan Command`
where
```go
type Command struct {
	Key   KeyCommand
	Value ValueCommand
}
```
And use this `Command chan` as a parameter to build your car.
```go
xboxController := input.NewXBoxInput(controller)
ch := xboxController.Listen()
car, err := service.NewQCar(*cipher, BLEcar.Devices(), ch)
```

## Key features
+ Scan, pair, and control the Brandbase Bluetooth RC battery car.
+ Capability to control RC car with:
  + **DualShock 4**: using github.com/mrasband/ps4 (evdev under the hood... **LINUX ONLY**).
  + **Keyboard**: using golang-evdev (**LINUX ONLY**), this is only for debug developed as a POC (*DUE TO ITS CONTENT IT SHOULD NOT BE VIEWED BY ANYONE* 😂)

## Todo list
+ Implement keyboard again.
+ Read car battery status from characteristic. (feature well documented in cheatsheet 😏)