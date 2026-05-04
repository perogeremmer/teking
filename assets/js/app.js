// Open Trip Marketplace - Platform Perantara
// Data Model: Gunung, Operator, Trip (terpisah)

// ============================================
// DUMMY DATA - GUNUNG (Pure Info)
// ============================================

const PROVINCES = [
  { id: 'jabar', name: 'Jawa Barat', image: 'https://images.unsplash.com/photo-1589308078059-be1415eab4c3?w=300&q=80', count: 3 },
  { id: 'jateng', name: 'Jawa Tengah', image: 'https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=300&q=80', count: 4 },
  { id: 'jatim', name: 'Jawa Timur', image: 'https://images.unsplash.com/photo-1464822759023-fed622ff2c3b?w=300&q=80', count: 3 },
  { id: 'diy', name: 'DI Yogyakarta', image: 'https://images.unsplash.com/photo-1519681393784-d120267933ba?w=300&q=80', count: 1 },
  { id: 'banten', name: 'Banten', image: 'https://images.unsplash.com/photo-1500534314209-a25ddb2bd429?w=300&q=80', count: 2 }
];

const MOUNTAINS = [
  { 
    id: 'ciremai', 
    name: 'Gunung Ciremai', 
    province: 'jabar', 
    height: '3.078 mdpl', 
    difficulty: 'Sulit', 
    image: 'https://images.unsplash.com/photo-1464822759023-fed622ff2c3b?w=600&q=80',
    description: 'Gunung tertinggi di Jawa Barat dengan puncak berbentuk kerucut sempurna dan savana yang luas. Menawarkan pemandangan matahari terbit yang spektakuler dari ketinggian 3.078 mdpl.',
    trending: true,
    location: { lat: -6.8921, lng: 108.4003, zoom: 13 }
  },
  { 
    id: 'papandayan', 
    name: 'Gunung Papandayan', 
    province: 'jabar', 
    height: '2.665 mdpl', 
    difficulty: 'Sedang', 
    image: 'https://images.unsplash.com/photo-1519681393784-d120267933ba?w=600&q=80',
    description: 'Gunung api aktif dengan kawah spektakuler, hutan mati yang instagramable, dan padang edelweiss yang memukau. Cocok untuk pendaki dengan pengalaman menengah.',
    trending: true,
    location: { lat: -7.3203, lng: 107.7303, zoom: 13 }
  },
  { 
    id: 'gede', 
    name: 'Gunung Gede', 
    province: 'jabar', 
    height: '2.958 mdpl', 
    difficulty: 'Sedang', 
    image: 'https://images.unsplash.com/photo-1589308078059-be1415eab4c3?w=600&q=80',
    description: 'Taman Nasional Gede Pangrango menawarkan keanekaragaman hayati yang luar biasa, air terjun, dan pemandangan alam yang masih sangat alami.',
    trending: false,
    location: { lat: -6.7893, lng: 106.9823, zoom: 13 }
  },
  { 
    id: 'merapi', 
    name: 'Gunung Merapi', 
    province: 'jateng', 
    height: '2.930 mdpl', 
    difficulty: 'Sedang', 
    image: 'https://images.unsplash.com/photo-1464822759023-fed622ff2c3b?w=600&q=80',
    description: 'Gunung api paling aktif di Indonesia dengan jalur pendakian yang menantang. Pemandangan kawah dan lautan pasir yang memukau dari ketinggian.',
    trending: true,
    location: { lat: -7.5407, lng: 110.4469, zoom: 13 }
  },
  { 
    id: 'merbabu', 
    name: 'Gunung Merbabu', 
    province: 'jateng', 
    height: '3.145 mdpl', 
    difficulty: 'Sedang', 
    image: 'https://images.unsplash.com/photo-1519681393784-d120267933ba?w=600&q=80',
    description: 'Sabana luas dengan pemandangan sunrise terbaik di Jawa Tengah. Jalur pendakian yang landai dengan panorama gunung Merapi sebagai latar belakang.',
    trending: true,
    location: { lat: -7.4550, lng: 110.4400, zoom: 13 }
  },
  { 
    id: 'sindoro', 
    name: 'Gunung Sindoro', 
    province: 'jateng', 
    height: '3.153 mdpl', 
    difficulty: 'Sulit', 
    image: 'https://images.unsplash.com/photo-1500534314209-a25ddb2bd429?w=600&q=80',
    description: 'Pendakian dengan jalur tebing dan pemandangan kawah yang menakjubkan. Dikenal dengan puncak yang sering diselimuti awan dan panorama yang dramatis.',
    trending: false,
    location: { lat: -7.3022, lng: 109.9922, zoom: 13 }
  },
  { 
    id: 'prau', 
    name: 'Gunung Prau', 
    province: 'jateng', 
    height: '2.565 mdpl', 
    difficulty: 'Mudah', 
    image: 'https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=600&q=80',
    description: 'Golden sunrise dengan lautan awan terindah di Indonesia. Cocok untuk pemula dengan durasi pendakian yang singkat dan jalur yang ramah.',
    trending: true,
    location: { lat: -7.1917, lng: 109.9000, zoom: 13 }
  },
  { 
    id: 'semeru', 
    name: 'Gunung Semeru', 
    province: 'jatim', 
    height: '3.676 mdpl', 
    difficulty: 'Sulit', 
    image: 'https://images.unsplash.com/photo-1464822759023-fed622ff2c3b?w=600&q=80',
    description: 'Gunung tertinggi di Jawa dengan puncak Mahameru yang legendaris. Tantangan sejati bagi pendaki dengan ranu kumbolo yang ikonik di perjalanannya.',
    trending: true,
    location: { lat: -8.1077, lng: 112.9223, zoom: 13 }
  },
  { 
    id: 'bromo', 
    name: 'Gunung Bromo', 
    province: 'jatim', 
    height: '2.329 mdpl', 
    difficulty: 'Mudah', 
    image: 'https://images.unsplash.com/photo-1519681393784-d120267933ba?w=600&q=80',
    description: 'Kawah aktif dengan lautan pasir yang luas dan sunrise yang ikonik. Salah satu destinasi pendakian paling populer dengan akses yang mudah.',
    trending: true,
    location: { lat: -7.9425, lng: 112.9530, zoom: 13 }
  },
  { 
    id: 'raung', 
    name: 'Gunung Raung', 
    province: 'jatim', 
    height: '3.344 mdpl', 
    difficulty: 'Sulit', 
    image: 'https://images.unsplash.com/photo-1500534314209-a25ddb2bd429?w=600&q=80',
    description: 'Pendakian ekstrem dengan jalur via Kalibaru yang penuh tantangan. Puncak tertinggi kedua di Jawa Timur dengan panorama yang luar biasa.',
    trending: false,
    location: { lat: -8.1250, lng: 114.0500, zoom: 13 }
  },
  { 
    id: 'wukir', 
    name: 'Bukit Wukir', 
    province: 'diy', 
    height: '1.200 mdpl', 
    difficulty: 'Mudah', 
    image: 'https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=600&q=80',
    description: 'Trekking ringan dengan pemandangan laut selatan dan sunset yang memukau. Cocok untuk keluarga dan pemula yang ingin menikmati alam.',
    trending: false,
    location: { lat: -8.1000, lng: 110.4000, zoom: 13 }
  },
  { 
    id: 'halimun', 
    name: 'Gunung Halimun Salak', 
    province: 'banten', 
    height: '1.929 mdpl', 
    difficulty: 'Mudah', 
    image: 'https://images.unsplash.com/photo-1464822759023-fed622ff2c3b?w=600&q=80',
    description: 'Taman nasional dengan hutan hujan tropis dan air terjun yang sejuk. Pendakian yang menyenangkan dengan keanekaragaman flora dan fauna.',
    trending: false,
    location: { lat: -6.7400, lng: 106.5300, zoom: 13 }
  },
  { 
    id: 'karang', 
    name: 'Gunung Karang', 
    province: 'banten', 
    height: '1.778 mdpl', 
    difficulty: 'Mudah', 
    image: 'https://images.unsplash.com/photo-1519681393784-d120267933ba?w=600&q=80',
    description: 'Gunung berapi purba dengan pemandungan lembah yang hijau dan asri. Jalur pendakian yang sejuk dan damai, cocok untuk weekend getaway.',
    trending: false,
    location: { lat: -6.2800, lng: 106.0800, zoom: 13 }
  }
];

// ============================================
// DUMMY DATA - OPERATOR
// ============================================

