#!/bin/bash

# ตรวจสอบว่ามี error หรือไม่
set -e

echo "เริ่มต้นกระบวนการทั้งหมด..."

# 1. รัน truffle compile
echo "กำลัง compile ด้วย Truffle..."
truffle compile

# 2. รัน truffle-to-abigen เพื่อแปลง artifacts
echo "กำลังแปลงไฟล์ artifacts เป็นรูปแบบที่ abigen รองรับ..."
truffle-to-abigen

# 3. สร้าง Go bindings สำหรับทุก contract
echo "กำลังสร้าง Go bindings..."

# สร้างไดเรกทอรีสำหรับเก็บไฟล์ Go
mkdir -p ./contracts-go

# วนลูปผ่านไฟล์ .abi ทั้งหมดในโฟลเดอร์ abigen
for abi_file in ./abigen/*.abi; do
	# ดึงชื่อไฟล์ไม่รวมนามสกุล
	base_name=$(basename "$abi_file" .abi)
	echo "กำลังสร้าง binding สำหรับ $base_name..."

	# ตรวจสอบว่ามีไฟล์ .bin หรือไม่
	bin_file="./abigen/$base_name.bin"
	if [ -f "$bin_file" ] && [ -s "$bin_file" ]; then
		# สร้าง Go binding พร้อม bytecode
		abigen --abi="$abi_file" --bin="$bin_file" --pkg=contracts --out="./contracts-go/$base_name.go"
	else
		# กรณีไม่มี bytecode หรือไฟล์ว่างเปล่า (เช่น interface หรือ abstract contract)
		abigen --abi="$abi_file" --pkg=contracts --out="./contracts-go/$base_name.go"
	fi
done

echo "เสร็จสิ้นกระบวนการทั้งหมด Go bindings ถูกสร้างไว้ที่ ./contracts-go/"
