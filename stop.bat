@echo off
chcp 65001 >nul 2>&1
echo Stopping all SE Platform services...

:: Kill by window title
taskkill /fi "WINDOWTITLE eq MinIO Service*" /f >nul 2>&1
taskkill /fi "WINDOWTITLE eq Go Backend*" /f >nul 2>&1
taskkill /fi "WINDOWTITLE eq Frontend*" /f >nul 2>&1

:: Kill by port (fallback)
for /f "tokens=5" %%a in ('netstat -aon ^| findstr :9000 ^| findstr LISTENING') do taskkill /pid %%a /f >nul 2>&1
for /f "tokens=5" %%a in ('netstat -aon ^| findstr :8080 ^| findstr LISTENING') do taskkill /pid %%a /f >nul 2>&1
for /f "tokens=5" %%a in ('netstat -aon ^| findstr :5173 ^| findstr LISTENING') do taskkill /pid %%a /f >nul 2>&1

echo Done. All services stopped.
timeout /t 2 /nobreak >nul
