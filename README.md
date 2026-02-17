# Killer Shell v1.0 - Unified Predator CLI

[**English**](#english) | [**Bahasa Indonesia**](#bahasa-indonesia)

---

## English

Killer Shell is a powerful Command & Control (C2) platform and polymorphic payload generator. It consolidates a scanner, builder, and controller into a single, standalone binary for highly monitored environments.

### Key Features
- **Standalone Binary**: Zero-dependency Go binary with embedded templates.
- **Polyglot Agents**: Generates hardened payloads for PHP, Node.js, and Python.
- **Zero-Knowledge Crypto**: AES-256-GCM encryption with keys derived only in memory.
- **LFI Resilience**: Optimized buffer clearing for PHP inclusion scenarios.
- **Instant Access**: Session persistence for rapid re-connection to active targets.
- **Kill Switch**: Manual `self-destruct` command to securely wipe agents from disk.

### Installation
```bash
go build -o shell main.go
./shell
```

### Usage
1. **BUILD SHELL**: Generate a payload tailored for your target.
2. **SCAN TARGET**: Fingerprint the backend technology (PHP/Node/Python).
3. **COMMAND CONTROL**: Connect to your agent and execute commands.

---

## Bahasa Indonesia

Killer Shell adalah platform Command & Control (C2) dan generator payload polimorfik yang tangguh. Alat ini menggabungkan fitur scanner, builder, dan controller ke dalam satu biner mandiri yang portabel.

### Fitur Utama
- **Standalone Binary**: Biner Go tanpa ketergantungan file eksternal (templates tersemat).
- **Polyglot Agents**: Menghasilkan payload yang diperkeras untuk PHP, Node.js, dan Python.
- **Zero-Knowledge Crypto**: Enkripsi AES-256-GCM dengan kunci yang hanya dibuat di memori.
- **LFI Resilience**: Pembersihan buffer otomatis untuk skenario PHP Inclusion (LFI).
- **Instant Access**: Persistensi sesi untuk penyambungan kembali yang cepat ke target.
- **Kill Switch**: Perintah `self-destruct` manual untuk menghapus agen secara aman.

### Instalasi
```bash
go build -o shell main.go
./shell
```

### Penggunaan
1. **BUILD SHELL**: Buat payload yang disesuaikan dengan target.
2. **SCAN TARGET**: Identifikasi teknologi backend (PHP/Node/Python).
3. **COMMAND CONTROL**: Terhubung ke agen dan jalankan perintah.

---

### Legal Disclaimer / Penyanggahan Hukum
> [!WARNING]
> This tool is developed for authorized security testing in controlled environments. The author is not responsible for any illegal use or damage caused by this tool.
>
> Alat ini dibuat untuk pengujian keamanan berizin di lingkungan yang terkontrol. Jika digunakan untuk tindakan ilegal, itu bukan tanggung jawab author.
