package ws

import (
    "auction-system/internal/app"
    "auction-system/internal/domain"
    "encoding/json"
    "github.com/gorilla/websocket"
    "net/http"
    "strconv"
)

type WebSocketHandler struct {
    auctionService *app.AuctionService
    clients        map[int]map[*websocket.Conn]bool
}

func NewWebSocketHandler(svc *app.AuctionService) *WebSocketHandler {
    return &WebSocketHandler{
        auctionService: svc,
        clients:        make(map[int]map[*websocket.Conn]bool),
    }
}

func (h *WebSocketHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
    up := websocket.Upgrader{
        CheckOrigin: func(r *http.Request) bool { return true },
    }
    conn, _ := up.Upgrade(w, r, nil)

    listingIDStr := r.URL.Query().Get("listing")
    listingID, _ := strconv.Atoi(listingIDStr)

    if h.clients[listingID] == nil {
        h.clients[listingID] = make(map[*websocket.Conn]bool)
    }
    h.clients[listingID][conn] = true

    go h.listen(conn, listingID)
}

func (h *WebSocketHandler) listen(conn *websocket.Conn, listingID int) {
    defer conn.Close()
    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            break
        }
        var bid domain.Bid
        if err := json.Unmarshal(msg, &bid); err == nil {
            h.auctionService.PlaceBid(bid)
        }
    }
}

func (h *WebSocketHandler) BroadcastBid(bid domain.Bid) {
    for conn := range h.clients[bid.ListingID] {
        _ = conn.WriteJSON(bid)
    }
}
