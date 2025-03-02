package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// TruffleArtifact แทนโครงสร้าง JSON ของ artifact ที่ Truffle สร้าง
type TruffleArtifact struct {
	ContractName string          `json:"contractName"`
	ABI          json.RawMessage `json:"abi"`
	Bytecode     string          `json:"bytecode"`
}

func main() {
	// รับ path จาก command line
	buildDir := flag.String("build", "./build/contracts", "Path to Truffle's build/contracts directory")
	outputDir := flag.String("output", "./abigen", "Output directory for abigen files")
	flag.Parse()

	fmt.Printf("จะแปลง Truffle artifacts จาก %s ไปยัง %s\n", *buildDir, *outputDir)

	// สร้างโฟลเดอร์ output ถ้ายังไม่มี
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Printf("ไม่สามารถสร้างโฟลเดอร์ output: %v\n", err)
		os.Exit(1)
	}

	// อ่านรายชื่อไฟล์ใน build directory
	files, err := os.ReadDir(*buildDir)
	if err != nil {
		fmt.Printf("ไม่สามารถอ่านโฟลเดอร์ build: %v\n", err)
		os.Exit(1)
	}

	// นับจำนวนไฟล์ที่แปลงแล้ว
	processed := 0

	// วนลูปผ่านแต่ละไฟล์
	for _, file := range files {
		// ข้ามถ้าไม่ใช่ไฟล์ .json
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		// อ่านไฟล์ artifact
		artifactPath := filepath.Join(*buildDir, file.Name())
		data, err := os.ReadFile(artifactPath)
		if err != nil {
			fmt.Printf("ไม่สามารถอ่านไฟล์ %s: %v\n", artifactPath, err)
			continue
		}

		// แปลง JSON เป็น TruffleArtifact
		var artifact TruffleArtifact
		if err := json.Unmarshal(data, &artifact); err != nil {
			fmt.Printf("ไม่สามารถแปลง JSON ของไฟล์ %s: %v\n", artifactPath, err)
			continue
		}

		// สร้างชื่อไฟล์ output
		baseName := artifact.ContractName
		abiPath := filepath.Join(*outputDir, baseName+".abi")
		binPath := filepath.Join(*outputDir, baseName+".bin")

		// บันทึกไฟล์ ABI
		if err := os.WriteFile(abiPath, artifact.ABI, 0644); err != nil {
			fmt.Printf("ไม่สามารถบันทึกไฟล์ ABI %s: %v\n", abiPath, err)
			continue
		}

		// ลบ "0x" นำหน้า bytecode (ถ้ามี)
		bytecode := artifact.Bytecode
		if strings.HasPrefix(bytecode, "0x") {
			bytecode = bytecode[2:]
		}

		// บันทึกไฟล์ bytecode
		if err := os.WriteFile(binPath, []byte(bytecode), 0644); err != nil {
			fmt.Printf("ไม่สามารถบันทึกไฟล์ bytecode %s: %v\n", binPath, err)
			continue
		}

		fmt.Printf("แปลงไฟล์ %s เป็น %s และ %s สำเร็จ\n", file.Name(), baseName+".abi", baseName+".bin")
		processed++
	}

	fmt.Printf("เสร็จสิ้น: แปลงไฟล์ทั้งหมด %d ไฟล์\n", processed)
}
