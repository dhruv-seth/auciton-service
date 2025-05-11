package app

import (
    "context"
    "time"
    "auction-system/internal/domain"
)

type AuctionService struct {
    repo domain.AuctionRepository
    live map[int]*domain.Auction
}

func NewAuctionService(repo domain.AuctionRepository) *AuctionService {
    return &AuctionService{
        repo: repo,
        live: make(map[int]*domain.Auction),
    }
}

func (s *AuctionService) StartAuction(ctx context.Context, listing *domain.Auction, duration time.Duration) {
    s.live[listing.ListingID] = listing
    listing.BroadcastCh = make(chan domain.Bid, 100)
    listing.CloseChan = make(chan struct{})

    go func() {
        select {
        case <-ctx.Done():
        case <-time.After(duration):
            s.CloseAuction(listing.ListingID)
        }
    }()

    go s.listenForBids(listing)
}

func (s *AuctionService) listenForBids(a *domain.Auction) {
    for bid := range a.BroadcastCh {
        if a.IsClosed {
            continue
        }
        if a.HighestBid == nil || bid.Amount > a.HighestBid.Amount {
            a.HighestBid = &bid
            s.repo.SaveBid(bid)
        }
    }
}

func (s *AuctionService) PlaceBid(bid domain.Bid) error {
    auction, ok := s.live[bid.ListingID]
    if !ok || auction.IsClosed {
        return nil
    }
    auction.BroadcastCh <- bid
    return nil
}

func (s *AuctionService) CloseAuction(listingID int) {
    if a, ok := s.live[listingID]; ok {
        a.IsClosed = true
        close(a.BroadcastCh)
        s.repo.CloseListing(listingID)
        delete(s.live, listingID)
    }
}
