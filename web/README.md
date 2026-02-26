# Frontend (UmiJS Max + ProComponents)

## Quick Start

```bash
npm install
npm run dev     # Development (http://localhost:8000, proxy to :8080)
npm run build   # Production build (output: ../internal/web/dist/)
```

## Structure

```
src/
├── app.tsx         # Runtime config (request interceptor, layout, auth)
├── access.ts       # Access control
├── constants.ts    # Constants (TOKEN_KEY, API_PREFIX)
├── services/       # API service calls
│   ├── auth.ts     # Login / Refresh token
│   └── example.ts  # Example CRUD
└── pages/
    ├── Login/      # Login page (ProForm LoginForm)
    ├── Dashboard/  # Dashboard (StatisticCard)
    └── Example/    # CRUD example (ProTable + ModalForm)
```

## Default Login

- Username: `admin`
- Password: `admin123`
