# 🗄️ Mini Object Store (S3-Inspired)

A minimal object storage system inspired by Amazon S3, built to understand how modern object storage works internally.

This project focuses on **chunking**, **distributed storage**, and **metadata-driven retrieval**, while keeping the implementation intentionally simple.

> ⚠️ This is an educational project. Not production-ready.

---

## ✨ Key Concepts

- Object storage (not a traditional filesystem)
- File chunking
- Distribution across multiple storage units (simulated disks)
- Metadata-driven file reconstruction
- FinOps-aware storage design

---

## 📁 Project Structure

```
├── storage_nodes/              # Simulated storage disks
│   ├── disk1/
│   ├── disk2/
│   └── disk3/
│       ├── buckets/            # Buckets (namespaces)
│       └── my-bucket/
│           └── {filename}.meta.json
│
├── downloaded_files/           # Reconstructed downloaded files
│
├── scripts/
│   ├── upload/
│   │   └── main.go             # Upload CLI
│   └── download/
│       └── main.go             # Download CLI
│
└── README.md
```Claude is AI and can make mistakes. Please double-check responses.

---

## ⚙️ How the System Works

### Upload Flow

1. A file is uploaded using the CLI.
2. The file is split into fixed-size chunks (1MB).
3. Chunks are distributed across multiple disks (round-robin).
4. Metadata is written describing where each chunk is stored.

### Download Flow

1. Metadata is read to locate all chunks.
2. Chunks are fetched in order from storage disks.
3. The original file is reconstructed and saved locally.

---

## 🚀 How to Run

### 1️⃣ Prepare directories

```bash
mkdir -p storage_chunks/disk1 storage_chunks/disk2 storage_chunks/disk3
mkdir -p buckets downloaded_files


2️⃣ Upload a file
go run .scrpts/upload my-bucket input_files/demo.txt


This will:

Split the file into chunks

Store chunks across multiple disks

Create buckets/my-bucket/metadata.json

3️⃣ Download a file
go run.scrpts/download my-bucket/demo.txt


The reconstructed file will be available at:

downloaded_files/demo.txt

⚠️ Known Limitations (V1)

No replication or redundancy

No checksum or corruption detection

No concurrent upload/download safety

No disk failure handling

Sequential chunk reads

These limitations are intentional to keep the system easy to understand.
