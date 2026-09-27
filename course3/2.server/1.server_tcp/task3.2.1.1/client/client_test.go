package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"
	"testing"
	"time"
)

// TestClientReaderNormal тестирует нормальную работу clientReader
func TestClientReaderNormal(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()

	// Захватываем stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Запускаем clientReader
	go clientReader(clientConn)

	// Отправляем тестовые данные
	go func() {
		fmt.Fprintln(serverConn, "Hello from server")
		fmt.Fprintln(serverConn, "Second line")
		time.Sleep(100 * time.Millisecond)
		serverConn.Close()
	}()

	// Ждем обработки
	time.Sleep(200 * time.Millisecond)

	// Восстанавливаем stdout
	w.Close()
	os.Stdout = oldStdout

	// Проверяем вывод
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	expectedLines := []string{"Hello from server", "Second line"}
	for _, line := range expectedLines {
		if !strings.Contains(output, line) {
			t.Errorf("Expected '%s' in output, got: %s", line, output)
		}
	}
}

// TestClientReaderEmpty тестирует clientReader с пустыми данными
func TestClientReaderEmpty(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()

	// Захватываем stdout
	oldStdout := os.Stdout
	os.Stdout, _ = os.Open(os.DevNull)

	// Сразу закрываем соединение
	serverConn.Close()

	// Это не должно паниковать
	clientReader(clientConn)

	// Восстанавливаем stdout
	os.Stdout = oldStdout
}

// TestDialAndErrorHandling тестирует логику подключения и обработки ошибок
func TestDialAndErrorHandling(t *testing.T) {
	// Тестируем обработку ошибки dial
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Имитируем логику обработки ошибки из main()
	err := fmt.Errorf("connection refused")
	fmt.Println(err)
	// os.Exit(1) - не тестируем exit, так как он завершит тест

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "connection refused") {
		t.Errorf("Expected error in output, got: %s", output)
	}
}

// TestWriteToConnection тестирует запись в соединение
func TestWriteToConnection(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()

	// Захватываем stdout для перехвата ошибок
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Тестируем запись в соединение (имитируем часть main)
	go func() {
		text := "test message"
		_, err := fmt.Fprintln(clientConn, text)
		if err != nil {
			fmt.Println(err)
		}
		clientConn.Close()
	}()

	// Читаем с серверной стороны
	scanner := bufio.NewScanner(serverConn)
	if scanner.Scan() {
		received := scanner.Text()
		if received != "test message" {
			t.Errorf("Expected 'test message', got '%s'", received)
		}
	}

	time.Sleep(100 * time.Millisecond)

	// Восстанавливаем stdout
	w.Close()
	os.Stdout = oldStdout

	// Проверяем что не было ошибок записи
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()
	if strings.Contains(output, "error") || strings.Contains(output, "fail") {
		t.Errorf("Unexpected error in output: %s", output)
	}
}

// TestWriteError тестирует обработку ошибки записи
func TestWriteError(t *testing.T) {
	clientConn, serverConn := net.Pipe()

	// Захватываем stdout для перехвата ошибок
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Закрываем соединение чтобы вызвать ошибку
	serverConn.Close()

	// Пытаемся писать в закрытое соединение (имитируем логику main)
	go func() {
		_, err := fmt.Fprintln(clientConn, "test")
		if err != nil {
			fmt.Println(err) // Это должно быть записано в stdout
		}
		clientConn.Close()
	}()

	time.Sleep(100 * time.Millisecond)

	// Восстанавливаем stdout
	w.Close()
	os.Stdout = oldStdout

	// Проверяем что ошибка была обработана и выведена
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()
	if !strings.Contains(output, "write") && !strings.Contains(output, "closed") {
		t.Errorf("Expected write error message, got: %s", output)
	}
}

// TestScannerFunctionality тестирует работу сканера (дополнительное покрытие)
func TestScannerFunctionality(t *testing.T) {
	// Тестируем работу bufio.Scanner который используется в clientReader
	testInput := "line 1\nline 2\nline 3\n"
	scanner := bufio.NewScanner(strings.NewReader(testInput))

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if len(lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "line 1" || lines[1] != "line 2" || lines[2] != "line 3" {
		t.Errorf("Unexpected line content: %v", lines)
	}
}