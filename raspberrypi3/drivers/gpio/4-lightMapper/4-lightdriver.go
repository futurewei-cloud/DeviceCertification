package lightdriver

import (
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"os"
)

//var (
//	// Use mcu pin 18, corresponds to physical pin 12 on the pi

//)

func main() {
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Unmap gpio memory when done
	defer rpio.Close()

	// Set pin to output mode
	pin1:=rpio.Pin(23)
	pin2:=rpio.Pin(24)
	pin3:=rpio.Pin(4)
	pin4:=rpio.Pin(22)

	pin1.Output()
	pin2.Output()
	pin3.Output()
	pin4.Output()

		//// Toggle pin 20 times
		//for x := 0; x < 20; x++ {
		////	pin.Toggle()
		//	time.Sleep(time.Second)
		//}
}

func TurnON(pinNo int) {
	// Open and map memory to access gpio, check for errors
	pin:=rpio.Pin(pinNo)
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Unmap gpio memory when done
	defer rpio.Close()

	// Set pin to output mode
	pin.Output()

	pin.High()
}


func TurnOff(pinNo int) {
	// Open and map memory to access gpio, check for errors
	pin:=rpio.Pin(pinNo)
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Unmap gpio memory when done
	defer rpio.Close()

	// Set pin to output mode
	pin.Output()
	pin.Low()

}

