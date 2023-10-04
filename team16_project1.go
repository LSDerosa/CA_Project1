package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

/*********************************************************/
/*********************************************************/
/***********      MATH       FUNCTIONS    ****************/
/*********************************************************/
/*********************************************************/

/***********************2s COMPLEMENT**********************/

func twosComplement(binaryStr string) (int, error) {
	var negFlag bool //Check for leading one (negative number)
	if binaryStr[0:1] == "1" {
		negFlag = true
	}
	inverted := ""

	// Invert each bit individually
	for _, bit := range binaryStr {
		if bit == '0' {
			inverted += "1"
		} else {
			inverted += "0"
		}
	}

	// Trim leading zeros
	binaryStr = strings.TrimLeft(inverted, "0")

	result := ""

	// Carry for addition
	carry := 1

	// Work through string from right to left
	for i := len(binaryStr) - 1; i >= 0; i-- {
		bit := int(binaryStr[i] - '0')
		sum := bit + carry

		// Update result and carry
		result = strconv.Itoa(sum%2) + result
		carry = sum / 2
	}

	// Add leftover carry to left
	if carry == 1 {
		result = "1" + result
	}

	// Add zeroes to return to original length
	for len(result) < len(binaryStr) {
		result = "0" + result
	}
	result2, err := binaryToInteger(result)
	if err != nil {
		fmt.Println("Error:", err)
	}
	if negFlag == true {
		result2 = -result2
	}
	return result2, nil
}

/***************BINARY TO INTEGER CONVERTER*************/
func binaryToInteger(binary string) (int, error) {
	result, err := strconv.ParseInt(binary, 2, 64)
	if err != nil {
		return 0, err
	}
	return int(result), nil
}

/*********************************************************/
/*********************************************************/
/*******      INSTRUCTION-TYPE FUNCTIONS      ************/
/***************    ALPHABETICALLY    ********************/
/*********************************************************/

/*****************ADD FUNCTION*********************/
func addInstruction(binaryInstruction string, lineNumber int) {
	firstSource := binaryInstruction[11:16]
	// Reg One Int Conversion
	firstSourceint, err := binaryToInteger(firstSource)
	if err != nil {
		fmt.Println("Error:", err)
	}
	secondSource := binaryInstruction[22:27]
	// Reg Two Int Conversion
	secondSourceint, err := binaryToInteger(secondSource)
	if err != nil {
		fmt.Println("Error:", err)
	}
	destinationReg := binaryInstruction[27:32]
	// Reg Three Int Conversion
	destInt, err := binaryToInteger(destinationReg)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(binaryInstruction[0:11], "\t", firstSource, "\t", binaryInstruction[16:22], secondSource,
		"\t", destinationReg, lineNumber, " ADD \t", "R", destInt, "R", firstSourceint, "R", secondSourceint)
	// binaryInstruction = ""    Maybe not needed now

	instructionType := binaryInstruction[0:11]
	valueShamt := binaryInstruction[16:22]
	// ***WRITING TO FILE***
	file, err := os.OpenFile("team16_out_dis.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println("Error creating the file:", err)
		return
	}
	defer file.Close()
	// Write the text to the file
	output := fmt.Sprintf("%s %s %s %s %s %d ADD R%d R%d R%d \n",
		instructionType, firstSource, valueShamt, secondSource, destinationReg, lineNumber, destInt, firstSourceint, secondSourceint)
	_, err = file.WriteString(output)
	if err != nil {
		fmt.Println("Error writing to the file:", err)
		// file.Close()
		return
	}

}

/***************AND FUNCTION*******************/
func andInstruction(binaryInstruction string, lineNumber int) {
	firstSource := binaryInstruction[11:16]
	// Reg One Int Conversion
	firstSourceint, err := binaryToInteger(firstSource)
	if err != nil {
		fmt.Println("Error:", err)
	}
	secondSource := binaryInstruction[22:27]
	// Reg Two Int Conversion
	secondSourceint, err := binaryToInteger(secondSource)
	if err != nil {
		fmt.Println("Error:", err)
	}
	destinationReg := binaryInstruction[27:32]
	// Reg Three Int Conversion
	destInt, err := binaryToInteger(destinationReg)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(binaryInstruction[0:11], "\t", firstSource, "\t", binaryInstruction[16:22], secondSource,
		"\t", destinationReg, lineNumber, "AND \t", "R", destInt, "R", firstSourceint, "R", secondSourceint)
	// binaryInstruction = ""    Maybe not needed now
	instructionType := binaryInstruction[0:11]
	valueShamt := binaryInstruction[16:22]
	// ***WRITING TO FILE***
	file, err := os.OpenFile("team16_out_dis.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println("Error creating the file:", err)
		return
	}
	defer file.Close()
	// Write the text to the file
	output := fmt.Sprintf("%s %s %s %s %s %d AND R%d R%d R%d \n",
		instructionType, firstSource, valueShamt, secondSource, destinationReg, lineNumber, destInt, firstSourceint, secondSourceint)
	_, err = file.WriteString(output)
	if err != nil {
		fmt.Println("Error writing to the file:", err)
		return
	}

}

