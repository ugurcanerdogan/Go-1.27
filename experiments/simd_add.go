//go:build goexperiment.simd

package main

import (
	"fmt"
	"simd"
)

// Go 1.27: Taşınabilir SIMD (Single Instruction, Multiple Data) — Deneysel
//
// ─── SIMD NEDİR VE NE SAĞLAR? (ASSEMBLY / CPU DÖNGÜSÜ ANALOJİSİ) ───
// Klasik programlamada (Scalar İşlem):
//   4 elemanlı iki diziyi toplamak için bir `for` döngüsü kurarsınız:
//   İşlemci 4 ayrı toplama işlemi yapar -> 4 kez Fetch, Decode, Execute (4 CPU döngüsü).
//
// SIMD ile (Vektörel İşlem):
//   SIMD = "Tek Komut, Çoklu Veri" (Single Instruction, Multiple Data) demektir.
//   Modern CPU'lardaki geniş vektör yazmaçları (register: ARM Neon 128-bit, Intel AVX2 256-bit, AVX-512)
//   aynı anda birden çok sayıyı (örneğin 4, 8 veya 16 adet float32) tek bir blokta tutar.
//   İşlemciye TEK BİR ASSEMBLY KOMUTU (örn: `vaddps`) verilir; donanım tüm bu sayıları
//   AYNI ANDA TEK BİR CPU DÖNGÜSÜNDE paralel olarak toplar!
//   -> 4-8 kat hızlanma doğrudan donanım seviyesinde elde edilir!
//
// ─── GO 1.27 İLE GELEN YENİLİK: TAŞINABİLİRLİK (PORTABILITY) ───
// Eskiden Go'da SIMD kullanmak için ya platforma özel CGO ya da karmaşık assembly
// yazmak gerekiyordu.
// Go 1.27'deki yeni `simd` paketi:
//   - Tamamen TAŞINABİLİRDİR: Kodunuzu hem Apple Silicon (M1/M2/M3 - ARM Neon)
//     hem Intel/AMD (AVX2/AVX-512) işlemcilerde aynı Go API'si ile çalıştırabilirsiniz.
//   - Donanım desteği yoksa runtime otomatik olarak saf Go döngüsüyle emüle eder.
//
// Çalıştırma:
//   GOEXPERIMENT=simd go run experiments/simd_add.go

func main() {
	// 16 elemanlı iki float32 dilimi
	a := []float32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	b := []float32{10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110, 120, 130, 140, 150, 160}

	// ─── 1. Dilimlerden SIMD Vektörlerine Yükleme ───
	// LoadFloat32s: İşlemcinin desteklediği vektör genişliği (Len) kadar elemanı
	// tek hamlede CPU'nun geniş SIMD vektör register'ına yükler.
	va := simd.LoadFloat32s(a)
	vb := simd.LoadFloat32s(b)

	// ─── 2. Paralel Toplama (Tek Bir CPU Komutuyla Paralel) ───
	// Döngü yok! Donanım seviyesinde tek bir makine komutu ile tüm elemanlar paralel toplanır:
	sum := va.Add(vb)

	// ─── 3. Sonucu Belleğe Geri Yazma ───
	out := make([]float32, sum.Len())
	sum.Store(out)

	fmt.Println("=== Taşınabilir SIMD Vektör Toplama ===")
	fmt.Printf("Bu donanımdaki vektör genişliği (lanes): %d eleman (tek komutta işlenen)\n", sum.Len())
	fmt.Printf("İlk 4 elemanın toplamı: %v\n", out[:4])
}
