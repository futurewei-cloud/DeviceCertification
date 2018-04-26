package switchdriver


import (
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"os"
)

var (
	// Use mcu pin 17, corresponds to physical pin 11 on the pi
	pin = rpio.Pin(17)
)

func ReadStatus() int64 {
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Unmap gpio memory when done
	defer rpio.Close()

	return int64(pin.Read())
}


func main() {
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Unmap gpio memory when done
	defer rpio.Close()

	fmt.Printf("Button state: %d\n", pin.Read())
}
