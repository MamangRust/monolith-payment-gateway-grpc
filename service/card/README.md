### Card GRPC Server
``CardService`` adalah layanan berbasis gRPC yang menangani seluruh operasi terkait data kartu pengguna dalam sistem, seperti pengelolaan kartu, statistik saldo dan transaksi, serta manajemen data aktif dan terhapus (soft delete).

#### Fitur Utama
1. CRUD Kartu: Buat, ubah, dan hapus kartu pengguna dengan sistem soft-delete.
2. Manajemen Data: Mendukung FindByUserId, FindByCardNumber, FindByActive, dan FindByTrashed.
3. Dashboard Analytics: Menampilkan statistik jumlah kartu aktif, terhapus, saldo, dan transaksi.
5. Analitik Transaksi: Statistik berdasarkan bulan/tahun untuk topup, withdraw, transfer (pengirim & penerima).
6. Granular Analytics by Card Number: Statistik individual berdasarkan nomor kartu.