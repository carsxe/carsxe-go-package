# 🚗 CarsXE API (Go Package)

[![Go Reference](https://pkg.go.dev/badge/github.com/carsxe/carsxe-go-package.svg)](https://pkg.go.dev/github.com/carsxe/carsxe-go-package)

**CarsXE** is a powerful and developer-friendly API that gives you instant access to a wide range of vehicle data. From VIN decoding and market value estimation to vehicle history, images, OBD code explanations, and plate recognition, CarsXE provides everything you need to build automotive applications at scale.

🌐 **Website:** [https://api.carsxe.com](https://api.carsxe.com)  
📄 **Docs:** [https://api.carsxe.com/docs](https://api.carsxe.com/docs)  
📦 **All Products:** [https://api.carsxe.com/all-products](https://api.carsxe.com/all-products)

To get started with the CarsXE API, follow these steps:

1. **Sign up for a CarsXE account:**

   - [Register here](https://api.carsxe.com/register)
   - Add a [payment method](https://api.carsxe.com/dashboard/billing#payment-methods) to activate your subscription and get your API key.

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

### `Specs` – Decode VIN & get full vehicle specifications

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

### `InternationalVINDecoder` – Decode VIN with worldwide support

**Required:**

- `vin`

**Optional:**

- None

**Example:**

```go
intvin := client.InternationalVINDecoder(map[string]string{"vin": "WF0MXXGBWM8R43240"})
```

---

### `PlateDecoder` – Decode license plate info (plate, country)

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

### `MarketValue` – Estimate vehicle market value based on VIN

**Required:**

- `vin`

**Optional:**

- `state`

**Example:**

```go
marketvalue := client.MarketValue(map[string]string{"vin": "WBAFR7C57CC811956"})
```

---

### `History` – Retrieve vehicle history

**Required:**

- `vin`

**Optional:**

- None

**Example:**

```go
history := client.History(map[string]string{"vin": "WBAFR7C57CC811956"})
```

---

### `Images` – Fetch images by make, model, year, trim

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

### `Recalls` – Get safety recall data for a VIN

**Required:**

- `vin`

**Optional:**

- None

**Example:**

```go
recalls := client.Recalls(map[string]string{"vin": "1C4JJXR64PW696340"})
```

---

### `PlateImageRecognition` – Read & decode plates from images

**Required:**

- `imageURL`

**Optional:**

- None

**Example:**

```go
plateimg := client.PlateImageRecognition("https://api.carsxe.com/img/apis/plate_recognition.JPG")
```

---

### `VinOCR` – Extract VINs from images using OCR

**Required:**

- `imageURL`

**Optional:**

- None

**Example:**

```go
vinocr := client.VinOCR("https://api.carsxe.com/img/apis/plate_recognition.JPG")
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

## Notes & Best Practices

- **Parameter requirements:** Each endpoint requires specific parameters—see the Required/Optional fields above.
- **Return values:** All responses are Go maps (`map[string]any`) for easy access and manipulation.
- **Error handling:** The client currently panics on network or JSON decode errors. Consider wrapping calls with `recover` or modifying the client to return errors for production use.
- **More info:** For advanced usage and full details, visit the [official API documentation](https://api.carsxe.com/docs).

---

## Overall

The CarsXE Go package provides a wide range of powerful, easy-to-use tools for accessing and integrating vehicle data into your applications and services. Whether you're a developer or a business owner, you can quickly get the information you need to take your projects to the next level—without hassle or inconvenience.