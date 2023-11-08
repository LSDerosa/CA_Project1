/*******  TO DO    *******/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

/***********************************************************************/
/*****************************B TYPE INSTRUCTIONS***********************/
/***********************************************************************/

type BranchInstruction struct {
	Opcode     string
	OffsetStr  string
	OffsetInt  int
	LineNumber int
}

func NewBranchInstruction(opcode string, offsetStr string, offsetInt int, lineNumber int) *BranchInstruction {
	return &BranchInstruction{
		Opcode:     opcode,
		OffsetStr:  offsetStr,
		OffsetInt:  offsetInt,
		LineNumber: lineNumber,
	}
}

func parseBranchInstruction(binaryInstruction string, lineNumber int) *BranchInstruction {
	opcode := binaryInstruction[:6]
	offsetStr := binaryInstruction[6:32]
	offsetInt, _ := twosComplement(offsetStr)
	return NewBranchInstruction(opcode, offsetStr, offsetInt, lineNumber)
}

/***********************************************************************/
/*****************************CB TYPE INSTRUCTIONS***********************/
/***********************************************************************/

type cBranchInstruction struct {
	Opcode      string
	OffsetStr   string
	RegistryStr string
	RegistryInt int
	OffsetInt   int
	LineNumber  int
}

func NewcBranchInstruction(opcode string, offsetStr string, registryStr string,
	registryInt int, offsetInt int, lineNumber int) *cBranchInstruction {
	return &cBranchInstruction{
		Opcode:      opcode,
		OffsetStr:   offsetStr,
		RegistryStr: registryStr,
		RegistryInt: registryInt,
		OffsetInt:   offsetInt,
		LineNumber:  lineNumber,
	}
}

func parsecBranchInstruction(binaryInstruction string, lineNumber int) *cBranchInstruction {
	opcode := binaryInstruction[:8]
	offsetStr := binaryInstruction[8:27]
	registryStr := binaryInstruction[27:32]
	registryInt, _ := strconv.ParseInt(registryStr, 2, 32)
	offsetInt, _ := twosComplement(offsetStr)
	return NewcBranchInstruction(opcode, offsetStr, registryStr, int(registryInt), offsetInt, lineNumber)
}

/***********************************************************************/
/*****************************IM TYPE INSTRUCTIONS***********************/
/***********************************************************************/
type imTypeInstruction struct {
	Opcode      string
	FieldValue  string
	DestStr     string
	Op2         string
	Destination int
	LineNumber  int
	FieldInt    uint
	BitPattern  int
}

func NewImTypeInstruction(opcode string, op2 string, destStr string, fieldValue string,
	bitPattern int, dest int, fieldInt uint, lineNumber int) *imTypeInstruction {
	return &imTypeInstruction{
		Opcode:      opcode,
		Op2:         op2,
		Destination: dest,
		FieldInt:    fieldInt,
		LineNumber:  lineNumber,
		FieldValue:  fieldValue,
		DestStr:     destStr,
		BitPattern:  bitPattern,
	}
}

func parseImTypeInstruction(binaryInstruction string, lineNumber int) *imTypeInstruction {
	opcode := binaryInstruction[:9]
	op2 := binaryInstruction[9:11]
	fieldValue := binaryInstruction[11:27]
	destStr := binaryInstruction[27:32]
	dest, _ := strconv.ParseInt(destStr, 2, 32)
	fieldInt, _ := strconv.ParseInt(fieldValue, 2, 32)
	bitPattern := 0
	if op2 == "00" {
		bitPattern = 0
	} else if op2 == "01" {
		bitPattern = 16
	} else if op2 == "10" {
		bitPattern = 32
	} else if op2 == "11" {
		bitPattern = 48
	}
	return NewImTypeInstruction(opcode, op2, destStr, fieldValue, bitPattern, int(dest), uint(fieldInt), lineNumber)
}

/***********************************************************************/
/*****************************I TYPE INSTRUCTIONS***********************/
/***********************************************************************/
type iTypeInstruction struct {
	Opcode      string
	Immediate   string
	Src1Str     string
	DestStr     string
	Destination int
	Source1     int
	LineNumber  int
	ImmedInt    int
}