const OPERATORS = [
  {
    id: 'tigadewa',
    name: 'Tiga Dewa Adventure',
    logo: 'https://images.unsplash.com/photo-1523987355523-c7b5b0dd90a7?w=200&q=80',
    rating: 4.9,
    trips_count: 156,
    verified: true,
    description: 'Operator profesional dengan pengalaman lebih dari 5 tahun. Menyediakan trip berkualitas dengan guide bersertifikat dan fasilitas lengkap.',
    contact: {
      phone: '0812-3456-7890',
      whatsapp: '6281234567890',
      instagram: '@tigadewaadventure'
    }
  },
  {
    id: 'tebet',
    name: 'Tebet Adventure',
    logo: 'https://images.unsplash.com/photo-1551632811-561732d1e306?w=200&q=80',
    rating: 4.8,
    trips_count: 134,
    verified: true,
    description: 'Komunitas pecinta alam yang menyediakan trip dengan konsep fun dan safety. Cocok untuk pendaki pemula hingga mahir.',
    contact: {
      phone: '0813-9876-5432',
      whatsapp: '6281398765432',
      instagram: '@tebetadventure'
    }
  },
  {
    id: 'komunitas',
    name: 'Komunitas Petualang',
    logo: 'https://images.unsplash.com/photo-1533240332313-0db49b459ad6?w=200&q=80',
    rating: 4.7,
    trips_count: 98,
    verified: true,
    description: 'Komunitas petualang yang aktif sejak 2015. Fokus pada edukasi alam dan konservasi lingkungan di setiap trip.',
    contact: {
      phone: '0815-1122-3344',
      whatsapp: '6281511223344',
      instagram: '@komunitaspetualang'
    }
  },
  {
    id: 'summit',
    name: 'Summit Seekers',
    logo: 'https://images.unsplash.com/photo-1527668752968-14dc70a27c95?w=200&q=80',
    rating: 4.8,
    trips_count: 112,
    verified: true,
    description: 'Specialist high-altitude expeditions. Guide berpengalaman untuk gunung-gunung ekstrem dengan peralatan standar internasional.',
    contact: {
      phone: '0816-5544-3322',
      whatsapp: '6281655443322',
      instagram: '@summitseekers'
    }
  },
  {
    id: 'hikingbuddies',
    name: 'Hiking Buddies',
    logo: 'https://images.unsplash.com/photo-1501555088652-021faa106b9b?w=200&q=80',
    rating: 4.6,
    trips_count: 87,
    verified: false,
    description: 'Komunitas hiking yang ramah dan welcoming. Trip dengan suasana kekeluargaan dan banyak aktivitas seru di perjalanan.',
    contact: {
      phone: '0817-7788-9900',
      whatsapp: '6281777889900',
      instagram: '@hikingbuddies'
    }
  }
];

// ============================================
// DUMMY DATA - TRIPS (Dibuat Operator untuk Gunung)
// ============================================

