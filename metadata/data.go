package metadata

// This document is Licensed under Creative Commons CC0.
// To the extent possible under law, the author(s) have dedicated all copyright and related and neighboring rights
// to this document to the public domain worldwide.
// This document is distributed without any warranty.
// You should have received a copy of the CC0 Public Domain Dedication along with this document.
// If not, see https://creativecommons.org/publicdomain/zero/1.0/legalcode.

var ActivationKey = "IISABPDENLCAEIKFLMBORDQQNLMLBLKJLQELSIJPIESFIAQAJINCMHLDLALMSLAM"

var SiteName = "WellWish\nCloud Decision Engine"

var SiteUrl = "http://127.0.0.1:7777"

var CompanyName = "Example Corporation (SAMPLE)"

var CompanyInfo = `Example Inc.
1010 Corporate Avenue, San Jose, CA, 55555, USA
TAX ID: 1234-56 Payment: ACH Routing# 12345 Account# 12345 https://cashbuddy.example.com/12345
`

var PaymentPattern = "https://example.com/%s"

// InvoicePattern Please update the pattern with your locally regulated purchase order format.
var OrderPattern = `
Company: %s
Billing address: %s
Billing email: %s
Our company places the following order.
The payment term is Net 30.                 
Ordering %s remoting vouchers for %s each.

The order total is %s.
The final amount includes Sales Tax of %s percent.
Satisfaction guarantee. Cancel or refund within 30 days.
Notes:
`

// InvoicePattern Please update the pattern with your locally regulated invoice format.
var InvoicePattern = `              INVOICE              

Payee: %s
Date: %s        Invoice Number: %s

Payer: %s
Billing address: %s
Billing email: %s
Please pay the following remoting vouchers.
The payment term is Net 30.
Satisfaction guarantee.
Order is cancellable within 30 days.
Order is refundable within 30 days, if paid.

Ordered %s cloud vouchers for %s each
               Invoiced Total %s.

The total amount includes sales tax of 0 percent.
Order Status:
%s.
`

var UnitPrice = "USD 1.03"

var VoucherPattern = `              SERVICE VOUCHER              

From: %s
Issue Date: %s
This voucher can be used at the servicing company listed above.
It is valid for 365 days from the time of issuance.
Invoice: %s
The voucher status is %s.
`
