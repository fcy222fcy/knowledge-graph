@echo off
chcp 65001 >nul 2>&1
title SE Platform - All Services

echo ============================================
echo   Software Engineering QA Platform
echo   Starting all services...
echo ============================================
echo.

:: Project root (with trailing backslash from %~dp0)
set "ROOT=%~dp0"

:: -------------------------------------------
:: 1. Start MinIO Service (port 9000/9001)
:: -------------------------------------------
echo [1/3] Starting MinIO Service (port 9000/9001)...
if exist "E:\minio\bin\minio.exe" (
    start "MinIO Service" cmd /k "E:\minio\bin\minio.exe server E:\minio\data --console-address :9001"
    timeout /t 2 /nobreak >nul
) else (
    echo       MinIO not found at E:\minio\bin\minio.exe
    echo       Please install MinIO first or update the path in start.bat
)

:: -------------------------------------------
:: 2. Start Go Backend (port 8080)
:: -------------------------------------------
echo [2/3] Starting Go Backend (port 8080)...
set "BACK_DIR=%ROOT%software engineering-back"
start "Go Backend" cmd /k "cd /d ""%BACK_DIR%"" && go run ./cmd/server"
timeout /t 3 /nobreak >nul

:: -------------------------------------------
:: 3. Start Vue Frontend (port 5173)
:: -------------------------------------------
echo [3/3] Starting Frontend (port 5173)...
set "FRONT_DIR=%ROOT%software engineering-fronted\frontend-vue-app"
if exist "%FRONT_DIR%\package.json" (
    start "Frontend" cmd /k "cd /d ""%FRONT_DIR%"" && npm run dev"
) else (
    echo       Frontend not initialized yet. Run the following to set it up:
    echo       cd "software engineering-fronted\frontend-vue-app" ^&^& npm install ^&^& npm run dev
)

echo.
echo ============================================
echo   All services started!
echo   - MinIO:        http://localhost:9000 (API) / http://localhost:9001 (Console)
echo   - Go Backend:   http://localhost:8080
echo   - Frontend:     http://localhost:5173
echo   - Ollama:       http://localhost:11434 (ensure it's running)
echo.
echo   MinIO Console Login:
echo     Username: minioadmin
echo     Password: minioadmin
echo ============================================
echo.
echo Press any key to STOP all services...
pause >nul

:: Kill all spawned windows
taskkill /fi "WINDOWTITLE eq MinIO Service*" /f >nul 2>&1
taskkill /fi "WINDOWTITLE eq Go Backend*" /f >nul 2>&1
taskkill /fi "WINDOWTITLE eq Frontend*" /f >nul 2>&1

echo All services stopped.
timeout /t 2 /nobreak >nul
