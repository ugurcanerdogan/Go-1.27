package main

import (
	"fmt"

	"uuid"
)

// Go 1.27: Resmi Standart Kütüphane `uuid` Paketi (RFC 9562)
//
// ─── NE DEĞİŞTİ? ───
// Yıllardır harici 'github.com/google/uuid' gibi paketlere bağımlıydık.
// Go 1.27 ile RFC 9562 standartlı resmi 'uuid' paketi doğrudan stdlib'e geldi!
//
// ─── UUID v4 vs v7 VE B-TREE İNDEKS BAĞLANTISI (NEDEN ÖNEMLİ?) ───
// İlişkisel ve NoSQL veritabanları (PostgreSQL, MySQL, Couchbase) birincil anahtarları
// (primary key) B-Tree indeks yapısında sıralı saklar.
//
// • UUID v4 (Tamamen Rastgele):
//   Yeni eklenen her kayıt indeks ağacında rastgele bir yaprak sayfaya (leaf page)
//   düşer. Sayfalar doldukça ortadan ikiye bölünür (B-Tree page split), bellek
//   önbelleği (cache) sürekli geçersiz kalır ve disk I/O tavan yapar!
//
// • UUID v7 (Zaman Damgalı - Time-Ordered):
//   İlk 48 biti milisaniye cinsinden zaman damgasıdır (timestamp).
//   Yeni kayıtlar zamana göre hep B-Tree indeksinin EN SAĞINA sırayla eklenir (append-only).
//   Sayfa bölünmeleri (page split) minimuma iner, önbellek verimliliği ve INSERT
//   performansı kat kat artar!
//
// ─── ÖNE ÇIKAN DİĞER ÖZELLİKLER ───
// 1. Kriptografik güvenli rastgelelik (CSPRNG).
// 2. UUID doğrudan [16]byte dizisidir; 'id1 == id2' ile sıfır maliyetle kıyaslanır.

func main() {
	// ─── 1. Parse Etme ve Karşılaştırma ───
	fmt.Println("=== 1. Parse Etme ve Sabitler ===")
	rawStr := "f81d4fae-7dec-11d0-a765-00a0c91e6bf6"

	id1 := uuid.MustParse(rawStr)
	id2 := uuid.MustParse(rawStr)

	fmt.Println("Ayrıştırılan ID:     ", id1)
	fmt.Println("Doğrudan eşitlik (==):", id1 == id2) // [16]byte olduğu için doğrudan çalışır!
	fmt.Println("Nil UUID:             ", uuid.Nil())
	fmt.Println("Max UUID:             ", uuid.Max())

	// ─── 2. UUID Üretimi (Generation) ───
	fmt.Println("\n=== 2. UUID Üretimi ===")

	// uuid.New(): Genel kullanım için önerilen varsayılan algoritma
	defaultID := uuid.New() // v4
	fmt.Println("uuid.New()   :", defaultID)

	// uuid.NewV4(): Tamamen rastgele üretilen UUID
	v4ID := uuid.NewV4()
	fmt.Println("uuid.NewV4() :", v4ID, "<- Rastgele (B-Tree indeks parçalanması riski)")

	// uuid.NewV7(): Zamana göre sıralı (time-ordered) UUID
	v7ID := uuid.NewV7()
	fmt.Println("uuid.NewV7() :", v7ID, "<- Zaman damgalı! (DB Primary Key için ideal)")
}
