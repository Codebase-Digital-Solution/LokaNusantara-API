package db

import (
	"api/entity"
	"log"
	"time"

	"gorm.io/gorm"
)

// Pawukon is the 210-day Balinese cycle (30 wuku of 7 days). Its holy days
// repeat exactly every cycle, so they can be computed from one known date.
const pawukonDays = 210

// pawukonEpoch is Redite (Sunday) of wuku Sinta, the first day of a cycle:
// 73 days before Galungan on 28 February 2024.
var pawukonEpoch = time.Date(2023, time.December, 17, 0, 0, 0, 0, time.Local)

type pawukonHoliday struct {
	name        string
	offset      int // day within the cycle: wuku index * 7 + weekday
	description string
}

var pawukonHolidays = []pawukonHoliday{
	{"Tumpek Landep", 13, "Hari untuk mengupacarai benda-benda dari besi dan logam seperti keris dan peralatan kerja, kini juga kendaraan, sebagai wujud syukur atas ketajaman pikiran."},
	{"Tumpek Wariga", 48, "Dikenal juga sebagai Tumpek Uduh atau Tumpek Bubuh. Umat mengupacarai pepohonan dengan bubur sumsum agar berbuah lebat menjelang Galungan, 25 hari kemudian."},
	{"Hari Raya Galungan", 73, "Hari kemenangan dharma melawan adharma. Penjor bambu melengkung menghiasi sepanjang jalan, dan umat Hindu bersembahyang ke pura serta merajan keluarga."},
	{"Hari Raya Kuningan", 83, "Sepuluh hari setelah Galungan, hari ketika leluhur dipercaya kembali ke alam dewata. Umat menghaturkan nasi kuning dan memasang hiasan tamiang serta endongan."},
	{"Hari Raya Saraswati", 209, "Hari turunnya ilmu pengetahuan, dipersembahkan kepada Dewi Saraswati. Buku dan lontar diupacarai, dan sekolah serta pura ramai oleh persembahyangan."},
	{"Hari Raya Pagerwesi", 3, "Hari untuk memperkuat benteng diri dengan ilmu pengetahuan. Dirayakan sangat meriah di Buleleng, hampir semeriah Galungan."},
}

type fixedEvent struct {
	name, description, category, region, location string
	start, end                                    string // YYYY-MM-DD; end optional
	estimate                                      bool
}

// Dates marked estimate follow a lunar calendar or an organiser's decision;
// update them once official dates are announced.
var fixedEvents = []fixedEvent{
	{"Hari Raya Natal", "Misa Natal di desa Katolik Palasari, dengan gereja berarsitektur Bali yang dihiasi penjor dan ornamen khas.", "gereja", "Jembrana", "Gereja Hati Kudus Yesus Palasari", "2026-12-25", "", false},
	{"Tahun Baru Imlek", "Perayaan Tahun Baru Imlek 2578 dengan sembahyang di klenteng, lampion merah, dan pertunjukan barongsai.", "klenteng", "Badung", "Vihara Dharmayana Kuta", "2027-02-06", "", false},
	{"Cap Go Meh", "Hari kelima belas setelah Imlek yang menutup rangkaian perayaan Tahun Baru, dengan sembahyang bersama umat Tionghoa dan Hindu.", "klenteng", "Denpasar", "Griya Kongco Dwipayana Tanah Kilap", "2027-02-20", "", false},
	{"Pawai Ogoh-ogoh", "Malam Pengerupukan sehari sebelum Nyepi. Patung raksasa ogoh-ogoh diarak keliling banjar sebagai simbol mengusir unsur negatif.", "tradisi", "Seluruh Bali", "", "2027-03-07", "", true},
	{"Hari Raya Nyepi", "Tahun Baru Saka. Seluruh Bali hening selama 24 jam: tanpa api, tanpa bepergian, tanpa hiburan, dan bandara ditutup.", "pura", "Seluruh Bali", "", "2027-03-08", "", true},
	{"Omed-omedan", "Tradisi sekaa teruna Banjar Kaja Sesetan saat Ngembak Geni, sehari setelah Nyepi, yang mempererat persaudaraan warga.", "tradisi", "Denpasar", "Banjar Kaja Sesetan", "2027-03-09", "", true},
	{"Hari Raya Idul Fitri", "Salat Id dan halalbihalal di kampung Muslim tertua di Bali, tempat umat Islam dan Hindu hidup berdampingan sejak abad ke-14.", "masjid", "Klungkung", "Masjid Nurul Huda Gelgel", "2027-03-10", "", true},
	{"Hari Raya Paskah", "Perayaan Paskah di gereja Katolik pertama di Bali, dengan ibadah bernuansa budaya Bali.", "gereja", "Badung", "Gereja Tritunggal Mahakudus Tuka", "2027-03-28", "", false},
	{"Hari Raya Idul Adha", "Salat Id dan penyembelihan hewan kurban di kompleks Puja Mandala yang menjadi simbol kerukunan umat beragama.", "masjid", "Badung", "Masjid Agung Ibnu Batutah", "2027-05-17", "", true},
	{"Hari Raya Waisak", "Peringatan kelahiran, pencerahan, dan wafatnya Buddha Gautama dengan puja bakti dan meditasi bersama.", "vihara", "Buleleng", "Brahmavihara-Arama", "2027-05-20", "", true},
	{"Perang Pandan (Mekare-kare)", "Ritual adat Tenganan untuk menghormati Dewa Indra: para pemuda saling bertarung dengan daun pandan berduri. Tanggal pasti ditentukan desa adat, biasanya sekitar Juni.", "tradisi", "Karangasem", "Desa Tenganan Pegringsingan", "2027-06-01", "2027-06-30", true},
	{"Pesta Kesenian Bali", "Festival seni tahunan selama sebulan dengan tari, gamelan, drama, dan pawai budaya dari seluruh kabupaten di Bali.", "tradisi", "Denpasar", "Taman Werdhi Budaya Art Centre", "2027-06-12", "2027-07-10", true},
	{"Maulid Nabi: Muludan Pegayaman", "Peringatan Maulid Nabi dengan tradisi sokok, rangkaian hiasan telur dan bunga mirip gebogan Bali, yang diarak ke masjid.", "masjid", "Buleleng", "Masjid Jami Safinatussalam Pegayaman", "2027-08-15", "", true},
	{"Hari Raya Natal", "Misa Natal di desa Katolik Palasari, dengan gereja berarsitektur Bali yang dihiasi penjor dan ornamen khas.", "gereja", "Jembrana", "Gereja Hati Kudus Yesus Palasari", "2027-12-25", "", false},
}

