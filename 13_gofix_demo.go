package main

import "fmt"

// Go 1.27: `go fix` Yeni Modernizer'lar (Kodunuzu Otomatik Güncelleyin)
//
// ─── NE DEĞİŞTİ? ───
// Go 1.26'da tamamen yeniden yazılan modernizer analiz altyapısına Go 1.27 ile
// 4 YENİ ANALYZER eklendi!
//
// Bu araç, eski Go sürümlerinden kalma kalıpları tespit eder ve tek komutla
// modern, güvenli ve performanslı Go 1.27 koduna dönüştürür.
//
// ─── YENİ ANALYZER'LAR ───
// 1. atomictypes   : sync/atomic fonksiyon çağrılarını strongly-typed atomic tiplere çevirir.
// 2. embedlit      : Struct literal içindeki gereksiz iç içe gömülü struct tanımlarını temizler.
// 3. slicesbackward: Tersten dönen manuel for döngülerini slices.Backward() iterator'ına çevirir.
// 4. unsafefuncs   : Manuel pointer aritmetiğini unsafe.Add / unsafe.Slice çağrılarına çevirir.
//
// ─── DİĞER GÜNCELLEMELER ───
// - 'waitgroup' analyzer'ı, 'waitgroupgo' olarak yeniden adlandırıldı.
// - 'fmtappendf' analyzer'ı kod stili endişeleri sebebiyle kaldırıldı.

func main() {
	fmt.Println("=== Go 1.27 `go fix` Modernizer Örnekleri ===")
	fmt.Println("Bu dosyadaki açıklamalar yeni eklenen 4 modernizer'ın")
	fmt.Println("önce/sonra kod dönüşümlerini gösterir.")
	fmt.Println()
	fmt.Println("Kendi projenizde denemek için:")
	fmt.Println("  go fix -diff .              # Değişiklikleri uygulamadan sadece diff göster")
	fmt.Println("  go fix -atomictypes .       # Sadece atomictypes kuralını uygula")
	fmt.Println("  go fix .                    # Tüm kuralları otomatik uygula")
}

// ─────────────────────────────────────────────────────────────────────────────
// 1. ANALYZER: atomictypes
// ─────────────────────────────────────────────────────────────────────────────
// Eski Yöntem: Düz ilkel tipler ve sync/atomic paket fonksiyonları.
//   var counter int32
//   atomic.AddInt32(&counter, 1)
//   val := atomic.LoadInt32(&counter)
//
// RİSK:
//   Değişken düz bir 'int32' olduğu için, başka bir geliştirici veya dikkatsiz
//   bir anınızda atomic fonksiyon yerine yanlışlıkla `counter++` veya `val = counter`
//   yazabilir! Bu da goroutine'ler arası veri yarışına (data race) ve tutarsızlığa yol açar.
//
// `go fix` Sonrası (atomic.Int32 kullanımı teşvik edilir):
//   var counter atomic.Int32
//   counter.Add(1)
//   val := counter.Load()
//   -> Artık değişkeni kazara atomik olmayan yolla değiştiremezsiniz!

// ─────────────────────────────────────────────────────────────────────────────
// 2. ANALYZER: embedlit
// ─────────────────────────────────────────────────────────────────────────────
// type Inner struct { ID int }
// type Outer struct { Inner; Name string }
//
// Eski Yöntem (İç içe gömülü struct tipi zorunluydu):
//   o := Outer{
//       Inner: Inner{ID: 10},
//       Name:  "Telefon",
//   }
//
// `go fix` Sonrası (Go 1.27 promoted field syntax):
//   o := Outer{
//       ID:   10,
//       Name: "Telefon",
//   }

// ─────────────────────────────────────────────────────────────────────────────
// 3. ANALYZER: slicesbackward
// ─────────────────────────────────────────────────────────────────────────────
// Eski Yöntem (Hata yapmaya müsait indeks tabanlı ters döngü):
//   for i := len(items) - 1; i >= 0; i-- {
//       process(items[i])
//   }
//
// `go fix` Sonrası (Go 1.23+ range iterator):
//   for _, v := range slices.Backward(items) {
//       process(v)
//   }

// ─────────────────────────────────────────────────────────────────────────────
// 4. ANALYZER: unsafefuncs
// ─────────────────────────────────────────────────────────────────────────────
// Eski Yöntem (Karmaşık ve okunması zor pointer aritmetiği):
//   nextPtr := unsafe.Pointer(uintptr(ptr) + uintptr(offset))
//
// `go fix` Sonrası (Standart ve güvenli yardımcı fonksiyon):
//   nextPtr := unsafe.Add(ptr, offset)