const TRIPS = [
  // --- CIREMAI (Jawa Barat) - Tiga Dewa Adventure ---
  {
    id: 'trip-ciremai-1',
    operator_id: 'tigadewa',
    mountain_id: 'ciremai',
    name: 'Open Trip Ciremai via Apuy',
    route: 'Apuy',
    duration: '3 hari 2 malam',
    price: 850000,
    meeting_point: 'Stasiun KA Cirebon, 06.00 WIB',
    meeting_map: { lat: -6.7058, lng: 108.5573 },
    includes: ['Transportasi PP dari Cirebon', 'Makan selama trip (5x)', 'Tenda & cooking set', 'Guide berpengalaman', 'P3K & asuransi', 'Dokumentasi'],
    addons: [
      { id: 'pole', name: 'Tracking Pole', price: 25000, icon: '🥾' },
      { id: 'carrier', name: 'Carrier 60L', price: 50000, icon: '🎒' },
      { id: 'sleeping_bag', name: 'Sleeping Bag', price: 35000, icon: '🛏️' }
    ],
    schedules: [
      // MEI 2026 - Banyak overlapping untuk demo
      { id: 'sched-c1-m1', date_start: '2026-05-01', date_end: '2026-05-03', quota_total: 15, quota_remaining: 8 },
      { id: 'sched-c1-m2', date_start: '2026-05-02', date_end: '2026-05-04', quota_total: 15, quota_remaining: 12 },
      { id: 'sched-c1-m3', date_start: '2026-05-05', date_end: '2026-05-07', quota_total: 15, quota_remaining: 5 },
      { id: 'sched-c1-m4', date_start: '2026-05-06', date_end: '2026-05-08', quota_total: 15, quota_remaining: 14 },
      { id: 'sched-c1-m5', date_start: '2026-05-09', date_end: '2026-05-11', quota_total: 15, quota_remaining: 10 },
      { id: 'sched-c1-m6', date_start: '2026-05-10', date_end: '2026-05-12', quota_total: 15, quota_remaining: 7 },
      // Triple overlap di minggu ke-3!
      { id: 'sched-c1-m7', date_start: '2026-05-13', date_end: '2026-05-15', quota_total: 15, quota_remaining: 9 },
      { id: 'sched-c1-m8', date_start: '2026-05-14', date_end: '2026-05-16', quota_total: 15, quota_remaining: 6 },
      { id: 'sched-c1-m9', date_start: '2026-05-15', date_end: '2026-05-17', quota_total: 15, quota_remaining: 11 },
      { id: 'sched-c1-m10', date_start: '2026-05-18', date_end: '2026-05-20', quota_total: 15, quota_remaining: 13 },
      { id: 'sched-c1-m11', date_start: '2026-05-19', date_end: '2026-05-21', quota_total: 15, quota_remaining: 4 },
      { id: 'sched-c1-m12', date_start: '2026-05-22', date_end: '2026-05-24', quota_total: 15, quota_remaining: 15 },
      { id: 'sched-c1-m13', date_start: '2026-05-23', date_end: '2026-05-25', quota_total: 15, quota_remaining: 8 },
      { id: 'sched-c1-m14', date_start: '2026-05-26', date_end: '2026-05-28', quota_total: 15, quota_remaining: 10 },
      { id: 'sched-c1-m15', date_start: '2026-05-27', date_end: '2026-05-29', quota_total: 15, quota_remaining: 12 },
      { id: 'sched-c1-m16', date_start: '2026-05-29', date_end: '2026-05-31', quota_total: 15, quota_remaining: 7 },
      // JUNI 2026
      { id: 'sched-c1-j1', date_start: '2026-06-04', date_end: '2026-06-06', quota_total: 15, quota_remaining: 12 },
      { id: 'sched-c1-j2', date_start: '2026-06-11', date_end: '2026-06-13', quota_total: 15, quota_remaining: 8 },
      { id: 'sched-c1-j3', date_start: '2026-06-18', date_end: '2026-06-20', quota_total: 15, quota_remaining: 14 },
      { id: 'sched-c1-j4', date_start: '2026-06-25', date_end: '2026-06-27', quota_total: 15, quota_remaining: 10 },
      // JULI 2026
      { id: 'sched-c1-ju1', date_start: '2026-07-02', date_end: '2026-07-04', quota_total: 15, quota_remaining: 15 },
      { id: 'sched-c1-ju2', date_start: '2026-07-09', date_end: '2026-07-11', quota_total: 15, quota_remaining: 7 }
    ]
  },
  // --- CIREMAI - Linggarjati (Tiga Dewa) ---
  {
    id: 'trip-ciremai-2',
    operator_id: 'tigadewa',
    mountain_id: 'ciremai',
    name: 'Open Trip Ciremai via Linggarjati',
    route: 'Linggarjati',
    duration: '3 hari 2 malam',
    price: 850000,
    meeting_point: 'Terminal Kuningan, 05.30 WIB',
    meeting_map: { lat: -7.0149, lng: 108.4833 },
    includes: ['Transportasi PP dari Kuningan', 'Makan selama trip (5x)', 'Tenda & cooking set', 'Guide berpengalaman', 'P3K & asuransi', 'Dokumentasi'],
    addons: [
      { id: 'pole', name: 'Tracking Pole', price: 25000, icon: '🥾' },
      { id: 'carrier', name: 'Carrier 60L', price: 50000, icon: '🎒' },
      { id: 'sleeping_bag', name: 'Sleeping Bag', price: 35000, icon: '🛏️' }
    ],
    schedules: [
      { id: 'sched-c2-1', date_start: '2026-05-17', date_end: '2026-05-19', quota_total: 12, quota_remaining: 4 },
      { id: 'sched-c2-2', date_start: '2026-05-24', date_end: '2026-05-26', quota_total: 12, quota_remaining: 8 },
      { id: 'sched-c2-3', date_start: '2026-05-31', date_end: '2026-06-02', quota_total: 12, quota_remaining: 6 },
      { id: 'sched-c2-4', date_start: '2026-06-07', date_end: '2026-06-09', quota_total: 12, quota_remaining: 10 },
      { id: 'sched-c2-5', date_start: '2026-06-14', date_end: '2026-06-16', quota_total: 12, quota_remaining: 3 },
      { id: 'sched-c2-6', date_start: '2026-06-21', date_end: '2026-06-23', quota_total: 12, quota_remaining: 11 },
      { id: 'sched-c2-7', date_start: '2026-06-28', date_end: '2026-06-30', quota_total: 12, quota_remaining: 9 },
      { id: 'sched-c2-8', date_start: '2026-07-05', date_end: '2026-07-07', quota_total: 12, quota_remaining: 7 },
      { id: 'sched-c2-9', date_start: '2026-07-12', date_end: '2026-07-14', quota_total: 12, quota_remaining: 5 },
      { id: 'sched-c2-10', date_start: '2026-07-19', date_end: '2026-07-21', quota_total: 12, quota_remaining: 12 }
    ]
  },
  // --- CIREMAI - Tebet Adventure ---
  {
    id: 'trip-ciremai-3',
    operator_id: 'tebet',
    mountain_id: 'ciremai',
    name: 'Adventure Ciremai via Apuy',
    route: 'Apuy',
    duration: '3 hari 2 malam',
    price: 800000,
    meeting_point: 'Stasiun KA Cirebon, 06.00 WIB',
    meeting_map: { lat: -6.7058, lng: 108.5573 },
    includes: ['Transportasi PP dari Cirebon', 'Makan selama trip (5x)', 'Tenda berkualitas', 'Guide ramah', 'Snack & mineral water', 'Dokumentasi drone'],
    addons: [
      { id: 'pole', name: 'Tracking Pole', price: 20000, icon: '🥾' },
      { id: 'tent', name: 'Tenda Single', price: 60000, icon: '⛺' },
      { id: 'headlamp', name: 'Headlamp', price: 15000, icon: '🔦' }
    ],
    schedules: [
      { id: 'sched-c3-1', date_start: '2026-05-15', date_end: '2026-05-17', quota_total: 20, quota_remaining: 14 },
      { id: 'sched-c3-2', date_start: '2026-05-22', date_end: '2026-05-24', quota_total: 20, quota_remaining: 18 },
      { id: 'sched-c3-3', date_start: '2026-05-29', date_end: '2026-05-31', quota_total: 20, quota_remaining: 16 },
      { id: 'sched-c3-4', date_start: '2026-06-05', date_end: '2026-06-07', quota_total: 20, quota_remaining: 12 },
      { id: 'sched-c3-5', date_start: '2026-06-12', date_end: '2026-06-14', quota_total: 20, quota_remaining: 19 },
      { id: 'sched-c3-6', date_start: '2026-06-19', date_end: '2026-06-21', quota_total: 20, quota_remaining: 15 },
      { id: 'sched-c3-7', date_start: '2026-06-26', date_end: '2026-06-28', quota_total: 20, quota_remaining: 20 },
      { id: 'sched-c3-8', date_start: '2026-07-03', date_end: '2026-07-05', quota_total: 20, quota_remaining: 17 },
      { id: 'sched-c3-9', date_start: '2026-07-10', date_end: '2026-07-12', quota_total: 20, quota_remaining: 8 },
      { id: 'sched-c3-10', date_start: '2026-07-17', date_end: '2026-07-19', quota_total: 20, quota_remaining: 11 }
    ]
  },

  // --- PAPANDAYAN (Jawa Barat) - Tiga Dewa ---
  {
    id: 'trip-papandayan-1',
    operator_id: 'tigadewa',
    mountain_id: 'papandayan',
    name: 'Open Trip Papandayan Camp David',
    route: 'Camp David',
    duration: '2 hari 1 malam',
    price: 650000,
    meeting_point: 'Terminal Garut, 07.00 WIB',
    meeting_map: { lat: -7.3920, lng: 107.9056 },
    includes: ['Transportasi PP dari Garut', 'Makan (3x)', 'Tenda & cooking set', 'Guide lokal', 'P3K', 'Dokumentasi'],
    addons: [
      { id: 'pole', name: 'Tracking Pole', price: 25000, icon: '🥾' },
      { id: 'carrier', name: 'Carrier 60L', price: 50000, icon: '🎒' }
    ],
    schedules: [
      { id: 'sched-p1-1', date_start: '2026-05-15', date_end: '2026-05-16', quota_total: 20, quota_remaining: 7 },
      { id: 'sched-p1-2', date_start: '2026-05-18', date_end: '2026-05-19', quota_total: 20, quota_remaining: 12 },
      { id: 'sched-p1-3', date_start: '2026-05-22', date_end: '2026-05-23', quota_total: 20, quota_remaining: 15 },
      { id: 'sched-p1-4', date_start: '2026-05-25', date_end: '2026-05-26', quota_total: 20, quota_remaining: 9 },
      { id: 'sched-p1-5', date_start: '2026-05-29', date_end: '2026-05-30', quota_total: 20, quota_remaining: 18 },
      { id: 'sched-p1-6', date_start: '2026-06-01', date_end: '2026-06-02', quota_total: 20, quota_remaining: 20 },
      { id: 'sched-p1-7', date_start: '2026-06-05', date_end: '2026-06-06', quota_total: 20, quota_remaining: 14 },
      { id: 'sched-p1-8', date_start: '2026-06-08', date_end: '2026-06-09', quota_total: 20, quota_remaining: 6 },
      { id: 'sched-p1-9', date_start: '2026-06-12', date_end: '2026-06-13', quota_total: 20, quota_remaining: 11 },
      { id: 'sched-p1-10', date_start: '2026-06-15', date_end: '2026-06-16', quota_total: 20, quota_remaining: 16 }
    ]
  },
  // --- PAPANDAYAN - Komunitas ---
  {
    id: 'trip-papandayan-2',
    operator_id: 'komunitas',
    mountain_id: 'papandayan',
    name: 'Explore Papandayan',
    route: 'Camp David',
    duration: '2 hari 1 malam',
    price: 600000,
    meeting_point: 'Terminal Garut, 07.00 WIB',
    meeting_map: { lat: -7.3920, lng: 107.9056 },
    includes: ['Transportasi PP', 'Makan (3x)', 'Tenda', 'Guide', 'Kegiatan edukasi alam', 'Dokumentasi'],
    addons: [
      { id: 'pole', name: 'Tracking Pole', price: 20000, icon: '🥾' },
      { id: 'sleeping_bag', name: 'Sleeping Bag', price: 30000, icon: '🛏️' }
    ],
    schedules: [
      { id: 'sched-p2-1', date_start: '2026-05-18', date_end: '2026-05-19', quota_total: 25, quota_remaining: 18 },
      { id: 'sched-p2-2', date_start: '2026-05-21', date_end: '2026-05-22', quota_total: 25, quota_remaining: 22 },
      { id: 'sched-p2-3', date_start: '2026-05-25', date_end: '2026-05-26', quota_total: 25, quota_remaining: 15 },
      { id: 'sched-p2-4', date_start: '2026-05-28', date_end: '2026-05-29', quota_total: 25, quota_remaining: 20 },
      { id: 'sched-p2-5', date_start: '2026-06-01', date_end: '2026-06-02', quota_total: 25, quota_remaining: 12 },
      { id: 'sched-p2-6', date_start: '2026-06-04', date_end: '2026-06-05', quota_total: 25, quota_remaining: 23 },
      { id: 'sched-p2-7', date_start: '2026-06-08', date_end: '2026-06-09', quota_total: 25, quota_remaining: 17 },
      { id: 'sched-p2-8', date_start: '2026-06-11', date_end: '2026-06-12', quota_total: 25, quota_remaining: 9 },
      { id: 'sched-p2-9', date_start: '2026-06-15', date_end: '2026-06-16', quota_total: 25, quota_remaining: 25 },
      { id: 'sched-p2-10', date_start: '2026-06-18', date_end: '2026-06-19', quota_total: 25, quota_remaining: 14 }
    ]
  },

  // --- MERAPI (Jawa Tengah) - Tiga Dewa ---
  {
    id: 'trip-merapi-1',
    operator_id: 'tigadewa',
    mountain_id: 'merapi',
    name: 'Sunrise Merapi via Selo',
    route: 'Selo',
    duration: '2 hari 1 malam',
    price: 550000,
    meeting_point: 'Basecamp Selo, 22.00 WIB',
    meeting_map: { lat: -7.5603, lng: 110.4425 },
    includes: ['Transportasi dari Jogja', 'Makan (3x)', 'Tenda', 'Guide berpengalaman', 'Headlamp', 'Dokumentasi'],
    addons: [
      { id: 'carrier', name: 'Carrier 60L', price: 45000, icon: '🎒' },
      { id: 'sleeping_bag', name: 'Sleeping Bag', price: 30000, icon: '🛏️' },
      { id: 'pole', name: 'Tracking Pole', price: 20000, icon: '🥾' }
    ],
    schedules: [
      { id: 'sched-m1-1', date_start: '2026-05-14', date_end: '2026-05-15', quota_total: 10, quota_remaining: 3 },
      { id: 'sched-m1-2', date_start: '2026-05-17', date_end: '2026-05-18', quota_total: 10, quota_remaining: 7 },
      { id: 'sched-m1-3', date_start: '2026-05-21', date_end: '2026-05-22', quota_total: 10, quota_remaining: 5 },
      { id: 'sched-m1-4', date_start: '2026-05-24', date_end: '2026-05-25', quota_total: 10, quota_remaining: 9 },
      { id: 'sched-m1-5', date_start: '2026-05-28', date_end: '2026-05-29', quota_total: 10, quota_remaining: 2 },
      { id: 'sched-m1-6', date_start: '2026-06-04', date_end: '2026-06-05', quota_total: 10, quota_remaining: 8 },
      { id: 'sched-m1-7', date_start: '2026-06-07', date_end: '2026-06-08', quota_total: 10, quota_remaining: 6 },
      { id: 'sched-m1-8', date_start: '2026-06-11', date_end: '2026-06-12', quota_total: 10, quota_remaining: 10 },
      { id: 'sched-m1-9', date_start: '2026-06-14', date_end: '2026-06-15', quota_total: 10, quota_remaining: 4 },
      { id: 'sched-m1-10', date_start: '2026-06-18', date_end: '2026-06-19', quota_total: 10, quota_remaining: 1 }
    ]
  },
  // --- MERAPI - Summit ---
  {
    id: 'trip-merapi-2',
    operator_id: 'summit',
    mountain_id: 'merapi',
    name: 'Merapi Expedition via Selo',
    route: 'Selo',
    duration: '2 hari 1 malam',
    price: 500000,
    meeting_point: 'Basecamp Selo, 22.00 WIB',
    meeting_map: { lat: -7.5603, lng: 110.4425 },
    includes: ['Transportasi dari Jogja', 'Makan (3x)', 'Tenda premium', 'Guide bersertifikat', 'P3K lengkap', 'Drone footage'],
    addons: [
      { id: 'carrier', name: 'Carrier 60L', price: 50000, icon: '🎒' },
      { id: 'sleeping_bag', name: 'Sleeping Bag -5°C', price: 40000, icon: '🛏️' }
    ],
    schedules: [
      { id: 'sched-m2-1', date_start: '2026-05-16', date_end: '2026-05-17', quota_total: 15, quota_remaining: 11 },
      { id: 'sched-m2-2', date_start: '2026-05-19', date_end: '2026-05-20', quota_total: 15, quota_remaining: 13 },
      { id: 'sched-m2-3', date_start: '2026-05-23', date_end: '2026-05-24', quota_total: 15, quota_remaining: 8 },
      { id: 'sched-m2-4', date_start: '2026-05-26', date_end: '2026-05-27', quota_total: 15, quota_remaining: 15 },
      { id: 'sched-m2-5', date_start: '2026-05-30', date_end: '2026-05-31', quota_total: 15, quota_remaining: 6 },
      { id: 'sched-m2-6', date_start: '2026-06-06', date_end: '2026-06-07', quota_total: 15, quota_remaining: 12 },
      { id: 'sched-m2-7', date_start: '2026-06-09', date_end: '2026-06-10', quota_total: 15, quota_remaining: 9 },
      { id: 'sched-m2-8', date_start: '2026-06-13', date_end: '2026-06-14', quota_total: 15, quota_remaining: 7 },
      { id: 'sched-m2-9', date_start: '2026-06-16', date_end: '2026-06-17', quota_total: 15, quota_remaining: 14 },
      { id: 'sched-m2-10', date_start: '2026-06-20', date_end: '2026-06-21', quota_total: 15, quota_remaining: 5 }
    ]
  },

  // --- MERBABU (Jawa Tengah) - Tebet ---
  {
    id: 'trip-merbabu-1',
    operator_id: 'tebet',
    mountain_id: 'merbabu',
    name: 'Open Trip Merbabu Savana',
    route: 'Suhat',
    duration: '3 hari 2 malam',
    price: 750000,
    meeting_point: 'Basecamp Suhat, 06.00 WIB',
    meeting_map: { lat: -7.4550, lng: 110.4400 },
    includes: ['Transportasi PP dari Jogja', 'Makan selama trip (5x)', 'Tenda & cooking set', 'Guide', 'P3K', 'Dokumentasi'],
    addons: [
      { id: 'pole', name: 'Tracking Pole', price: 25000, icon: '🥾' },
      { id: 'carrier', name: 'Carrier 60L', price: 50000, icon: '🎒' },
      { id: 'tent', name: 'Tenda (2P)', price: 70000, icon: '⛺' }
    ],
    schedules: [
      { id: 'sched-mb1-1', date_start: '2026-05-15', date_end: '2026-05-17', quota_total: 12, quota_remaining: 6 },
      { id: 'sched-mb1-2', date_start: '2026-05-22', date_end: '2026-05-24', quota_total: 12, quota_remaining: 9 },
      { id: 'sched-mb1-3', date_start: '2026-05-29', date_end: '2026-05-31', quota_total: 12, quota_remaining: 12 },
      { id: 'sched-mb1-4', date_start: '2026-06-05', date_end: '2026-06-07', quota_total: 12, quota_remaining: 4 },
      { id: 'sched-mb1-5', date_start: '2026-06-12', date_end: '2026-06-14', quota_total: 12, quota_remaining: 8 },
      { id: 'sched-mb1-6', date_start: '2026-06-19', date_end: '2026-06-21', quota_total: 12, quota_remaining: 11 },
      { id: 'sched-mb1-7', date_start: '2026-06-26', date_end: '2026-06-28', quota_total: 12, quota_remaining: 7 },
      { id: 'sched-mb1-8', date_start: '2026-07-03', date_end: '2026-07-05', quota_total: 12, quota_remaining: 10 },
      { id: 'sched-mb1-9', date_start: '2026-07-10', date_end: '2026-07-12', quota_total: 12, quota_remaining: 3 },
      { id: 'sched-mb1-10', date_start: '2026-07-17', date_end: '2026-07-19', quota_total: 12, quota_remaining: 5 }
    ]
  },

  // --- PRAU (Jawa Tengah) - Hiking Buddies ---
  {
    id: 'trip-prau-1',
    operator_id: 'hikingbuddies',
    mountain_id: 'prau',
    name: 'Golden Sunrise Prau',
    route: 'Patak Banteng',
    duration: '2 hari 1 malam',
    price: 350000,
    meeting_point: 'Basecamp Patak Banteng, 11.00 WIB',
    meeting_map: { lat: -7.1917, lng: 109.9000 },
    includes: ['Transportasi dari Dieng', 'Makan (3x)', 'Tenda', 'Guide', 'Api unggun', 'Games & fun'],
    addons: [
      { id: 'sleeping_bag', name: 'Sleeping Bag', price: 25000, icon: '🛏️' },
      { id: 'mat', name: 'Matras', price: 10000, icon: '🛏️' }
    ],
    schedules: [
      { id: 'sched-pr1-1', date_start: '2026-05-14', date_end: '2026-05-15', quota_total: 30, quota_remaining: 20 },
      { id: 'sched-pr1-2', date_start: '2026-05-18', date_end: '2026-05-19', quota_total: 30, quota_remaining: 25 },
      { id: 'sched-pr1-3', date_start: '2026-05-21', date_end: '2026-05-22', quota_total: 30, quota_remaining: 18 },
      { id: 'sched-pr1-4', date_start: '2026-05-25', date_end: '2026-05-26', quota_total: 30, quota_remaining: 28 },
      { id: 'sched-pr1-5', date_start: '2026-05-28', date_end: '2026-05-29', quota_total: 30, quota_remaining: 15 },
      { id: 'sched-pr1-6', date_start: '2026-06-04', date_end: '2026-06-05', quota_total: 30, quota_remaining: 22 },
      { id: 'sched-pr1-7', date_start: '2026-06-08', date_end: '2026-06-09', quota_total: 30, quota_remaining: 12 },
      { id: 'sched-pr1-8', date_start: '2026-06-11', date_end: '2026-06-12', quota_total: 30, quota_remaining: 30 },
      { id: 'sched-pr1-9', date_start: '2026-06-15', date_end: '2026-06-16', quota_total: 30, quota_remaining: 8 },
      { id: 'sched-pr1-10', date_start: '2026-06-18', date_end: '2026-06-19', quota_total: 30, quota_remaining: 17 }
    ]
  },
  // --- PRAU - Komunitas ---
  {
    id: 'trip-prau-2',
    operator_id: 'komunitas',
    mountain_id: 'prau',
    name: 'Prau Dieng Adventure',
    route: 'Dieng',
    duration: '2 hari 1 malam',
    price: 400000,
    meeting_point: 'Dieng Plateau, 10.00 WIB',
    meeting_map: { lat: -7.2000, lng: 109.9000 },
    includes: ['Transportasi', 'Makan (3x)', 'Tenda', 'Guide', 'Edukasi geologi', 'Dokumentasi'],
    addons: [
      { id: 'sleeping_bag', name: 'Sleeping Bag', price: 20000, icon: '🛏️' },
      { id: 'pole', name: 'Tracking Pole', price: 15000, icon: '🥾' }
    ],
    schedules: [
      { id: 'sched-pr2-1', date_start: '2026-05-15', date_end: '2026-05-16', quota_total: 25, quota_remaining: 22 },
      { id: 'sched-pr2-2', date_start: '2026-05-19', date_end: '2026-05-20', quota_total: 25, quota_remaining: 20 },
      { id: 'sched-pr2-3', date_start: '2026-05-22', date_end: '2026-05-23', quota_total: 25, quota_remaining: 18 },
      { id: 'sched-pr2-4', date_start: '2026-05-26', date_end: '2026-05-27', quota_total: 25, quota_remaining: 24 },
      { id: 'sched-pr2-5', date_start: '2026-05-29', date_end: '2026-05-30', quota_total: 25, quota_remaining: 16 },
      { id: 'sched-pr2-6', date_start: '2026-06-05', date_end: '2026-06-06', quota_total: 25, quota_remaining: 23 },
      { id: 'sched-pr2-7', date_start: '2026-06-09', date_end: '2026-06-10', quota_total: 25, quota_remaining: 14 },
      { id: 'sched-pr2-8', date_start: '2026-06-12', date_end: '2026-06-13', quota_total: 25, quota_remaining: 25 },
      { id: 'sched-pr2-9', date_start: '2026-06-16', date_end: '2026-06-17', quota_total: 25, quota_remaining: 10 },
      { id: 'sched-pr2-10', date_start: '2026-06-19', date_end: '2026-06-20', quota_total: 25, quota_remaining: 19 }
    ]
  },

  // --- SEMERU (Jawa Timur) - Tiga Dewa ---
  {
    id: 'trip-semeru-1',
    operator_id: 'tigadewa',
    mountain_id: 'semeru',
    name: 'Expedition Mahameru',
    route: 'Ranu Pane',
    duration: '5 hari 4 malam',
    price: 1350000,
    meeting_point: 'Stasiun Malang / Surabaya, 07.00 WIB',
    meeting_map: { lat: -8.1500, lng: 112.6000 },
    includes: ['Transportasi PP', 'Makan selama trip (10x)', 'Tenda & cooking set', 'Guide berpengalaman', 'Porter', 'P3K & asuransi', 'Dokumentasi'],
    addons: [
      { id: 'carrier', name: 'Carrier 80L', price: 75000, icon: '🎒' },
      { id: 'sleeping_bag', name: 'Sleeping Bag -10°C', price: 50000, icon: '🛏️' },
      { id: 'pole', name: 'Tracking Pole (2 pcs)', price: 40000, icon: '🥾' },
      { id: 'tent', name: 'Tenda Single', price: 80000, icon: '⛺' }
    ],
    schedules: [
      { id: 'sched-s1-1', date_start: '2026-05-17', date_end: '2026-05-21', quota_total: 10, quota_remaining: 4 },
      { id: 'sched-s1-2', date_start: '2026-05-24', date_end: '2026-05-28', quota_total: 10, quota_remaining: 7 },
      { id: 'sched-s1-3', date_start: '2026-05-31', date_end: '2026-06-04', quota_total: 10, quota_remaining: 9 },
      { id: 'sched-s1-4', date_start: '2026-06-07', date_end: '2026-06-11', quota_total: 10, quota_remaining: 6 },
      { id: 'sched-s1-5', date_start: '2026-06-14', date_end: '2026-06-18', quota_total: 10, quota_remaining: 3 },
      { id: 'sched-s1-6', date_start: '2026-06-21', date_end: '2026-06-25', quota_total: 10, quota_remaining: 8 },
      { id: 'sched-s1-7', date_start: '2026-06-28', date_end: '2026-07-02', quota_total: 10, quota_remaining: 10 },
      { id: 'sched-s1-8', date_start: '2026-07-05', date_end: '2026-07-09', quota_total: 10, quota_remaining: 5 },
      { id: 'sched-s1-9', date_start: '2026-07-12', date_end: '2026-07-16', quota_total: 10, quota_remaining: 2 },
      { id: 'sched-s1-10', date_start: '2026-07-19', date_end: '2026-07-23', quota_total: 10, quota_remaining: 7 }
    ]
  },
  // --- SEMERU - Summit ---
  {
    id: 'trip-semeru-2',
    operator_id: 'summit',
    mountain_id: 'semeru',
    name: 'Semeru Extreme',
    route: 'Ranu Pane',
    duration: '5 hari 4 malam',
    price: 1450000,
    meeting_point: 'Bandara Abdul Rachman Saleh, 06.00 WIB',
    meeting_map: { lat: -8.1500, lng: 112.6000 },
    includes: ['Transportasi dari bandara', 'Makan premium (10x)', 'Tenda premium', 'Guide bersertifikat', 'Porter profesional', 'P3K lengkap', 'Drone & GoPro'],
    addons: [
      { id: 'carrier', name: 'Carrier 80L Osprey', price: 100000, icon: '🎒' },
      { id: 'sleeping_bag', name: 'Sleeping Bag -15°C', price: 75000, icon: '🛏️' },
      { id: 'headlamp', name: 'Headlamp Fenix', price: 25000, icon: '🔦' }
    ],
    schedules: [
      { id: 'sched-s2-1', date_start: '2026-05-24', date_end: '2026-05-28', quota_total: 8, quota_remaining: 6 },
      { id: 'sched-s2-2', date_start: '2026-05-31', date_end: '2026-06-04', quota_total: 8, quota_remaining: 8 },
      { id: 'sched-s2-3', date_start: '2026-06-07', date_end: '2026-06-11', quota_total: 8, quota_remaining: 5 },
      { id: 'sched-s2-4', date_start: '2026-06-14', date_end: '2026-06-18', quota_total: 8, quota_remaining: 3 },
      { id: 'sched-s2-5', date_start: '2026-06-21', date_end: '2026-06-25', quota_total: 8, quota_remaining: 7 },
      { id: 'sched-s2-6', date_start: '2026-06-28', date_end: '2026-07-02', quota_total: 8, quota_remaining: 4 },
      { id: 'sched-s2-7', date_start: '2026-07-05', date_end: '2026-07-09', quota_total: 8, quota_remaining: 2 },
      { id: 'sched-s2-8', date_start: '2026-07-12', date_end: '2026-07-16', quota_total: 8, quota_remaining: 8 },
      { id: 'sched-s2-9', date_start: '2026-07-19', date_end: '2026-07-23', quota_total: 8, quota_remaining: 6 },
      { id: 'sched-s2-10', date_start: '2026-07-26', date_end: '2026-07-30', quota_total: 8, quota_remaining: 1 }
    ]
  },

  // --- BROMO (Jawa Timur) - Tebet ---
  {
    id: 'trip-bromo-1',
    operator_id: 'tebet',
    mountain_id: 'bromo',
    name: 'Bromo Midnight Tour',
    route: 'Sunrise Point',
    duration: '1 hari',
    price: 450000,
    meeting_point: 'Hotel Probolinggo / Malang, 00.00 WIB',
    meeting_map: { lat: -7.9425, lng: 112.9530 },
    includes: ['Jeep 4WD', 'Breakfast', 'Guide lokal', 'Masker', 'Dokumentasi'],
    addons: [
      { id: 'horse', name: 'Sewa Kuda Bromo', price: 150000, icon: '🐴' },
      { id: 'jacket', name: 'Jacket Tebal', price: 30000, icon: '🧥' }
    ],
    schedules: [
      { id: 'sched-b1-1', date_start: '2026-05-14', date_end: '2026-05-14', quota_total: 40, quota_remaining: 15 },
      { id: 'sched-b1-2', date_start: '2026-05-16', date_end: '2026-05-16', quota_total: 40, quota_remaining: 28 },
      { id: 'sched-b1-3', date_start: '2026-05-18', date_end: '2026-05-18', quota_total: 40, quota_remaining: 35 },
      { id: 'sched-b1-4', date_start: '2026-05-21', date_end: '2026-05-21', quota_total: 40, quota_remaining: 22 },
      { id: 'sched-b1-5', date_start: '2026-05-23', date_end: '2026-05-23', quota_total: 40, quota_remaining: 38 },
      { id: 'sched-b1-6', date_start: '2026-05-25', date_end: '2026-05-25', quota_total: 40, quota_remaining: 19 },
      { id: 'sched-b1-7', date_start: '2026-05-28', date_end: '2026-05-28', quota_total: 40, quota_remaining: 30 },
      { id: 'sched-b1-8', date_start: '2026-05-30', date_end: '2026-05-30', quota_total: 40, quota_remaining: 25 },
      { id: 'sched-b1-9', date_start: '2026-06-01', date_end: '2026-06-01', quota_total: 40, quota_remaining: 33 },
      { id: 'sched-b1-10', date_start: '2026-06-04', date_end: '2026-06-04', quota_total: 40, quota_remaining: 17 }
    ]
  },
  // --- BROMO - Hiking Buddies ---
  {
    id: 'trip-bromo-2',
    operator_id: 'hikingbuddies',
    mountain_id: 'bromo',
    name: 'Bromo Sunrise Fun Trip',
    route: 'Sunrise Point',
    duration: '1 hari',
    price: 400000,
    meeting_point: 'Hotel Malang, 00.00 WIB',
    meeting_map: { lat: -7.9425, lng: 112.9530 },
    includes: ['Jeep 4WD', 'Breakfast', 'Guide', 'Masker', 'Snack', 'Dokumentasi'],
    addons: [
      { id: 'horse', name: 'Sewa Kuda', price: 125000, icon: '🐴' },
      { id: 'jacket', name: 'Jacket', price: 25000, icon: '🧥' }
    ],
    schedules: [
      { id: 'sched-b2-1', date_start: '2026-05-15', date_end: '2026-05-15', quota_total: 35, quota_remaining: 28 },
      { id: 'sched-b2-2', date_start: '2026-05-17', date_end: '2026-05-17', quota_total: 35, quota_remaining: 30 },
      { id: 'sched-b2-3', date_start: '2026-05-19', date_end: '2026-05-19', quota_total: 35, quota_remaining: 22 },
      { id: 'sched-b2-4', date_start: '2026-05-22', date_end: '2026-05-22', quota_total: 35, quota_remaining: 15 },
      { id: 'sched-b2-5', date_start: '2026-05-24', date_end: '2026-05-24', quota_total: 35, quota_remaining: 32 },
      { id: 'sched-b2-6', date_start: '2026-05-26', date_end: '2026-05-26', quota_total: 35, quota_remaining: 20 },
      { id: 'sched-b2-7', date_start: '2026-05-29', date_end: '2026-05-29', quota_total: 35, quota_remaining: 27 },
      { id: 'sched-b2-8', date_start: '2026-06-02', date_end: '2026-06-02', quota_total: 35, quota_remaining: 18 },
      { id: 'sched-b2-9', date_start: '2026-06-05', date_end: '2026-06-05', quota_total: 35, quota_remaining: 35 },
      { id: 'sched-b2-10', date_start: '2026-06-07', date_end: '2026-06-07', quota_total: 35, quota_remaining: 12 }
    ]
  },

  // --- GEDE (Jawa Barat) - Komunitas ---
  {
    id: 'trip-gede-1',
    operator_id: 'komunitas',
    mountain_id: 'gede',
    name: 'Explore Gede Pangrango',
    route: 'Cibodas',
    duration: '2 hari 1 malam',
    price: 650000,
    meeting_point: 'Kebun Raya Cibodas, 07.00 WIB',
    meeting_map: { lat: -6.7300, lng: 106.9800 },
    includes: ['Transportasi dari Bogor', 'Makan (3x)', 'Tenda', 'Guide', 'Edukasi flora fauna', 'Dokumentasi'],
    addons: [
      { id: 'pole', name: 'Tracking Pole', price: 20000, icon: '🥾' },
      { id: 'carrier', name: 'Carrier 60L', price: 45000, icon: '🎒' }
    ],
    schedules: [
      { id: 'sched-g1-1', date_start: '2026-05-20', date_end: '2026-05-21', quota_total: 15, quota_remaining: 12 },
      { id: 'sched-g1-2', date_start: '2026-05-23', date_end: '2026-05-24', quota_total: 15, quota_remaining: 14 },
      { id: 'sched-g1-3', date_start: '2026-05-27', date_end: '2026-05-28', quota_total: 15, quota_remaining: 9 },
      { id: 'sched-g1-4', date_start: '2026-05-30', date_end: '2026-05-31', quota_total: 15, quota_remaining: 15 },
      { id: 'sched-g1-5', date_start: '2026-06-03', date_end: '2026-06-04', quota_total: 15, quota_remaining: 11 },
      { id: 'sched-g1-6', date_start: '2026-06-06', date_end: '2026-06-07', quota_total: 15, quota_remaining: 7 },
      { id: 'sched-g1-7', date_start: '2026-06-10', date_end: '2026-06-11', quota_total: 15, quota_remaining: 13 },
      { id: 'sched-g1-8', date_start: '2026-06-13', date_end: '2026-06-14', quota_total: 15, quota_remaining: 8 },
      { id: 'sched-g1-9', date_start: '2026-06-17', date_end: '2026-06-18', quota_total: 15, quota_remaining: 10 },
      { id: 'sched-g1-10', date_start: '2026-06-20', date_end: '2026-06-21', quota_total: 15, quota_remaining: 6 }
    ]
  },

  // --- WUKIR (DIY) - Hiking Buddies ---
  {
    id: 'trip-wukir-1',
    operator_id: 'hikingbuddies',
    mountain_id: 'wukir',
    name: 'Sunset Bukit Wukir',
    route: 'Puncak Wukir',
    duration: '1 hari',
    price: 250000,
    meeting_point: 'Kampus UGM, 14.00 WIB',
    meeting_map: { lat: -7.7700, lng: 110.3800 },
    includes: ['Transportasi dari Jogja', 'Snack & minum', 'Guide', 'Dokumentasi', 'Games'],
    addons: [
      { id: 'mat', name: 'Matras', price: 5000, icon: '🛏️' }
    ],
    schedules: [
      { id: 'sched-w1-1', date_start: '2026-05-18', date_end: '2026-05-18', quota_total: 25, quota_remaining: 20 },
      { id: 'sched-w1-2', date_start: '2026-05-21', date_end: '2026-05-21', quota_total: 25, quota_remaining: 22 },
      { id: 'sched-w1-3', date_start: '2026-05-25', date_end: '2026-05-25', quota_total: 25, quota_remaining: 18 },
      { id: 'sched-w1-4', date_start: '2026-05-28', date_end: '2026-05-28', quota_total: 25, quota_remaining: 15 },
      { id: 'sched-w1-5', date_start: '2026-06-01', date_end: '2026-06-01', quota_total: 25, quota_remaining: 23 },
      { id: 'sched-w1-6', date_start: '2026-06-04', date_end: '2026-06-04', quota_total: 25, quota_remaining: 19 },
      { id: 'sched-w1-7', date_start: '2026-06-08', date_end: '2026-06-08', quota_total: 25, quota_remaining: 12 },
      { id: 'sched-w1-8', date_start: '2026-06-11', date_end: '2026-06-11', quota_total: 25, quota_remaining: 25 },
      { id: 'sched-w1-9', date_start: '2026-06-15', date_end: '2026-06-15', quota_total: 25, quota_remaining: 9 },
      { id: 'sched-w1-10', date_start: '2026-06-18', date_end: '2026-06-18', quota_total: 25, quota_remaining: 17 }
    ]
  }
];

