@echo off
chcp 65001 >nul 2>&1
echo 正在停止 SE智图问答平台所有服务...

:: 停止 Neo4j
echo 正在停止 Neo4j...
D:\neo4j-community-5.26.0\bin\neo4j.bat stop >nul 2>&1

:: 按窗口标题停止
taskkill /fi "WINDOWTITLE eq Go Backend*" /f >nul 2>&1
taskkill /fi "WINDOWTITLE eq Frontend*" /f >nul 2>&1
taskkill /fi "WINDOWTITLE eq Admin*" /f >nul 2>&1

:: 按端口停止（备用方案）
for /f "tokens=5" %%a in ('netstat -aon ^| findstr :8080 ^| findstr LISTENING') do taskkill /pid %%a /f >nul 2>&1
for /f "tokens=5" %%a in ('netstat -aon ^| findstr :4173 ^| findstr LISTENING') do taskkill /pid %%a /f >nul 2>&1
for /f "tokens=5" %%a in ('netstat -aon ^| findstr :4174 ^| findstr LISTENING') do taskkill /pid %%a /f >nul 2>&1

echo.
echo 所有服务已停止。
timeout /t 2 /nobreak >nul
