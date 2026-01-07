package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

var CRLF = "\r\n"

func buildCommand(cmd []string) string {
	var result string
	result += fmt.Sprintf("*%d%s", len(cmd), CRLF)
	for _, arg := range cmd {
		result += fmt.Sprintf("$%d%s%s%s", len(arg), CRLF, arg, CRLF)
	}
	return result
}

// readRESPResponse 读取 RESP 协议响应
func readRESPResponse(reader *bufio.Reader) (string, error) {
	firstByte, err := reader.ReadByte()
	if err != nil {
		return "", err
	}

	switch firstByte {
	case '+': // Simple String
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		return strings.TrimSuffix(line, "\r\n"), nil

	case '-': // Error
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		return strings.TrimSuffix(line, "\r\n"), nil

	case ':': // Integer
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		return strings.TrimSuffix(line, "\r\n"), nil

	case '$': // Bulk String
		lenLine, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		lenStr := strings.TrimSuffix(lenLine, "\r\n")
		strLen, err := strconv.Atoi(lenStr)
		if err != nil {
			return "", err
		}
		if strLen == -1 {
			return "(nil)", nil
		}
		data := make([]byte, strLen)
		if _, err := io.ReadFull(reader, data); err != nil {
			return "", err
		}
		// 读取末尾的 \r\n
		crlf := make([]byte, 2)
		if _, err := io.ReadFull(reader, crlf); err != nil {
			return "", err
		}
		return string(data), nil

	case '*': // Array
		lenLine, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		lenStr := strings.TrimSuffix(lenLine, "\r\n")
		arrayLen, err := strconv.Atoi(lenStr)
		if err != nil {
			return "", err
		}
		if arrayLen == -1 {
			return "(nil)", nil
		}
		var result strings.Builder
		result.WriteString("[")
		for i := 0; i < arrayLen; i++ {
			if i > 0 {
				result.WriteString(", ")
			}
			item, err := readRESPResponse(reader)
			if err != nil {
				return "", err
			}
			result.WriteString(item)
		}
		result.WriteString("]")
		return result.String(), nil

	default:
		return "", fmt.Errorf("unknown RESP type: %c", firstByte)
	}
}

var addr string

func main() {
	flag.StringVar(&addr, "addr", "127.0.0.1:6380", "address to connect to")
	flag.Parse()
	fmt.Println("Dial to", addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "quit" {
			break
		}
		commands := strings.Fields(line)
		if len(commands) == 0 {
			continue
		}
		commandString := buildCommand(commands)
		_, err := conn.Write([]byte(commandString))
		if err != nil {
			fmt.Println("write error:", err)
			continue
		}
		response, err := readRESPResponse(reader)
		if err != nil {
			fmt.Println("read error:", err)
			continue
		}
		fmt.Println(response)
	}
}
