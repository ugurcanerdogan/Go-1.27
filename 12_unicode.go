package main

import (
	"fmt"
	"unicode"
)

// Go 1.27: Unicode 17 Standardına Geçiş
//
// ─── NE DEĞİŞTİ? ───
// Go'nun standart kütüphanesindeki `unicode` paketi ve dildeki karakter tabloları
// Go 1.26'da Unicode 15.0 seviyesindeydi.
// Go 1.27 ile birlikte iki büyük sürüm birden atlanarak doğrudan
// Unicode 17.0 standardına güncellendi!
//
// Bu sayede son iki yılda eklenen yüzlerce yeni emoji, sembol, para birimi ve
// alfabe karakteri artık Go tarafından yerel olarak tanınıyor.

func main() {
	// Standart kütüphanedeki güncel Unicode sürümü
	fmt.Printf("Mevcut Unicode Sürümü: %s\n\n", unicode.Version)

	// Unicode 16.0 ile eklenen kök sebze (root vegetable) emojisi:
	// Kod noktası: U+1FADC (🫜)
	r := '\U0001FADC'

	// Go 1.26'da bu kod çalıştırıldığında:
	//   IsSymbol = false
	//   IsGraphic = false dönerdi çünkü Go 1.26 henüz Unicode 15'teydi.
	//
	// Go 1.27'de Unicode 17 ile birlikte:
	//   Her iki fonksiyon da 'true' döner!
	isSymbol := unicode.IsSymbol(r)
	isGraphic := unicode.IsGraphic(r)

	fmt.Printf("Karakter: %c (Unicode: %#U)\n", r, r)
	fmt.Printf("unicode.IsSymbol (Sembol mü?) : %v\n", isSymbol)
	fmt.Printf("unicode.IsGraphic(Görsel mi?) : %v\n", isGraphic)
}
