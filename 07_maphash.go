package main

import (
	"fmt"
	"hash/maphash"
	"strings"
)

// Go 1.27: hash/maphash.Hasher[T] (Özelleştirilebilir Eşitlik ve Hashing Interface'i)
//
// ─── BU INTERFACE İLE HAYATIMIZA NE GİRDİ? ───
// Go'nun yerleşik `map[K]V` yapısı, anahtarları (key) karşılaştırırken SADECE '=='
// operatörünü kullanır.
//
// Peki ya iki string'i büyük/küçük harf duyarsız (case-insensitive) eşleştirmek
// isterseniz? Örneğin "Go" ile "GO" aynı kabul edilsin istiyorsanız standart map yetmez!
//
// İşte `maphash.Hasher[T]` tam olarak bu ihtiyacı çözer:
// Bir tipin nasıl eşit sayılacağını (`Equal`) ve bu eşitliğe uygun hash değerinin
// nasıl üretileceğini (`Hash`) tek bir standart sözleşmede toplar.
//
// Bu interface tipini implemente edenler iki yöntemi açıkça belirtir:
//   1. Equal(x, y T) bool   -> Hangi iki değer birbirine eşittir? (Örn: ToLower(x) == ToLower(y))
//   2. Hash(h *Hash, val T) -> Bu değerin hash'i nasıl hesaplanır? (Örn: ToLower(val) hash'e beslenir)
//
// Mantık çok basittir: Eğer iki değer `Equal` metoduna göre eşit sayılıyorsa,
// `Hash` metodunun da bu iki değer için aynı tohumla AYNI sayıyı üretmesi gerekir!

// Büyük/küçük harfe duyarsız (case-insensitive) Hasher implementasyonu
type caseInsensitiveHasher struct{}

// Hash yöntemi: Metni küçük harfe çevirerek hash'e besler.
// Böylece "Go" ve "GO" için aynı hash hesaplanır.
func (caseInsensitiveHasher) Hash(h *maphash.Hash, s string) {
	h.WriteString(strings.ToLower(s))
}

// Equal yöntemi: İki metni küçük harfe çevirerek eşitliği kontrol eder.
func (caseInsensitiveHasher) Equal(x, y string) bool {
	return strings.ToLower(x) == strings.ToLower(y)
}

func main() {
	// Hasher interface'ini tanımlıyoruz
	var hasher maphash.Hasher[string] = caseInsensitiveHasher{}

	// ─── 1. Mantıksal Eşitlik Kontrolü (Equal) ───
	fmt.Println("=== 1. Mantıksal Eşitlik (Equal) ===")
	fmt.Println("'Go' == 'GO'   (eşit mi?):", hasher.Equal("Go", "GO"))   // true
	fmt.Println("'Go' == 'Rust' (eşit mi?):", hasher.Equal("Go", "Rust")) // false, because Go >> Rust :P

	// ─── 2. Hash Çıktılarının Eşitliğini Doğrulama (Hash) ───
	fmt.Println("\n=== 2. Hash Çıktısının Doğrulanması ===")
	// İki ayrı Hash nesnesine aynı rastgele tohumu (seed) atıyoruz:
	sharedSeed := maphash.MakeSeed()

	var hashA, hashB maphash.Hash
	hashA.SetSeed(sharedSeed)
	hashB.SetSeed(sharedSeed)

	hasher.Hash(&hashA, "Go")
	hasher.Hash(&hashB, "GO")

	hashValueA := hashA.Sum64()
	hashValueB := hashB.Sum64()

	fmt.Printf("Hash('Go'): %d\n", hashValueA)
	fmt.Printf("Hash('GO'): %d\n", hashValueB)
	fmt.Println("Hash'ler eşit çıktı mı?:", hashValueA == hashValueB) // true!
}