/********************BRANCH INSTRUCTION************************/
func branchInstruction(binaryInstruction string, lineNumber int) {
	instructionType := binaryInstruction[0:6]
	bOffset := binaryInstruction[6:32]
	bOffsetInt, err := twosComplement(bOffset)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(instructionType, bOffset, lineNumber, "B", "#", bOffsetInt)
	// ***WRITING TO FILE***
	file, err := os.OpenFile("team16_out_dis.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println("Error creating the file:", err)
		return
	}
	defer file.Close()
	// Write the text to the file
	output := fmt.Sprintf("%s %s %d B #%d \n",
		instructionType, bOffset, lineNumber, bOffsetInt)
	_, err = file.WriteString(output)
	if err != nil {
		fmt.Println("Error writing to the file:", err)
		return
	}
}

/*******************CONDITIONAL BRANCH NZERO**********************/
func conditionalBranchNz(binaryInstruction string, lineNumber int) {
	instructionType := binaryInstruction[0:8]
	branchOffset := binaryInstruction[8:27]
	bRegistry := binaryInstruction[27:32]
	branchOffsetInt, err := twosComplement(branchOffset)
	bRegistryInt, err := binaryToInteger(bRegistry)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(instructionType, branchOffset, bRegistry, lineNumber, "CBNZ", "R", bRegistryInt, "#", branchOffsetInt)
	// ***WRITING TO FILE***
	file, err := os.OpenFile("team16_out_dis.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println("Error creating the file:", err)
		return
	}
	defer file.Close()
	// Write the text to the file
	output := fmt.Sprintf("%s %s %s %d CBZ R%d #%d \n",
		instructionType, branchOffset, bRegistry, lineNumber, bRegistryInt, branchOffsetInt)
	_, err = file.WriteString(output)
	if err != nil {
		fmt.Println("Error writing to the file:", err)
		return
	}
}

/*******************CONDITIONAL BRANCH ZERO**********************/
func conditionalBranch(binaryInstruction string, lineNumber int) {
	instructionType := binaryInstruction[0:8]
	branchOffset := binaryInstruction[8:27]
	bRegistry := binaryInstruction[27:32]
	bRegistryInt, err := binaryToInteger(bRegistry)
	branchOffsetInt, err := twosComplement(branchOffset)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(instructionType, branchOffset, bRegistry, lineNumber, "CBZ", "R", bRegistryInt, "#", branchOffsetInt)
	// ***WRITING TO FILE***
	file, err := os.OpenFile("team16_out_dis.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println("Error creating the file:", err)
		return
	}
	defer file.Close()
	// Write the text to the file
	output := fmt.Sprintf("%s %s %s %d CBNZ R%d #%d \n",
		instructionType, branchOffset, bRegistry, lineNumber, bRegistryInt, branchOffsetInt)
	_, err = file.WriteString(output)
	if err != nil {
		fmt.Println("Error writing to the file:", err)
		return
	}
}

/*****************ORR FUNCTION*********************/
func orrInstruction(binaryInstruction string, lineNumber int) {
	firstSource := binaryInstruction[11:16]
	// Reg One Int Conversion
	firstSourceint, err := binaryToInteger(firstSource)
	if err != nil {
		fmt.Println("Error:", err)
	}
	secondSource := binaryInstruction[22:27]
	// Reg Two Int Conversion
	secondSourceint, err := binaryToInteger(secondSource)
	if err != nil {
		fmt.Println("Error:", err)
	}
	destinationReg := binaryInstruction[27:32]
	// Reg Three Int Conversion
	destInt, err := binaryToInteger(destinationReg)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(binaryInstruction[0:11], "\t", firstSource, "\t", binaryInstruction[16:22], secondSource,
		"\t", destinationReg, lineNumber, "ORR ", " R", destInt, "R", firstSourceint, "R", secondSourceint)
	// binaryInstruction = ""    Maybe not needed now
	instructionType := binaryInstruction[0:11]
	valueShamt := binaryInstruction[16:22]
	// ***WRITING TO FILE***
	file, err := os.OpenFile("team16_out_dis.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println("Error creating the file:", err)
		return
	}
	defer file.Close()
	// Write the text to the file
	output := fmt.Sprintf("%s %s %s %s %s %d ORR R%d R%d R%d \n",
		instructionType, firstSource, valueShamt, secondSource, destinationReg, lineNumber, destInt, firstSourceint, secondSourceint)
	_, err = file.WriteString(output)
	if err != nil {
		fmt.Println("Error writing to the file:", err)
		return
	}
}