func NewITypeInstruction(opcode string, immediate string,
	src1Str string, destStr string, dest int, src1 int, immedInt int, lineNumber int) *iTypeInstruction {
	return &iTypeInstruction{
		Opcode:      opcode,
		Destination: dest,
		Source1:     src1,
		ImmedInt:    immedInt,
		LineNumber:  lineNumber,
		Immediate:   immediate,
		Src1Str:     src1Str,
		DestStr:     destStr,
	}
}

func parseITypeInstruction(binaryInstruction string, lineNumber int) *iTypeInstruction {
	opcode := binaryInstruction[:10]
	immediate := binaryInstruction[10:22]
	src1Str := binaryInstruction[22:27]
	destStr := binaryInstruction[27:32]
	dest, _ := strconv.ParseInt(destStr, 2, 32)
	src1, _ := strconv.ParseInt(src1Str, 2, 32)
	immedInt, _ := strconv.ParseInt(immediate, 2, 32)
	return NewITypeInstruction(opcode, immediate, src1Str, destStr, int(dest), int(src1), int(immedInt), lineNumber)
}

/***********************************************************************/
/*****************************R TYPE INSTRUCTIONS***********************/
/***********************************************************************/
type rTypeInstruction struct {
	Opcode      string
	Destination int
	Source1     int
	Source2     int
	LineNumber  int
	ShamtInt    int
	Shamt       string
	Src2Str     string
	Src1Str     string
	DestStr     string
}

func NewrTypeInstruction(opcode string, shamt string, src2Str string,
	src1Str string, destStr string, dest int, src1 int, src2 int, shamtInt int, lineNumber int) *rTypeInstruction {
	return &rTypeInstruction{
		Opcode:      opcode,
		Destination: dest,
		Source1:     src1,
		Source2:     src2,
		ShamtInt:    shamtInt,
		LineNumber:  lineNumber,
		Shamt:       shamt,
		Src2Str:     src2Str,
		Src1Str:     src1Str,
		DestStr:     destStr,
	}
}

func parserTypeInstruction(binaryInstruction string, lineNumber int) *rTypeInstruction {
	opcode := binaryInstruction[:11]
	src2Str := binaryInstruction[11:16]
	shamt := binaryInstruction[16:22]
	src1Str := binaryInstruction[22:27]
	destStr := binaryInstruction[27:32]
	dest, _ := strconv.ParseInt(destStr, 2, 32)
	src1, _ := strconv.ParseInt(src1Str, 2, 32)
	src2, _ := strconv.ParseInt(src2Str, 2, 32)
	shamtInt, _ := strconv.ParseInt(shamt, 2, 32)
	return NewrTypeInstruction(opcode, shamt, src2Str, src1Str, destStr, int(dest), int(src1), int(src2), int(shamtInt), lineNumber)
}

/*******************************************************************/
/**************************D TYPE INSTRUCTIONS***********************/
/*******************************************************************/
type dTypeInstruction struct {
	Opcode     string
	AddressStr string
	Op2        string
	Src2Str    string
	Src1Str    string
	AddressInt int
	Source1    int
	Source2    int
	LineNumber int
}

func NewDTypeInstruction(opcode string, addressStr string, op2 string,
	src1Str string, src2Str string, addressInt int, src1 int, src2 int, lineNumber int) *dTypeInstruction {
	return &dTypeInstruction{
		Opcode:     opcode,
		AddressStr: addressStr,
		Source1:    src1,
		Source2:    src2,
		LineNumber: lineNumber,
		Op2:        op2,
		Src2Str:    src2Str,
		Src1Str:    src1Str,
		AddressInt: addressInt,
	}
}

func parseDTypeInstruction(binaryInstruction string, lineNumber int) *dTypeInstruction {
	opcode := binaryInstruction[0:11]
	addressStr := binaryInstruction[11:20]
	op2 := binaryInstruction[20:22]
	src2Str := binaryInstruction[22:27]
	src1Str := binaryInstruction[27:32]
	addressInt, _ := strconv.ParseInt(addressStr, 2, 32)
	src1, _ := strconv.ParseInt(src1Str, 2, 32)
	src2, _ := strconv.ParseInt(src2Str, 2, 32)
	return NewDTypeInstruction(opcode, addressStr, op2, src1Str, src2Str, int(addressInt), int(src1), int(src2), lineNumber)
}

