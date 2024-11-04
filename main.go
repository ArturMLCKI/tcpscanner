package main

import (
	"fmt"
	"net"
	"sort"
	"sync"
)

func worker(ports, results chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for p := range ports {
		address := fmt.Sprintf("localhost:%d", p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

func main() {
	//ZA POMOCĄ TEGO KOD SKANUJEMY TYLKO JEDEN PORT
	//_, err := net.Dial("tcp", "scanme.nmap.org:80")
	//if err == nil {
	//	fmt.Println("Connection successful")
	//} else {
	//	fmt.Println("Connection denied")
	//}

	//SKANER SPRAWDZAJĄCY WSZYSTKIE PORTY PO KOLEI I INFORMUJE NAS CZY JEST OTWARTY
	//for i := 1; i <= 1024; i++ {
	//	address := fmt.Sprintf("scanme.nmap.org:%d", i)
	//	conn, err := net.Dial("tcp", address)
	//	if err != nil {
	//		//port is closed or filtered
	//		continue
	//	}
	//	conn.Close()
	//	fmt.Printf("%d open\n", i)
	//}

	//SZYBKI SKANER UŻYWAJĄCY GORUTYN ORAZ SYNC
	//var wg sync.WaitGroup
	//for i := 1; i <= 65535; i++ {
	//	wg.Add(1)
	//	go func(j int) {
	//		defer wg.Done()
	//		address := fmt.Sprintf("127.0.0.1:%d", j)
	//		conn, err := net.Dial("tcp", address)
	//		if err != nil {
	//			return
	//		}
	//
	//		conn.Close()
	//		fmt.Printf("%d open\n", j)
	//	}(i)
	//}
	//wg.Wait()
	ports := make(chan int, 100)
	results := make(chan int)
	var openports []int
	var wg sync.WaitGroup

	for i := 0; i < cap(ports); i++ {
		wg.Add(1)
		go worker(ports, results, &wg)
	}

	go func() {
		for i := 1; i <= 65535; i++ {
			ports <- i
		}
		close(ports)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for port := range results {
		if port != 0 {
			openports = append(openports, port)
		}
	}
	sort.Ints(openports)
	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}
}
