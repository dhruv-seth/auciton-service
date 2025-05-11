package domain

type Bid struct {
    Bidder    string
    Amount    int
    ListingID int
}

type Auction struct {
    ListingID   int
    ItemName    string
    HighestBid  *Bid
    IsClosed    bool
    CloseChan   chan struct{}
    BroadcastCh chan Bid
}

type AuctionRepository interface {
    GetListingByID(id int) (*Auction, error)
    SaveBid(bid Bid) error
    CloseListing(id int) error
}
