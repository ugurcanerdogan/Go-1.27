# Introduction to Go 1.27

Go 1.27, Ağustos 2026'da yayınlandı. Dil değişiklikleri, standart kütüphaneye eklenen yeni paketler, performans iyileştirmeleri ve araç güncellemeleri içeriyor. Bu depodaki tüm örnekler **go1.27.1** ile test edilmiş ve doğrulanmıştır.

Go 1.25'te deneysel sunulan `encoding/json/v2` ile Go 1.26'da deneysel sunulan goroutine leak profili bu sürümle birlikte kararlı hale gelmiştir.

## Contents

### Dil Değişiklikleri
* **Generic Metotlar**: Metotlar artık kendi tip parametrelerini (`func (b Box[T]) Map[U any]`) tanımlayabiliyor.
* **Struct Literal Field Selectors**: Gömülü (promoted) alanlar, iç içe struct literal yazmadan doğrudan atanabiliyor.
* **Genelleştirilmiş Fonksiyon Tip Çıkarımı**: Generic fonksiyonlar composite literal, tip dönüşümü ve channel send işlemlerinde açık tip parametresi verilmeden kullanılabiliyor.

### Yeni Standart Kütüphane Özellikleri
* **`uuid` Paketi**: Standart kütüphaneye RFC 9562 uyumlu UUID üretimi ve ayrıştırması eklendi (`NewV4`, `NewV7`).
* **`encoding/json/v2`**: Deneysel bayrak kalktı. v2 katı kurallarla (duplicate key ve bozuk UTF-8 reddi) çalışır, v1 API ise v2 motoruna geçirildi.
* **`strings.CutLast` / `bytes.CutLast`**: Ayracın son geçtiği yerden tek hamlede bölme.
* **`hash/maphash.Hasher`**: Eşitlik ve hash kuralını birlikte tanımlayan tip-güvenli sözleşme.
* **`math/big.Int.Divide`**: Yuvarlama modlu (`Ceil`, `Floor`, `Trunc`, `Round`) büyük sayı bölme işlemi.
* **`crypto/mldsa`**: Kuantum sonrası (post-quantum) dijital imza desteği (FIPS 204). TLS 1.3 ve X.509 entegrasyonu.

### Çalışma Zamanı (Runtime) ve Performans
* **Traceback Başlıklarında Goroutine Labels**: `pprof.Labels` etiketleri panic ve traceback başlıklarında otomatik listelenir.
* **Goroutine Leak Profile**: Kararlı hale geldi (`pprof.Lookup("goroutineleak")`), deneysel bayrak kaldırıldı.
* **Boyuta Özel Bellek Tahsisi (Size-specialized Allocation)**: 80 byte altı küçük nesnelerin tahsis maliyeti %30'a varan oranda azaldı.
* **Unicode 17 Desteği**: İki sürüm birden güncellenerek Unicode 17.0 karakter setine geçildi.

### Test ve Araç Güncellemeleri
* **`testing/synctest.Sleep`**: Sanal saat balonu içinde `time.Sleep` ve `synctest.Wait` adımlarını birleştiren deterministik test yardımcısı.
* **`net/http/httptest.NewTestServer`**: Gerçek TCP portu açmayan, in-memory ve otomatik temizlenen HTTP test sunucusu.
* **`go fix` Modernizers**: `atomictypes`, `embedlit`, `slicesbackward` ve `unsafefuncs` olmak üzere 4 yeni analiz kuralı.
* **`go test`**: `stdversion` vet kontrolü artık varsayılan olarak çalışır.
* **`go mod tidy`**: Dağınık `require` bloklarını otomatik olarak direct/indirect olmak üzere 2 blokta toplar.

### Deneysel Özellikler
* **`simd`**: Donanım ve mimariden bağımsız taşınabilir vektörel işlem paketi (`GOEXPERIMENT=simd`).

## Örnekler ve Kullanım

