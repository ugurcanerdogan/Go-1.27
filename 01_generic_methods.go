package main

import (
	"fmt"
	"math/rand/v2"
	"time"
)

// Go 1.27: Generic Methods (Generic Metotlar)
//
// ─── GO 1.18 vs GO 1.27 FARKI (BASİT ANALOJİ) ───
// Go 1.18'de Generics geldiğinde struct ve fonksiyonlar generic yapılabildi,
// ancak METOTLAR alıcıdan (receiver) bağımsız YENİ bir tip parametresi alamazdı!
//
// Örneğin elinizde `Person[T]` olsun (T: Name tipi, string vs.).
// Person nesnesine bir `Age` (U tipi: int, float64, string) ekleyip yeni bir nesne
// üretmek isteseydiniz:
//
// ❌ Go 1.18 - 1.26 Arası (Zorunlu Paket Fonksiyonu):
//    func WithAge[T, U any](p Person[T], age U) PersonWithAge[T, U] { ... }
//    // Kullanım: p2 := WithAge(person, 30) -> Metot zincirleme (fluent API) imkansızdı!
//
// ✅ Go 1.27 İle Gelen (Doğrudan Metot):
//    func (p Person[T]) WithAge[U any](age U) PersonWithAge[T, U] { ... }
//    // Kullanım: person.WithAge(30) -> Tıpkı modern dillerdeki gibi metot zinciri!
//
// ─── ÖNEMLİ KISIT ───
// Interface'ler generic metot tanımlayamaz ve generic metotlar bir interface'i
// implemente etmek için kullanılamaz (vtable boyutunun derleme anında sabit kalması için).
//
//	type Mapper interface {
//	    Map[U any](f func(int) U) any // DERLENMEZ: interface method must have no type parameters
//	}

// ─── 1. Analoji Örneği: Person ve WithAge Metodu ───
type Person[T any] struct {
	Name T
}

type PersonWithAge[T, U any] struct {
	Name T
	Age  U
}

// WithAge metodu kendi bağımsız [U any] tip parametresini tanımlar (Go 1.27+):
func (p Person[T]) WithAge[U any](age U) PersonWithAge[T, U] {
	return PersonWithAge[T, U]{Name: p.Name, Age: age}
}

// ─── 2. Container Örneği: Box ve Map Metodu ───
type Box[T any] struct {
	v T
}

// Map metodu kendi tip parametresi olan [U any]'yi tanımlar.
// Receiver Box[T] iken, dönen değer Box[U] olur.
func (b Box[T]) Map[U any](f func(T) U) Box[U] {
	return Box[U]{v: f(b.v)}
}

func main() {
	// ─── 1. Person - WithAge Analojisi ───
	fmt.Println("=== 1. Person - Age Analojisi (1.18 vs 1.27) ===")
	user := Person[string]{Name: "Uğurcan"}

	// Go 1.27 sayesinde doğrudan metot üzerinden farklı bir tip (int) ekliyoruz:
	userWithIntAge := user.WithAge(28)
	fmt.Printf("Kullanıcı: %s, Yaş (int): %d\n", userWithIntAge.Name, userWithIntAge.Age)

	// Aynı metot farklı bir tiple (string) de çağrılabilir:
	userWithStringAge := user.WithAge("Yirmi Sekiz")
	fmt.Printf("Kullanıcı: %s, Yaş (string): %s\n", userWithStringAge.Name, userWithStringAge.Age)

	// ─── 2. Box.Map Fonksiyonel Dönüşüm Örneği ───
	fmt.Println("\n=== 2. Generic Metot Örneği: Box.Map ===")
	intBox := Box[int]{v: 21}

	// int -> int dönüşümü
	doubled := intBox.Map(func(n int) int { return n * 2 })

	// int -> string dönüşümü (U tipi string oldu)
	stringBox := doubled.Map(func(n int) string {
		return fmt.Sprintf("value=%d", n)
	})
	fmt.Println("Sonuç:", stringBox.v) // "value=42"

	// ─── 3. Standart Kütüphaneden Örnek: (*rand.Rand).N ───
	fmt.Println("\n=== 3. Standart Kütüphane: (*rand.Rand).N ===")
	// Go 1.22'de paket seviyesinde generic rand.N() fonksiyonu gelmişti.
	// Go 1.27 ile generic metot desteği sayesinde bu özellik doğrudan
	// *rand.Rand instance'ı üzerine de eklendi: func (r *Rand) N[Int intType](n Int) Int
	// Permuted Congruential Generator (PCG) tohumu ile:
	r := rand.New(rand.NewPCG(1, 2))

	fmt.Println("int rastgele:     ", r.N(100))
	fmt.Println("int32 rastgele:   ", r.N(int32(10)))
	fmt.Println("duration rastgele:", r.N(5*time.Second))
}
