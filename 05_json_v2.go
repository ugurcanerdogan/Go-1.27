package main

import (
	"fmt"

	jsonv1 "encoding/json"
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
)

// Go 1.27: encoding/json/v2 Artık Kararlı (Deneysel Aşama Bitti!)
//
// ─── GO 1.25 SUNUMUMUZA ATIF: NELER GÖRMÜŞTÜK? ───
// Go 1.25 sunumumuzda bu paketi deneysel olarak incelemiş, GOEXPERIMENT=jsonv2
// bayrağını açarak test etmiştik. O sunumda bahsettiğimiz gibi:
//   - Eski json paketi 10+ yıldır dilin en çok şikayet edilen, yavaş kısımlarındandı.
//   - Go 1.25'te ön izlemesini gördüğümüz bu yeni mimari, Go 1.27 ile artık
//     HERHANGİ BİR BAYRAĞA İHTİYAÇ OLMADAN doğrudan standart kütüphaneye girdi!
//   - Dahası, mevcut projelerinizdeki eski 'encoding/json' (v1) kodları da
//     arka planda artık bu v2 motoruyla çalışıyor (%20-30 hızlanma!).
//
// ─── v2'NİN DAHA GÜVENLİ VE KATI VARSAYILANLARI ───
// v1 bazı hatalı durumları sessizce yutarken, v2 RFC uyumluluğu ve veri güvenliği
// için daha katı varsayılanlarla gelir:
//   - Yinelenen anahtarlar (duplicate keys): v1 sonuncuyu kabul eder, v2 reddeder.
//   - Geçersiz UTF-8: v1 sessizce bozar (\uFFFD), v2 reddeder.
//   - Map sıralaması: v1 her zaman sıralar (yavaş), v2 varsayılan olarak sıralamaz (hızlı).

func main() {
	// ─── 1. Yinelenen Anahtar (Duplicate Object Keys) Kontrolü ───
	// Go 1.25 sunumunda da değinmiştik: JSON standartlarına göre aynı anahtar iki kez olamaz!
	fmt.Println("=== 1. Yinelenen Anahtar (Duplicate Key) Davranışı ===")
	duplicateJSON := []byte(`{"info":"IPHONE-15","info":"IPHONE-16"}`)

	// v1: Hatasız unmarshal eder, ikinci değer ilkinin üzerine yazar (sessiz risk!)
	var mapV1 map[string]string
	_ = jsonv1.Unmarshal(duplicateJSON, &mapV1)
	fmt.Printf("v1 sonucu (sessizce sonuncu alındı): %v\n", mapV1)

	// v2: Veri bütünlüğü için varsayılan olarak hata fırlatır!
	var mapV2 map[string]string
	errV2 := jsonv2.Unmarshal(duplicateJSON, &mapV2)
	fmt.Printf("v2 sonucu (varsayılan hata): %v\n", errV2)

	// İstenirse v2'de seçenekle eski davranışa izin verilebilir:
	_ = jsonv2.Unmarshal(duplicateJSON, &mapV2, jsontext.AllowDuplicateNames(true))
	fmt.Printf("v2 + AllowDuplicateNames(true): %v\n", mapV2)

	// ─── 2. Geçersiz UTF-8 Kontrolü ───
	fmt.Println("\n=== 2. Geçersiz UTF-8 Davranışı ===")
	invalidUTF8 := []byte("{\"name\":\"caf\xffe\"}") // \xff geçerli UTF-8 değildir

	// v1: Hatasız geçer, bozuk karakter yerine Unicode replacement char koyar
	var nameMapV1 map[string]string
	_ = jsonv1.Unmarshal(invalidUTF8, &nameMapV1)
	fmt.Printf("v1 sonucu (karakter bozuldu): %q\n", nameMapV1["name"])

	// v2: Veri bozulmasını önlemek için doğrudan hata döndürür
	var nameMapV2 map[string]string
	errUTF := jsonv2.Unmarshal(invalidUTF8, &nameMapV2)
	fmt.Printf("v2 sonucu (varsayılan hata): %v\n", errUTF)

	// ─── 3. Map Key Sıralaması ve Performans ───
	fmt.Println("\n=== 3. Map Key Sıralaması ve Performans ===")
	// v1 her marshal işleminde map anahtarlarını sıralar (CPU harcar).
	// v2 varsayılan olarak sıralamaz (çok daha hızlı!).
	// Golden testlerde veya imza gerektiren yerlerde sabit sıra için Deterministic(true) verilir.
	stock := map[string]int{"zebra": 10, "apple": 20, "mango": 30}

	// Varsayılan v2 marshal (hızlı, sırasız)
	defaultData, _ := jsonv2.Marshal(stock)
	fmt.Println("v2 varsayılan çıktı:   ", string(defaultData))

	// Deterministic v2 marshal (anahtarlar sıralı, testler için sabit)
	detData, _ := jsonv2.Marshal(stock, jsonv2.Deterministic(true))
	fmt.Println("v2 Deterministic çıktı:", string(detData))
}
