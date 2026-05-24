package database

import (
	"log"

	"github.com/ayomendaki/ayomendaki-admin/internal/auth"
)

func Seed() error {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM operators").Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		log.Println("Database already seeded, skipping")
		return nil
	}

	log.Println("Seeding database...")

	hash, err := auth.HashPassword("admin123")
	if err != nil {
		return err
	}
	userHash, err := auth.HashPassword("user123")
	if err != nil {
		return err
	}

	users := []struct {
		username, passwordHash, name, role, description string
		rating                                          float64
		verified                                        int
	}{
		{"superadmin", hash, "Super Admin", "superadmin", "Full akses ke seluruh fitur sistem.", 5.0, 1},
		{"admin", hash, "Admin Operator", "admin", "Open trip operator terpercaya sejak 2020.", 4.9, 1},
		{"user", userHash, "User Report", "user", "Hanya bisa melihat laporan.", 4.0, 0},
	}

	for _, u := range users {
		_, err = DB.Exec(`INSERT INTO operators (username, password_hash, name, role, rating, verified, description, phone, whatsapp, instagram) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			u.username, u.passwordHash, u.name, u.role, u.rating, u.verified, u.description, "08123456789", "08123456789", "@ayomendaki")
		if err != nil {
			return err
		}
	}

	provinces := []struct{ id, name, image string }{
		{"jabar", "Jawa Barat", "https://images.unsplash.com/photo-1555400038-63f5ba517a47"},
		{"jateng", "Jawa Tengah", "https://images.unsplash.com/photo-1585409677983-0f6c41ca9c3b"},
		{"jatim", "Jawa Timur", "https://images.unsplash.com/photo-1596900778232-0f6c6d4b9cf5"},
		{"diy", "DI Yogyakarta", "https://images.unsplash.com/photo-1580624474891-7a1f5b8e0d5a"},
		{"banten", "Banten", "https://images.unsplash.com/photo-1555400038-63f5ba517a47"},
	}
	for _, p := range provinces {
		_, err = DB.Exec("INSERT INTO provinces (id, name, image) VALUES (?, ?, ?)", p.id, p.name, p.image)
		if err != nil {
			return err
		}
	}

	mountains := []struct {
		id, name, provinceID string
		height              int
		difficulty          string
		image               string
		description         string
		trending            bool
		lat, lng            float64
		zoom                int
	}{
		{"ciremai", "Gunung Ciremai", "jabar", 3078, "Sulit", "https://images.unsplash.com/photo-1555400038-63f5ba517a47", "Gunung tertinggi di Jawa Barat dengan pemandangan sunrise yang memukau.", true, -6.8947, 108.4, 12},
		{"papandayan", "Gunung Papandayan", "jabar", 2665, "Sedang", "https://images.unsplash.com/photo-1585409677983-0f6c41ca9c3b", "Gunung dengan kawah aktif dan Taman Wisata Alam yang indah.", true, -7.3194, 107.7294, 13},
		{"merapi", "Gunung Merapi", "diy", 2968, "Sulit", "https://images.unsplash.com/photo-1596900778232-0f6c6d4b9cf5", "Gunung berapi paling aktif di Indonesia dengan pemandangan spektakuler.", false, -7.5407, 110.4443, 12},
		{"merbabu", "Gunung Merbabu", "jateng", 3145, "Sulit", "https://images.unsplash.com/photo-1580624474891-7a1f5b8e0d5a", "Gunung dengan padang savana luas dan pemandangan 5 gunung sekaligus.", true, -7.45, 110.4333, 13},
		{"sindoro", "Gunung Sindoro", "jateng", 3136, "Sedang", "https://images.unsplash.com/photo-1555400038-63f5ba517a47", "Gunung dengan jalur pendakian yang indah dan kawah yang eksotis.", true, -7.3, 109.9833, 12},
		{"prau", "Gunung Prau", "jateng", 2565, "Mudah", "https://images.unsplash.com/photo-1585409677983-0f6c41ca9c3b", "Gunung favorit pemula dengan pemandangan sunrise terbaik di Dieng.", true, -7.1833, 109.9, 13},
		{"semeru", "Gunung Semeru", "jatim", 3676, "Sulit", "https://images.unsplash.com/photo-1596900778232-0f6c6d4b9cf5", "Gunung tertinggi di Pulau Jawa dengan panorama Mahameru.", true, -8.1075, 112.92, 11},
		{"bromo", "Gunung Bromo", "jatim", 2329, "Mudah", "https://images.unsplash.com/photo-1580624474891-7a1f5b8e0d5a", "Gunung ikonik dengan lautan pasir dan sunrise yang legendaris.", true, -7.9425, 112.9533, 13},
		{"raung", "Gunung Raung", "jatim", 3344, "Sulit", "https://images.unsplash.com/photo-1555400038-63f5ba517a47", "Gunung dengan kawah terbesar di Pulau Jawa.", false, -8.125, 114.0417, 12},
		{"gede", "Gunung Gede", "jabar", 2958, "Sedang", "https://images.unsplash.com/photo-1585409677983-0f6c41ca9c3b", "Gunung dengan air panas alami dan padang edelweiss.", false, -6.78, 106.9833, 12},
		{"wukir", "Bukit Wukir", "jateng", 500, "Mudah", "https://images.unsplash.com/photo-1596900778232-0f6c6d4b9cf5", "Bukit dengan pemandangan perbukitan hijau yang asri.", false, -7.5, 110.2, 14},
		{"halimun", "Gunung Halimun Salak", "jabar", 1929, "Sedang", "https://images.unsplash.com/photo-1580624474891-7a1f5b8e0d5a", "Kawasan hutan hujan tropis dengan keanekaragaman hayati tinggi.", false, -6.72, 106.52, 12},
		{"karang", "Gunung Karang", "banten", 1778, "Mudah", "https://images.unsplash.com/photo-1555400038-63f5ba517a47", "Gunung dengan kawah mati dan pemandangan Selat Sunda.", false, -6.27, 106.05, 13},
	}
	for _, m := range mountains {
		trending := 0
		if m.trending {
			trending = 1
		}
		_, err = DB.Exec(`INSERT INTO mountains (id, name, province_id, height, difficulty, image, description, trending, lat, lng, zoom) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			m.id, m.name, m.provinceID, m.height, m.difficulty, m.image, m.description, trending, m.lat, m.lng, m.zoom)
		if err != nil {
			return err
		}
	}

	_, err = DB.Exec("UPDATE provinces SET count = (SELECT COUNT(*) FROM mountains WHERE mountains.province_id = provinces.id)")
	if err != nil {
		return err
	}

	// Seed meeting points
	mps := []struct {
		mpType, name, address string
		lat, lng              float64
	}{
		{"titik_jemput", "Cawang UKI", "Jl. Dewi Sartika, Cawang, Jakarta Timur", -6.2447, 106.8585},
		{"titik_jemput", "Pool Cikarang", "Mega Regency, Cikarang, Bekasi", -6.3174, 107.1676},
		{"basecamp", "Basecamp Apuy", "Desa Apuy, Argapura, Majalengka", -6.8947, 108.4000},
		{"basecamp", "Basecamp Papandayan", "Desa Sirnajaya, Cisurupan, Garut", -7.3194, 107.7294},
		{"titik_jemput", "Terminal Pulo Gadung", "Terminal Pulo Gadung, Jakarta Timur", -6.1833, 106.9087},
		{"titik_jemput", "Stasiun Malang", "Stasiun Malang Kota, Jawa Timur", -7.9774, 112.6370},
		{"titik_jemput", "Terminal Probolinggo", "Terminal Bayuangga, Probolinggo", -7.7572, 113.2171},
	}
	for _, mp := range mps {
		_, err = DB.Exec(`INSERT INTO meeting_points (operator_id, type, name, address, lat, lng) VALUES (1, ?, ?, ?, ?, ?)`,
			mp.mpType, mp.name, mp.address, mp.lat, mp.lng)
		if err != nil {
			return err
		}
	}

	// Seed packages with facilities JSON
	pkgs := []struct {
		name, description, facilities string
	}{
		{
			"Paket A - Ekonomis", "Paket ekonomis, semua perlengkapan bawa sendiri",
			`[{"name":"Transportasi PP","detail":"Jakarta - Basecamp","type":"include"},{"name":"Simaksi / Entry Ticket","detail":"","type":"include"},{"name":"Tenda","detail":"Bawa sendiri","type":"exclude"},{"name":"Logistik Team","detail":"Bawa sendiri","type":"exclude"},{"name":"Makan Selama Pendakian","detail":"Bawa sendiri","type":"exclude"},{"name":"Perlengkapan Pribadi","detail":"Shoes, jaket, celana","type":"exclude"},{"name":"Surat Sehat","detail":"","type":"exclude"},{"name":"Asuransi Perjalanan","detail":"+Rp 50.000","type":"exclude"}]`,
		},
		{
			"Paket B - Standard", "Paket lengkap dengan fasilitas tim pendakian",
			`[{"name":"Transportasi PP","detail":"Jakarta - Basecamp","type":"include"},{"name":"Simaksi / Entry Ticket","detail":"","type":"include"},{"name":"Guide","detail":"Berpengalaman","type":"include"},{"name":"Porter Tenda","detail":"Membawa tenda & alat masak","type":"include"},{"name":"Porter Logistik","detail":"Membawa logistik pendakian","type":"include"},{"name":"Porter Air Team","detail":"Air untuk masak & camp","type":"include"},{"name":"Tenda","detail":"Kapasitas 4-5, max isi 4 orang","type":"include"},{"name":"Perlengkapan Masak","detail":"Kompor, nesting, gas","type":"include"},{"name":"Tenda Toilet","detail":"","type":"include"},{"name":"Makan Selama Pendakian","detail":"4x (Siang, Malam, Summit, Sebelum Turun)","type":"include"},{"name":"Makan Pagi Sebelum Pendakian","detail":"","type":"include"},{"name":"Perlengkapan Makan & Minum","detail":"Piring, sendok, gelas","type":"include"},{"name":"Logistik Team","detail":"","type":"include"},{"name":"Welcome Drink","detail":"","type":"include"},{"name":"Buah","detail":"","type":"include"},{"name":"P3K / First Aid Kit","detail":"Standar","type":"include"},{"name":"Handy Talky","detail":"","type":"include"},{"name":"Dokumentasi","detail":"Foto & video","type":"include"},{"name":"Sertifikat Pendakian","detail":"Format PDF","type":"include"},{"name":"Kebersihan Basecamp","detail":"","type":"include"},{"name":"Rumah Singgah","detail":"","type":"include"},{"name":"Makan di Perjalanan / Rest Area","detail":"","type":"exclude"},{"name":"Porter Pribadi","detail":"","type":"exclude"},{"name":"Perlengkapan Pribadi","detail":"Shoes, jaket, celana","type":"exclude"},{"name":"Cemilan Pribadi","detail":"","type":"exclude"},{"name":"Air Mineral","detail":"","type":"exclude"},{"name":"Surat Sehat","detail":"","type":"exclude"},{"name":"Pickup","detail":"","type":"exclude"},{"name":"Asuransi Perjalanan","detail":"+Rp 50.000","type":"exclude"}]`,
		},
		{
			"Paket C - Premium", "Paket all-inclusive dengan perlengkapan pribadi lengkap",
			`[{"name":"Transportasi PP","detail":"Jakarta - Basecamp","type":"include"},{"name":"Simaksi / Entry Ticket","detail":"","type":"include"},{"name":"Guide","detail":"Berpengalaman","type":"include"},{"name":"Porter Tenda","detail":"Membawa tenda & alat masak","type":"include"},{"name":"Porter Logistik","detail":"Membawa logistik pendakian","type":"include"},{"name":"Porter Air Team","detail":"Air untuk masak & camp","type":"include"},{"name":"Tenda","detail":"Kapasitas 4-5, max isi 4 orang","type":"include"},{"name":"Perlengkapan Masak","detail":"Kompor, nesting, gas","type":"include"},{"name":"Tenda Toilet","detail":"","type":"include"},{"name":"Sleeping Bag","detail":"","type":"include"},{"name":"Matras","detail":"","type":"include"},{"name":"Trekking Pole","detail":"1 pasang","type":"include"},{"name":"Makan Selama Pendakian","detail":"4x (Siang, Malam, Summit, Sebelum Turun)","type":"include"},{"name":"Makan Pagi Sebelum Pendakian","detail":"","type":"include"},{"name":"Makan Sebelum & Sesudah Pendakian","detail":"2x","type":"include"},{"name":"Perlengkapan Makan & Minum","detail":"Piring, sendok, gelas","type":"include"},{"name":"Logistik Team","detail":"","type":"include"},{"name":"Welcome Drink","detail":"","type":"include"},{"name":"Buah","detail":"","type":"include"},{"name":"Air Mineral","detail":"1.5L 2 botol + 600ml 1 botol","type":"include"},{"name":"P3K / First Aid Kit","detail":"Standar","type":"include"},{"name":"Handy Talky","detail":"","type":"include"},{"name":"Dokumentasi","detail":"Foto & video","type":"include"},{"name":"Sertifikat Pendakian","detail":"Format PDF","type":"include"},{"name":"Kaos","detail":"","type":"include"},{"name":"Pickup","detail":"","type":"include"},{"name":"Surat Sehat","detail":"","type":"include"},{"name":"Kebersihan Basecamp","detail":"","type":"include"},{"name":"Rumah Singgah","detail":"","type":"include"},{"name":"Stiker & Gantungan Kunci","detail":"","type":"include"},{"name":"Makan di Perjalanan / Rest Area","detail":"","type":"exclude"},{"name":"Porter Pribadi","detail":"","type":"exclude"},{"name":"Perlengkapan Pribadi","detail":"Shoes, jaket, celana","type":"exclude"},{"name":"Cemilan Pribadi","detail":"","type":"exclude"},{"name":"Asuransi Perjalanan","detail":"+Rp 50.000","type":"exclude"}]`,
		},
	}
	for _, pkg := range pkgs {
		_, err = DB.Exec(`INSERT INTO packages (operator_id, name, description, facilities) VALUES (1, ?, ?, ?)`,
			pkg.name, pkg.description, pkg.facilities)
		if err != nil {
			return err
		}
	}

	// Seed addon templates
	addons := []struct{ name, icon string; price int64 }{
		{"Tracking Pole", "bx-walk", 25000},
		{"Carrier", "bx-backpack", 30000},
		{"Sleeping Bag", "bx-bed", 20000},
		{"Tenda", "bx-building-house", 45000},
		{"Headlamp", "bx-bulb", 15000},
		{"Jaket Gunung", "bx-tshirt", 35000},
		{"Sarung Tangan", "bx-hand", 10000},
		{"Matras", "bx-layer", 15000},
	}
	for _, a := range addons {
		_, err = DB.Exec(`INSERT INTO addon_templates (operator_id, name, price, icon) VALUES (1, ?, ?, ?)`,
			a.name, a.price, a.icon)
		if err != nil {
			return err
		}
	}

	log.Println("Seed complete: 3 operators, 5 provinces, 13 mountains, 7 meeting points, 3 packages, 8 addon templates")
	return nil
}