| Dosya | Konu |
|-------|------|
| `01_generic_methods.go` | Generic metotlar ve `(*rand.Rand).N` |
| `02_struct_literals.go` | Gömülü struct alanlarına doğrudan değer atama |
| `03_type_inference.go` | Composite literal, dönüşüm ve kanallarda tip çıkarımı |
| `04_uuid.go` | Standart `uuid` paketi (`NewV4`, `NewV7`) |
| `05_json_v2.go` | `encoding/json/v2` katı varsayılanlar ve `Deterministic` |
| `06_cutlast.go` | `strings.CutLast` ve `bytes.CutLast` ile sondan kesme |
| `07_maphash.go` | `maphash.Hasher[T]` ile özel eşitlik ve hashleme |
| `08_big_divide.go` | `math/big.Int.Divide` yuvarlama modları |
| `09_mldsa.go` | Post-quantum ML-DSA-65 imzalama ve doğrulama |
| `10_goroutine_labels.go` | Traceback başlıklarında goroutine etiketleri |
| `11_goroutine_leak.go` | `goroutineleak` profili ile sızıntı tespiti |
| `12_unicode.go` | Unicode 17 desteği ve yeni karakterler |
| `13_gofix_demo.go` | `go fix` için 4 yeni modernizer |
| `tests/14_synctest_sleep_test.go` | `synctest.Sleep` ile deterministik eşzamanlılık testi |
| `tests/15_httptest_server_test.go` | `httptest.NewTestServer` bellek içi test sunucusu |
| `experiments/simd_add.go` | Taşınabilir SIMD vektör toplama (deneysel) |

```bash
# Tek bir örneği çalıştırma
go run 01_generic_methods.go

# Testleri çalıştırma
go test -v -race ./tests/

# go fix değişikliklerini dosyaya yazmadan diff olarak inceleme
go fix -diff .

# SIMD örneğini çalıştırma (deneysel)
GOEXPERIMENT=simd go run experiments/simd_add.go
```

## Performans Karşılaştırmaları

### Boyuta Özel Bellek Tahsisi (Size-specialized Malloc)
* <80 byte nesne tahsislerinde: **~%30'a varan maliyet düşüşü**
* Gerçek dünya tahsis-yoğun programlarda: **~%1 genel performans artışı**
* Binary boyut etkisi: **+~60 KB**
* Devre dışı bırakma seçeneği: `GOEXPERIMENT=nosizespecializedmalloc`

### JSON v2 Motoru
* Unmarshal işlemlerinde: **%20-30 hızlanma**
* Map serialization: Varsayılan olarak anahtar sıralaması yapılmadığı için belirgin CPU tasarrufu

## Geçiş ve Uyumluluk Notları

* **`asynctimerchan` Kaldırıldı**: `time.After`, `NewTimer` ve `NewTicker` kanalları artık daima unbuffered (senkron) çalışır. Eski tamponlu davranışa dönme bayrağı silinmiştir.
* **macOS 13 Ventura**: macOS desteği için minimum sürüm 13 Ventura olarak güncellendi.
* **`encoding/json` Hata Mesajları**: v1 motoru v2 altyapısına geçtiği için hata metinlerinde küçük değişiklikler olabilir (davranış aynıdır). Gerekirse geçici geri dönüş: `GOEXPERIMENT=nojsonv2`.
* **HTTP/1 Keep-Alive**: `Response.Body.Close()`, bağlantının yeniden kullanılabilmesi için okunmamış gövdeyi belirli bir sınıra kadar otomatik tahliye eder (drain).
* **HTTP/2 Client Priority**: Sunucu artık RFC 9218 öncelik sinyallerini varsayılan olarak dikkate alır.

## Kaynaklar

* [Go 1.27 Release Notes](https://go.dev/doc/go1.27)
* [Go 1.27 Blog Post](https://go.dev/blog/go1.27)
* [Go 1.27 Interactive Tour (VictoriaMetrics)](https://victoriametrics.com/blog/go-1-27/)

## Sunum Detayları

* **Presenter**: Uğurcan Erdoğan
* **Contact**: [ugurcanerdogan3306@gmail.com](mailto:ugurcanerdogan3306@gmail.com)
