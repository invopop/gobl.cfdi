# GOBL to CFDI Conversion

Convert GOBL documents in the Mexican CFDI (Comprobante Fiscal Digital por Internet) format.

Copyright [Invopop Ltd.](https://invopop.com) 2023. Released publicly under the [Apache License Version 2.0](LICENSE). For commercial licenses please contact the [dev team at invopop](mailto:dev@invopop.com). In order to accept contributions to this library we will require transferring copyrights to Invopop Ltd.

[![Lint](https://github.com/invopop/gobl.cfdi/actions/workflows/lint.yaml/badge.svg)](https://github.com/invopop/gobl.cfdi/actions/workflows/lint.yaml)
[![Test Go](https://github.com/invopop/gobl.cfdi/actions/workflows/test.yaml/badge.svg)](https://github.com/invopop/gobl.cfdi/actions/workflows/test.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/invopop/gobl.cfdi)](https://goreportcard.com/report/github.com/invopop/gobl.cfdi)
[![GoDoc](https://godoc.org/github.com/invopop/gobl.cfdi?status.svg)](https://godoc.org/github.com/invopop/gobl.cfdi)
![Latest Tag](https://img.shields.io/github/v/tag/invopop/gobl.cfdi)

## Usage

### Go Package

Usage of the CFDI conversion library is quite straight forward. You must first have a GOBL Envelope including an invoice for Mexico ready to convert. There are some samples here in the test/data directory.

```go
package main

import (
    "os"

    "github.com/invopop/gobl"
    "github.com/invopop/gobl.cfdi"
)

func main {
    data, _ := os.ReadFile("./test/data/invoice.json")

    env := new(gobl.Envelope)
    if err := json.Unmarshal(data, env); err != nil {
        panic(err)
    }

    // Prepare the CFDI document
    doc, err := cfdi.NewDocument(env)
    if err != nil {
        panic(err)
    }

    // Create the XML output
    out, err := doc.Bytes()
    if err != nil {
        panic(err)
    }

    // TODO: do something with the output
}
```

### Testing

This package uses [lestrrat-go/libxml2](https://github.com/lestrrat-go/libxml2) for testing purporses, which in turn depends on the libxml-2.0 C library. Ensure you have the development dependency installed. In linux this implies:

```bash
sudo apt-get install libxml2-dev
```

Tests can take a while to run as they download the complete XML documents to test against, please be patient.

## Addendas

For details on support for converting Addendas, please see the [addendas package](addendas).

## Command Line

The GOBL to CFDI tool also includes a command line helper. You can find pre-built [gobl.cfdi binaries](https://github.com/invopop/gobl.cfdi/releases) in the github repository, or install manually in your Go environment with:

```bash
go install github.com/invopop/gobl.cfdi
```

Usage is very straightforward:

```bash
gobl.cfdi convert ./test/data/invoice.json
```

Which should produce something like:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<cfdi:Comprobante xmlns:cfdi="http://www.sat.gob.mx/cfd/4" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.sat.gob.mx/cfd/4 http://www.sat.gob.mx/sitio_internet/cfd/4/cfdv40.xsd" Version="4.0" TipoDeComprobante="I" Serie="LMC" Folio="0010" Fecha="2023-05-29T00:00:00" LugarExpedicion="26015" SubTotal="400.40" Descuento="200.20" Total="232.23" Moneda="MXN" Exportacion="01" MetodoPago="PUE" FormaPago="03" CondicionesDePago="Pago a 30 días." Sello="" NoCertificado="00000000000000000000" Certificado="">
  <cfdi:Emisor Rfc="EKU9003173C9" Nombre="ESCUELA KEMPER URGATE" RegimenFiscal="601"></cfdi:Emisor>
  <cfdi:Receptor Rfc="URE180429TM6" Nombre="UNIVERSIDAD ROBOTICA ESPAÑOLA" DomicilioFiscalReceptor="65000" RegimenFiscalReceptor="601" UsoCFDI="G01"></cfdi:Receptor>
  <cfdi:Conceptos>
    <cfdi:Concepto ClaveProdServ="50211502" Cantidad="2" ClaveUnidad="H87" Descripcion="Cigarros" ValorUnitario="200.2020" Importe="400.4040" Descuento="200.2020" ObjetoImp="02">
      <cfdi:Impuestos>
        <cfdi:Traslados>
          <cfdi:Traslado Base="200.2020" Importe="32.0323" Impuesto="002" TasaOCuota="0.160000" TipoFactor="Tasa"></cfdi:Traslado>
        </cfdi:Traslados>
      </cfdi:Impuestos>
    </cfdi:Concepto>
  </cfdi:Conceptos>
  <cfdi:Impuestos TotalImpuestosTrasladados="32.03">
    <cfdi:Traslados>
      <cfdi:Traslado Base="200.20" Importe="32.03" Impuesto="002" TasaOCuota="0.160000" TipoFactor="Tasa"></cfdi:Traslado>
    </cfdi:Traslados>
  </cfdi:Impuestos>
</cfdi:Comprobante>
```
