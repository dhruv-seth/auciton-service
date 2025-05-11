package main

import (
    "auction-system/internal/app"
    "auction-system/internal/infra/db"
    "auction-system/internal/infra/ws"
    "database/sql"
    "log"
    "net/http"
    "time"

    _ "github.com/lib/pq"
)

func main() {
    dbConn, err := sql.Open("postgres", "user=postgres dbname=auction sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }

    repo := db.NewSQLAuctionRepo(dbConn)
    svc := app.NewAuctionService(repo)
    wsHandler := ws.NewWebSocketHandler(svc)

    listing, _ := repo.GetListingByID(1)
    svc.StartAuction(nil, listing, 60*time.Second)

    http.HandleFunc("/ws", wsHandler.ServeWS)

    log.Println("Listening on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
