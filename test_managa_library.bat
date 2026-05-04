@echo off
set TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImpvaG5AZ21haWwuY29tIiwiZXhwIjoxNzc3ODY2NzMxLCJpYXQiOjE3Nzc4NjYxMzEsInVzZXJfaWQiOiJ1c3JfMTc3Nzg2NjEzMTExNDUxOTQwMCIsInVzZXJuYW1lIjoiam9obmRvZSJ9.p3jx606Xi4xbVDo63Fqzelc3UVUnA2kQ4ZuoGF2iblg
echo TEST MANGA (NEW SCHEMA)
echo ========================================

echo.
echo [1] Add manga success (multi genres)
curl -X POST http://localhost:8080/manga ^
-H "Content-Type: application/json" ^
-H "Authorization: Bearer %TOKEN%" ^
-d "{\"title\":\"One Piece\",\"author\":\"Oda Eiichiro\",\"genres\":[\"action\",\"adventure\"],\"status\":\"ongoing\",\"total_chapters\":1100,\"description\":\"Pirate king\"}"

echo.
echo [2] Add manga duplicate
curl -X POST http://localhost:8080/manga ^
-H "Content-Type: application/json" ^
-H "Authorization: Bearer %TOKEN%" ^
-d "{\"title\":\"One Piece\",\"author\":\"Oda Eiichiro\",\"genres\":[\"action\"],\"status\":\"ongoing\",\"total_chapters\":1100,\"description\":\"duplicate\"}"

echo.
echo [3] Add manga without token (should fail)
curl -X POST http://localhost:8080/manga ^
-H "Content-Type: application/json" ^
-d "{\"title\":\"Naruto\",\"author\":\"Kishimoto\",\"genres\":[\"action\"],\"status\":\"completed\",\"total_chapters\":700}"

echo.
echo [4] Add second manga
curl -X POST http://localhost:8080/manga ^
-H "Content-Type: application/json" ^
-H "Authorization: Bearer %TOKEN%" ^
-d "{\"title\":\"Naruto\",\"author\":\"Kishimoto\",\"genres\":[\"action\",\"adventure\"],\"status\":\"completed\",\"total_chapters\":700}"

echo.
echo ========================================
echo SEARCH TEST
echo ========================================

echo.
echo [5] Search by name
curl "http://localhost:8080/manga?query=one"

echo.
echo [6] Search by genre (single)
curl "http://localhost:8080/manga?genre=action"

echo.
echo [7] Search by genre (multiple)
curl "http://localhost:8080/manga?genre=action&genre=adventure"

echo.
echo [8] Search by status
curl "http://localhost:8080/manga?status=ongoing"

echo.
echo [9] Search combined
curl "http://localhost:8080/manga?query=one&genre=action&status=ongoing"

echo.
echo [10] Search not found
curl "http://localhost:8080/manga?query=notexist"

echo.
echo ========================================
echo GET MANGA
echo ========================================

set /p MANGA_ID=Enter manga_id from search result: 

echo.
echo [11] Get manga detail
curl "http://localhost:8080/manga/%MANGA_ID%"

echo.
echo [12] Get manga not found
curl "http://localhost:8080/manga/not-exist"

echo.
echo ========================================
echo TEST LIBRARY
echo ========================================

echo.
echo [13] Add to library success
curl -X POST http://localhost:8080/users/library ^
-H "Content-Type: application/json" ^
-H "Authorization: Bearer %TOKEN%" ^
-d "{\"manga_id\":\"%MANGA_ID%\",\"status\":\"reading\",\"current_chapter\":1}"

echo.
echo [14] Add duplicate
curl -X POST http://localhost:8080/users/library ^
-H "Content-Type: application/json" ^
-H "Authorization: Bearer %TOKEN%" ^
-d "{\"manga_id\":\"%MANGA_ID%\",\"status\":\"reading\",\"current_chapter\":1}"

echo.
echo [15] Add manga not exist
curl -X POST http://localhost:8080/users/library ^
-H "Content-Type: application/json" ^
-H "Authorization: Bearer %TOKEN%" ^
-d "{\"manga_id\":\"not-exist\",\"status\":\"reading\",\"current_chapter\":0}"

echo.
echo [16] View library
curl http://localhost:8080/users/library ^
-H "Authorization: Bearer %TOKEN%"

echo.
echo [17] Update progress success
curl -X PUT http://localhost:8080/users/progress ^
-H "Content-Type: application/json" ^
-H "Authorization: Bearer %TOKEN%" ^
-d "{\"manga_id\":\"%MANGA_ID%\",\"current_chapter\":50,\"status\":\"reading\"}"

echo.
echo [18] Update invalid chapter
curl -X PUT http://localhost:8080/users/progress ^
-H "Content-Type: application/json" ^
-H "Authorization: Bearer %TOKEN%" ^
-d "{\"manga_id\":\"%MANGA_ID%\",\"current_chapter\":-1}"

echo.
echo ========================================
echo TEST RATING (NEW)
echo ========================================

echo.
echo [19] Rate manga success
curl -X POST http://localhost:8080/manga/rate ^
-H "Content-Type: application/json" ^
-H "Authorization: Bearer %TOKEN%" ^
-d "{\"manga_id\":\"%MANGA_ID%\",\"rating\":9}"

echo.
echo [20] Rate invalid (out of range)
curl -X POST http://localhost:8080/manga/rate ^
-H "Content-Type: application/json" ^
-H "Authorization: Bearer %TOKEN%" ^
-d "{\"manga_id\":\"%MANGA_ID%\",\"rating\":11}"

echo.
echo [21] Rate manga not exist
curl -X POST http://localhost:8080/manga/rate ^
-H "Content-Type: application/json" ^
-H "Authorization: Bearer %TOKEN%" ^
-d "{\"manga_id\":\"not-exist\",\"rating\":8}"

echo.
echo ========================================
echo DONE
echo ========================================
pause