// ============================================
// BACKWARD COMPATIBILITY: Add single schedule properties
// ============================================

TRIPS.forEach(trip => {
  const firstSchedule = trip.schedules[0];
  if (firstSchedule) {
    trip.date_start = firstSchedule.date_start;
    trip.date_end = firstSchedule.date_end;
    trip.quota_total = firstSchedule.quota_total;
    trip.quota_remaining = firstSchedule.quota_remaining;
  }
});

// ============================================
// CALENDAR HELPERS
// ============================================

function getTripSchedules(tripId) {
  const trip = getTripById(tripId);
  return trip ? trip.schedules : [];
}

function getScheduleById(tripId, scheduleId) {
  const schedules = getTripSchedules(tripId);
  return schedules.find(s => s.id === scheduleId);
}

function getSchedulesByMonth(tripId, year, month) {
  const schedules = getTripSchedules(tripId);
  return schedules.filter(s => {
    const date = new Date(s.date_start);
    return date.getFullYear() === year && date.getMonth() === month;
  });
}

function hasScheduleOnDate(tripId, dateStr) {
  const schedules = getTripSchedules(tripId);
  const checkDate = new Date(dateStr + 'T00:00:00');
  
  return schedules.some(s => {
    const start = new Date(s.date_start + 'T00:00:00');
    const end = new Date(s.date_end + 'T00:00:00');
    return checkDate >= start && checkDate <= end;
  });
}

