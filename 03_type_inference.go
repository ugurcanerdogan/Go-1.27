package main

import "fmt"

// Go 1.27: Generalized Function Type Inference (Genişletilmiş Fonksiyon Tip Çıkarımı)
//
// ─── NE DEĞİŞTİ? ───
// Go 1.21'den bu yana, generic bir fonksiyonu doğrudan bir değişkene atarken
// derleyici tipi çıkarabiliyordu (`var fn func(int) int = first`).
// Ancak diğer atama durumlarında derleyici tipi otomatik çıkaramıyor ve bizi
// `first[int]` gibi tip parametrelerini açıkça yazmaya zorluyordu:
//   1. Composite literal elemanlarında (`[]func([]int) int{first}`)
//   2. Açık tip dönüşümlerinde (`MyFuncType(genericFunc)`)
//   3. Channel'a veri gönderirken (`ch <- genericFunc`)
//
// Go 1.27 ile tip çıkarımı genelleştirildi: Hedef fonksiyon tipinin bilindiği
// TÜM atama bağlamlarında generic fonksiyonlar ek tip argümanı olmadan kullanılabilir.

// İki basit generic yardımcı fonksiyon
func first[T any](s []T) T { return s[0] }
func last[T any](s []T) T  { return s[len(s)-1] }

func genericFormatter[T any](v T) string {
	return fmt.Sprintf("değer: %v", v)
}

type IntFormatter func(int) string

func main() {
	// ─── 1. Composite Literal İçinde Tip Çıkarımı ───
	// Go 1.26'da: []func([]int) int{first[int], last[int]} yazmak zorunluydu.
	// Go 1.27'de: Slice tipi func([]int) int olduğu için T = int otomatik çıkarılır.
	fmt.Println("=== 1. Composite Literal İçinde ===")
	ops := []func([]int) int{first, last}
	numbers := []int{10, 20, 30}
	for _, op := range ops {
		fmt.Println("Sonuç:", op(numbers))
	}

	// ─── 2. Tip Dönüşümü ve Channel Send İşlemleri ───
	// Go 1.26'da: IntFormatter(genericFormatter[int]) yazmak gerekiyordu.
	// Go 1.27'de: Hedef tip IntFormatter (func(int) string) olduğundan T=int otomatik anlaşılır.
	fmt.Println("\n=== 2. Tip Dönüşümü ve Channel Send ===")

	// Tip dönüşümü
	fn := IntFormatter(genericFormatter)
	fmt.Println("Dönüştürülen fonksiyon:", fn(42))

	// Channel'a gönderme
	ch := make(chan IntFormatter, 1)
	ch <- genericFormatter // Otomatik T=int çıkarımı
	receivedFn := <-ch
	fmt.Println("Kanaldan gelen fonksiyon:", receivedFn(99)) // -> string
}
