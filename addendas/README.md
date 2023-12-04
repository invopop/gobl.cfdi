# CFDI Addendas

"Addendas" add functionality to regular CFDI documents so that private companies can leverage existing infrastructure around the CFDI format and SAT to extract additional structured data.

Each of the addendas currently supported are listed below, with instructions on the mappings and key fields.

## MABE

Most of the MABE Addenda fields are determined automatically from the base GOBL Invoice, with the exception of the following:

| MABE Field                          | GOBL Invoice Property          | GOBL Invoice Value                                    | Description                                                               |
| ----------------------------------- | ------------------------------ | ----------------------------------------------------- | ------------------------------------------------------------------------- |
| Purchase Order (Orden de Compra)    | `ordering.identities`          | `[{"key":"mx-mabe-purchase-order", "code":"-CODE-"}]` | Provided by Mabe for the order                                            |
| Provider Code (Código de Proveedor) | `supplier.identities`          | `[{"key":"mx-mabe-provider-code", "code":"-CODE-"}]`  | Code issued by Mabe to identify the supplier                              |
| Delivery Plant (Planta de Entrega)  | `delivery.receiver.identities` | `[{"key":"mx-mabe-delivery-plant","code":"-CODE-"}]`  | Delivery Plant Code                                                       |
| Article Code (Código de Artículo)   | `lines[i].item.identities`     | `[{"key":"mx-mabe-article-code","code":"-CODE-"}]`    | Article code provided by Mabe                                             |
| Unit (Unidad)                       | `lines[i].item.identities`     | `[{"key":"mx-mabe-unit","code":"-CODE-"}]`            | Item unit code provided by Mabe (falls back to CFDI/UNECE unit)           |
| Reference 1                         | `ordering.identities`          | `[{"key":"mx-mabe-reference1","code":"-CODE-"}]`      | Additional code required by Mabe in certain circumstances while ordering. |
| Reference 2                         | `ordering.identities`          | `[{"key":"mx-mabe-reference2","code":"-CODE-"}]`      | Set to `NA` by default as not currently used by Mabe.                     |
