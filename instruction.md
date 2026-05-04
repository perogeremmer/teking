Role: Senior Frontend Developer & Product Architect.
Task: Buatkan struktur folder dan file HTML/JS terpisah untuk prototype Mobile Web App "Open Trip Marketplace". Gunakan Tailwind CSS (via CDN).

Struktur File yang Diminta:

index.html (Home Dashboard):

Header: Sticky navigation dengan search bar.

Section Provinsi: Horizontal scroll kartu provinsi (Jawa Tengah, Jawa Barat, dsb).

Section Trending: Daftar gunung ramai daki 7 hari terakhir (misal: Ciremai, Papandayan).

Section Top Operators: Daftar 10 operator terbaik (Tiga Dewa, Tebet Adventure, dsb).

Bottom Nav: Menu bar (Home, My Trips, Profile).

region-detail.html (Daftar Gunung per Provinsi):

Halaman yang muncul saat provinsi diklik.

List kartu gunung di provinsi tersebut (Contoh: Jawa Tengah -> Merapi, Merbabu, Sindoro, Sumbing).

Tampilkan info singkat: Ketinggian, Level Kesulitan, dan Jumlah Trip Aktif.

trip-detail.html (Detail Trip & Booking):

Foto hero gunung, deskripsi trip, dan fasilitas.

Cek Slot: Tabel sisa kuota (seperti data Tiga Dewa: "Ciremai Apuy 14-16 Mei: Sisa 9").

Add-on Inventori: Checkbox untuk sewa alat (Tracking Pole, Tenda, Carrier).

Sticky Footer: Tombol "Booking Sekarang".

assets/js/app.js (Shared Logic):

Script sederhana untuk handle navigasi antar halaman (dummy).

Logic untuk filter gunung berdasarkan provinsi.

Kriteria Desain:

Warna: Emerald-900 (Forest Green), Slate-50 (Background), Orange-500 (CTA/Highlight).

Mobile-First: Maksimalkan penggunaan flexbox dan grid Tailwind untuk kenyamanan jempol user.
