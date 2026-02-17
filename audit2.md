# AUDIT SKALABILITAS SISTEM: KILLER SHELL v1.0
## ANALISIS INFRASTRUKTUR & BATASAN OPERASIONAL MULTI-TARGET

Laporan ini menganalisis kesiapan arsitektur Killer Shell dalam menangani banyak agen atau target secara simultan. Fokus utama adalah identifikasi hambatan struktural yang muncul jika sistem dipaksa beralih dari model 1-operator-1-target ke model kontrol terdistribusi.

### 1. MODEL EKSEKUSI CONTROLLER

* **Sifat Eksekusi**: Controller bersifat pemblokir (Blocking) dan serial.
* **Penanganan Sesi**: Sistem hanya menangani satu sesi aktif pada satu waktu dalam loop utama. Tidak ada penggunaan goroutine untuk pengiriman perintah paralel atau pemrosesan respons di latar belakang.
* **Proses Perintah**: Perintah diproses secara berurutan (Sequential). Operator harus menunggu respons dari target saat ini sebelum dapat mengirim perintah berikutnya atau berinteraksi dengan target lain.

### 2. MODEL MANAJEMEN SESI

* **Registri Sesi**: Killer Shell tidak memiliki registri sesi atau daftar target aktif dalam memori.
* **Struktur Data**: Sistem hanya menggunakan satu struct `ConfigSession` yang menyimpan data satu target terakhir.
* **Batasan**: Karena tidak adanya map atau slice untuk menyimpan daftar target, operator kehilangan konteks target sebelumnya begitu target baru dimasukkan, kecuali melalui file persistensi sesi terakhir.

### 3. DESAIN PERSISTENSI STATE

* **File Sesi**: Menggunakan satu file `.kshell_session` dalam format JSON yang menyimpan satu objek tunggal.
* **Risiko Skalabilitas**: Jika terdapat 10+ agen, sistem tidak memiliki mekanisme untuk membedakan atau mengelola kredensial (Salt, Trigger, Header) untuk masing-masing agen tersebut secara otomatis. Setiap pergantian target memerlukan input manual atau penimpaan file sesi yang sama.

### 4. MODEL KOMUNIKASI AGEN

* **Inisiator**: Komunikasi sepenuhnya diinisiasi oleh Controller (Controller-to-Agent). Agen bersifat pasif dan hanya merespons jika dipicu (Trigger-based).
* **Sinkronisasi**: Komunikasi bersifat sinkron. Controller tertahan (hang) selama durasi request HTTP (Timeout 30 detik). Jika satu agen mengalami latency tinggi, seluruh antarmuka operator akan membeku.
* **Dampak Skalabilitas**: Model ini mencegah pengawasan real-time terhadap banyak agen karena tidak ada mekanisme beaconing (Agen mengirim sinyal ke Controller).

### 5. BATASAN PENGIRIMAN PERINTAH (DISPATCH)

* **Interaksi Multi-Target**: Pengiriman perintah ke 20 agen secara sekaligus (Batching) tidak dimungkinkan dalam arsitektur saat ini.
* **Hambatan**: Tidak adanya antrean tugas (Task Queue) atau logika batching berarti operator harus masuk ke setiap sesi secara manual satu per satu.

### 6. TEKANAN MEMORI & CPU

* **10-25 Agen**: Tekanan pada mesin operator tetap rendah karena operasi dilakukan secara serial.
* **50+ Agen**: Titik tekanan pertama bukan pada sumber daya komputasi, melainkan pada efisiensi operator (Human Bottleneck). Dari sisi teknis, tekanan akan muncul pada penanganan I/O jika timeout terjadi pada banyak target secara berurutan. Arsitektur saat ini sangat hemat sumber daya karena tidak ada pemrosesan paralel.

### 7. BOTTLENECK ARSITEKTURAL

* **Sentralisasi main.go**: Logika C2 terikat erat dalam satu fungsi loop tunggal tanpa abstraksi session manager.
* **Ketiadaan Concurrency Layer**: Tidak ada pemanfaatan fitur konkurensi Go untuk mengelola multiple connection stream.
* **Tight Coupling**: Modul perintah dan modul enkripsi dijalankan secara sinkron dalam thread utama CLI.

### 8. POLA PERILAKU AGEN

* **Passive-Trigger Pattern**: Agen tidak melakukan polling atau beaconing. Ini sangat baik untuk OPSEC (siluman), tetapi sangat buruk untuk manajemen skala besar. Controller tidak tahu apakah agen masih hidup atau sudah dihancurkan kecuali jika operator mencoba melakukan handshake secara manual.

### 9. SKOR KEMATANGAN SKALABILITAS (1-10)

* **Multi-session readiness (1/10)**: Sistem tidak dirancang untuk menangani lebih dari satu sesi dalam satu waktu.
* **Concurrency readiness (2/10)**: Meskipun dibangun dengan Go, potensi konkurensi tidak dimanfaatkan di level aplikasi.
* **Controller scalability (3/10)**: Terbatas pada kecepatan input operator manusia.
* **Agent scalability (8/10)**: Agen sangat ringan dan dapat disebar ke ribuan target tanpa membebani infrastruktur target.

### 10. ESTIMASI BATAS ATAS SKALABILITAS

* **Batas Aman**: 1 - 3 target (melalui manajemen manual yang intens).
* **Titik Kegagalan**: Di atas 5 target, risiko kesalahan input (Salt/Trigger tertukar) meningkat drastis karena tidak adanya manajemen identitas agen.
* **Konsekuensi Skala Besar**: Operasi akan melambat secara eksponensial. Latency pada satu target akan mematikan produktivitas pada target lainnya karena model pemblokiran I/O.

### KESIMPULAN AUDIT

Killer Shell v1.0 adalah alat "Bedah Taktis" yang sangat efisien untuk operasi presisi tunggal. Arsitekturnya saat ini tidak memiliki fondasi untuk menjadi platform "Serangan Massal" (Mass Deployment). Untuk mendukung operasi skala besar, sistem memerlukan perombakan total pada model manajemen sesi (dari Single-Struct ke Registry-Map) dan pengenalan lapisan asinkron untuk komunikasi agen.