function getScheduleTypeOnDate(tripId, dateStr) {
  const schedules = getTripSchedules(tripId);
  const checkDate = new Date(dateStr + 'T00:00:00');
  
  for (const s of schedules) {
    const start = new Date(s.date_start + 'T00:00:00');
    const end = new Date(s.date_end + 'T00:00:00');
    
    if (checkDate >= start && checkDate <= end) {
      if (start.getTime() === end.getTime()) return { type: 'single', schedule: s };
      if (checkDate.getTime() === start.getTime()) return { type: 'start', schedule: s };
      if (checkDate.getTime() === end.getTime()) return { type: 'end', schedule: s };
      return { type: 'mid', schedule: s };
    }
  }
  
  return null;
}

function getScheduleByDate(tripId, dateStr) {
  const schedules = getTripSchedules(tripId);
  const checkDate = new Date(dateStr + 'T00:00:00');
  
  return schedules.find(s => {
    const start = new Date(s.date_start + 'T00:00:00');
    const end = new Date(s.date_end + 'T00:00:00');
    return checkDate >= start && checkDate <= end;
  });
}

function getScheduleDates(tripId) {
  const schedules = getTripSchedules(tripId);
  return schedules.map(s => s.date_start);
}

function getAllSchedulesOnDate(tripId, dateStr) {
  const schedules = getTripSchedules(tripId);
  const checkDate = new Date(dateStr + 'T00:00:00');
  
  return schedules.filter(s => {
    const start = new Date(s.date_start + 'T00:00:00');
    const end = new Date(s.date_end + 'T00:00:00');
    return checkDate >= start && checkDate <= end;
  }).map(s => {
    const start = new Date(s.date_start + 'T00:00:00');
    const end = new Date(s.date_end + 'T00:00:00');
    let type = 'mid';
    if (start.getTime() === end.getTime()) type = 'single';
    else if (checkDate.getTime() === start.getTime()) type = 'start';
    else if (checkDate.getTime() === end.getTime()) type = 'end';
    
    return { ...s, position_type: type };
  });
}

