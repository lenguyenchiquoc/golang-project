@echo off
echo ========================================
echo TEST AUTH
echo ========================================

echo.
echo [1] Register success
curl -X POST http://localhost:8080/auth/register -H "Content-Type: application/json" -d "{\"username\": \"johndoe1\", \"email\": \"john@gmail.com\", \"password\": \"12345678910\", \"RePassword\": \"12345678910\"}"

echo.
echo [2] Register username duplicate
curl -X POST http://localhost:8080/auth/register -H "Content-Type: application/json" -d "{\"username\": \"johndoe\", \"email\": \"other@gmail.com\", \"password\": \"12345678910\", \"RePassword\": \"12345678910\"}"

echo.
echo [3] Register email duplicate
curl -X POST http://localhost:8080/auth/register -H "Content-Type: application/json" -d "{\"username\": \"janedoe\", \"email\": \"john@gmail.com\", \"password\": \"12345678910\", \"RePassword\": \"12345678910\"}"

echo.
echo [4] Register weak password
curl -X POST http://localhost:8080/auth/register -H "Content-Type: application/json" -d "{\"username\": \"janedoe\", \"email\": \"jane@gmail.com\", \"password\": \"123\", \"RePassword\": \"123\"}"

echo.
echo [5] Register password mismatch
curl -X POST http://localhost:8080/auth/register -H "Content-Type: application/json" -d "{\"username\": \"janedoe\", \"email\": \"jane@gmail.com\", \"password\": \"12345678910\", \"RePassword\": \"different\"}"

echo.
echo [6] Login success
curl -X POST http://localhost:8080/auth/login -H "Content-Type: application/json" -d "{\"username\": \"johndoe\", \"password\": \"12345678910\"}"

echo.
echo [7] Login wrong password
curl -X POST http://localhost:8080/auth/login -H "Content-Type: application/json" -d "{\"username\": \"johndoe\", \"password\": \"wrongpass\"}"

echo.
echo [8] Login user not found
curl -X POST http://localhost:8080/auth/login -H "Content-Type: application/json" -d "{\"username\": \"nobody\", \"password\": \"12345678910\"}"

echo.
echo ========================================
echo DONE
echo ========================================
pause