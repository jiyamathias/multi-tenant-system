# This file contains the request samples for each endpoints

## Tenant
- Tenant signup

method: **POST**

endpoint: **localhost:5002/api/v1/tenant**

```json
{
    "businessName": "Myce",
    "email": "myce@gmail.com",
    "password": "123456"
}
```
- Tenant login
  
method: **POST**

endpoint: **localhost:5002/api/v1/tenant/login**

```json
{
    "email": "myce@gmail.com",
    "password": "123456"
}
```

- Get users by tenent ID

method: **GET**

endpoint: **localhost:5002/api/v1/tenant**

## User
- User signup - pass in the tenant access token to the auth header inother to create a user

method: **POST**

endpoint: **localhost:5002/api/v1/auth/signup**

```json
{
    "firstName": "john",
    "lastName": "doe",
    "email": "johndoe@gmail.com",
    "password": "123456"
}
```

- User login

method: **POST**

endpoint: **localhost:5002/api/v1/auth/login**

```json
{
    "email": "johndoe@gmail.com",
    "password": "123456"
}
```

- Get user by ID

method: **GET**

endpoint: **localhost:5002/api/v1/auth/user/{id}**

- Update user

method: **PATCH**

endpoint: **localhost:5002/api/v1/auth/user**

```json
{
    "firstName": "jonny",
    "lastName": "drake",
}
```

## Wallet
- Get user wallet

method: **GET**

endpoint: **localhost:5002/api/v1/wallet**

## Payment
- Deposit

method: **POST**

endpoint: **localhost:5002/api/v1/payment/deposit**

```json
{
    "amount": 5000
}
```

- Transfer

method: **POST**

endpoint: **localhost:5002/api/v1/payment/transfer**

```json
{
    "bankNumber": "052",
    "accountNumber": "5376661243",
    "amount": 5000
}
```

- Bank transfer

method: **POST**

endpoint: **localhost:5002/api/v1/payment/bank-transfer**

```json
{
    "fullName": "john doe",
    "bankName": "zenith bank",
}
```

## Wallet
- Get user wallet

method: **GET**

endpoint: **localhost:5002/api/v1/wallet**

## Transaction
- Get transaction by ID

method: **GET**

endpoint: **localhost:5002/api/v1/transaction/{id}**

- Get transactions

method: **GET**

endpoint: **localhost:5002/api/v1/transaction**

- Get transactions by flow e.g revenue or withdrawal

method: **GET**

endpoint: **localhost:5002/api/v1/transaction/flow?flow=revenue**

## Audit Log
- Get all audit logs by transaction ID

method: **GET**

endpoint: **localhost:5002/api/v1/audit-log/transaction/{id}**

- Get audit log by ID

method: **GET**

endpoint: **localhost:5002/api/v1/audit-log**

## Webhook
- webhook simulation a payment provider

**NOTE:** In the request header, use the key `auth` with value `payment`. This is supposed to search as a authentication mechanism to verify if the incoming webhook request is coming from the set payment provider or not. So as to avoid making updated to the database on false.

method: **POST**

endpoint: **localhost:5002/api/v1/webhook/payment**

```json
{
    "event": "success",
    "data": {
        "status": "success",
        "reference": "crt_81cb0b68-f980-4d56-9d02-3b54919e99af",
        "amount": 5000,
        "currency": "NGN",
        "fees": 250
    }
}
```