function getScheduleCountOnDate(tripId, dateStr) {
  return getAllSchedulesOnDate(tripId, dateStr).length;
}

// ============================================
// UTILITIES
// ============================================

function formatPrice(price) {
  return 'Rp ' + price.toLocaleString('id-ID');
}

function formatDate(dateStr) {
  const options = { day: 'numeric', month: 'long', year: 'numeric' };
  return new Date(dateStr).toLocaleDateString('id-ID', options);
}

function formatShortDate(dateStr) {
  const options = { day: 'numeric', month: 'short' };
  return new Date(dateStr).toLocaleDateString('id-ID', options);
}

function getURLParam(param) {
  const urlParams = new URLSearchParams(window.location.search);
  return urlParams.get(param);
}

function navigateTo(page, params = {}) {
  const queryString = new URLSearchParams(params).toString();
  const separator = queryString ? '?' : '';
  window.location.href = page + separator + queryString;
}

function getMountainById(id) {
  return MOUNTAINS.find(m => m.id === id);
}

function getOperatorById(id) {
  return OPERATORS.find(o => o.id === id);
}

function getTripById(id) {
  return TRIPS.find(t => t.id === id);
}

function getTripsByMountain(mountainId) {
  return TRIPS.filter(t => t.mountain_id === mountainId);
}

