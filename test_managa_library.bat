@echo off
set TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImpvaG5AZ21haWwuY29tIiwiZXhwIjoxNzc3MzA0MTcyLCJpYXQiOjE3NzczMDM1NzIsInVzZXJfaWQiOiJ1c3JfMTc3NzMwMjE2NzQ0NzI5MzIwMCIsInVzZXJuYW1lIjoiam9obmRvZSJ9.SEVuND8odyeGpD5khkQkqcKg83Iegf9JRLoBakMRifI

echo ========================================
echo TEST MANGA
echo ========================================

echo.
echo [1] Add manga success
curl -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"One Piece\", \"author\": \"Oda Eiichiro\", \"genres\": \"Action,Adventure\", \"status\": \"ongoing\", \"total_chapters\": 1100, \"description\": \"Luffy wants to be pirate king\"}"

echo.
echo [2] Add manga duplicate
curl -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"One Piece\", \"author\": \"Oda Eiichiro\", \"genres\": \"Action\", \"status\": \"ongoing\", \"total_chapters\": 1100, \"description\": \"duplicate\"}"

echo.
echo [3] Add manga without token
curl -X POST http://localhost:8080/manga -H "Content-Type: application/json" -d "{\"title\": \"Naruto\", \"author\": \"Kishimoto\", \"genres\": \"Action\", \"status\": \"completed\", \"total_chapters\": 700, \"description\": \"Ninja story\"}"

echo.
echo [4] Add second manga
curl -X POST http://localhost:8080/manga -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"title\": \"Naruto\", \"author\": \"Kishimoto\", \"genres\": \"Action,Adventure\", \"status\": \"completed\", \"total_chapters\": 700, \"description\": \"Ninja story\"}"

echo.
echo [5] Search by name
curl "http://localhost:8080/manga?query=one"

echo.
echo [6] Search by genre
curl "http://localhost:8080/manga?genre=Action"

echo.
echo [7] Search by status
curl "http://localhost:8080/manga?status=ongoing"

echo.
echo [8] Search combined
curl "http://localhost:8080/manga?query=one&genre=Action&status=ongoing"

echo.
echo [9] Search not found
curl "http://localhost:8080/manga?query=notexist"

echo.
echo [10] Get all manga
curl "http://localhost:8080/manga"

echo.
echo ========================================
echo GET MANGA ID FROM SEARCH RESULT ABOVE
echo THEN UPDATE MANGA_ID BELOW
echo ========================================
echo.

set /p MANGA_ID=Enter manga_id from search result: 

echo ========================================
echo TEST LIBRARY
echo ========================================

echo.
echo [11] Get manga detail
curl "http://localhost:8080/manga/%MANGA_ID%"

echo.
echo [12] Get manga not found
curl "http://localhost:8080/manga/not-exist"

echo.
echo [13] Add to library success
curl -X POST http://localhost:8080/users/library -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"manga_id\": \"%MANGA_ID%\", \"status\": \"reading\", \"current_chapter\": 1}"

echo.
echo [14] Add to library duplicate
curl -X POST http://localhost:8080/users/library -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"manga_id\": \"%MANGA_ID%\", \"status\": \"reading\", \"current_chapter\": 1}"

echo.
echo [15] Add manga not exist to library
curl -X POST http://localhost:8080/users/library -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"manga_id\": \"not-exist\", \"status\": \"reading\", \"current_chapter\": 0}"

echo.
echo [16] Add invalid status
curl -X POST http://localhost:8080/users/library -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"manga_id\": \"%MANGA_ID%\", \"status\": \"invalid\", \"current_chapter\": 0}"

echo.
echo [17] View library
curl http://localhost:8080/users/library -H "Authorization: Bearer %TOKEN%"

echo.
echo [18] Update progress success
curl -X PUT http://localhost:8080/users/progress -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"manga_id\": \"%MANGA_ID%\", \"current_chapter\": 50, \"status\": \"reading\"}"

echo.
echo [19] Update progress negative chapter
curl -X PUT http://localhost:8080/users/progress -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"manga_id\": \"%MANGA_ID%\", \"current_chapter\": -1, \"status\": \"reading\"}"

echo.
echo [20] Update progress manga not in library
curl -X PUT http://localhost:8080/users/progress -H "Content-Type: application/json" -H "Authorization: Bearer %TOKEN%" -d "{\"manga_id\": \"not-exist\", \"current_chapter\": 50, \"status\": \"reading\"}"

echo.
echo [21] Update progress no token
curl -X PUT http://localhost:8080/users/progress -H "Content-Type: application/json" -d "{\"manga_id\": \"%MANGA_ID%\", \"current_chapter\": 50}"

echo.
echo [22] View library no token
curl http://localhost:8080/users/library

echo.
echo ========================================
echo DONE
echo ========================================
pause