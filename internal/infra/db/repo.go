package db

import (
	"auction-system/internal/domain"
	"database/sql"
)

type SQLAuctionRepo struct {
	db *sql.DB
}

func NewSQLAuctionRepo(db *sql.DB) *SQLAuctionRepo {
	return &SQLAuctionRepo{db: db}
}

func (r *SQLAuctionRepo) GetListingByID(id int) (*domain.Auction, error) {
	row := r.db.QueryRow("SELECT id, item_name FROM listings WHERE id=$1 AND is_closed=FALSE", id)

	var listing domain.Auction
	if err := row.Scan(&listing.ListingID, &listing.ItemName); err != nil {
		return nil, err
	}
	return &listing, nil
}

func (r *SQLAuctionRepo) SaveBid(bid domain.Bid) error {
	_, err := r.db.Exec("INSERT INTO bids (bidder, amount, listing_id) VALUES ($1, $2, $3)",
		bid.Bidder, bid.Amount, bid.ListingID)
	return err
}

func (r *SQLAuctionRepo) CloseListing(id int) error {
	_, err := r.db.Exec("UPDATE listings SET is_closed=TRUE WHERE id=$1", id)
	return err
}