function getTripsByOperator(operatorId) {
  return TRIPS.filter(t => t.operator_id === operatorId);
}

function getOperatorByTrip(tripId) {
  const trip = getTripById(tripId);
  return trip ? getOperatorById(trip.operator_id) : null;
}

function getMountainByTrip(tripId) {
  const trip = getTripById(tripId);
  return trip ? getMountainById(trip.mountain_id) : null;
}

// ============================================
// WISHLIST
// ============================================

function getWishlist() {
  return JSON.parse(localStorage.getItem('wishlist') || '[]');
}

function toggleWishlist(mountainId) {
  let wishlist = getWishlist();
  if (wishlist.includes(mountainId)) {
    wishlist = wishlist.filter(id => id !== mountainId);
  } else {
    wishlist.push(mountainId);
  }
  localStorage.setItem('wishlist', JSON.stringify(wishlist));
  return wishlist.includes(mountainId);
}

function isInWishlist(mountainId) {
  return getWishlist().includes(mountainId);
}

// ============================================
// BOOKINGS
// ============================================

// Dummy bookings for demo
const BOOKINGS = [
  {
    id: 'book-001',
    trip_id: 'trip-ciremai-1',
    schedule_id: 'sched-c1-m1',
    lead_name: 'Budi Santoso',
    lead_phone: '0812-3456-7890',
    lead_email: 'budi@email.com',
    participants: [
      { name: 'Budi Santoso', ktp: '3171234567890001' },
      { name: 'Ani Wulandari', ktp: '3171234567890002' },
      { name: 'Dedi Pratama', ktp: '3171234567890003' }
    ],
    trip_price: 850000,
    addons: [],
    total: 2550000,
    status: 'confirmed',
    created_at: '2026-04-15'
  },
  {
    id: 'book-002',
    trip_id: 'trip-ciremai-1',
    schedule_id: 'sched-c1-m2',
    lead_name: 'Rina Susanti',
    lead_phone: '0813-9876-5432',
    lead_email: 'rina@email.com',
    participants: [
      { name: 'Rina Susanti', ktp: '3171234567890004' },
      { name: 'Agus Wijaya', ktp: '3171234567890005' }
    ],
    trip_price: 850000,
    addons: [
      { id: 'pole', name: 'Tracking Pole', price: 25000 }
    ],
    total: 1750000,
    status: 'pending',
    created_at: '2026-04-20'
  },
  {
    id: 'book-003',
    trip_id: 'trip-ciremai-1',
    schedule_id: 'sched-c1-m7',
    lead_name: 'Dewi Kusuma',
    lead_phone: '0815-1122-3344',
    lead_email: 'dewi@email.com',
    participants: [
      { name: 'Dewi Kusuma', ktp: '3171234567890006' },
      { name: 'Bayu Aji', ktp: '3171234567890007' },
      { name: 'Citra Lestari', ktp: '3171234567890008' },
      { name: 'Fajar Maulana', ktp: '3171234567890009' }
    ],
    trip_price: 850000,
    addons: [
      { id: 'carrier', name: 'Carrier 60L', price: 50000 }
    ],
    total: 3450000,
    status: 'completed',
    created_at: '2026-04-25'
  },
  {
    id: 'book-004',
    trip_id: 'trip-ciremai-2',
    schedule_id: 'sched-c2-1',
    lead_name: 'Eko Prasetyo',
    lead_phone: '0816-5544-3322',
    lead_email: 'eko@email.com',
    participants: [
      { name: 'Eko Prasetyo', ktp: '3171234567890010' }
    ],
    trip_price: 850000,
    addons: [],
    total: 850000,
    status: 'pending',
    created_at: '2026-04-28'
  },
  {
    id: 'book-005',
    trip_id: 'trip-papandayan-1',
    schedule_id: 'sched-p1-1',
    lead_name: 'Fitriani',
    lead_phone: '0817-7788-9900',
    lead_email: 'fitri@email.com',
    participants: [
      { name: 'Fitriani', ktp: '3171234567890011' },
      { name: 'Hendra Gunawan', ktp: '3171234567890012' }
    ],
    trip_price: 650000,
    addons: [],
    total: 1300000,
    status: 'confirmed',
    created_at: '2026-04-22'
  },
  {
    id: 'book-006',
    trip_id: 'trip-merapi-1',
    schedule_id: 'sched-m1-1',
    lead_name: 'Irfan Hakim',
    lead_phone: '0818-1234-5678',
    lead_email: 'irfan@email.com',
    participants: [
      { name: 'Irfan Hakim', ktp: '3171234567890013' },
      { name: 'Joko Widodo', ktp: '3171234567890014' },
      { name: 'Kartini Sari', ktp: '3171234567890015' }
    ],
    trip_price: 550000,
    addons: [
      { id: 'pole', name: 'Tracking Pole', price: 20000 }
    ],
    total: 1710000,
    status: 'pending',
    created_at: '2026-05-01'
  },
  {
    id: 'book-007',
    trip_id: 'trip-semeru-1',
    schedule_id: 'sched-s1-1',
    lead_name: 'Lina Marlina',
    lead_phone: '0819-8765-4321',
    lead_email: 'lina@email.com',
    participants: [
      { name: 'Lina Marlina', ktp: '3171234567890016' },
      { name: 'Maman Supratman', ktp: '3171234567890017' },
      { name: 'Nani Wijaya', ktp: '3171234567890018' },
      { name: 'Omar Faruq', ktp: '3171234567890019' },
      { name: 'Putri Ayu', ktp: '3171234567890020' }
    ],
    trip_price: 1350000,
    addons: [
      { id: 'carrier', name: 'Carrier 80L', price: 75000 }
    ],
    total: 6825000,
    status: 'confirmed',
    created_at: '2026-04-30'
  },
  {
    id: 'book-008',
    trip_id: 'trip-ciremai-1',
    schedule_id: 'sched-c1-m8',
    lead_name: 'Qori Amalia',
    lead_phone: '0821-1111-2222',
    lead_email: 'qori@email.com',
    participants: [
      { name: 'Qori Amalia', ktp: '3171234567890021' },
      { name: 'Rudi Hartono', ktp: '3171234567890022' }
    ],
    trip_price: 850000,
    addons: [],
    total: 1700000,
    status: 'completed',
    created_at: '2026-04-18'
  },
  {
    id: 'book-009',
    trip_id: 'trip-prau-1',
    schedule_id: 'sched-pr1-1',
    lead_name: 'Sinta Dewi',
    lead_phone: '0822-3333-4444',
    lead_email: 'sinta@email.com',
    participants: [
      { name: 'Sinta Dewi', ktp: '3171234567890023' },
      { name: 'Toni Sujono', ktp: '3171234567890024' },
      { name: 'Umi Kalsum', ktp: '3171234567890025' }
    ],
    trip_price: 350000,
    addons: [],
    total: 1050000,
    status: 'cancelled',
    created_at: '2026-04-10'
  },
  {
    id: 'book-010',
    trip_id: 'trip-ciremai-1',
    schedule_id: 'sched-c1-m3',
    lead_name: 'Vino Bastian',
    lead_phone: '0823-5555-6666',
    lead_email: 'vino@email.com',
    participants: [
      { name: 'Vino Bastian', ktp: '3171234567890026' }
    ],
    trip_price: 850000,
    addons: [
      { id: 'sleeping_bag', name: 'Sleeping Bag', price: 35000 }
    ],
    total: 885000,
    status: 'confirmed',
    created_at: '2026-05-02'
  }
];

