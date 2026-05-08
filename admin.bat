@echo off
setlocal enabledelayedexpansion

:: --- THÔNG TIN CẤU HÌNH ---
set "BASE_URL=http://localhost:8080"
set "USER=adminins"
set "PASS=12345678910"
set "EMAIL=admin2@gmail.com"

echo ========================================
echo 1. DANG KY TAI KHOAN...
curl -X POST "%BASE_URL%/auth/register" ^
     -H "Content-Type: application/json" ^
     -d "{\"username\": \"%USER%\", \"email\": \"%EMAIL%\", \"password\": \"%PASS%\", \"RePassword\": \"%PASS%\"}"

echo.
echo.
echo ========================================
echo 2. DANG NHAP VA LAY TOKEN...

:: Ghi response vao file tam login_res.txt de tranh loi syntax CMD
curl -s -X POST "%BASE_URL%/auth/login" ^
     -H "Content-Type: application/json" ^
     -d "{\"username\": \"%USER%\", \"password\": \"%PASS%\"}" > login_res.txt

:: Doc noi dung tu file vao bien response
set /p response=<login_res.txt

:: Xử lý tách Token
:: Cắt từ chữ token":" trở đi
set "search_str=token\":\""
set "tmp=!response:*token\":\"=!"

:: Lấy chuỗi cho đến dấu ngoặc kép tiếp theo
for /f "delims=\"" %%a in ("!tmp!") do (
    set "TOKEN=%%a"
)

:: Xóa file tạm sau khi lấy xong token
if exist login_res.txt del login_res.txt

if "!TOKEN!"=="" (
    echo [LOI] Khong lay duoc Token. Response nhan duoc:
    echo !response!
    pause
    exit /b
)

echo Dang nhap thanh cong.
echo Token: !TOKEN:~0,25!...
echo.

echo ========================================
echo 3. THEM TRUYEN MANGA...
curl -s -X POST "%BASE_URL%/manga" ^
     -H "Content-Type: application/json" ^
     -H "Authorization: Bearer !TOKEN!" ^
     -d "{\"title\": \"Naruto\", \"author\": \"Kishimoto Masashi\", \"genres\": \"Action,Adventure,Shounen\", \"status\": \"completed\", \"total_chapters\": 700, \"description\": \"A young ninja seeks recognition and dreams of becoming Hokage\"}"

echo.
echo.
echo ========================================
echo DA HOAN THANH!
pause