// SeedEvents inserts Pawukon holidays for 2026 and 2027 plus the fixed list
// above. Rows are matched by name and start date, so reruns are safe.
func SeedEvents(db *gorm.DB) {
	from := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.Local)
	to := time.Date(2028, time.January, 1, 0, 0, 0, 0, time.Local)

	var events []entity.Event
	for _, h := range pawukonHolidays {
		for _, day := range pawukonDates(h.offset, from, to) {
			events = append(events, entity.Event{
				Name:        h.name,
				Description: h.description,
				Category:    "pura",
				StartDate:   day,
				Region:      "Seluruh Bali",
			})
		}
	}

	for _, f := range fixedEvents {
		start, err := time.ParseInLocation("2006-01-02", f.start, time.Local)
		if err != nil {
			log.Printf("Invalid date for event '%s': %v", f.name, err)
			continue
		}
		event := entity.Event{
			Name:        f.name,
			Description: f.description,
			Category:    f.category,
			StartDate:   start,
			IsEstimate:  f.estimate,
			Region:      f.region,
		}
		if f.end != "" {
			end, err := time.ParseInLocation("2006-01-02", f.end, time.Local)
			if err == nil {
				event.EndDate = &end
			}
		}
		if f.location != "" {
			var loc entity.Location
			if err := db.Select("id").Where("name = ?", f.location).First(&loc).Error; err == nil {
				event.LocationID = &loc.ID
			} else {
				log.Printf("Location '%s' for event '%s' not found", f.location, f.name)
			}
		}
		events = append(events, event)
	}

	for _, e := range events {
		var count int64
		db.Model(&entity.Event{}).
			Where("name = ? AND start_date = ?", e.Name, e.StartDate).
			Count(&count)
		if count > 0 {
			continue
		}
		if err := db.Create(&e).Error; err != nil {
			log.Printf("Failed to seed event '%s': %v", e.Name, err)
		}
	}
}

// pawukonDates returns every date in [from, to) that falls on the given day
// of the Pawukon cycle.
func pawukonDates(offset int, from, to time.Time) []time.Time {
	first := pawukonEpoch.AddDate(0, 0, offset)
	// Step back or forward whole cycles to the first occurrence >= from.
	days := int(from.Sub(first).Hours() / 24)
	cycles := days / pawukonDays
	if days < 0 {
		cycles--
	}
	day := first.AddDate(0, 0, cycles*pawukonDays)
	for day.Before(from) {
		day = day.AddDate(0, 0, pawukonDays)
	}

	var dates []time.Time
	for ; day.Before(to); day = day.AddDate(0, 0, pawukonDays) {
		dates = append(dates, day)
	}
	return dates
}
