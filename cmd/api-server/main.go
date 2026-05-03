package apiserver

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"managahub/internal/auth"
	grpcserver "managahub/internal/grpc"
	"managahub/internal/library"
	"managahub/internal/manga"
	"managahub/internal/tcp"
	"managahub/internal/udp"
	wsocket "managahub/internal/websocket"
)

func Run(tcpServer *tcp.ProgressSyncServer, udpServer *udp.NotificationServer, ws *wsocket.ChatHub, db *sql.DB) {
	jwtSecret := "ITITIU22134_LENGUYENCHIQUOC"
	grpcClient := grpcserver.NewMangaGRPCClient("localhost:9092")
	authService := auth.NewAuthService(db, jwtSecret)
	mangaService := manga.NewMangaService(db)
	libraryService := library.NewLibraryService(db, tcpServer)

	authHandler := auth.NewAuthHandler(authService)
	mangaHandler := manga.NewMangaHandler(mangaService)
	libraryHandler := library.NewLibraryHandler(libraryService, grpcClient)
	notifyHandler := udp.NewNotificationHandler(udpServer)
	r := gin.Default()

	public := r.Group("/auth")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
	}

	r.GET("/manga", mangaHandler.SearchManga)
	r.GET("/manga/:id", mangaHandler.GetManga)
	r.GET("/ws", func(c *gin.Context) {
		ws.HandleWS(c.Writer, c.Request)
	})
	protected := r.Group("/")
	protected.Use(auth.JWTMiddleware(authService))
	{

		protected.POST("/auth/logout", authHandler.Logout)
		protected.GET("/me", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			username, _ := c.Get("username")
			c.JSON(http.StatusOK, gin.H{
				"user_id":  userID,
				"username": username,
			})
		})

		protected.POST("/manga", mangaHandler.CreateManga)

		protected.POST("/users/library", libraryHandler.AddToLibrary)
		protected.GET("/users/library", libraryHandler.GetLibrary)
		protected.PUT("/users/progress", libraryHandler.UpdateProgress)
		protected.POST("/notify", notifyHandler.SendNotification)
	}
	log.Println("HTTP API Server running at http://localhost:8080")
	log.Println("WebSocket Chat running at ws://localhost:8080/ws")
	r.Run(":8080")
}
