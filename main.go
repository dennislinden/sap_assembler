package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

const (
	NOP = "0000"
	LDA = "0001"
	ADD = "0010"
	SUB = "0011"
	STA = "0100"
	LDI = "0101"
	JMP = "0110"
	JC  = "0111"
	JZ  = "1000"
	OUT = "1110"
	HLT = "1111"
)

var opcodes = []string{"NOP", "LDA", "ADD", "SUB", "STA", "LDI", "JMP", "JC", "JZ", "OUT", "HLT"}

var flagFilevar string

func init() {
	flag.StringVar(&flagFilevar, "i", "", "input file")
}

func main() {
	flag.Parse()

	scanner, file, err := readAsmFile(flagFilevar)
	if err != nil {
		log.Fatal("cannot read file", err)
	}

	finalOutput, err := convertOpcodes(scanner)
	if err != nil {
		log.Fatal("cannot convert opcodes: ", err)
	}

	defer file.Close()

	if err := writeAsmFile(finalOutput); err != nil {
		log.Fatal("cannot write file", err)
	}

}

func writeAsmFile(file string) error {
	os.WriteFile("out.bin", []byte(file), 0644)
	return nil
}

func readAsmFile(filename string) (fscanner *bufio.Scanner, file *os.File, err error) {
	asmFile, err := os.Open(filename)
	if err != nil {
		fmt.Println("cannot read file")
	}

	scanner := bufio.NewScanner(asmFile)
	return scanner, file, nil
}

func convertOpcodes(scanner *bufio.Scanner) (opcodes string, err error) {

	finalArray := []string{}
	finalOutput := ""

	// convert mnemonics to machinecode
	for scanner.Scan() {
		scanLine := scanner.Text()

		// basic check for invalid characters
		if err := sanityCheck(scanLine); err != nil {
			log.Fatal("Invalid char", err)
		}

		scanLine = cleanComments(scanLine)

		// remove multispaces
		space := regexp.MustCompile(`\s+`)
		scanLine = space.ReplaceAllString(scanLine, " ")

		lineArray := strings.Split(scanLine, " ")

		if len(lineArray) <= 2 {
			for key, line := range lineArray {

				// check for empty lines
				if line == "" {
					// delete item
					lineArray = append(lineArray[:key], lineArray[key+1:]...)
					continue
				}

				// get all the first words and parse them to machinecodes
				if key%2 == 0 {
					if line != "" {
						mCode, err := parse(line)
						if err != nil {
							log.Fatal(err)
						}
						lineArray[key] = mCode
					}
				}

			}
		} else if strings.Contains(scanLine, "//") == true {
			// remove inline comments
			for key, line := range lineArray {
				if line == "//" {
					lineArray = append(lineArray[:key])

				}
			}
			for key, line := range lineArray {

				if key%2 == 0 {
					if line != "" {
						mCode, err := parse(line)
						if err != nil {
							log.Fatal(err)
						}
						lineArray[key] = mCode
					}
				}
			}
		} else if strings.Contains(scanLine, "#") {
			fmt.Println(scanLine) //CONTINUE HERE
			for _, line := range lineArray {
				val := convertDecBin(line)
				fmt.Println(val)
			}
		} else {
			return "", fmt.Errorf("assembly error to many words in line: %v", lineArray)
		}

		if len(lineArray) != 0 {
			finalArray = append(finalArray, strings.Join(lineArray, ""))
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	finalOutput = strings.Join(finalArray, "\n")
	return finalOutput, nil
}

func cleanComments(line string) string {
	if strings.HasPrefix(line, "//") == true {
		return ""
	}
	return line
}

func convertDecBin(dec string) string {
	// remove #
	dec = dec[1:]
	bin := fmt.Sprintf("%b", dec)
	return bin
}
func sanityCheck(line string) error {

	valid := false
	validchars := []string{"%", "#", "//"}
	validPrefix := append(validchars, opcodes...)

	for _, prefix := range validPrefix {
		if strings.HasPrefix(line, prefix) == true {
			valid = true
			break
		}
	}
	if valid == false {
		return fmt.Errorf("Unvalid char in line:%v", line)
	}
	return nil
}

func parse(opcode string) (string, error) {
	switch opcode {
	case "NOP":
		return NOP, nil
	case "LDA":
		return LDA, nil
	case "ADD":
		return ADD, nil
	case "SUB":
		return SUB, nil
	case "STA":
		return STA, nil
	case "LDI":
		return LDI, nil
	case "JMP":
		return JMP, nil
	case "JC":
		return JC, nil
	case "JZ":
		return JZ, nil
	case "OUT":
		return OUT, nil
	case "HLT":
		return HLT, nil
	default:
		return "", fmt.Errorf("Cannot use opcode: %v", opcode)

	}
}