/******************SUB FUNCTION****************/
func subInstruction(binaryInstruction string, lineNumber int) {
	firstSource := binaryInstruction[11:16]
	// Reg One Int Conversion
	firstSourceint, err := binaryToInteger(firstSource)
	if err != nil {
		fmt.Println("Error:", err)
	}
	secondSource := binaryInstruction[22:27]
	// Reg Two Int Conversion
	secondSourceint, err := binaryToInteger(secondSource)
	if err != nil {
		fmt.Println("Error:", err)
	}
	destinationReg := binaryInstruction[27:32]
	instructionType := binaryInstruction[0:11]
	valueShamt := binaryInstruction[16:22]
	// Reg Three Int Conversion
	destInt, err := binaryToInteger(destinationReg)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(instructionType, "\t", firstSource, "\t", valueShamt, secondSource,
		"\t", destinationReg, lineNumber, "SUB \t", "R", destInt, "R", firstSourceint, "R", secondSourceint)
	// binaryInstruction = ""    Maybe not needed now

	// ***WRITING TO FILE***
	file, err := os.OpenFile("team16_out_dis.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println("Error creating the file:", err)
		return
	}
	defer file.Close()
	// Write the text to the file
	output := fmt.Sprintf("%s %s %s %s %s %d SUB R%d R%d R%d \n",
		instructionType, firstSource, valueShamt, secondSource, destinationReg, lineNumber, destInt, firstSourceint, secondSourceint)
	_, err = file.WriteString(output)
	if err != nil {
		fmt.Println("Error writing to the file:", err)
		return
	}

}

/*******************************************************************/
/***********PRIMARY FUNCTION TO READ IN INSTRUCTIONS***************/
/******************************************************************/

func readAndProcessInstructions(binaryInstruction string, lineNumber int) {
	// Check if it matches any known instruction type
	switch binaryInstruction[:6] {
	case "000101":
		branchInstruction(binaryInstruction, lineNumber)
	default:
		switch binaryInstruction[:8] {
		case "10110100":
			conditionalBranch(binaryInstruction, lineNumber)
		case "10110101":
			conditionalBranchNz(binaryInstruction, lineNumber)
		default:
			switch binaryInstruction[:9] {
			case "110100101":
				println("MOVK")
			case "111100101":
				println("MOVK")
			default:
				switch binaryInstruction[:10] {
				case "1101000100":
					println("SUBI")
				case "1001000100":
					println("ADDI")
				default:
					switch binaryInstruction[:11] {
					case "10001011000":
						addInstruction(binaryInstruction, lineNumber)
					case "10001010000":
						andInstruction(binaryInstruction, lineNumber)
					case "10101010000":
						orrInstruction(binaryInstruction, lineNumber)
					case "11001011000":
						subInstruction(binaryInstruction, lineNumber)
					case "11111000000":
						println("STUR")
					case "11111000010":
						println("LDUR")
					case "11101010000":
						println("EOR")
					case "11010011100":
						println("ASR")
					case "11010011011":
						println("LSL")
					case "11010011010":
						println("LSR")
					default:
						calc, err := twosComplement(binaryInstruction)
						if err != nil {
							println(calc)
						}
						println(calc)
						//NEED BREAK CASE
					}
				}
			}
		}
	}
}

/************   MAIN   ****************/

func main() {
	// Open the input file
	inputFile, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening input file:", err)
		return
	}
	defer inputFile.Close()

	// Create a reader to read from the input file
	reader := bufio.NewReader(inputFile)
	lineNumber := 96
	for {
		// Create a new buffer for each loop
		buffer := make([]byte, 32)
		// Read 32 characters into new buffer
		var bytesRead int
		for bytesRead < 32 {
			char, err := reader.ReadByte()
			if err != nil {
				if err == io.EOF {
					// End of file reached, break the loop
					break
				} else {
					fmt.Println("Error reading from file:", err)
					return
				}
			}
			// Skip newline characters and empty spaces
			if char == '\n' {
				continue
			}
			if unicode.IsSpace(rune(char)) {
				continue
			}
			buffer[bytesRead] = char
			bytesRead++
		}
		if bytesRead < 32 {
			// End of file reached, break the loop
			break
		}
		// Convert  to string
		binaryInstruction := string(buffer[:bytesRead])
		readAndProcessInstructions(binaryInstruction, lineNumber)
		lineNumber += 4
	}

}
