# City Traffic Simulator

## Build Requirements
- Go version 1.10 or above

## How to compile

```sh 
$ make
```

## How to run 

```sh 
$ ./city-traffic-sim -width <size of the map> -cars <number-of-cars> -semaphores <number-of-semaphores>
```
### Notes

- width must be multiple of 7 + 2. For example: 9, 16, 30, 44, 51, ...

- Cars must be less than the width