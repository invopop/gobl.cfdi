# CFDI Addendas

"Addendas" add functionality to regular CFDI documents so that private companies can leverage existing infrastructure around the CFDI format and SAT to extract additional structured data.

Each of the addendas currently supported are listed below, with instructions on the mappings and key fields.

## MABE

Most of the MABE Addenda fields are determined automatically from the base GOBL Invoice, with the exception of the following:

| MABE Field     | GOBL Invoice Value                                                          | Description                                                               |
| -------------- | --------------------------------------------------------------------------- | ------------------------------------------------------------------------- |
| Order Code     | `ordering.code = "-CODE-"`                                                  | Provided by Mabe for the order                                            |
| Provider Code  | `supplier.identities = [{"type":"MABE", "code":"-CODE-"}]`                  | Code issued by Mabe to identify the supplier                              |
| Delivery Plant | `delivery.receiver.identities = [{"type":"MABE-PLANT-ID","code":"-CODE-"}]` | Delivery Plant ID                                                         |
| Item Code      | `lines[i].item.identities = [{"type":"MABE","code":"-CODE-"}]`              | Article code provided by Mabe                                             |
| Reference 1    | `ordering.identities = [{"type":"MABE-REF1","code":"-CODE-"}]`              | Additional code required by Mabe in certain circumstances while ordering. |
| Reference 2    | NA                                                                          | Always empty as not currently used by Mabe.                               |
