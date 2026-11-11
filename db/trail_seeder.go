package db

import (
	"api/entity"
	"log"

	"gorm.io/gorm"
)

type trailSeed struct {
	slug, title, subtitle, story, color string
	stops                               []stopSeed
}

type stopSeed struct {
	location, note string
}

var trailSeeds = []trailSeed{
	{
		slug:     "jejak-nirartha",
		title:    "Jejak Dang Hyang Nirartha",
		subtitle: "Menapak perjalanan pendeta besar dari Jawa",
		color:    "#A87618",
		story:    "Pada akhir abad ke-15, Dang Hyang Nirartha menyeberang dari Blambangan ke Bali dan membawa pembaruan besar dalam kehidupan keagamaan Hindu. Jalur ini mengikuti tempat-tempat yang dalam tradisi dikaitkan dengan perjalanannya, dari pantai barat hingga tebing selatan, tempat ia mencapai moksa.",
		stops: []stopSeed{
			{"Pura Gede Perancak", "Titik yang dipercaya sebagai tempat Dang Hyang Nirartha pertama kali menginjakkan kaki di Bali."},
			{"Pura Rambut Siwi", "Di sini ia menyerahkan sehelai rambutnya kepada warga setelah membebaskan desa dari wabah."},
			{"Pura Tanah Lot", "Karang tempat ia beristirahat dan menyarankan warga membangun pura untuk memuja Dewa Baruna."},
			{"Pura Uluwatu", "Pura di ujung selatan Bali, tempat ia dipercaya mencapai moksa."},
			{"Desa Mas", "Desa tempat keturunannya menetap dan tradisi ukir kayu berkembang."},
		},
	},
	{
		slug:     "islam-pulau-dewata",
		title:    "Jejak Islam di Pulau Dewata",
		subtitle: "Kampung Muslim yang hidup berdampingan sejak berabad lalu",
		color:    "#1F7A5C",
		story:    "Islam hadir di Bali sejak abad ke-14 dan tumbuh lewat prajurit, pedagang, dan pelaut yang diterima oleh puri-puri Bali. Jalur ini menyusuri kampung-kampung Muslim tertua, tempat budaya Bali dan Islam berbaur dalam nama, tradisi, dan arsitektur.",
		stops: []stopSeed{
			{"Masjid Nurul Huda Gelgel", "Kampung Muslim tertua di Bali, berawal dari prajurit Majapahit yang mengiringi Raja Gelgel."},
			{"Masjid Al-Muhajirin Kepaon", "Kampung Bugis yang diberi tempat oleh Puri Pemecutan dan tetap dekat dengan pihak puri."},
			{"Masjid Jami Safinatussalam Pegayaman", "Warga memakai nama Wayan dan Nyoman di depan nama Islam, dan merayakan Maulid dengan sokok."},
			{"Masjid Agung Jami Singaraja", "Masjid di kota pelabuhan dengan pintu ukir pemberian Puri Buleleng."},
			{"Masjid Baitul Qadim Loloan", "Komunitas Melayu di tepi sungai yang masih bertutur bahasa Melayu dan tinggal di rumah panggung."},
		},
	},
	{
		slug:     "desa-bali-aga",
		title:    "Desa Bali Aga",
		subtitle: "Bertemu masyarakat Bali asli dan adat tertuanya",
		color:    "#8E3B5E",
		story:    "Bali Aga adalah masyarakat Bali asli yang mempertahankan adat dari masa sebelum pengaruh Majapahit. Jalur ini mengunjungi desa-desa yang menjaga awig-awig, tata ruang, dan ritual leluhur yang tidak ditemukan di tempat lain.",
		stops: []stopSeed{
			{"Desa Tenganan Pegringsingan", "Penjaga kain gringsing ikat ganda dan tradisi Perang Pandan."},
			{"Desa Adat Penglipuran", "Tata ruang Tri Mandala dan hutan bambu yang dijaga aturan adat."},
			{"Desa Trunyan", "Tradisi mepasah di bawah pohon Taru Menyan di tepi Danau Batur."},
			{"Pura Ulun Danu Tamblingan", "Pura yang dijaga masyarakat adat Catur Desa sebagai hulu air kehidupan."},
		},
	},
	{
		slug:     "rumah-ibadah-berdampingan",
		title:    "Rumah Ibadah Berdampingan",
		subtitle: "Toleransi lintas agama di Bali Selatan",
		color:    "#2B5C8A",
		story:    "Di Bali, rumah ibadah berbagai agama berdiri berdampingan dan umatnya saling menjaga. Jalur ini dimulai dari Puja Mandala, tempat lima rumah ibadah berdiri dalam satu kawasan, lalu berlanjut ke klenteng tempat umat Tionghoa dan Hindu bersembahyang bersama.",
		stops: []stopSeed{
			{"Masjid Agung Ibnu Batutah", "Masjid di Puja Mandala, kawasan lima rumah ibadah yang diresmikan pada 1997."},
			{"Gereja Maria Bunda Segala Bangsa", "Gereja Katolik yang berdiri tepat di samping masjid."},
			{"GKPB Jemaat Bukit Doa", "Gereja Protestan di kompleks yang sama, bagian dari tradisi gereja bernuansa Bali."},
			{"Vihara Buddha Guna", "Vihara yang melengkapi kerukunan di Puja Mandala."},
			{"Griya Kongco Dwipayana Tanah Kilap", "Klenteng tempat umat Tionghoa dan Hindu bersembahyang dalam satu kompleks."},
			{"Vihara Dharmayana Kuta", "Salah satu klenteng tertua di Bali, saksi panjang komunitas Tionghoa di Kuta."},
		},
	},
	{
		slug:     "jejak-rasa-bali",
		title:    "Jejak Rasa Bali",
		subtitle: "Sehari mencicipi kuliner warisan keluarga",
		color:    "#6B8E23",
		story:    "Kuliner Bali lahir dari dapur upacara dan resep keluarga yang diwariskan turun-temurun, dengan base genep sebagai jiwanya. Jalur ini mengajakmu berkeliling dari Ubud hingga Jimbaran, dari sarapan nasi ayam sampai makan malam ikan bakar di tepi pantai.",
		stops: []stopSeed{
			{"Nasi Ayam Kedewatan Ibu Mangku", "Mulai pagi dengan nasi ayam, sate lilit, dan lawar khas Kedewatan."},
			{"Warung Babi Guling Ibu Oka", "Babi guling legendaris di dekat Puri Ubud untuk makan siang."},
			{"Pasar Malam Gianyar", "Sore hari, berburu jaje Bali dan babi guling Gianyar di pasar malam."},
			{"Pasar Badung", "Pasar tradisional terbesar Bali untuk mengenal bumbu dan bahan banten."},
			{"Warung Mak Beng", "Sup ikan dan ikan goreng sambal yang bertahan sejak 1941."},
			{"Ikan Bakar Jimbaran", "Tutup hari dengan ikan bakar bumbu Bali sambil menikmati senja."},
		},
	},
	{
		slug:     "jalur-seni-rupa",
		title:    "Jalur Seni Rupa",
		subtitle: "Dari perak Celuk hingga lukisan wayang Kamasan",
		color:    "#A9432A",
		story:    "Seni rupa Bali tumbuh di desa-desa perajin yang melayani pura dan puri. Jalur ini menelusuri kerajinan perak, ukir kayu, gerakan seni Pita Maha di Ubud, hingga lukisan wayang klasik yang menghiasi balai pengadilan kerajaan Klungkung.",
		stops: []stopSeed{
			{"Desa Celuk", "Sentra perak dengan teknik granulasi dan filigri yang dikerjakan tangan."},
			{"Desa Mas", "Desa ukir kayu yang melahirkan pematung-pematung terkenal."},
			{"Puri Saren Agung Ubud", "Pusat gerakan seni Pita Maha pada 1930-an."},
			{"Kertha Gosa Klungkung", "Langit-langit balai pengadilan yang dihiasi lukisan wayang Kamasan."},
			{"Desa Kamasan", "Tempat lukisan wayang klasik masih dibuat secara turun-temurun."},
		},
	},
}

// SeedTrails inserts the curated trails once; existing slugs are skipped.
func SeedTrails(db *gorm.DB) {
	for _, t := range trailSeeds {
		var count int64
		db.Model(&entity.Trail{}).Where("slug = ?", t.slug).Count(&count)
		if count > 0 {
			continue
		}

		trail := entity.Trail{
			Slug:     t.slug,
			Title:    t.title,
			Subtitle: t.subtitle,
			Story:    t.story,
			Color:    t.color,
		}
		for i, s := range t.stops {
			var loc entity.Location
			if err := db.Select("id").Where("name = ?", s.location).First(&loc).Error; err != nil {
				log.Printf("Trail '%s': location '%s' not found, skipping stop", t.slug, s.location)
				continue
			}
			trail.Stops = append(trail.Stops, entity.TrailStop{
				Position:   i + 1,
				LocationID: loc.ID,
				Note:       s.note,
			})
		}
		if err := db.Create(&trail).Error; err != nil {
			log.Printf("Failed to seed trail '%s': %v", t.slug, err)
		}
	}
}
