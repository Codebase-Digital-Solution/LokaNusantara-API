package db

import (
	"api/entity"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SeedLocations inserts sample LokaNusantara locations in Bali, one batch
// per category. Existing rows (matched by name) are left untouched.
func SeedLocations(db *gorm.DB) {
	seedTradisi(db)
	seedPura(db)
	seedMasjid(db)
	seedGereja(db)
	seedVihara(db)
	seedKlenteng(db)
	seedKuliner(db)
}

func findCategory(db *gorm.DB, name string) (entity.Category, bool) {
	var category entity.Category
	if err := db.Where("name = ?", name).First(&category).Error; err != nil {
		log.Printf("Category '%s' not found, skipping", name)
		return category, false
	}
	return category, true
}

func seedTradisi(db *gorm.DB) {
	category, ok := findCategory(db, "tradisi")
	if !ok {
		return
	}

	locations := []entity.Location{
		{ID: uuid.NewString(), Name: "Desa Adat Penglipuran", Description: "Desa adat Bali Aga dengan tata ruang tradisional Tri Hita Karana, deretan rumah berpintu angkul-angkul seragam, dan hutan bambu sakral.", Latitude: -8.4222, Longitude: 115.3592, Address: "Kubu, Bangli, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Desa Tenganan Pegringsingan", Description: "Desa Bali Aga penjaga tenun ikat ganda gringsing dan tradisi Perang Pandan (Mekare-kare) setiap tahun.", Latitude: -8.4785, Longitude: 115.5664, Address: "Tenganan, Manggis, Karangasem, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Desa Trunyan", Description: "Desa Bali Aga di tepi Danau Batur dengan tradisi pemakaman mepasah di bawah pohon Taru Menyan.", Latitude: -8.2374, Longitude: 115.4134, Address: "Trunyan, Kintamani, Bangli, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Subak Jatiluwih", Description: "Lanskap sawah terasering dengan sistem irigasi tradisional subak yang diakui sebagai Warisan Budaya Dunia UNESCO.", Latitude: -8.3705, Longitude: 115.1319, Address: "Jatiluwih, Penebel, Tabanan, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Puri Saren Agung Ubud", Description: "Istana kerajaan Ubud yang menggelar pertunjukan tari Legong dan Barong setiap malam di pelatarannya.", Latitude: -8.5069, Longitude: 115.2625, Address: "Jl. Raya Ubud, Ubud, Gianyar, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Kertha Gosa Klungkung", Description: "Balai pengadilan Kerajaan Klungkung dengan langit-langit lukisan wayang gaya Kamasan.", Latitude: -8.5355, Longitude: 115.4039, Address: "Semarapura, Klungkung, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Desa Kamasan", Description: "Desa seniman lukis wayang klasik gaya Kamasan yang diwariskan turun-temurun sejak zaman Kerajaan Gelgel.", Latitude: -8.555, Longitude: 115.405, Address: "Kamasan, Klungkung, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Taman Werdhi Budaya Art Centre", Description: "Pusat kesenian Bali dan tuan rumah Pesta Kesenian Bali setiap tahun, dengan panggung terbuka Ardha Candra.", Latitude: -8.6536, Longitude: 115.2336, Address: "Jl. Nusa Indah, Denpasar, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Museum Bali", Description: "Museum etnografi tertua di Bali dengan koleksi pakaian adat, keris, topeng, dan perlengkapan upacara.", Latitude: -8.6567, Longitude: 115.2183, Address: "Jl. Mayor Wisnu, Denpasar, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Garuda Wisnu Kencana", Description: "Taman budaya dengan patung Garuda Wisnu Kencana setinggi 121 meter dan pertunjukan tari Bali harian.", Latitude: -8.8104, Longitude: 115.1676, Address: "Ungasan, Kuta Selatan, Badung, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Desa Celuk", Description: "Sentra kerajinan perak dan emas Bali dengan teknik filigri tradisional.", Latitude: -8.6065, Longitude: 115.2697, Address: "Celuk, Sukawati, Gianyar, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Desa Mas", Description: "Desa seniman ukir kayu yang karyanya menghiasi pura dan puri di seluruh Bali.", Latitude: -8.547, Longitude: 115.262, Address: "Mas, Ubud, Gianyar, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Desa Sidemen", Description: "Desa penenun kain songket dan endek Bali di lembah sawah dengan latar Gunung Agung.", Latitude: -8.485, Longitude: 115.442, Address: "Sidemen, Karangasem, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Banjar Kaja Sesetan", Description: "Tempat tradisi Omed-omedan, ritual saling tarik muda-mudi sehari setelah Nyepi.", Latitude: -8.688, Longitude: 115.2195, Address: "Sesetan, Denpasar Selatan, Bali", CategoryID: category.ID},
	}
	seedBatch(db, locations)
}

func seedPura(db *gorm.DB) {
	category, ok := findCategory(db, "pura")
	if !ok {
		return
	}

	locations := []entity.Location{
		// Bali Selatan & Tengah
		{ID: uuid.NewString(), Name: "Pura Besakih", Description: "Pura terbesar dan paling suci di Bali di lereng Gunung Agung. Dikenal sebagai Pura Ibu Bali.", Latitude: -8.3743, Longitude: 115.4517, Address: "Besakih, Rendang, Karangasem, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Pura Tanah Lot", Description: "Pura Hindu ikonik di atas batu karang di tengah laut. Paling difoto saat matahari terbenam.", Latitude: -8.6210, Longitude: 115.0866, Address: "Beraban, Kediri, Tabanan, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Pura Uluwatu", Description: "Pura di tebing karang setinggi 70 meter di ujung selatan Bali. Salah satu Sad Kahyangan Bali.", Latitude: -8.8291, Longitude: 115.0849, Address: "Pecatu, Kuta Selatan, Badung, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Pura Tirta Empul", Description: "Pura dengan mata air suci untuk ritual melukat, dibangun tahun 962 M di Tampaksiring.", Latitude: -8.4153, Longitude: 115.3150, Address: "Manukaya, Tampaksiring, Gianyar, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Pura Ulun Danu Beratan", Description: "Pura yang tampak mengapung di Danau Beratan Bedugul. Didedikasikan untuk Dewi Danu.", Latitude: -8.2748, Longitude: 115.1677, Address: "Candikuning, Baturiti, Tabanan, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Pura Taman Ayun", Description: "Pura kerajaan Mengwi yang dikelilingi kolam. Dibangun tahun 1634 dan masuk Warisan Budaya Dunia UNESCO.", Latitude: -8.5424, Longitude: 115.1757, Address: "Mengwi, Badung, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Pura Gunung Kawi", Description: "Kompleks candi purbakala abad ke-11 yang dipahat langsung dari tebing batu di lembah Sungai Pakerisan.", Latitude: -8.4197, Longitude: 115.3122, Address: "Tampaksiring, Gianyar, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Pura Kehen", Description: "Pura kerajaan Bangli yang megah dengan gapura bertingkat, pohon beringin raksasa, dan relief indah.", Latitude: -8.4443, Longitude: 115.3566, Address: "Cempaga, Bangli, Bali", CategoryID: category.ID},

		// Bali Utara
		{ID: uuid.NewString(), Name: "Pura Maduwe Karang", Description: "Pura unik di Bali Utara yang terkenal dengan ukiran relief sepeda kuno pada dindingnya. Dibangun tahun 1890.", Latitude: -8.0851, Longitude: 115.1693, Address: "Kubutambahan, Buleleng, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Pura Pulaki", Description: "Pura suci di tepi pantai Bali Utara yang dikelilingi ratusan monyet. Terletak di antara dua tebing laut.", Latitude: -8.1649, Longitude: 114.7842, Address: "Banyupoh, Gerokgak, Buleleng, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Pura Ponjok Batu", Description: "Pura yang menjorok ke laut di ujung tebing Bali Utara. Salah satu pura tepi laut paling dramatis di Bali.", Latitude: -8.0674, Longitude: 115.4047, Address: "Julah, Tejakula, Buleleng, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Pura Ulun Danu Tamblingan", Description: "Pura sakral di tepi Danau Tamblingan yang hanya dapat dikunjungi dengan perahu tradisional melalui hutan.", Latitude: -8.2524, Longitude: 115.1105, Address: "Munduk, Banjar, Buleleng, Bali", CategoryID: category.ID},

		// Bali Barat
		{ID: uuid.NewString(), Name: "Pura Rambut Siwi", Description: "Pura suci di tepi tebing pantai Bali Barat. Berdiri kokoh menghadap Samudra Hindia, didirikan oleh Dang Hyang Nirartha.", Latitude: -8.3570, Longitude: 114.8208, Address: "Yehembang, Mendoyo, Jembrana, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Pura Gede Perancak", Description: "Pura bersejarah di muara sungai Jembrana tempat Dang Hyang Nirartha pertama kali menginjakkan kaki di Bali.", Latitude: -8.4079, Longitude: 114.5738, Address: "Perancak, Negara, Jembrana, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Pura Luhur Batukaru", Description: "Pura suci di lereng Gunung Batukaru, dikelilingi hutan tropis lebat. Salah satu Sad Kahyangan Bali.", Latitude: -8.3872, Longitude: 115.1245, Address: "Wongaya Gede, Penebel, Tabanan, Bali", CategoryID: category.ID},

		// Bali Timur
		{ID: uuid.NewString(), Name: "Pura Lempuyang Luhur", Description: "Pura tertinggi di Bali di lereng Gunung Lempuyang. Terkenal dengan Gerbang Surga (Gates of Heaven) yang ikonik.", Latitude: -8.3921, Longitude: 115.6287, Address: "Lemukih, Abang, Karangasem, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Pura Goa Lawah", Description: "Pura gua kelelawar di tepi pantai timur Bali. Dihuni ribuan kelelawar, salah satu Sad Kahyangan Bali.", Latitude: -8.5445, Longitude: 115.4556, Address: "Pesinggahan, Dawan, Klungkung, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Pura Pasar Agung Besakih", Description: "Pura di lereng Gunung Agung yang menjadi gerbang pendakian spiritual. Terletak di ketinggian 1700 mdpl.", Latitude: -8.3219, Longitude: 115.5063, Address: "Selat, Karangasem, Bali", CategoryID: category.ID},

		// Nusa Penida
		{ID: uuid.NewString(), Name: "Pura Goa Giri Putri", Description: "Pura di dalam gua batu raksasa di Nusa Penida. Salah satu pura paling sakral dengan stalagmit dan stalagtit.", Latitude: -8.6940, Longitude: 115.4857, Address: "Suana, Nusa Penida, Klungkung, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Pura Dalem Penataran Ped", Description: "Pura utama di Nusa Penida yang dipercaya sebagai rumah Ratu Gede Mecaling, dewa penjaga pulau.", Latitude: -8.7001, Longitude: 115.5341, Address: "Ped, Nusa Penida, Klungkung, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Pura Segara Kidul Nusa Penida", Description: "Pura di tepi laut Nusa Penida yang menghadap Samudra Hindia dengan pemandangan alam yang dramatis.", Latitude: -8.7694, Longitude: 115.4291, Address: "Nusa Penida, Klungkung, Bali", CategoryID: category.ID},
	}
	seedBatch(db, locations)
}

func seedMasjid(db *gorm.DB) {
	category, ok := findCategory(db, "masjid")
	if !ok {
		return
	}

	locations := []entity.Location{
		{ID: uuid.NewString(), Name: "Masjid Agung Ibnu Batutah", Description: "Masjid di kompleks Puja Mandala Nusa Dua yang berdiri berdampingan dengan gereja, vihara, dan pura sebagai simbol toleransi.", Latitude: -8.8035, Longitude: 115.2155, Address: "Puja Mandala, Kuta Selatan, Badung, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Masjid Nurul Huda Gelgel", Description: "Masjid tertua di Bali, berdiri sejak masa Kerajaan Gelgel di kampung muslim Gelgel.", Latitude: -8.5613, Longitude: 115.405, Address: "Gelgel, Klungkung, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Masjid Jami Safinatussalam Pegayaman", Description: "Masjid di kampung muslim Pegayaman yang memadukan budaya Bali dan Islam, terkenal dengan tradisi Muludan dan sokok.", Latitude: -8.238, Longitude: 115.089, Address: "Pegayaman, Sukasada, Buleleng, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Masjid Baitul Qadim Loloan", Description: "Masjid bersejarah komunitas Melayu Loloan di tepi Sungai Ijo Gading dengan rumah panggung khas di sekitarnya.", Latitude: -8.364, Longitude: 114.627, Address: "Loloan Timur, Jembrana, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Masjid Agung Jami Singaraja", Description: "Masjid bersejarah di Kampung Kajanan dengan pintu ukir hadiah dari Puri Buleleng.", Latitude: -8.1137, Longitude: 115.087, Address: "Kampung Kajanan, Singaraja, Buleleng, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Masjid Al-Hidayah Bedugul", Description: "Masjid di perbukitan Bedugul yang menghadap Danau Beratan, tidak jauh dari Pasar Candikuning.", Latitude: -8.2808, Longitude: 115.166, Address: "Candikuning, Baturiti, Tabanan, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Masjid Al-Muhajirin Kepaon", Description: "Masjid di Kampung Islam Kepaon yang warganya hidup rukun dengan masyarakat Hindu sejak era Puri Pemecutan.", Latitude: -8.697, Longitude: 115.201, Address: "Kepaon, Denpasar Selatan, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Masjid Agung Sudirman", Description: "Masjid besar di pusat kota Denpasar yang menjadi pusat kegiatan umat muslim.", Latitude: -8.6727, Longitude: 115.2265, Address: "Jl. Jend. Sudirman, Denpasar, Bali", CategoryID: category.ID},
	}
	seedBatch(db, locations)
}

func seedGereja(db *gorm.DB) {
	category, ok := findCategory(db, "gereja")
	if !ok {
		return
	}

	locations := []entity.Location{
		{ID: uuid.NewString(), Name: "Gereja Hati Kudus Yesus Palasari", Description: "Gereja Katolik berarsitektur perpaduan gotik dan Bali di desa Katolik Palasari, dikelilingi bukit dan danau.", Latitude: -8.249, Longitude: 114.596, Address: "Palasari, Melaya, Jembrana, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Gereja Pniel Blimbingsari", Description: "Gereja Kristen Protestan bergaya wantilan Bali di desa Kristen Blimbingsari yang asri.", Latitude: -8.253, Longitude: 114.617, Address: "Blimbingsari, Melaya, Jembrana, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Gereja Tritunggal Mahakudus Tuka", Description: "Gereja Katolik pertama di Bali yang berdiri sejak 1936, dengan ornamen ukiran Bali.", Latitude: -8.6, Longitude: 115.164, Address: "Tuka, Dalung, Kuta Utara, Badung, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Gereja Maria Bunda Segala Bangsa", Description: "Gereja Katolik di kompleks Puja Mandala Nusa Dua, berdampingan dengan rumah ibadah lain.", Latitude: -8.8037, Longitude: 115.2151, Address: "Puja Mandala, Kuta Selatan, Badung, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "GKPB Jemaat Bukit Doa", Description: "Gereja Kristen Protestan di Bali di kompleks Puja Mandala Nusa Dua.", Latitude: -8.8039, Longitude: 115.2147, Address: "Puja Mandala, Kuta Selatan, Badung, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Gereja Katedral Roh Kudus Denpasar", Description: "Katedral Keuskupan Denpasar dengan sentuhan ornamen Bali pada fasad dan interiornya.", Latitude: -8.6736, Longitude: 115.2337, Address: "Jl. Tukad Musi, Renon, Denpasar, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Gereja Santo Fransiskus Xaverius Kuta", Description: "Gereja Katolik di jantung Kuta yang melayani umat lokal dan wisatawan.", Latitude: -8.7224, Longitude: 115.1774, Address: "Jl. Kartika Plaza, Kuta, Badung, Bali", CategoryID: category.ID},
	}
	seedBatch(db, locations)
}

func seedVihara(db *gorm.DB) {
	category, ok := findCategory(db, "vihara")
	if !ok {
		return
	}

	locations := []entity.Location{
		{ID: uuid.NewString(), Name: "Brahmavihara-Arama", Description: "Vihara Buddha terbesar di Bali dengan replika stupa Borobudur, taman meditasi, dan pemandangan laut utara.", Latitude: -8.212, Longitude: 114.964, Address: "Banjar Tegeha, Banjar, Buleleng, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Vihara Buddha Guna", Description: "Vihara di kompleks Puja Mandala Nusa Dua, simbol kerukunan lima agama di Bali.", Latitude: -8.8033, Longitude: 115.2159, Address: "Puja Mandala, Kuta Selatan, Badung, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Vihara Dharma Giri", Description: "Vihara di perbukitan Pupuan dengan patung Buddha berbaring dan suasana pegunungan yang sejuk.", Latitude: -8.344, Longitude: 115.037, Address: "Pupuan, Tabanan, Bali", CategoryID: category.ID},
	}
	seedBatch(db, locations)
}

func seedKlenteng(db *gorm.DB) {
	category, ok := findCategory(db, "klenteng")
	if !ok {
		return
	}

	locations := []entity.Location{
		{ID: uuid.NewString(), Name: "Vihara Dharmayana Kuta", Description: "Salah satu klenteng tertua di Bali, tempat ibadah Tri Dharma yang ramai saat Imlek dan Cap Go Meh.", Latitude: -8.7186, Longitude: 115.1741, Address: "Jl. Blambangan, Kuta, Badung, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Griya Kongco Dwipayana Tanah Kilap", Description: "Klenteng unik tempat umat Tionghoa dan Hindu bersembahyang bersama, cerminan akulturasi budaya.", Latitude: -8.702, Longitude: 115.201, Address: "Tanah Kilap, Denpasar Selatan, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Klenteng Ling Gwan Kiong", Description: "Klenteng bersejarah di kawasan pelabuhan tua Buleleng, saksi jejak pedagang Tionghoa di Bali Utara.", Latitude: -8.1124, Longitude: 115.0883, Address: "Pelabuhan Buleleng, Singaraja, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Vihara Satya Dharma Benoa", Description: "Klenteng Caow Eng Bio di tepi laut Tanjung Benoa, dipercaya didirikan oleh pelaut Tionghoa.", Latitude: -8.759, Longitude: 115.2225, Address: "Tanjung Benoa, Kuta Selatan, Badung, Bali", CategoryID: category.ID},
	}
	seedBatch(db, locations)
}

func seedKuliner(db *gorm.DB) {
	category, ok := findCategory(db, "kuliner")
	if !ok {
		return
	}

	locations := []entity.Location{
		{ID: uuid.NewString(), Name: "Warung Babi Guling Ibu Oka", Description: "Warung babi guling legendaris di Ubud dengan kulit renyah dan bumbu base genep khas Bali.", Latitude: -8.5067, Longitude: 115.2638, Address: "Jl. Tegal Sari No.2, Ubud, Gianyar, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Ayam Betutu Men Tempeh", Description: "Ayam betutu pedas khas Gilimanuk yang dimasak dengan bumbu base genep, legendaris sejak 1978.", Latitude: -8.1715, Longitude: 114.4362, Address: "Gilimanuk, Melaya, Jembrana, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Warung Mak Beng", Description: "Warung sup ikan dan ikan goreng sambal matah di Sanur yang berdiri sejak 1941.", Latitude: -8.6756, Longitude: 115.2634, Address: "Jl. Hang Tuah, Sanur, Denpasar, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Nasi Ayam Kedewatan Ibu Mangku", Description: "Nasi campur ayam khas Bali dengan sate lilit, lawar, dan sambal matah.", Latitude: -8.4953, Longitude: 115.249, Address: "Kedewatan, Ubud, Gianyar, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Warung Wardani", Description: "Nasi campur Bali legendaris di Denpasar dengan lauk sate lilit, ayam suwir, dan sambal khas.", Latitude: -8.6591, Longitude: 115.2185, Address: "Jl. Yudistira, Denpasar, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Pasar Malam Gianyar", Description: "Pasar kuliner malam dengan babi guling, sate, lawar, dan aneka jaje Bali.", Latitude: -8.5395, Longitude: 115.3267, Address: "Jl. Ngurah Rai, Gianyar, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Pasar Badung", Description: "Pasar tradisional terbesar di Bali, tempat berburu jaje Bali, bumbu, dan perlengkapan banten.", Latitude: -8.656, Longitude: 115.2129, Address: "Jl. Gajah Mada, Denpasar, Bali", CategoryID: category.ID},
		{ID: uuid.NewString(), Name: "Ikan Bakar Jimbaran", Description: "Deretan warung ikan bakar bumbu Bali di tepi Pantai Jimbaran, disantap sambil menikmati senja.", Latitude: -8.7748, Longitude: 115.1607, Address: "Pantai Jimbaran, Kuta Selatan, Badung, Bali", CategoryID: category.ID},
	}
	seedBatch(db, locations)
}

func seedBatch(db *gorm.DB, locations []entity.Location) {
	for _, loc := range locations {
		loc.History = locationHistories[loc.Name]

		var count int64
		db.Model(&entity.Location{}).Where("name = ?", loc.Name).Count(&count)
		if count > 0 {
			// Backfill history for rows seeded before the column existed.
			if loc.History != "" {
				db.Model(&entity.Location{}).
					Where("name = ? AND (history IS NULL OR history = '')", loc.Name).
					Update("history", loc.History)
			}
			continue
		}

		if err := db.Create(&loc).Error; err != nil {
			log.Printf("Failed to seed '%s': %v", loc.Name, err)
		} else {
			log.Printf("✓ Seeded: %s", loc.Name)
		}
	}
}
