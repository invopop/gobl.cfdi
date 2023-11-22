# CFDI Addendas

"Addendas" add functionality to regular CFDI documents so that private companies can leverage existing infrastructure around the CFDI format and SAT to extract additional structured data.

Each of the addendas currently supported are listed below, with instructions on the mappings and key fields.

## MABE

Most of the MABE Addenda fields are determined automatically from the base GOBL Invoice, with the exception of the following:

| MABE Field     | GOBL Invoice Property          | GOBL Invoice Value                              | Description                                                               |
| -------------- | ------------------------------ | ----------------------------------------------- | ------------------------------------------------------------------------- |
| Order Code     | `ordering.code`                | `"-CODE-"`                                      | Provided by Mabe for the order                                            |
| Provider Code  | `supplier.identities`          | `[{"key":"mx-mabe-provider", "code":"-CODE-"}]` | Code issued by Mabe to identify the supplier                              |
| Delivery Plant | `delivery.receiver.identities` | `[{"key":"mx-mabe-plant","code":"-CODE-"}]`     | Delivery Plant Code                                                       |
| Item Code      | `lines[i].item.identities`     | `[{"key":"mx-mabe-item","code":"-CODE-"}]`      | Article code provided by Mabe                                             |
| Reference 1    | `ordering.identities`          | `[{"key":"mx-mabe-ref1","code":"-CODE-"}]`      | Additional code required by Mabe in certain circumstances while ordering. |
| Reference 2    | `ordering.identities`          | `[{"key":"mx-mabe-ref2","code":"-CODE-"}]`      | Set to `NA` by default as not currently used by Mabe.                     |