/*********************************************************/
/*********************************************************/
/***********      MAIN DRIVER FUNCTION   ****************/
/*********************************************************/
/*********************************************************/

func readAndProcessInstructions(binaryInstruction string, lineNumber int, outputFile *string, output *string) {

	/************B TYPE INSTRUCTIONS****************/
	switch binaryInstruction[:6] {
	case "000101":
		branchInst := parseBranchInstruction(binaryInstruction, lineNumber)
		*outputFile += fmt.Sprintf("%.6s %.26s\t\t%.1d\tB\t#%d\n",
			branchInst.Opcode, branchInst.OffsetStr, lineNumber, branchInst.OffsetInt)
		//fmt.Printf(output)
		//writeToFile(*outputFile, output)
	default:
		/************CB TYPE INSTRUCTIONS****************/

		switch binaryInstruction[:8] {
		case "10110100":
			cbTypeinst := parsecBranchInstruction(binaryInstruction, lineNumber)
			*output += fmt.Sprintf("%.8s %.19s %.5s\t\t%.1d\tCBZ\tR%.1d, #%.1d\n",
				cbTypeinst.Opcode, cbTypeinst.OffsetStr, cbTypeinst.RegistryStr, lineNumber, cbTypeinst.RegistryInt,
				cbTypeinst.OffsetInt)
			//fmt.Printf(output)
			//writeToFile(*outputFile, output)

		case "10110101":
			cbTypeinst := parsecBranchInstruction(binaryInstruction, lineNumber)
			*output += fmt.Sprintf("%.8s %.19s %.5s\t%.1d\tCBNZ\tR%.1d, #%.1d\n",
				cbTypeinst.Opcode, cbTypeinst.OffsetStr, cbTypeinst.RegistryStr, lineNumber, cbTypeinst.RegistryInt,
				cbTypeinst.OffsetInt)
			//fmt.Printf(output)
			//writeToFile(*outputFile, output)
		default:
			/************IM TYPE INSTRUCTIONS****************/
			switch binaryInstruction[:9] {
			case "110100101":
				imTypeinst := parseImTypeInstruction(binaryInstruction, lineNumber)
				*output += fmt.Sprintf("%.9s %.2s %.16s %.5s\t%.1d\tMOVZ\tR%.1d, %.1d, LSL %.1d\n",
					imTypeinst.Opcode, imTypeinst.Op2, imTypeinst.FieldValue, imTypeinst.DestStr, lineNumber,
					imTypeinst.Destination, imTypeinst.FieldInt, imTypeinst.BitPattern)
				//fmt.Printf(output)
				//writeToFile(*outputFile, output)
			case "111100101":
				imTypeinst := parseImTypeInstruction(binaryInstruction, lineNumber)
				*output += fmt.Sprintf("%.9s %.2s %.16s %.5s\t%.1d\tMOVK\tR%.1d, %.1d, LSL %.1d\n",
					imTypeinst.Opcode, imTypeinst.Op2, imTypeinst.FieldValue, imTypeinst.DestStr, lineNumber,
					imTypeinst.Destination, imTypeinst.FieldInt, imTypeinst.BitPattern)
				//fmt.Printf(output)
				//writeToFile(*outputFile, output)
			default:
				/***********I TYPE INSTRUCTIONS******************/

				/***********SUBI TYPE INSTRUCTIONS******************/
				switch binaryInstruction[:10] {
				case "1101000100":
					iTypeInst := parseITypeInstruction(binaryInstruction, lineNumber)
					*output += fmt.Sprintf("%.10s %.12s %.5s %.5s\t\t%.1d\tSUBI\tR%.1d, R%.1d, #%.1d\n",
						iTypeInst.Opcode, iTypeInst.Immediate, iTypeInst.Src1Str, iTypeInst.DestStr, lineNumber,
						iTypeInst.Destination, iTypeInst.Source1, iTypeInst.ImmedInt)
					//fmt.Printf(output)
					//writeToFile(*outputFile, output)
					/***********ADDI TYPE INSTRUCTIONS******************/
				case "1001000100":
					iTypeInst := parseITypeInstruction(binaryInstruction, lineNumber)
					*output += fmt.Sprintf("%.10s %.12s %.5s %.5s \t%.1d ADDI R%.1d, R%.1d, #%.1d \n",
						iTypeInst.Opcode, iTypeInst.Immediate, iTypeInst.Src1Str, iTypeInst.DestStr, lineNumber,
						iTypeInst.Source1, iTypeInst.Destination, iTypeInst.ImmedInt)
					//fmt.Printf(output)
					//writeToFile(*outputFile, output)
				default:
					/*************D TYPE INSTRUCTIONS****************/
					switch binaryInstruction[:11] {
					//*******STUR******//
					case "11111000000":
						dTypeInst := parseDTypeInstruction(binaryInstruction, lineNumber)
						*output += fmt.Sprintf("%.11s %.9s %.2s %.5s %.5s\t%.1d\tSTUR\tR%.1d, [R%.1d, #%.1d]\n",
							dTypeInst.Opcode, dTypeInst.AddressStr, dTypeInst.Op2, dTypeInst.Src2Str,
							dTypeInst.Src1Str, lineNumber, dTypeInst.Source1, dTypeInst.Source2, dTypeInst.AddressInt)
						//fmt.Printf(output)
						//writeToFile(*outputFile, output)
						//*******LDUR******//
					case "11111000010": //
						dTypeInst := parseDTypeInstruction(binaryInstruction, lineNumber)
						*output += fmt.Sprintf("%.11s %.9s %.2s %.5s %.5s\t%.1d\tLDUR\tR%.1d, [R%.1d, #%.1d]\n",
							dTypeInst.Opcode, dTypeInst.AddressStr, dTypeInst.Op2, dTypeInst.Src2Str,
							dTypeInst.Src1Str, lineNumber, dTypeInst.Source1, dTypeInst.Source2, dTypeInst.AddressInt)
						//fmt.Printf(output)
						//writeToFile(*outputFile, output)
						/*******************RTYPE INSTRUCTIONS******************/
					default:
						switch binaryInstruction[:11] {
						//*******ADD******//
						case "10001011000":
							rTypeInst := parserTypeInstruction(binaryInstruction, lineNumber)
							*output += fmt.Sprintf("%.11s %.5s %.6s %.5s %.5s\t%.1d\tADD\tR%.1d, R%.1d, R%.1d\n",
								rTypeInst.Opcode, rTypeInst.Src2Str, rTypeInst.Shamt, rTypeInst.Src1Str, rTypeInst.DestStr, lineNumber,
								rTypeInst.Destination, rTypeInst.Source1, rTypeInst.Source2)
							//fmt.Printf(output)
							//writeToFile(*outputFile, output)
							//*******AND******//
						case "10001010000":
							rTypeInst := parserTypeInstruction(binaryInstruction, lineNumber)
							*output += fmt.Sprintf("%.11s %.5s %.6s %.5s %.5s\t%.1d\tAND\tR%.1d, R%.1d, R%.1d\n",
								rTypeInst.Opcode, rTypeInst.Src2Str, rTypeInst.Shamt, rTypeInst.Src1Str, rTypeInst.DestStr, lineNumber,
								rTypeInst.Destination, rTypeInst.Source1, rTypeInst.Source2)
							//fmt.Printf(output)
							//writeToFile(*outputFile, output)
							//*******ORR******//
						case "10101010000":
							rTypeInst := parserTypeInstruction(binaryInstruction, lineNumber)
							*output += fmt.Sprintf("%.11s %.5s %.6s %.5s %.5s\t%.1d\tORR\tR%.1d, R%.1d, R%.1d\n",
								rTypeInst.Opcode, rTypeInst.Src2Str, rTypeInst.Shamt, rTypeInst.Src1Str, rTypeInst.DestStr, lineNumber,
								rTypeInst.Destination, rTypeInst.Source1, rTypeInst.Source2)
							//fmt.Printf(output)
							//writeToFile(*outputFile, output)
							//****************SUB***************//
						case "11001011000":
							rTypeInst := parserTypeInstruction(binaryInstruction, lineNumber)
							*output += fmt.Sprintf("%.11s %.5s %.6s %.5s %.5s\t%.1d\tSUB\tR%.1d, R%.1d, R%.1d\n",
								rTypeInst.Opcode, rTypeInst.Src2Str, rTypeInst.Shamt, rTypeInst.Src1Str, rTypeInst.DestStr, lineNumber,
								rTypeInst.Destination, rTypeInst.Source1, rTypeInst.Source2)
							//fmt.Printf(output)
							//writeToFile(*outputFile, output)
							//***************EOR****************//
						case "11101010000":
							rTypeInst := parserTypeInstruction(binaryInstruction, lineNumber)
							*output += fmt.Sprintf("%.11s %.5s %.6s %.5s %.5s\t%.1d\tEOR\tR%.1d, R%.1d, R%.1d\n",
								rTypeInst.Opcode, rTypeInst.Src2Str, rTypeInst.Shamt, rTypeInst.Src1Str, rTypeInst.DestStr, lineNumber,
								rTypeInst.Destination, rTypeInst.Source1, rTypeInst.Source2)
							//fmt.Printf(output)
							//writeToFile(*outputFile, output)
							//***********ASR*****************//
						case "11010011100":
							rTypeInst := parserTypeInstruction(binaryInstruction, lineNumber)
							*output += fmt.Sprintf("%.11s %.5s %.6s %.5s %.5s\t%.1d\tASR\tR%.1d, R%.1d, #%.1d\n",
								rTypeInst.Opcode, rTypeInst.Src2Str, rTypeInst.Shamt, rTypeInst.Src1Str, rTypeInst.DestStr, lineNumber,
								rTypeInst.Destination, rTypeInst.Source1, rTypeInst.ShamtInt)
							//fmt.Printf(output)
							//writeToFile(*outputFile, output)
							//**************LSL******************//
						case "11010011011":
							rTypeInst := parserTypeInstruction(binaryInstruction, lineNumber)
							*output += fmt.Sprintf("%.11s %.5s %.6s %.5s %.5s\t%.1d\tLSL\tR%.1d, R%.1d, #%.1d\n",
								rTypeInst.Opcode, rTypeInst.Src2Str, rTypeInst.Shamt, rTypeInst.Src1Str, rTypeInst.DestStr, lineNumber,
								rTypeInst.Destination, rTypeInst.Source1, rTypeInst.ShamtInt)
							//fmt.Printf(output)
							//writeToFile(*outputFile, output)
							//**************LSR*******************//
						case "11010011010":
							rTypeInst := parserTypeInstruction(binaryInstruction, lineNumber)
							*output += fmt.Sprintf("%.11s %.5s %.6s %.5s %.5s\t%.1d\tLSR\tR%.1d, R%.1d, #%.1d\n",
								rTypeInst.Opcode, rTypeInst.Src2Str, rTypeInst.Shamt, rTypeInst.Src1Str, rTypeInst.DestStr, lineNumber,
								rTypeInst.Destination, rTypeInst.Source1, rTypeInst.ShamtInt)
							//fmt.Printf(output)
							//writeToFile(*outputFile, output)
						default:
							switch binaryInstruction[:32] {
							case "11111110110111101111111111100111":
								*output += fmt.Sprintf("%.8s %.3s %.5s %.5s %.5s %.6s\t%.1d\tBREAK\n",
									binaryInstruction[:8], binaryInstruction[8:11], binaryInstruction[11:16], binaryInstruction[16:21],
									binaryInstruction[21:26], binaryInstruction[26:32], lineNumber)
								//fmt.Printf(output)
								//writeToFile(*outputFile, output)
							case "00000000000000000000000000000000":
								*output += fmt.Sprintf("%.8s %.3s %.5s %.5s %.5s %.6s\t%.1d\tNOP\n",
									binaryInstruction[:8], binaryInstruction[8:11], binaryInstruction[11:16], binaryInstruction[16:21],
									binaryInstruction[21:26], binaryInstruction[26:32], lineNumber)
								//fmt.Printf(output)
								//writeToFile(*outputFile, output)
							default:
								calc, err := twosComplement(binaryInstruction)
								calcString := strconv.Itoa(calc)
								*output += fmt.Sprintf("%.8s %.3s %.5s %.5s %.5s %.5s\t%.1d\t%.5s\n",
									binaryInstruction[:8], binaryInstruction[8:11], binaryInstruction[11:16], binaryInstruction[16:21],
									binaryInstruction[21:26], binaryInstruction[26:32], lineNumber, calcString)
								if err != nil {
								}
								//fmt.Printf(output)
								//writeToFile(*outputFile, output)
							}
						}
					}
				}

			}
		}
	}
}

