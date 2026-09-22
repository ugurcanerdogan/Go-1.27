package main

import "fmt"

// Go 1.27: Struct Literal Field Selectors (Gömülü Alanlara Doğrudan Erişim)
//
// ─── NE DEĞİŞTİ? ───
// Go'da struct embedding (gömme) yapıldığında içteki alanlar dışarı "promote"
// edilir ve nesne üzerinden `u.UpdatedBy` şeklinde doğrudan okunabilir.
// Ancak struct literal ile nesne oluştururken bu alanlara değer atamak için
// mutlaka iç içe struct literal yazmak gerekiyordu.
//
// Go 1.27 ile artık gömülü (promoted) alanlar struct literal içinde doğrudan
// anahtar (key) olarak kullanılabiliyor!
//
// ─── DİKKAT EDİLECEK NOKTA ───
// Eğer birden fazla gömülü struct aynı isimde alan içeriyorsa (çakışma),
// derleyici belirsizlik hatası verir ve eski yöntemle açıkça belirtmenizi ister.

type AuditInfo struct {
	CreatedBy string
	UpdatedBy string
}

type ProductListing struct {
	AuditInfo // Gömülü (embedded) struct
	Info      string
	Price     float64
}

func main() {
	// ─── ESKİ YÖNTEM (Go 1.26 ve öncesi) ───
	// Gömülü alanlara değer vermek için iç içe AuditInfo{...} yazmak zorunluydu:
	fmt.Println("=== Eski Yöntem (İç İçe Literal) ===")
	oldItem := ProductListing{
		AuditInfo: AuditInfo{
			CreatedBy: "system",
			UpdatedBy: "catalog-team",
		},
		Info:  "TSHIRT-BLK-M",
		Price: 199.99,
	}
	fmt.Printf("Info: %s, Güncelleyen: %s\n", oldItem.Info, oldItem.UpdatedBy)

	// ─── YENİ YÖNTEM (Go 1.27) ───
	// Gömülü struct adını yazmaya gerek yok! Doğrudan alan isimlerini yazıyoruz:
	fmt.Println("\n=== Yeni Yöntem (Go 1.27 — Doğrudan Atama) ===")
	newItem := ProductListing{
		CreatedBy: "system",
		UpdatedBy: "catalog-team", // Doğrudan AuditInfo.UpdatedBy alanına atanır
		Info:      "TSHIRT-BLK-M",
		Price:     199.99,
	}
	fmt.Printf("Info: %s, Güncelleyen: %s\n", newItem.Info, newItem.UpdatedBy)

	// İpucu: 'go fix -embedlit .' komutu eski stildeki kodları otomatik olarak
	// bu yeni ve sade syntax'a dönüştürür.
}
