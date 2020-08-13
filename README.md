# Simple email validator service in Go

## Installation

`go get github.com/tony2001/go-mail-validate`

## Usage

`go-mail-validate` runtime flags:
- -d - enable debug logging
- -s - SMTP timeout in milliseconds (1000ms by default)
- -c - [Clearout API](https://clearout.io/) timeout in milliseconds (1000ms by default)

Environment variables used:
- `PORT` - port to listen (8080 by default)
- `CLEAROUT_TOKEN` - Clearout API token (empty by default, which means Clearout API is disabled)

## Available Endpoints

*/email/validate* - accepts request in the following form:
```
{ "email": "address@here.com" }
```

Response example:
```javascript
{
 "Valid":false,
 "ValidatorResults":
 [
  {
   "Name":"rfc5322regex",
   "Valid":true,
   "Reason":""
  },
  {
   "Name":"domainNotReserved",
   "Valid":true,
   "Reason":""
  },
  {
   "Name":"mxRecord",
   "Valid":true,
   "Reason":""
  },
  {
   "Name":"smtpTest",
   "Valid":false,
   "Reason":"550 5.7.606 Access denied, banned sending IP [X.X.X.X]. To request removal from this list please visit https://sender.office.com/ and follow the directions. For more information please go to  http://go.microsoft.com/fwlink/?LinkID=526655 AS(1430) [HE1EUR02FT008.eop-EUR02.prod.protection.outlook.com]"
  }
 ]
}
```

## Currently supported validators
- RFC5322-compliant email address regex parser
- Common (simpler) email regex parser
- International email regex parser
- Reserved domain name check
- MX domain check
- SMTP server check
- Clearout API instant check
