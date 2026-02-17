# Killer Shell v1.0 - UNIFIED PREDATOR CLI

Killer Shell adalah platform pusat kendali (Command & Control) dan generator payload polimorfik yang dirancang untuk lingkungan dengan tingkat pengawasan tinggi. Versi ini mengonsolidasikan seluruh kekuatan scanner, builder, dan controller ke dalam satu biner mandiri yang portabel.

## Struktur Proyek v1.0

- shell (Binary): Biner tunggal hasil kompilasi.
- main.go: Source code konsol terpadu dengan template tersemat.
- output/: Folder default untuk penyimpanan payload hasil generate.

## Alur Kerja (Workflow) Operasional

Operasi dengan Killer Shell v1.0 mengikuti protokol sistematis berikut:

1. RECONNAISSANCE: Menggunakan fitur SCAN TARGET untuk menentukan teknologi backend target (PHP, Node.js, atau Python).
2. FORGING: Menghasilkan payload yang hardened berdasarkan hasil scan menggunakan fitur BUILD SHELL.
3. DEPLOYMENT: Mengunggah payload ke target. Agen akan menetap hingga diperintah mati secara manual.
4. ASCENSION: Menggunakan fitur COMMAND CONTROL (C2) untuk masuk ke sesi kontrol yang terenkripsi penuh.

## Detail Fitur dan Penggunaan

### 1. Build Shell (Forging the Ghost)
Fitur ini menghasilkan agen predator yang telah diperkeras (hardened).
- Cara Penggunaan: Pilih opsi 1, masukkan tipe (php/node/python).
- **Support Custom Header**: Anda dapat menentukan kunci header sendiri (misal: `X-Origin-Token`). Default: `X-Shield-Key`.
- **LFI Resilience (PHP)**: Agent memiliki fitur `ob_end_clean()` otomatis untuk membersihkan sampah HTML dari halaman host saat dieksekusi melalui LFI.
- **Persistence Control (v1.0)**: Fitur auto-kill telah dihapus. Shell akan menetap di server target sampai Anda mengirimkan perintah pemusnahan secara eksplisit.
- Mekanisme KDF: Kunci dekripsi (AES-256-GCM) hanya tercipta saat operator mengirimkan header trigger yang tepat.
- Anti-Forensics: Versi Python menggunakan mekanisme Double-Fork. Penghapusan file dilakukan secara manual via C2.

### 2. Scan Target (Fingerprinting)
Melakukan identifikasi otomatis terhadap web stack target.
- Cara Penggunaan: Pilih opsi 2 dan masukkan URL target.
- Output: Sistem akan mendeteksi teknologi backend untuk merekomendasikan tipe agen yang paling stabil.

### 3. Command Control (C2 Mode)
Mode interaktif untuk mengendalikan target yang telah terinfeksi.
- Cara Penggunaan: Pilih opsi 3, masukkan URL endpoint agen.
- **Custom Header**: Masukkan nama Header Key yang Anda gunakan saat build.
- **Trigger**: Masukkan nilai Trigger (secret value) yang sesuai.
- **self-destruct**: Ketik perintah ini di prompt `pwn@ghost>` untuk menghapus agen secara permanen.
- **Debugging**: Ketik `debug` untuk melihat raw response dari target jika terjadi error format.

## Spesifikasi Teknis Hardening

### Cryptography (Zero-Knowledge)
Kunci sesi diturunkan di dalam memori agen menggunakan algoritma SHA256 terhadap kombinasi Trigger Header dan Salt. Tanpa salah satu komponen ini, biner agen tidak dapat dianalisis untuk mendapatkan kunci akses.

### Memory-Only Persistence
- PHP: Memanfaatkan output buffering rahasia dan exit-on-reply untuk stabilitas LFI.
- Node.js: Pembersihan FD dan global error handling.
- Python: Daemonisasi tingkat kernel (Double-Fork).

## Cara Instalasi

```bash
go build -o shell main.go
./shell
```

CAST THE FINAL CURSE. 🔥🩸
