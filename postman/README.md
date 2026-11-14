# Postman — LokaNusantara API

Folder ini berada di luar `api/` dan memuat collection Postman v2.1 untuk endpoint backend.

## Import

1. Import `LokaNusantara API.postman_collection.json` ke Postman.
2. Import `LokaNusantara Local.postman_environment.json`, lalu pilih environment **LokaNusantara Local**.
3. Pastikan API berjalan di `http://localhost:8080` (ubah `baseUrl` jika host/port berbeda).
4. Untuk endpoint terlindungi, isi `loginEmail` dan `loginPassword`, lalu jalankan **Login**. Token respons otomatis disimpan sebagai `jwt_token`.
5. Untuk endpoint admin, login dengan akun admin. Setelah membuat resource, isi variable ID yang sesuai (`locationId`, `categoryId`, `eventId`, `trailId`, `reviewId`, atau `savedPlanId`) dari respons.
6. Untuk terjemahan, set `lang` ke `id`, `en`, `ja`, atau `ko`. Endpoint konten menambahkan `include_translations=true` pada contoh yang sesuai.

## Isi

Collection mencakup endpoint bahasa dan auth; baca lokasi, kategori, agenda, jalur; AI; ulasan, bookmark, dan rencana pengguna; CRUD admin untuk lokasi, kategori, agenda, dan jalur; dashboard, pengguna, pengaturan AI, moderasi ulasan; unggah gambar; serta pengambilan file pada `/uploads`.

Semua nilai login, token, API key, ID, dan Google ID token berupa variable kosong. Jangan commit credential asli ke collection atau environment yang dibagikan. Endpoint AI memerlukan konfigurasi layanan AI di server; Google login memerlukan audience token yang dikonfigurasi.
