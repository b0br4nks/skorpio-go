package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var subscriptCounter = 0

func subscript(reset bool) int {
	if reset {
		subscriptCounter = 0
	}
	result := subscriptCounter
	subscriptCounter++
	return result
}

var (
	OP_PUSH   = subscript(true)
	OP_PLUS   = subscript(false)
	OP_MINUS  = subscript(false)
	OP_DUMP   = subscript(false)
	COUNT_OPS = subscript(false)
)

type operation struct {
	opcode int
	value  int
}

func push(x int) operation {
	return operation{OP_PUSH, x}
}

func plus() operation {
	return operation{OP_PLUS, 0}
}

func minus() operation {
	return operation{OP_MINUS, 0}
}

func dump() operation {
	return operation{OP_DUMP, 0}
}

func simulateProgram(program []operation) {
	var stack []int
	for _, op := range program {
		switch op.opcode {
		case OP_PUSH:
			stack = append(stack, op.value)
		case OP_PLUS:
			a := stack[len(stack)-1]
			b := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			stack = append(stack, a+b)
		case OP_MINUS:
			a := stack[len(stack)-1]
			b := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			stack = append(stack, b-a)
		case OP_DUMP:
			a := stack[len(stack)-1]
			fmt.Println(a)
			stack = stack[:len(stack)-1]
		default:
			panic("unreachable")
		}
	}
}

func compileProgram(program []operation, outFilePath string) {
	out, err := os.Create(outFilePath)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	out.WriteString("segment .text\n")
	out.WriteString("dump:\n")
	out.WriteString("    ; Implementation of dump function\n")
	out.WriteString("    mov     r9, -3689348814741910323\n")
	out.WriteString("    sub     rsp, 40\n")
	out.WriteString("    mov     BYTE [rsp+31], 10\n")
	out.WriteString("    lea     rcx, [rsp+30]\n")
	out.WriteString(".L2:\n")
	out.WriteString("    mov     rax, rdi\n")
	out.WriteString("    lea     r8, [rsp+32]\n")
	out.WriteString("    mul     r9\n")
	out.WriteString("    mov     rax, rdi\n")
	out.WriteString("    sub     r8, rcx\n")
	out.WriteString("    shr     rdx, 3\n")
	out.WriteString("    lea     rsi, [rdx+rdx*4]\n")
	out.WriteString("    add     rsi, rsi\n")
	out.WriteString("    sub     rax, rsi\n")
	out.WriteString("    add     eax, 48\n")
	out.WriteString("    mov     BYTE [rcx], al\n")
	out.WriteString("    mov     rax, rdi\n")
	out.WriteString("    mov     rdi, rdx\n")
	out.WriteString("    mov     rdx, rcx\n")
	out.WriteString("    sub     rcx, 1\n")
	out.WriteString("    cmp     rax, 9\n")
	out.WriteString("    ja      .L2\n")
	out.WriteString("    lea     rax, [rsp+32]\n")
	out.WriteString("    mov     edi, 1\n")
	out.WriteString("    sub     rdx, rax\n")
	out.WriteString("    xor     eax, eax\n")
	out.WriteString("    lea     rsi, [rsp+32+rdx]\n")
	out.WriteString("    mov     rdx, r8\n")
	out.WriteString("    mov     rax, 1\n")
	out.WriteString("    syscall\n")
	out.WriteString("    add     rsp, 40\n")
	out.WriteString("    ret\n")

	out.WriteString("global _start\n")
	out.WriteString("_start:\n")

	for _, op := range program {
		switch op.opcode {
		case OP_PUSH:
			out.WriteString(fmt.Sprintf("    ;; -- push %d --\n", op.value))
			out.WriteString(fmt.Sprintf("    push %d\n", op.value))
		case OP_PLUS:
			out.WriteString("    ;; -- plus --\n")
			out.WriteString("    pop rax\n")
			out.WriteString("    pop rbx\n")
			out.WriteString("    add rax, rbx\n")
			out.WriteString("    push rax\n")
		case OP_MINUS:
			out.WriteString("    ;; -- minus --\n")
			out.WriteString("    pop rax\n")
			out.WriteString("    pop rbx\n")
			out.WriteString("    sub rbx, rax\n")
			out.WriteString("    push rbx\n")
		case OP_DUMP:
			out.WriteString("    ;; -- dump --\n")
			out.WriteString("    pop rdi\n")
			out.WriteString("    call dump\n")
		default:
			panic("unreachable")
		}
	}

	out.WriteString("    mov rax, 60\n")
	out.WriteString("    xor edi, edi\n")
	out.WriteString("    syscall\n")
}

func parseWordAsOp(word string) operation {
	assertOps()
	switch word {
	case "+":
		return plus()
	case "-":
		return minus()
	case "=>":
		return dump()
	default:
		num, err := strconv.Atoi(word)
		if err != nil {
			panic("Invalid word in program")
		}
		return push(num)
	}
}

func loadProgramFromFile(filePath string) []operation {
	content, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	words := strings.Fields(string(content))
	program := make([]operation, len(words))
	for i, word := range words {
		program[i] = parseWordAsOp(word)
	}
	return program
}

func usage(programName string) {
	fmt.Println("Usage: musc <SUBCOMMAND> [ARGS]")
	fmt.Println("SUBCOMMANDS:")
	fmt.Println("    -s <file>      Simulate the program")
	fmt.Println("    -c <file>      Compile the program")
}

func callCmd(cmd []string) {
	fmt.Println(cmd)
	command := exec.Command(cmd[0], cmd[1:]...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Run()
	if err != nil {
		panic(err)
	}
}

func uncons(xs []string) (string, []string) {
	return xs[0], xs[1:]
}

func assertOps() {
	if COUNT_OPS != 4 {
		panic("Exhaustive op handling in parseWordAsOp")
	}
}

func main() {
	args := os.Args
	if len(args) < 2 {
		usage(args[0])
		fmt.Println("[!] ERROR: no subcommand is provided")
		os.Exit(1)
	}

	programName, remainingArgs := uncons(args)
	subcommand, args := uncons(remainingArgs)

	switch subcommand {
	case "-s":
		if len(args) < 1 {
			usage(programName)
			fmt.Println("[!] ERROR: no input file is provided")
			os.Exit(1)
		}
		programPath, _ := uncons(args)
		program := loadProgramFromFile(programPath)
		simulateProgram(program)

	case "-c":
		if len(args) < 1 {
			usage(programName)
			fmt.Println("[!] ERROR: no input file is provided")
			os.Exit(1)
		}
		programPath, _ := uncons(args)
		program := loadProgramFromFile(programPath)
		compileProgram(program, "skorpio.asm")
		callCmd([]string{"nasm", "-felf64", "skorpio.asm"})
		callCmd([]string{"ld", "-o", "skorpio", "skorpio.o"})

	default:
		usage(programName)
		fmt.Printf("ERROR: unknown subcommand %s\n", subcommand)
		os.Exit(1)
	}
}
