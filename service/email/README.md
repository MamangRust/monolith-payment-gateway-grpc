### Email(Worker Service)
Layanan ini bertanggung jawab untuk mengirim email secara asynchronous berdasarkan event yang diterima dari berbagai topik Kafka. Ini adalah komponen worker/service yang berjalan secara terus-menerus dan mendengarkan pesan dari Kafka, lalu mengirim email menggunakan SMTP.
#### Fitur Utama
- Mendengarkan berbagai topik Kafka terkait email.
- Mengirimkan email menggunakan SMTP (saat ini menggunakan Ethereal Email untuk testing).
- Mendaftarkan endpoint /metrics menggunakan Prometheus untuk monitoring.
- Struktur handler yang modular dan dapat dikembangkan untuk berbagai jenis email.