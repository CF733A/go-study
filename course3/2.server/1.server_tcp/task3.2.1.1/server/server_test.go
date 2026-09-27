package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

func TestBroadcaster(t *testing.T) {

	oldEntering := entering
	oldLeaving := leaving
	oldMessages := messages

	entering = make(chan client, 10)
	leaving = make(chan client, 10)
	messages = make(chan string, 10)

	defer func() {
		entering = oldEntering
		leaving = oldLeaving
		messages = oldMessages
	}()

	go broadcaster()

	time.Sleep(10 * time.Millisecond)

	conn1, _ := net.Pipe()
	conn2, _ := net.Pipe()
	defer conn1.Close()
	defer conn2.Close()

	ch1 := make(chan string, 10)
	ch2 := make(chan string, 10)

	cli1 := client{conn1, "client1", ch1}
	cli2 := client{conn2, "client2", ch2}

	entering <- cli1
	entering <- cli2
	time.Sleep(10 * time.Millisecond)

	testMessage := "Hello everyone!"
	messages <- testMessage

	select {
	case msg := <-ch1:
		if msg != testMessage {
			t.Errorf("Client1 expected %q, got %q", testMessage, msg)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Client1 didn't receive message")
	}

	select {
	case msg := <-ch2:
		if msg != testMessage {
			t.Errorf("Client2 expected %q, got %q", testMessage, msg)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Client2 didn't receive message")
	}

	leaving <- cli1
	time.Sleep(10 * time.Millisecond)

	messages <- "Second message"

	select {
	case msg := <-ch2:
		if msg != "Second message" {
			t.Errorf("Client2 expected 'Second message', got %q", msg)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Client2 didn't receive second message")
	}

	select {
	case _, ok := <-ch1:
		if ok {
			t.Error("Client1 received message after leaving")
		}
	case <-time.After(50 * time.Millisecond):
	}

	leaving <- cli2
	time.Sleep(10 * time.Millisecond)

	select {
	case _, ok := <-ch2:
		if ok {
			t.Error("Client2 channel should be closed after leaving")
		}
	case <-time.After(50 * time.Millisecond):
	}
}


func TestHandleConn(t *testing.T) {

    oldEntering := entering
    oldLeaving := leaving
    oldMessages := messages

    entering = make(chan client, 10)
    leaving = make(chan client, 10)
    messages = make(chan string, 10)

    defer func() {
        entering = oldEntering
        leaving = oldLeaving
        messages = oldMessages
    }()

    // Создаем in-memory соединение
    clientConn, serverConn := net.Pipe()
    defer clientConn.Close()
    defer serverConn.Close()

    // Запускаем handleConn
    go handleConn(clientConn)

    // Тест 1: Проверяем приветственное сообщение
    scanner := bufio.NewScanner(serverConn)
    if scanner.Scan() {
        msg := scanner.Text()
        if !strings.Contains(msg, "You are") {
            t.Errorf("Expected welcome message, got %q", msg)
        }
    } else {
        t.Error("Failed to read welcome message")
    }

    // Тест 2: Проверяем сообщение о подключении
    select {
    case msg := <-messages:
        if !strings.Contains(msg, "has arrived") {
            t.Errorf("Expected arrival message, got %q", msg)
        }
    case <-time.After(100 * time.Millisecond):
        t.Error("No arrival message received")
    }


    select {
    case cli := <-entering:
        if cli.conn == nil {
            t.Error("Client connection is nil")
        }
        if cli.name == "" {
            t.Error("Client name is empty")
        }
    case <-time.After(100 * time.Millisecond):
        t.Error("No client added to entering channel")
    }


    testMessage := "Hello from client"
    fmt.Fprintln(serverConn, testMessage)


    select {
    case msg := <-messages:
        if !strings.Contains(msg, "Hello from client") {
            t.Errorf("Expected client message, got %q", msg)
        }
    case <-time.After(100 * time.Millisecond):
        t.Error("No client message sent to messages channel")
    }


    serverConn.Close()


    time.Sleep(50 * time.Millisecond)


    select {
    case msg := <-messages:
        if !strings.Contains(msg, "has left") {
            t.Errorf("Expected leave message, got %q", msg)
        }
    case <-time.After(100 * time.Millisecond):
        t.Error("No leave message received")
    }

    select {
    case cli := <-leaving:
        if cli.conn == nil {
            t.Error("Leaving client connection is nil")
        }
    case <-time.After(100 * time.Millisecond):
        t.Error("No client removed from leaving channel")
    }
}


func TestClientWriter(t *testing.T) {
    clientConn, serverConn := net.Pipe()
    defer clientConn.Close()
    defer serverConn.Close()

    ch := make(chan string, 3)
    
    go clientWriter(clientConn, ch)


    ch <- "Hello"
    ch <- "World"
    ch <- "Test message"
    close(ch)

    scanner := bufio.NewScanner(serverConn)
    var messages []string

    for scanner.Scan() {
        messages = append(messages, scanner.Text())
        if len(messages) == 3 {
            break
        }
    }
    expected := []string{"Hello", "World", "Test message"}
    if len(messages) != len(expected) {
        t.Fatalf("Expected %d messages, got %d: %v", len(expected), len(messages), messages)
    }

    for i, msg := range messages {
        if msg != expected[i] {
            t.Errorf("Expected %q, got %q", expected[i], msg)
        }
    }
}