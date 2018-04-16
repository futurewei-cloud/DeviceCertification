package wrtnodedriver

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"strconv"
	"time"
)

// at most 8 devices per edge box
var macsInt [8]uint64              // mac has 8 bytes = 64 bits
var macsString [8]string
var macToIndex map[uint64]int

var currentGPIOState [8]uint64      // status has 4 bytes
var latestGPIOInput [8]uint64       // input data has 4 x 8 = 32 bits
var currentDTH11State [8]uint64     // status has 4 bytes

func printCommand(cmd *exec.Cmd) {
  fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func hexArrayToInt64(items []string) (uint64, error) {
	var mac uint64
	mac = 0
	//fmt.Printf("length of items: %d\n", len(items))
	//fmt.Printf("items: %+v\n", items)
	cnt := 0
	for _, str := range items {
		fmt.Printf("item[%d]: %+v\n", cnt, str)
		cur, err := strconv.ParseUint("0x" + str, 0, 8)
		if err != nil {
			return 0, err
		}
		mac = mac << 8 + cur
		cnt = cnt + 1
	}
	if cnt < 4 || (cnt > 4 && cnt != 8) {
		return 0, fmt.Errorf("Input string is malformed")
	}

	return mac, nil
}

func convertToCmdArray(newCmdArgsInt uint64) [8]string {
	var cmds [8]string
	var cur uint8

	cur = uint8(newCmdArgsInt & 0xFF)
	cmds[6] = fmt.Sprintf("%02X", cur)
	cmds[7] = fmt.Sprintf("%02X", ^cur)
	cur = uint8((newCmdArgsInt >> 8) & 0xFF)
	cmds[4] = fmt.Sprintf("%02X", cur)
	cmds[5] = fmt.Sprintf("%02X", ^cur)
	cur = uint8((newCmdArgsInt >> 16) & 0xFF)
	cmds[2] = fmt.Sprintf("%02X", cur)
	cmds[3] = fmt.Sprintf("%02X", ^cur)
	cur = uint8((newCmdArgsInt >> 24) & 0xFF)
	cmds[0] = fmt.Sprintf("%02X", cur)
	cmds[1] = fmt.Sprintf("%02X", ^cur)
	return cmds
}

func InitDevice() error {
	cmd := exec.Command("/bin/send", "57", "3E", "02", "01", "00", "00")
	printCommand(cmd)
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
		return err
	}
	out := strings.Split(cmdOutput.String(), "\n")
        ret_prefix := "buf: "
        startIndex := strings.Index(out[3], ret_prefix)
	if startIndex < 0 {
		fmt.Printf("==> Output: %v\n", out[3])
		return err
	}

	out[3] = out[3] + " 00"
	items := strings.Split(out[3], " ")
	l := len(items)
	fmt.Printf("==> len, items: %d\n%+v\n", l, items)

	macToIndex = make(map[uint64]int)
	low := 8
	high := 16
	i := 0
	for {
		if low >= l {
			break;
		}
		fmt.Printf("==> low, high: %d\n%d\n", low, high)
		macsString[i] = strings.Join(items[low:high:high], " ")
		macsInt[i], _ = hexArrayToInt64(items[low:high:high])
		macToIndex[macsInt[i]] = i
		latestGPIOInput[i] = 0
		low = high + 1
		high = low + 8
		i = i + 1
	}
	fmt.Printf("==> MAC string: %+v\n", macsString)
	fmt.Printf("==> MAC int: %+v\n", macsInt)
	fmt.Printf("==> latestGPIOInput: %+v\n", latestGPIOInput)
	return nil
}

func SetGPIO(deviceIndex int, cmdArgs string, onOff int) error {
	fmt.Printf("====> index: %d\n", deviceIndex)
	fmt.Printf("====> m[0]: %+v\n", macsString[0])
	fmt.Printf("====> m[1]: %+v\n", macsString[1])
	fmt.Printf("====> m[index]: %+v\n", macsString[deviceIndex])
	// first, check against latest input
	cmdArgsInt, err := hexArrayToInt64(strings.Split(cmdArgs, " "))
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
		return err
	}
	if onOff == 1 {
		latestGPIOInput[deviceIndex] = latestGPIOInput[deviceIndex] | cmdArgsInt
	} else {
		latestGPIOInput[deviceIndex] = latestGPIOInput[deviceIndex] &^ cmdArgsInt
	}
	
	c := convertToCmdArray(latestGPIOInput[deviceIndex])

	// different input, set
	m := strings.Split(macsString[deviceIndex], " ")
	fmt.Printf("==> m: %+v\n", m)
	fmt.Printf("==> c: %+v\n", c)
	cmd := exec.Command("/bin/send", "57", "3E", "08", "0D", "10", "00", m[0], m[1], m[2], m[3], m[4], m[5], m[6], m[7], c[0], c[1], c[2], c[3], c[4], c[5], c[6], c[7])
	printCommand(cmd)

	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err = cmd.Run()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
		return err
	}

	out := strings.Split(cmdOutput.String(), "\n")
	fmt.Printf("==> Output: %v\n%v\n%v\n%v\n", out[0], out[1], out[2], out[3])

	time.Sleep(time.Second)

	// update currentGPIOState
	err = readGPIO(deviceIndex)
	if err != nil {
		return err
	}

	return nil
}

func readGPIO(deviceIndex int) error {
	var out []string
	var count int
	count = 2

	for count > 1 {
		time.Sleep(time.Second)
		cmd := exec.Command("/bin/send", "57", "3E", "83", "01", "00", "00")
		printCommand(cmd)

		cmdOutput := &bytes.Buffer{}
		cmd.Stdout = cmdOutput
		err := cmd.Run()
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
			return err
		}

		out = strings.Split(cmdOutput.String(), "\n")

		out2 := strings.Split(out[2], " ");
		if out2[1] != "19" {
			os.Stderr.WriteString(fmt.Sprintf("==> Error: return state string: %s\n", out[3]))
			return fmt.Errorf("response mismatched")
		}

		ret_prefix := "buf: "
		startIndex := strings.Index(out[3], ret_prefix)
		if startIndex != 0 {
			os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", out[3]))
			return fmt.Errorf("response mismatched")
		}
		out[3] = out[3] + " 00"
		out3 := strings.Split(out[3][5:], " ")
		cnt, err := strconv.ParseInt("0x" + out3[6], 0, 64)
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
			return err
		}
		fmt.Printf("==> Count: %v\n%d\n", out3, cnt)
		count = int(cnt)
		if count == 0 {
			break
		}

		// get mac address from response
		curMac, err := hexArrayToInt64(out3[7:15])
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
			return err
		}
		state, err := hexArrayToInt64(out3[15:19])
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
			return err
		}
		fmt.Printf(" ==> curMac, macToIndex, state: %08X, %d, %08X\n", curMac, macToIndex[curMac], state)
		currentGPIOState[macToIndex[curMac]] = state
	}

	return nil
}

func main() {
	err1 := InitDevice()
	if err1 != nil {
		return
	}
	fmt.Printf(" mac: %+v\n", macsString)

	// turn on motor
	SetGPIO(0, "00 10 00 00", 1)
	// turn off motor
	SetGPIO(0, "00 10 00 00", 0)

	fmt.Printf("==> State: %08X\n", currentGPIOState[0])

	// turn on motor
	SetGPIO(1, "00 10 00 00", 1)
	// turn off motor
	SetGPIO(1, "00 10 00 00", 0)

	fmt.Printf("==> State: %08X\n", currentGPIOState[1])
}

func ReadGPIO(deviceIndex int) uint64 {
	return currentGPIOState[deviceIndex]
}
