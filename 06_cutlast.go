package main

import (
	"bytes"
	"fmt"
	"strings"
)

// Go 1.27: strings.CutLast ve bytes.CutLast (Sondan Arama ve Kesme)
//
// ─── NE DEĞİŞTİ? ───
// Go 1.18'de gelen `strings.Cut` fonksiyonu string'i ayracın İLK geçtiği yerden
// ikiye bölüyordu. Ancak dosya yolu, uzantı, namespace gibi sondan bölme
// gereken senaryolarda `strings.LastIndex` bulup manuel slice (`s[:i]`, `s[i+1:]`)
// yapmak zorunda kalıyorduk.
//
// Go 1.27 ile hem `strings` hem de `bytes` paketlerine `CutLast` fonksiyonu eklendi!
// Ayracın SON geçtiği yerden tek satırda bölme sağlar.

func main() {
	path := "images/products/pdp/cover.jpg"

	// ─── 1. strings.CutLast — Sondaki Ayraca Göre Bölme ───
	fmt.Println("=== 1. strings.CutLast (Go 1.27) ===")
	// '/' ayracının son geçtiği yerden böler -> dizin yolu ve dosya adı ayrılır
	dir, file, found := strings.CutLast(path, "/")
	fmt.Printf("Kaynak:    %q\n", path)
	fmt.Printf("Dizin:     %q\n", dir)
	fmt.Printf("Dosya:     %q\n", file)
	fmt.Printf("Bulundu mu: %v\n", found)

	// Dosya uzantısını ayırma örneği:
	name, ext, hasExt := strings.CutLast(file, ".")
	fmt.Printf("Ad: %q, Uzantı: %q, Uzantı var mı: %v\n", name, ext, hasExt)

	// ─── 2. strings.Cut ile Farkı (Go 1.18'den beri olan ilk ayraç) ───
	fmt.Println("\n=== 2. strings.Cut Karşılaştırması (İlk Ayraç) ===")
	firstPart, rest, _ := strings.Cut(path, "/")
	fmt.Printf("Cut (ilk '/'):  before=%q, after=%q\n", firstPart, rest)
	lastPart, fileOnly, _ := strings.CutLast(path, "/")
	fmt.Printf("CutLast (son '/'): before=%q, after=%q\n", lastPart, fileOnly)

	// ─── 3. Ayraç Bulunamadığında Ne Olur? ───
	fmt.Println("\n=== 3. Ayraç Bulunamazsa ===")
	// Ayraç yoksa: tüm metin 'before' olur, 'after' boş döner, 'found' false olur.
	noSepBefore, noSepAfter, noSepFound := strings.CutLast("tekkelime", "/")
	fmt.Printf("Sonuç: before=%q, after=%q, found=%v\n", noSepBefore, noSepAfter, noSepFound)

	// ─── 4. bytes.CutLast — []byte Desteği ───
	fmt.Println("\n=== 4. bytes.CutLast ===")
	rawBytes := []byte("catalogs/fashion/dresses")
	bBefore, bAfter, bFound := bytes.CutLast(rawBytes, []byte("/"))
	fmt.Printf("bytes.CutLast: before=%q, after=%q, found=%v\n", bBefore, bAfter, bFound)
}
