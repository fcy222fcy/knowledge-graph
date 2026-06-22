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
:: 1. Start Python AI Service (port 5000)
:: -------------------------------------------
echo [1/3] Starting AI Service (port 5000)...
set "AI_DIR=%ROOT%software engineering-back\python-ai-service"
start "AI Service" cmd /k "cd /d ""%AI_DIR%"" && venv\Scripts\activate && python main.py"
timeout /t 3 /nobreak >nul

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
set "FRONT_DIR=%ROOT%software engineering-frontend"
if exist "%FRONT_DIR%\package.json" (
    start "Frontend" cmd /k "cd /d ""%FRONT_DIR%"" && npm run dev"
) else (
    echo       Frontend not initialized yet. Run the following to set it up:
    echo       cd "software engineering-frontend" ^&^& npm create vite@latest . -- --template vue ^&^& npm install ^&^& npm run dev
)

echo.
echo ============================================
echo   All services started!
echo   - AI Service:   http://localhost:5000
echo   - Go Backend:   http://localhost:8080
echo   - Frontend:     http://localhost:5173
echo ============================================
echo.
echo Press any key to STOP all services...
pause >nul

:: Kill all spawned windows
taskkill /fi "WINDOWTITLE eq AI Service*" /f >nul 2>&1
taskkill /fi "WINDOWTITLE eq Go Backend*" /f >nul 2>&1
taskkill /fi "WINDOWTITLE eq Frontend*" /f >nul 2>&1

echo All services stopped.
timeout /t 2 /nobreak >nul