function getBookings() {
  const stored = localStorage.getItem('bookings');
  if (stored) {
    return JSON.parse(stored);
  }
  // Return dummy data for demo
  return BOOKINGS;
}

function addBooking(booking) {
  const bookings = getBookings();
  booking.id = Date.now();
  booking.status = 'Menunggu Konfirmasi';
  booking.createdAt = new Date().toISOString();
  bookings.push(booking);
  localStorage.setItem('bookings', JSON.stringify(bookings));
  return booking;
}

// ============================================
// TOAST NOTIFICATION
// ============================================

function showToast(message, type = 'success') {
  const toast = document.createElement('div');
  toast.className = `fixed top-4 left-4 right-4 z-50 p-4 rounded-xl text-center text-sm font-medium transform transition-all duration-300 translate-y-[-100%] opacity-0 ${
    type === 'success' ? 'bg-emerald-600 text-white' : 'bg-red-500 text-white'
  }`;
  toast.textContent = message;
  document.body.appendChild(toast);
  
  setTimeout(() => {
    toast.classList.remove('translate-y-[-100%]', 'opacity-0');
  }, 100);
  
  setTimeout(() => {
    toast.classList.add('translate-y-[-100%]', 'opacity-0');
    setTimeout(() => toast.remove(), 300);
  }, 3000);
}

// ============================================
// BOTTOM NAVIGATION
// ============================================

function renderBottomNav(activePage) {
  const nav = document.getElementById('bottom-nav');
  if (!nav) return;
  
  const items = [
    { page: 'index.html', icon: '🏠', label: 'Home' },
    { page: 'my-trips.html', icon: '🏔️', label: 'My Trips' },
    { page: 'profile.html', icon: '👤', label: 'Profile' }
  ];
  
  nav.innerHTML = `
    <div class="flex justify-around items-center w-full h-full">
      ${items.map(item => {
        const isActive = window.location.pathname.includes(item.page) || 
                        (item.page === 'index.html' && window.location.pathname.endsWith('/')) ||
                        (item.page === 'index.html' && activePage === 'home');
        return `
          <button onclick="navigateTo('${item.page}')" 
                  class="flex-1 flex flex-col items-center justify-center py-2 ${isActive ? 'text-orange-500' : 'text-slate-400'}">
            <span class="text-xl mb-0.5">${item.icon}</span>
            <span class="text-[10px] font-medium">${item.label}</span>
          </button>
        `;
      }).join('')}
    </div>
  `;
}

// ============================================
// SEARCH & FILTER MOUNTAINS
// ============================================

function filterMountains(query = '', province = '', difficulty = '') {
  return MOUNTAINS.filter(m => {
    const matchQuery = !query || m.name.toLowerCase().includes(query.toLowerCase());
    const matchProvince = !province || m.province === province;
    const matchDifficulty = !difficulty || m.difficulty === difficulty;
    return matchQuery && matchProvince && matchDifficulty;
  });
}

function renderSearchBar(containerId, onSearch, placeholder = 'Cari gunung...') {
  const container = document.getElementById(containerId);
  if (!container) return;
  
  container.innerHTML = `
    <div class="relative">
      <input type="text" 
             id="search-input"
             placeholder="${placeholder}" 
             class="w-full bg-white/20 backdrop-blur-sm text-white placeholder-white/70 rounded-xl px-4 py-3 pl-11 text-sm focus:outline-none focus:ring-2 focus:ring-orange-500/50"
             ${onSearch ? '' : 'readonly'}
             ${onSearch ? '' : 'onclick="navigateTo(\'search-results.html\')"'}
             style="${onSearch ? '' : 'cursor: pointer;'}"
      >
      <span class="absolute left-3.5 top-1/2 -translate-y-1/2 text-white/70">🔍</span>
    </div>
  `;
  
  if (onSearch) {
    const input = document.getElementById('search-input');
    let debounceTimer;
    input.addEventListener('input', (e) => {
      clearTimeout(debounceTimer);
      debounceTimer = setTimeout(() => onSearch(e.target.value), 300);
    });
  }
}

// ============================================
// MAP EMBED (Google Maps Static)
// ============================================

function renderMap(containerId, location, zoom = 13) {
  const container = document.getElementById(containerId);
  if (!container) return;
  
  const mapUrl = `https://maps.google.com/maps?q=${location.lat},${location.lng}&z=${zoom}&output=embed`;
  container.innerHTML = `
    <iframe 
      src="${mapUrl}" 
      width="100%" 
      height="200" 
      style="border:0; border-radius: 12px;" 
      allowfullscreen="" 
      loading="lazy"
      class="w-full rounded-xl"
    ></iframe>
  `;
}

// Initialize on DOM ready
document.addEventListener('DOMContentLoaded', () => {
  renderBottomNav();
});
