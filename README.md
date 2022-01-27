# Fully Isolated System Architecture (Microservices)

## Arsitektur

```
  Request
    |
    |
    |
Api Gateway --- Auth Provider
    |\__________________________
    |           |               |
    |           |               |
 Service-1   Service-2      Service-3
```

## File Konfigurasi

Setiap service terdapat file konfigurasi (`.env`) yang dapat disesuaikan dengan runtime environment

## Cara Install

- Buat database `microservices`

- Ubah file `.env` yang terdapat pada service security (`security/.env`), sesuaikan dengan konfigurasi database

- Build dan jalankan `docker-compose up` atau jika menggunakan task cukup dengan `task start`

- Seeding data melalui endpoint `POST http://localhost:2727/api/seed`

- Login melalui endpoint `POST http://localhost:2727/api/login`

```json
{
    "email": "surya.iksanudin@gmail.com",
    "password": "admin"
}
```

Maka akan mengembalikan response

```json
{
    "token": "[TOKEN]"
}
```

- Panggil endpoint `GET http://localhost:2727/api/hello` dengan header `Authorization: Bearer [TOKEN]` maka akan mengembalikan response

```json
{
    "message": "Hello Service 1"
}
```

- Panggil endpoint `GET http://localhost:2727/api/hey` dengan header `Authorization: Bearer [TOKEN]` untuk simulasi internal service call yaitu call service 2 dari service 3, maka akan mengembalikan response

```json
{
	"message": "Service 3 called service 2 with response: Hello Service 2"
}
```

## Endpoint

- Publik

```bash
POST   /login 
POST   /seed
```

- Privat (Butuh header authorization)

```bash
POST   /api/users
PUT    /api/users/:id
DELETE /api/users/:id
GET    /api/users/:id
GET    /api/users 
POST   /api/validate
GET    /api/hello
GET    /api/hi
GET    /api/hey
```

## Pengembangan Independen

Setiap service dapat dikembangkan secara independen tanpa perlu terkoneksi dengan Api Gateway atau Arsitektur secara utuh. Untuk melakukan simulasi privat endpoint (Api) cukup dengan mengirimkan header `X-User-ID` dan `X-User-Email` pada request.

Ini karena Api Gateway dirancang untuk melakukan validasi token dan kemudian mengubah token menjadi header `X-User-ID` dan `X-User-Email` melalui endpoint `/api/validate` yang mengembalikan informasi user. Header (`X-User-ID` dan `X-User-Email`) kemudian diteruskan ke endpoint service (sesuai dengan mapping) sebagai tanda bahwa request telah divalidasi oleh Api Gateway