/*********************************************************/
/*********************************************************/
/***********      MATH       FUNCTIONS    ****************/
/*********************************************************/
/*********************************************************/

/***********************2s COMPLEMENT**********************/

func twosComplement(binaryStr string) (int, error) {
	var negFlag bool //Check for leading one (negative number)
	if binaryStr[0:1] != "1" {
		negFlag = false
		result2, err := binaryToInteger(binaryStr)
		return result2, err
	} else {
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

/************   MAIN   ****************/

func main() {
	// String flag with default value "input.txt"
	inputFile := flag.String("i", "input.txt", "Enter input file")

	// String flag with default value "team16_out_dis.txt"
	outputFile := flag.String("o", "team16_out", "Enter desired name for the output file")

	// Enable command-line parsing
	flag.Parse()
	//simulator output file
	outputSimFile := *outputFile + "_sim.txt"
	//disassembler output file
	*outputFile += "_dis.txt"
	openInputFile, err := os.Open(*inputFile)
	if err != nil {
		fmt.Println("Error opening input file:", err)
		return
	}
	defer func(inputFile *os.File) {
		err := inputFile.Close()
		if err != nil {
		}
	}(openInputFile)

	// Create a reader to read from the input file
	reader := bufio.NewReader(openInputFile)
	//EDITED FROM WRITE TO FILE FUNCTION

	writeOutputFile, err := os.OpenFile(*outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println("Error creating the file:", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(writeOutputFile)
	lineNumber := 96
	output := ""
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
		readAndProcessInstructions(binaryInstruction, lineNumber, outputFile, &output)
		lineNumber += 4
	}
	//writing to file
	_, err = writeOutputFile.WriteString(output)
	if err != nil {
		fmt.Println("Error writing to the file:", err)
		// file.Close()
		return
	}
	err = writeOutputFile.Close()
	if err != nil {
		return
	}
	//WRITE TO SIMULATOR FILE
	writeOutputSimFile, err := os.OpenFile(outputSimFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println("Error creating the file:", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(writeOutputSimFile)
	//Initializing each variable used
	instructionAddress := 96
	cycle := 1
	outputSim := ""
	R0 := 0
	R1 := 0
	R2 := 0
	R3 := 0
	R4 := 0
	R5 := 0
	R6 := 0
	R7 := 0
	R8 := 0
	R9 := 0
	R10 := 0
	R11 := 0
	R12 := 0
	R13 := 0
	R14 := 0
	R15 := 0
	R16 := 0
	R17 := 0
	R18 := 0
	R19 := 0
	R20 := 0
	R21 := 0
	R22 := 0
	R23 := 0
	R24 := 0
	R25 := 0
	R26 := 0
	R27 := 0
	R28 := 0
	R29 := 0
	R30 := 0
	R31 := 0
	//write simulator format
	for instructionAddress <= lineNumber {
		outputSim += "====================\n"
		outputSim += fmt.Sprintf("cycle:%d \t %d \t instruction string\n\n", cycle, instructionAddress)
		outputSim += "registers:\n"
		outputSim += fmt.Sprintf("r00:\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\n",
			R0, R1, R2, R3, R4, R5, R6, R7)
		outputSim += fmt.Sprintf("r08:\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\n",
			R8, R9, R10, R11, R12, R13, R14, R15)
		outputSim += fmt.Sprintf("r16:\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\n",
			R16, R17, R18, R19, R20, R21, R22, R23)
		outputSim += fmt.Sprintf("r24:\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\n",
			R24, R25, R26, R27, R28, R29, R30, R31)
		outputSim += "data:"
		cycle++
		instructionAddress += 4
	}
	//WRITE TO SIMULATOR FILE
	_, err = writeOutputSimFile.WriteString(outputSim)
	if err != nil {
		fmt.Println("Error writing to the file:", err)
		// file.Close()
		return
	}
	err = writeOutputFile.Close()
	if err != nil {
		return
	}

}
