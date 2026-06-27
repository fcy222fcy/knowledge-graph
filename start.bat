@echo off
chcp 65001 >nul 2>&1
title SE Platform - All Services

echo ============================================
echo   SE智图问答平台 - 启动所有服务
echo ============================================
echo.

:: 项目根目录
set "ROOT=%~dp0"

:: -------------------------------------------
:: 1. 启动 Neo4j 服务 (端口 7474/7687)
:: -------------------------------------------
echo [1/4] 启动 Neo4j 服务 (端口 7474/7687)...
D:\neo4j-community-5.26.0\bin\neo4j.bat console >nul 2>&1
timeout /t 3 /nobreak >nul

:: -------------------------------------------
:: 2. 启动 Go 后端 (端口 8080)
:: -------------------------------------------
echo [2/4] 启动 Go 后端 (端口 8080)...
set "BACK_DIR=%ROOT%software engineering-back"
start "Go Backend" cmd /k "cd /d ""%BACK_DIR%"" && go run ./cmd/server"
timeout /t 3 /nobreak >nul

:: -------------------------------------------
:: 3. 启动前台学生端 (端口 4173)
:: -------------------------------------------
echo [3/4] 启动前台学生端 (端口 4173)...
set "FRONT_DIR=%ROOT%software engineering-fronted\frontend-vue-app"
if exist "%FRONT_DIR%\package.json" (
    start "Frontend" cmd /k "cd /d ""%FRONT_DIR%"" && npm run dev"
) else (
    echo       前端项目未初始化，请先运行: npm install
)
timeout /t 2 /nobreak >nul

:: -------------------------------------------
:: 4. 启动后台管理端 (端口 4174)
:: -------------------------------------------
echo [4/4] 启动后台管理端 (端口 4174)...
set "ADMIN_DIR=%ROOT%software engineering-fronted\admin-vue-app"
if exist "%ADMIN_DIR%\package.json" (
    start "Admin" cmd /k "cd /d ""%ADMIN_DIR%"" && npm run dev"
) else (
    echo       后台项目未初始化，请先运行: npm install
)

echo.
echo ============================================
echo   所有服务已启动！
echo   - Neo4j:        http://localhost:7474 (Web) / bolt://localhost:7687
echo   - Go Backend:   http://localhost:8080
echo   - 前台学生端:   http://localhost:4173
echo   - 后台管理端:   http://localhost:4174
echo   - Ollama:       http://localhost:11434 (请确保已启动)
echo.
echo   Neo4j 登录:    neo4j / neo4j
echo ============================================
echo.
echo 按任意键停止所有服务...
pause >nul

:: 停止所有服务
taskkill /fi "WINDOWTITLE eq Go Backend*" /f >nul 2>&1
taskkill /fi "WINDOWTITLE eq Frontend*" /f >nul 2>&1
taskkill /fi "WINDOWTITLE eq Admin*" /f >nul 2>&1

echo 所有服务已停止。
timeout /t 2 /nobreak >nul
