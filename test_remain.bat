@echo off
set TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImpvaG5AZ21haWwuY29tIiwiZXhwIjoxNzc3MzA0MTcyLCJpYXQiOjE3NzczMDM1NzIsInVzZXJfaWQiOiJ1c3JfMTc3NzMwMjE2NzQ0NzI5MzIwMCIsInVzZXJuYW1lIjoiam9obmRvZSJ9.SEVuND8odyeGpD5khkQkqcKg83Iegf9JRLoBakMRifI
set MANGA_ID=manga_1777302746617845700

echo ========================================
echo TEST JWT MIDDLEWARE
echo ========================================

echo.
echo [1] No token
curl http://localhost:8080/users/library

echo.
echo [2] Wrong format - missing Bearer
curl http://localhost:8080/users/library -H "Authorization: %TOKEN%"

echo.
echo [3] Fake token
curl http://localhost:8080/users/library -H "Authorization: Bearer faketoken123"

echo.
echo [4] Valid token
curl http://localhost:8080/users/library -H "Authorization: Bearer %TOKEN%"

echo.
echo ========================================
echo TEST UDP NOTIFICATION
echo ========================================

echo.
echo [5] Send notification - success
curl -X POST http://localhost:8080/notify -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"manga_id\": \"%MANGA_ID%\", \"manga_title\": \"One Piece\", \"message\": \"Chapter 1101 is out!\"}"

echo.
echo [6] Send notification - no token
curl -X POST http://localhost:8080/notify -H "Content-Type: application/json" -d "{\"manga_id\": \"%MANGA_ID%\", \"manga_title\": \"One Piece\", \"message\": \"Chapter 1101 is out!\"}"

echo.
echo [7] Send notification - missing fields
curl -X POST http://localhost:8080/notify -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"manga_id\": \"%MANGA_ID%\"}"

echo.
echo ========================================
echo TEST gRPC via HTTP
echo ========================================

echo.
echo [8] Search manga via gRPC - by name
curl "http://localhost:8080/manga?query=one"

echo.
echo [9] Search manga via gRPC - by genre
curl "http://localhost:8080/manga?genre=Action"

echo.
echo [10] Search manga via gRPC - pagination
curl "http://localhost:8080/manga?page=1&limit=1"

echo.
echo [11] Get manga via gRPC - success
curl "http://localhost:8080/manga/%MANGA_ID%"

echo.
echo [12] Get manga via gRPC - not found
curl "http://localhost:8080/manga/not-exist-id"

echo.
echo [13] Update progress via gRPC - success
curl -X PUT http://localhost:8080/users/progress -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"manga_id\": \"%MANGA_ID%\", \"current_chapter\": 100, \"status\": \"reading\"}"

echo.
echo [14] Update progress via gRPC - negative chapter
curl -X PUT http://localhost:8080/users/progress -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"manga_id\": \"%MANGA_ID%\", \"current_chapter\": -1, \"status\": \"reading\"}"

echo.
echo [15] Update progress via gRPC - not in library
curl -X PUT http://localhost:8080/users/progress -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"manga_id\": \"not-exist-id\", \"current_chapter\": 50, \"status\": \"reading\"}"

echo.
echo ========================================
echo TEST WEBSOCKET
echo ========================================
echo WebSocket must be tested manually via client:
echo   go run cmd/client/main.go
echo   Login -^> Choose 5. Chat -^> room: general
echo   Test commands: /users, /rooms, /join ^<room^>, /dm ^<user^> ^<msg^>, /quit
echo.

echo ========================================
echo TEST TCP BROADCAST
echo ========================================
echo TCP must be tested manually:
echo   Terminal 1: go run cmd/main.go
echo   Terminal 2: go run cmd/client/main.go -^> Login johndoe
echo   Terminal 3: go run cmd/client/main.go -^> Login janedoe
echo   Terminal 2: Choose 4. Update progress
echo   Terminal 3: Should receive broadcast
echo.

echo ========================================
echo DONE
echo ========================================
pause