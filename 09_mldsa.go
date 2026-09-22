package main

import (
	"crypto/mldsa"
	"crypto/tls"
	"fmt"
)

// Go 1.27: crypto/mldsa (Kuantum Sonrası / Post-Quantum Dijital İmza)
//
// ─── TLS NEDİR VE POST-QUANTUM İLE NE İLGİSİ VAR? ───
// TLS (Transport Layer Security), web'de HTTPS üzerinden gezinirken tarayıcınız ile
// sunucu arasındaki iletişimi şifreleyen ve sunucunun kimliğini doğrulayan protokoldür.
//
// Günümüzde TLS bağlantılarında kullanılan RSA ve ECDSA gibi imza algoritmalarının,
// gelecekte yeterince güçlü kuantum bilgisayarlar tarafından kolayca kırılabileceği biliniyor.
// Bu yüzden NIST (ABD Ulusal Standartlar Enstitüsü), kuantum bilgisayarlara dayanıklı
// FIPS 204 standardı ML-DSA (Module-Lattice Digital Signature Algorithm) algoritmasını yayımladı.
//
// Go 1.27, bu algoritmayı doğrudan standart kütüphaneye dahil etti: `crypto/mldsa`.
// Ve en önemlisi, doğrudan `crypto/tls` paketine TLS 1.3 imza şeması olarak entegre edildi!
//
// ─── 3 PARAMETRE SETİ ───
//   - MLDSA44: Hızlı ve hafif, genel uygulamalar için önerilen varsayılan (NIST Güvenlik Seviyesi 2)
//   - MLDSA65: Orta güvenlik seviyesi (NIST Seviye 3 - bu örnekte kullanıldı)
//   - MLDSA87: En yüksek güvenlik seviyesi (NIST Seviye 5)

func main() {
	// ─── 1. ML-DSA-65 Anahtar Çifti Üretimi ───
	fmt.Println("=== 1. ML-DSA Anahtar Üretimi ve İmzalama ===")
	params := mldsa.MLDSA65()
	privKey, err := mldsa.GenerateKey(params)
	if err != nil {
		panic(err)
	}

	// İmzalanacak mesaj
	message := []byte("trendyol-order-id:458921")

	// İmzala: ML-DSA imzalama sırasında dışarıdan io.Reader gerektirmez (nil geçilebilir)
	signature, err := privKey.Sign(nil, message, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Şema adı:         ", params.String())
	fmt.Println("İmza boyutu:      ", len(signature), "byte") // ML-DSA-65 için 3309 byte!
	fmt.Println("Public key boyutu:", len(privKey.PublicKey().Bytes()), "byte")

	// ─── 2. İmza Doğrulama (Verification) ───
	fmt.Println("\n=== 2. İmza Doğrulama ===")

	// Orijinal mesaj ile doğrulama -> Başarılı (err == nil)
	errVerify := mldsa.Verify(privKey.PublicKey(), message, signature, nil)
	fmt.Println("Orijinal mesaj doğrulandı mı? :", errVerify == nil)

	// Değiştirilmiş mesaj ile doğrulama -> Başarısız
	tamperedMessage := []byte("trendyol-order-id:999999")
	errTampered := mldsa.Verify(privKey.PublicKey(), tamperedMessage, signature, nil)
	fmt.Println("Tahrif edilmiş mesaj doğrulandı mı?:", errTampered == nil)

	// ─── 3. TLS 1.3 Desteği ───
	fmt.Println("\n=== 3. TLS 1.3 Entegrasyonu ===")
	// crypto/tls paketi artık HTTPS bağlantılarında bu şemaları yerel olarak tanır:
	fmt.Printf("TLS 1.3 Şeması: %s (Değer: 0x%04x)\n", tls.MLDSA65, uint16(tls.MLDSA65))
}
