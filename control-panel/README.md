# Control Panel

Administrative UI for issuing tickets through the issuer-service.

## Development

### Frontend (Vue/Vite)

```bash
cd frontend
npm install
npm run dev
```

**Note:** If you encounter import errors for `qrcode`, make sure to run `npm install` to install all dependencies.

The frontend will be available at `http://localhost:3000`

### Backend Server

```bash
go run main.go
```

The server will serve the built frontend on port 3001 (or PORT environment variable).

### Building for Production

```bash
cd frontend
npm run build
cd ..
go build -o control-panel
```

## Testing

### Frontend Tests

```bash
cd frontend
npm test
```

### Backend Tests

```bash
go test -v
```

## Configuration

- `PORT` - Server port (default: 3001)
- `FRONTEND_DIR` - Directory containing built frontend (default: ./frontend/dist)

## API Integration

The control panel communicates with the issuer-service at `http://localhost:8080/generate` to issue tickets.

