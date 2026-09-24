# 🚗 CarsXE API (Go Package)

[![Go Reference](https://pkg.go.dev/badge/github.com/carsxe/carsxe-go-package.svg)](https://pkg.go.dev/github.com/carsxe/carsxe-go-package)

**[CarsXE](https://carsxe.com)** is a powerful and developer-friendly API that gives you instant access to a wide range of vehicle data. From [VIN decoding](https://carsxe.com/vehicle-specifications) and [market value](https://carsxe.com/vehicle-market-value) estimation to [vehicle history](https://carsxe.com/vehicle-history), [images](https://carsxe.com/vehicle-images), OBD code explanations, and [plate recognition](https://carsxe.com/vehicle-plate-decoder), CarsXE provides everything you need to build automotive applications at scale.

🌐 **Website:** [https://carsxe.com](https://carsxe.com)  
📄 **Docs:** [https://carsxe.com/docs](https://carsxe.com/docs)  
📦 **All Products:** [https://carsxe.com/all-products](https://carsxe.com/all-products)

**Product pages:** [Vehicle History](https://carsxe.com/vehicle-history) · [Plate Decoder](https://carsxe.com/vehicle-plate-decoder) · [Vehicle Specifications](https://carsxe.com/vehicle-specifications) · [International VIN Decoder](https://carsxe.com/international-vin-decoder) · [Vehicle Images](https://carsxe.com/vehicle-images) · [Vehicle Recalls](https://carsxe.com/vehicle-recalls) · [Market Value](https://carsxe.com/vehicle-market-value)

To get started with the CarsXE API, follow these steps:

1. **Sign up for a CarsXE account:**
   - [Register here](https://carsxe.com/register)
   - Add a [payment method](https://carsxe.com/dashboard/billing#payment-methods) to activate your subscription and get your API key.

2. **Install the CarsXE Go package:**

   Run this command in your terminal:

   ```bash
   go get -u github.com/carsxe/carsxe-go-package
   ```

3. **Import the CarsXE API into your code:**

   ```go
   import "github.com/carsxe/carsxe-go-package"
   ```

4. **Initialize the API with your API key:**

   ```go
   client := carsxe.New("YOUR_API_KEY")
   ```

5. **Use the various endpoint methods provided by the API to access the data you need.**

## Usage

```go
package main

import (
	"fmt"

	"github.com/carsxe/carsxe-go-package"
)

func main() {
	client := carsxe.New("YOUR_API_KEY")
	vin := "WBAFR7C57CC811956"

	vehicle := client.Specs(map[string]string{"vin": vin})
	fmt.Println(vehicle["input"].(map[string]interface{})["vin"])
}
```

---

## 📚 Endpoints

The CarsXE API provides the following endpoint methods:

### `Specs` – Decode VIN & get full [vehicle specifications](https://carsxe.com/vehicle-specifications)

**Required:**

- `vin`

**Optional:**

- `deepdata`
- `disableIntVINDecoding`

**Example:**

```go
vehicle := client.Specs(map[string]string{"vin": "WBAFR7C57CC811956"})
```

---

### `InternationalVINDecoder` – Decode VIN with [worldwide support](https://carsxe.com/international-vin-decoder)

**Required:**

- `vin`

**Optional:**

- None

**Example:**

```go
intvin := client.InternationalVINDecoder(map[string]string{"vin": "WF0MXXGBWM8R43240"})
```

---

### `PlateDecoder` – Decode [license plate](https://carsxe.com/vehicle-plate-decoder) info (plate, country)

**Required:**

- `plate`
- `country` (always required except for US, where it is optional and defaults to 'US')

**Optional:**

- `state` (required for some countries, e.g. US, AU, CA)
- `district` (required for Pakistan)

> **Note:**
>
> - The `state` parameter is required only when applicable (for
>   specific countries such as US, AU, CA, etc.).
> - For Pakistan (`country='pk'`), both `state` and `district`
>   are required.

**Example:**

```go
decodedPlate := client.PlateDecoder(map[string]string{"plate": "7XER187", "state": "CA", "country": "US"})
```

---

### `MarketValue` – Estimate vehicle [market value](https://carsxe.com/vehicle-market-value) based on VIN

**Required:**

- `vin`

**Optional:**

- `state` — US state code for regional pricing (e.g. `CA`, `TX`)
- `mileage` — current mileage to adjust the value
- `condition` — vehicle condition: `excellent` | `clean` | `average` | `rough`

**Example:**

```go
marketvalueDetailed := client.MarketValue(map[string]string{"vin": "WBAFR7C57CC811956", "state": "CA", "mileage": "45000", "condition": "clean"})
```

---

### `History` – Retrieve [vehicle history](https://carsxe.com/vehicle-history)

**Required:**

- `vin`

**Optional:**

- None

**Example:**

```go
history := client.History(map[string]string{"vin": "WBAFR7C57CC811956"})
```

---

### `Images` – Fetch [images](https://carsxe.com/vehicle-images) by make, model, year, trim

**Required:**

- `make`
- `model`

**Optional:**

- `year`
- `trim`
- `color`
- `transparent`
- `angle`
- `photoType`
- `size`
- `license`

**Example:**

```go
images := client.Images(map[string]string{"make": "BMW", "model": "X5", "year": "2019"})
```

---

### `Recalls` – Get safety [recall](https://carsxe.com/vehicle-recalls) data for a VIN

**Required:**

- `vin`

**Optional:**

- None

**Example:**

```go
recalls := client.Recalls(map[string]string{"vin": "1C4JJXR64PW696340"})
```

---

### `PlateImageRecognition` – Read & decode [plates](https://carsxe.com/vehicle-plate-decoder) from images

**Required:**

- `imageURL`

**Optional:**

- None

**Example:**

```go
plateimg := client.PlateImageRecognition("https://imagedelivery.net/moyiiSImjJPI_EZVxNMBBw/f49aed53-d736-4370-f3f4-97418841c800/public")
```

![Plate recognition sample](https://imagedelivery.net/moyiiSImjJPI_EZVxNMBBw/f49aed53-d736-4370-f3f4-97418841c800/public)

---

### `VinOCR` – Extract VINs from images using OCR

**Required:**

- `imageURL`

**Optional:**

- None

**Example:**

```go
vinocr := client.VinOCR("https://imagedelivery.net/moyiiSImjJPI_EZVxNMBBw/f49aed53-d736-4370-f3f4-97418841c800/public")
```

---

### `YearMakeModel` – Query vehicle by year, make, model and trim (optional)

**Required:**

- `year`
- `make`
- `model`

**Optional:**

- `trim`

**Example:**

```go
yymm := client.YearMakeModel(map[string]string{"year": "2012", "make": "BMW", "model": "5 Series"})
```

---

### `ObdCodesDecoder` – Decode OBD error/diagnostic codes

**Required:**

- `code`

**Optional:**

- None

**Example:**

```go
obdcode := client.ObdCodesDecoder(map[string]string{"code": "P0115"})
```

---

### `LienAndTheft` – Check lien and theft records by VIN

**Required:**

- `vin`

**Optional:**

- None

**Example:**

```go
lienTheft := client.LienAndTheft(map[string]string{"vin": "2C3CDXFG1FH762860"})
```

---

### `RecallsYmm` – Get safety [recall](https://carsxe.com/vehicle-recalls) data by year, make, and model

**Required:**

- `year`
- `make`
- `model`

**Optional:**

- None

**Example:**

```go
recallsYmm := client.RecallsYmm(map[string]string{"year": "2026", "make": "toyota", "model": "corolla"})
```

---

### `SubmitRecallsBatch` – Submit VINs for async bulk [recall](https://carsxe.com/vehicle-recalls) checking

POST JSON to `/v1/recalls-batch/submit`. Provide at least one of `vins`, `csv`, or `csvUrl` (they can be combined). Max 10,000 unique VINs.

**Required (at least one):**

- `vins` — array of 17-character VIN strings
- `csv` — inline CSV text
- `csvUrl` — HTTPS URL to a CSV file

**Optional:**

- `webhookUrl` — HTTPS URL notified when the batch finishes

**Example:**

```go
batch := client.SubmitRecallsBatch(map[string]any{
	"vins": []string{"1HGBH41JXMN109186", "5YJSA1E26HF000001", "1C4JJXR64PW696340"},
})
```

---

### `RecallsBatchStatus` – Poll a [recall](https://carsxe.com/vehicle-recalls) batch job

**Required:**

- `batchId`

**Optional:**

- None

**Example:**

```go
status := client.RecallsBatchStatus(map[string]string{"batchId": "brb_mnablbn7_wvbaqv"})
```

---

### `RecallsBatchResults` – Fetch completed [recall](https://carsxe.com/vehicle-recalls) batch results as JSON

**Required:**

- `batchId`

**Optional:**

- None

**Example:**

```go
results := client.RecallsBatchResults(map[string]string{"batchId": "brb_mnablbn7_wvbaqv"})
```

---

### `RecallsBatchDownload` – Download completed [recall](https://carsxe.com/vehicle-recalls) batch results as CSV

**Required:**

- `batchId`

**Optional:**

- None

CSV responses are returned as `map[string]any{"csv": "<csv text>"}`. JSON error bodies are decoded as usual.

**Example:**

```go
download := client.RecallsBatchDownload(map[string]string{"batchId": "brb_mnablbn7_wvbaqv"})
```

---

### `YmmOptions` – List years, makes, models, variants, or trims for dropdowns

**Required:**

- None (no filters lists years)

**Optional:**

- `dimension` — `years` | `makes` | `models` | `trims` | `variants`
- `year`
- `make` (required for `dimension=models`)
- `model` (required for `dimension=trims`, and for `dimension=variants` unless both `year` and `make` are set)
- `trim` — substring filter on trim/variant names

**Example:**

```go
years := client.YmmOptions(map[string]string{})
makes := client.YmmOptions(map[string]string{"year": "2026"})
models := client.YmmOptions(map[string]string{"make": "Toyota"})
variants := client.YmmOptions(map[string]string{"year": "2026", "make": "Toyota", "model": "Tacoma"})
```

---

### `OwnershipVin` – Look up registered owner(s) by VIN

Enterprise only. Billed per owner returned; a `404` with `error: "no_data"` is not billed.

**Required:**

- `vin`

**Optional:**

- `include` — comma-separated subset of `demographics,emails,phones,vehicle_history`

**Example:**

```go
owners := client.OwnershipVin(map[string]string{"vin": "1FT8X3BT0BEA61538"})
```

---

### `OwnershipPerson` – Resolve contact details by name and address

Enterprise only.

**Required:**

- `first_name`
- `last_name`
- `address` — street address only
- `zip` — 5-digit US ZIP, optionally ZIP+4

**Optional:**

- `include` — comma-separated subset of `demographics,emails,phones,vehicle_history`

**Example:**

```go
person := client.OwnershipPerson(map[string]string{
	"first_name": "John",
	"last_name":  "Sample",
	"address":    "123 Example St",
	"zip":        "90210",
})
```

---

### `OwnershipAddress` – Find residents at a street address

Enterprise only.

**Required:**

- `address` — street address only
- `zip` — 5-digit US ZIP, optionally ZIP+4

**Optional:**

- `include` — comma-separated subset of `demographics,emails,phones,vehicle_history`
- `variant` — legacy alias; prefer `include`

**Example:**

```go
residents := client.OwnershipAddress(map[string]string{"address": "123 Example St", "zip": "90210"})
```

---

### `OwnershipZip` – Search people in a ZIP code with optional filters

Enterprise only. Paginated; each returned record is billed.

**Required:**

- `zip` — exactly 5 digits

**Optional:**

- `gender` — `M` or `F`
- `min_age`
- `max_age`
- `income` — letter code or full label (e.g. `F` or `$50,000-$59,999`)
- `page` — default `1`
- `limit` — default `15`, max `100`
- `include` — comma-separated subset of `demographics,emails,phones,vehicle_history`
- `variant` — legacy alias; prefer `include`

**Example:**

```go
area := client.OwnershipZip(map[string]string{"zip": "90210", "gender": "f", "min_age": "45"})
```

## Notes & Best Practices

- **Parameter requirements:** Each endpoint requires specific parameters—see the Required/Optional fields above.
- **Return values:** All responses are Go maps (`map[string]any`) for easy access and manipulation.
- **Error handling:** The client currently panics on network or JSON decode errors. Consider wrapping calls with `recover` or modifying the client to return errors for production use.
- **More info:** For advanced usage and full details, visit the [official API documentation](https://carsxe.com/docs).

---

## Overall

The CarsXE Go package provides a wide range of powerful, easy-to-use tools for accessing and integrating vehicle data into your applications and services. Whether you're a developer or a business owner, you can quickly get the information you need to take your projects to the next level—without hassle or inconvenience.
