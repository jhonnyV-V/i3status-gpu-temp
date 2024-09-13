package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type Element struct {
	Name     string `json:"name"`
	Instance string `json:"instance"`
	Markup   string `json:"markup"`
	FullText string `json:"full_text"`
}

func isVersion(line string) bool {
	if len(line) < 10 {
		return false
	}
	if line[:10] == "{\"version\"" {
		return true
	}

	return false
}

func isEmpty(line string) bool {
	if len(line) == 0 {
		return true
	}

	if len(line) == 1 && line[0] == '[' {
		return true
	}

	return false
}

func main() {
	cmd := exec.Command("i3status")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating stdout pipe:", err)
		return
	}
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting i3status:", err)
		return
	}

	gpuTemps := Element{
		Name:     "cpu_temperature",
		Instance: "2",
		Markup:   "pango",
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		var elements []Element
		line := scanner.Text()

		if isVersion(line) {
			fmt.Printf("%s\n[[]\n", line)
			continue
		}

		if isEmpty(line) {
			continue
		}

		if line[0] == ',' {
			line = line[1:]
		}

		err := json.Unmarshal([]byte(line), &elements)
		if err != nil {
			fmt.Printf("Sometime Went Really Wrong \n")
			fmt.Println(line)
			fmt.Println(err)
			continue
		}

		nvidiaCmd := exec.Command("nvidia-smi", "--query-gpu=temperature.gpu", "--format=csv,noheader")
		nvidiaOut, err := nvidiaCmd.Output()
		if err != nil {
			fmt.Println("Error executing nvidia-smi:", err)
			return
		}
		temp := strings.TrimSpace(string(nvidiaOut))

		gpuTemps.FullText = fmt.Sprintf(
			"<span background='#d12f2c'> 󱎓 </span><span background='#bfbaac'> GPU %s °C  </span>",
			temp,
		)

		elements = append([]Element{gpuTemps}, elements...)

		formattedOutput, err := json.Marshal(elements)
		if err != nil {
			fmt.Printf("failed to marshall\n")
			fmt.Println(err)
			continue
		}

		elements = []Element{}

		fmt.Printf(",%s\n", formattedOutput)
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("Error waiting for i3status:", err)
	}
}
