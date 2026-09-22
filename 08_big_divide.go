package main

import (
	"fmt"
	"math/big"
)

// Go 1.27: math/big.Int.Divide (Yuvarlama Modu Destekli Bölme)
//
// ─── NE DEĞİŞTİ? ───
// Klasik `big.Int.Quo` ve `Rem` fonksiyonları bölme sonucunu her zaman sıfıra
// doğru kesiyordu (truncation). Ancak finans, muhasebe, faiz ve kriptografi
// hesaplamalarında yukarı yuvarlama (Ceil) veya aşağı yuvarlama (Floor)
// kritik önem taşır.
//
// Go 1.27 ile `Int.Divide(x, y, r, mode)` metodu eklendi.
// Bölüm (q) ve kalan (r) seçilen yuvarlama moduna göre BİRLİKTE hesaplanır:
//   - big.Ceil  : Tavana yuvarlar (+sonsuza doğru)
//   - big.Floor : Tabana yuvarlar (-sonsuza doğru)
//   - big.Trunc : Sıfıra doğru keser (klasik tamsayı bölmesi)
//   - big.Round : En yakın tam sayıya yuvarlar
//
// ─── FORMÜL KURALI ───
// Her modda matematiksel olarak `x = y*q + r` eşitliği korunur.
// Dolayısıyla bölüm yukarı yuvarlandığında kalanın negatif çıkması doğaldır!

func main() {
	// Örnek: 7 bölü 2 (normal bölmede 3.5)
	x := big.NewInt(7)
	y := big.NewInt(2)

	q := new(big.Int) // Bölüm (quotient)
	r := new(big.Int) // Kalan (remainder)

	fmt.Println("İşlem: 7 / 2")

	// ─── 1. big.Ceil (Yukarı Yuvarlama) ───
	// 3.5 yukarı yuvarlanır -> q = 4
	// Eşitliğin (7 = 2*4 + r) sağlanması için kalan r = -1 olur
	q.Divide(x, y, r, big.Ceil)
	fmt.Printf("Ceil (Tavan) : Bölüm (q)=%s, Kalan (r)=%s  [Sağlama: 2*4 + (%s) = 7]\n", q, r, r)

	// ─── 2. big.Floor (Aşağı Yuvarlama) ───
	// 3.5 aşağı yuvarlanır -> q = 3
	// Eşitliğin (7 = 2*3 + r) sağlanması için kalan r = 1 olur
	q.Divide(x, y, r, big.Floor)
	fmt.Printf("Floor (Taban): Bölüm (q)=%s, Kalan (r)=%s   [Sağlama: 2*3 + %s = 7]\n", q, r, r)

	// ─── 3. big.Trunc (Sıfıra Doğru Kesme - Klasik Yöntem) ───
	q.Divide(x, y, r, big.Trunc)
	fmt.Printf("Trunc (Kesme): Bölüm (q)=%s, Kalan (r)=%s   [Klasik Quo/Rem ile aynı]\n", q, r)
}
