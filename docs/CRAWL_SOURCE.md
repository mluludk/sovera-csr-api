Berikut daftar sumber data publik yang valid, terstruktur, dan dapat di-*crawl* atau diakses melalui feed publik untuk memetakan sinyal kemitraan CSR, TJSL BUMN, dan dana sosial di Indonesia:

---

### 1. Keterbukaan Informasi & Laporan Emiten Bursa (BEI / IDX)

Kanal ini merupakan sumber data paling valid dengan nilai akurasi finansial dan komitmen ESG tertinggi:

* **Keterbukaan Informasi Publik IDX:**
* *Endpoint Web:* Halaman Pengumuman / Keterbukaan Informasi di situs resmi Bursa Efek Indonesia (`idx.co.id/id/berita/pengumuman`).
* *Target Data:* Pengumuman dividen, RUPS, realisasi alokasi laba bersih untuk dana sosial/lingkungan, dan aksi korporasi emiten terbuka.


* **Laporan Keberlanjutan & Tahunan Emiten (*Sustainability & Annual Reports*):**
* *Pola URL:* Halaman Investor Relations (Hubungan Investor) emiten Tbk (misal: Bank Mandiri, BRI, Telkom, Astra International, Indofood, Adaro).
* *Target Data:* Berkas PDF *Sustainability Report* tahunan yang diwajibkan oleh regulasi OJK (POJK 51/2017). Bagian ini memuat bab alokasi anggaran CSR, wilayah binaan, dan realisasi program kemasyarakatan.



---

### 2. Portal Berita Finansial, Bisnis, dan CSR (RSS Feed & Webhook)

Menggunakan RSS Feed resmi portal media nasional memungkinkan *crawler* berjalan ringan tanpa risiko pemblokiran IP (*IP blocking*):

* **Google News RSS - Search Query CSR Indonesia (`https://news.google.com/rss/search?q=CSR+perusahaan+Indonesia&hl=id&gl=ID&ceid=ID:id`):**
  * Memantau seluruh artikel berita terbaru dari media Indonesia terkait pengumuman CSR dan alokasi dana TJSL korporasi.

* **Bisnis.com (Kategori CSR, Industri, & Finansial):**
* Menyediakan feed berita berkala terkait alokasi hibah dan kegiatan TJSL perusahaan swasta maupun multinasional.


* **Kontan.co.id (Kategori Nasional, Keuangan, & Industri):**
* Fokus pada rilis neraca emiten, pembagian laba, dan liputan program pemberdayaan masyarakat oleh korporasi.


* **Antara News (Kategori Warta Ekonomi & Humaniora):**
* Sumber resmi rilis pers kementerian dan lembaga negara yang memuat liputan bantuan tanggap bencana atau kemitraan sosial korporat di daerah 3T.


* **Republika Online & Kumparan:**
* Memiliki rubrik khusus filantropi, zakat korporasi, serta inisiatif pemberdayaan sosial lembaga kemanusiaan.



---

### 3. Portal Resmi BUMN & Kementerian Terkait

Program Tanggung Jawab Sosial dan Lingkungan (TJSL) BUMN memiliki regulasi wajib untuk mempublikasikan laporan programnya:

* **Situs Resmi Kementerian BUMN / Info TJSL:**
* Siaran pers berkala mengenai penugasan TJSL BUMN dalam klaster pendidikan, lingkungan, dan pengembangan UMKM.


* **Newsroom / Press Release Resmi Korporasi BUMN:**
* Subdomain rilis pers situs korporat BUMN besar (seperti Pertamina, PLN, Pelindo, Pupuk Indonesia, dan BUMN Perbankan Himbara). Korporasi ini secara rutin merilis artikel pengumuman pembukaan program binaan atau penyaluran bantuan kemanusiaan.



---

### 4. Portal Pengadaan & Hibah Lembaga Donor / Kedutaan (Grants)

Untuk melengkapi jangkauan lembaga kemanusiaan non-zakat, portal hibah resmi menyediakan panggilan proposal (*Call for Proposals*):

* **UNDP Indonesia & UN Agencies Procurement Notices:**
* Pengumuman hibah dan tender kemitraan implementasi program lapangan (pendidikan, kesetaraan gender, mitigasi iklim).


* **Kedutaan Besar Asing di Indonesia (Direct Aid Program / Hibah Grassroots):**
* Contoh: Program DAP Kedutaan Australia, GGP Kedutaan Jepang, atau hibah USAID/UKAID yang dipublikasikan terbuka di situs resmi perwakilan mereka untuk mendanai NGO lokal.



---

### Strategi Implementasi Crawler (Rekomendasi Teknis)

1. **Prioritaskan Format RSS / XML Sitemap:** Mulai *crawler* dari RSS Feed portal media untuk meminimalkan beban komputasi dan menghindari batasan scraping HTML murni.
2. **Ekstraksi PDF Parser Terjadwal:** Untuk laporan tahunan/keberlanjutan emiten, jadwalkan bot mingguan/bulanan yang hanya mengunduh dokumen saat terdeteksi versi tahun buku baru.
3. **Filter Kata Kunci Relevan:** Konfigurasikan bot agar hanya memproses artikel yang mengandung kata kunci pemicu: *"TJSL"*, *"CSR"*, *"Bina Lingkungan"*, *"Zakat Perusahaan"*, *"Pemberdayaan Masyarakat"*, *"Bantuan Kemanusiaan"*, *"Hibah"*, atau *"Sustainability"*.