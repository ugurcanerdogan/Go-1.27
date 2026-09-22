//go:build !goexperiment.simd

package main

import "fmt"

// Bu dosya GOEXPERIMENT=simd bayrağı verilmediğinde derlenir.
// Amacı: 'simd' paketi deneysel olduğu için normal derlemede derleme hatası
// almamak ve geliştiriciyi doğru bayrakla yönlendirmektir.

func main() {
	fmt.Println("=== Taşınabilir SIMD Paketi Bilgilendirmesi ===")
	fmt.Println("'simd' paketi Go 1.27'de deneyseldir ve standart derlemede kapalıdır.")
	fmt.Println("Örneği çalıştırmak için lütfen şu komutu kullanın:")
	fmt.Println("\n  GOEXPERIMENT=simd go run experiments/simd_add.go")
}
