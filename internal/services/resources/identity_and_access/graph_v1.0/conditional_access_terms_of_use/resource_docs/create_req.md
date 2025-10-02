Request URL
https://graph.microsoft.com/v1.0/agreements
Request Method
POST

{
  "displayName":"test",
  "isViewingBeforeAcceptanceRequired":true,
  "isPerDeviceAcceptanceRequired":true,
  "termsExpiration":{
    "startDateTime":"2025-10-30T01:00:00.000Z",
    "frequency":"P365D"
    },
  "userReacceptRequiredFrequency":"P90D",
  "file":{
    "localizations":
    [
      {
        "fileName":"Apple-File-System-Reference.pdf",
        "displayName":"test",
        "language":"af",
        "isDefault":true,
        "fileData":{
          "data":"JVBERi0xLjUKJeTw7fgKNiAwIG9iago8PC9GaWx0ZXIvRmxhdGVEZWNvZGUvTGVuZ3RoIDIwNj4+CnN0cmVhbQp42m1Qu2oDMRDs8xX7A1Z2R284VAScQLoEdSGFdWe7SpH8f+GVfBeMMZIWMTs7zCz9khDrEYrQyzT/0Eul51chK8bDC9WT9neAST5TXb4mZnaMHK/PN4a4soOz2nEnhsM9OnvGMZYwMQ6+KILGRWTqHEYYPCivtcHpII4Fdv3nViRv5MRgP8S/6zvtK7EJVnp1sddg6fPtAfh33oJJMD5jCxZcuqZKh/lfklSBHg1IMjnfDDGW4cWpvSUXiavNnl/XOTa0Rli1P54uoI1RNQplbmRzdHJlYW0KZW5kb2JqCjU1IDAgb2JqCjw8L0ZpbHRlci9GbGF0ZURlY29kZS9MZW5ndGggMjA5Mz4+CnN0cmVhbQp42u1cS2/cNhC+91foD5jhkMMhCSz2UKAp0Fta34oeLD96aQ7ppX+/fI1Eaje1tHKR9YZI1tJSMjmUP46+eXH4MsAgwz8YrAr/5fD4efjxfvjwEQYlBaGk4f4lXL9TKDzAcP/0+0FKRVKFmxU5qUZ5BHmQCnD6bsLXB3P84/6X4af72L3A8HHDrz9XX/7+Mw+kBiccxFFQoDTDHYFw3pSRZOhVQh5tfEi93ymF5Yq3+WPGOD5fwRepUC1bH41Uz/bokmypI2SpcRY1SngqGClhjWpFQ/CXztChMMqVGRqYn108+vEI6iDhWacZoPZSjjbPPly7g0Oec3wicSLzDJ/kUetDnGW+I/ZW/xm+OrcgQysUwrhybjp2ZNFVs/NOSGtzR8r7IF2cDaajkk/zMcrq83ekh9w+qrPywuCFdKjiMEYKg9iOg8Yc70jKM2ewo+32ermNWSpp5U6UAoAgtYApw3GCa4RlhKmG/B3MDFe9Eq5OGNCLATvSvoNeAkphL0qVEWCpoHR8LEqyHFE3ynSLEkUljMXFAB2V3+kstyL1hNEAWuGkL5TGq0zajM4kRJpMW0a5hpAYK4xZ9omgw0ftWEmBu2lT+I2yTEdUPhps9Tq/D3Q5Pvvc/nvr+UyNDg0NwolJUVPRgo="
          }
        }
      ]
    }
  }

  update

  Request URL
https://graph.microsoft.com/v1.0/agreements('ba66d872-67f2-4438-a64f-0264b13db06a')
Request Method
PATCH

{"id":"ba66d872-67f2-4438-a64f-0264b13db06a","displayName":"test-2","isViewingBeforeAcceptanceRequired":true,"isPerDeviceAcceptanceRequired":true}

update existing file

Request URL
https://graph.microsoft.com/v1.0/agreements('ba66d872-67f2-4438-a64f-0264b13db06a')/file/localizations
Request Method
POST

{
  "fileName":"Apple-File-System-Reference.pdf",
  "displayName":"test-2",
  "language":"af",
  "isDefault":true,
  "isMajorVersion":true,
  "fileData":{
    "data":"JVBERi0xLjUKJeTw7fgKNiAwIG9iago8PC9GaWx0ZXIvRmxhdGVEZWNvZGUvTGVuZ3RoIDIwNj4+CnN0cmVhbQp42m1Qu2oDMRDs8xX7A1Z2R284VAScQLoEdSGFdWe7SpH8f+GVfBeMMZIWMTs7zCz9khDrEYrQyzT/0Eul51chK8bDC9WT9neAST5TXb4mZnaMHK/PN4a4soOz2nEnhsM9OnvGMZYwMQ6+KILGRWTqHEYYPCivtcHpII4Fdv3nViRv5MRgP8S/6zvtK7EJVnp1sddg6fPtAfh33oJJMD5jCxZcuqZKh/lfklSBHg1IMjnfDDGW4cWpvSUXiavNnl/XOTa0Rli1P54uoI1RNQplbmRzdHJlYW0KZW5kb2JqCjU1IDAgb2JqCjw8L0ZpbHRlci9GbGF0ZURlY29kZS9MZW5ndGggMjA5Mz4+CnN0cmVhbQp42u1cS2/cNhC+91NvZGUvTGVuZ3RoIDE4NjYVPRgo="
  }
}

add second language

Request URL
https://graph.microsoft.com/v1.0/agreements('ba66d872-67f2-4438-a64f-0264b13db06a')/file/localizations
Request Method
POST

{
  "fileName":"Your Policy Schedule.pdf",
  "displayName":"second-language",
  "language":"ar-SA",
  "isDefault":false,
  "isMajorVersion":false,
  "fileData":{
    "data":"JVBERi0xLjQNCiW0tba3DQolDQoxIDAgb2JqDQo8PA0KL0Rlc3RzIDIgMCBSDQovUGFnZXMgMyAwIFINCi9QYWdlTGF5b3V0IC9PbmVDb2x1bW4NCi9UeXBlIC9DYXRhbG9nDQovVmlld2VyUHJlZmVyZW5jZXMgNCAwIFINCi9QYWdlTW9kZSAvVXNlTm9uZQ0KPj4NCmVuZG9iag0KOSAwIG9iag0KPDwNCi9GaWx0ZXIgL1N0YW5kYXJkDQovUCAxOTINCi9VIDxDNjdCNzA0RUoNCg0Kc3RhcnR4cmVmDQozODQyOA0KJSVFT0YNCg=="
  }
}


Delete

Request URL
https://graph.microsoft.com/v1.0/agreements('ba66d872-67f2-4438-a64f-0264b13db06a')
Request Method
DELETE