# Connect

Личный кабинет и шлюз SMS. Клиент говорит с Connect по REST или gRPC. Шлюз проверяет свои правила и только после этого вызывает HTTPS API EasySMS: `POST /api/send-sms`, `/api/create-otp`, `/api/validate-otp`.

## Состав

- `apps/web` — кабинет на Next.js и TypeScript. К базе ходит через SQL.
- `services/gateway` — Go: HTTP `:8080`, gRPC `:50051` и миграции goose при старте.
- `proto/connect/v1/gateway.proto` — контракт gRPC.
- Postgres, по умолчанию порт `5434`.

## Запуск

```powershell
docker compose up -d postgres
cd services/gateway
go run ./cmd/seed
go run ./cmd/gateway
```

В другом окне:

```powershell
cd apps/web
npm install
npm run dev
```

Кабинет: http://localhost:3000

Демо-клиент: `demo@connect.local` / `demo12345`  
Администратор: `admin@connect.local` / `admin12345`  
Ключ API: `cnt_demo_7f3a9c2e1b`

Пока `EASYSMS_API_KEY` пустой, шлюз не звонит поставщику и отмечает отправку локально. Боевой ключ и адрес задаются в `.env`.
