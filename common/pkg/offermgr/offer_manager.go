// FCROfferMgr manages offer storage
package offermgr

import (
	"encoding/hex"

	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-common/pkg/database"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
)

var inUnitTest bool = false

var offerDBInitialized bool = false

type FCROfferMgr struct {
	db *database.Database
}

// NewFCROfferMgr returns
func NewFCROfferMgr() *FCROfferMgr {
	db, _ := database.NewDatabase()

	if !offerDBInitialized {
		db.Exec(`PRAGMA journal_mode=wal`)
		if inUnitTest {
			db.Exec(`drop table if exists offer`)
			db.Exec(`drop table if exists content`)
		} else {
			offerDBInitialized = true
		}
		db.Exec(`create table if not exists offer (digest blob, provider_id blob, expiry blob, price blob, qos blob, signature blob, primary key (digest))`)
		db.Exec(`create table if not exists content (content_id blob, content_no int, digest blob, primary key (digest, content_no))`)
		db.Exec(`create index if not exists content_cid_idx on content (content_id)`)
	}
	return &FCROfferMgr{
		db: db,
	}
}

// AddGroupOffer stores a group offer
func (mgr *FCROfferMgr) AddGroupOffer(offer *cidoffer.CIDOffer) error {
	return mgr.insertOffer(offer)
}

// AddDHTOffer stores a dht offer
func (mgr *FCROfferMgr) AddDHTOffer(offer *cidoffer.CIDOffer) error {
	return mgr.insertOffer(offer)
}

// GetGroupOffers returns a list of group offers that contain the given cid
func (mgr *FCROfferMgr) GetGroupOffers(c *cid.ContentID) ([]cidoffer.CIDOffer, bool) {
	return mgr.selectOffersSingle(c)
}

// GetDHTOffers returns a list of dht offers that contain the given cid
func (mgr *FCROfferMgr) GetDHTOffers(c *cid.ContentID) ([]cidoffer.CIDOffer, bool) {
	return mgr.selectOffersSingle(c)
}

// GetDHTOffersWithinRange returns a list of dht offers contains a cid within the given range
func (mgr *FCROfferMgr) GetDHTOffersWithinRange(cidMin, cidMax *cid.ContentID, maxOffers int) ([]cidoffer.CIDOffer, bool) {
	return mgr.selectOffersRange(maxOffers, cidMin, cidMax)
}

// GetOffers returns a list of all offers (group or dht) that contain the given cid
func (mgr *FCROfferMgr) GetOffers(c *cid.ContentID) ([]cidoffer.CIDOffer, bool) {
	return mgr.selectOffersSingle(c)
}

// GetOfferByDigest allows a gateway to be able to respond to a query to search for an offer by the offer digest
func (mgr *FCROfferMgr) GetOfferByDigest(digest [cidoffer.CIDOfferDigestSize]byte) (*cidoffer.CIDOffer, bool) {
	return mgr.selectOffersDigest(digest[:])
}

// insertOffer inserts records to offer and content table, if primary keys exist do nothing
func (mgr *FCROfferMgr) insertOffer(offer *cidoffer.CIDOffer) (err error) {
	// insert new if not exists, do nothing if exists
	sqlInsertOffer :=
		`insert into offer (digest, provider_id, expiry, price, qos, signature) 
		select ? digest, ? provider_id, ? expiry, ? price, ? qos, ? signature
		from (values (1)) a 
		where not exists (select digest from offer where digest = ?)`
	sqlInsertContent :=
		`insert into content (content_id, content_no, digest) 
		select ? content_id, ? content_no, ? digest
		from (values (1)) a 
		where not exists (select 1 from content where content_no = ? and digest = ?)`
	digestArr := offer.GetMessageDigest()
	digestHex := hex.EncodeToString(digestArr[:])
	_, err = mgr.db.Exec(sqlInsertOffer, digestHex,
		offer.GetProviderID().ToString(),
		offer.GetExpiry(),
		offer.GetPrice(),
		offer.GetQoS(),
		offer.GetSignature(), digestHex)
	if err != nil {
		return err
	}
	i := 0
	for _, contentID := range offer.GetCIDs() {
		_, err = mgr.db.Exec(sqlInsertContent, contentID.ToString(), i, digestHex, i, digestHex)
		i++
	}
	return err
}

// selectOffersSingle retrieves offers by a CID
// TODO: expiry
func (mgr *FCROfferMgr) selectOffersSingle(c *cid.ContentID) (res []cidoffer.CIDOffer, find bool) {
	sqlSelectOffer :=
		`select o.digest, o.provider_id, o.expiry, o.price, o.qos, o.signature
		from offer o join content c on o.digest=c.digest where c.content_id=?`

	return mgr.selectOffers(sqlSelectOffer, -1, c.ToString())
}

// selectOffersRange retrieves offers by a range of CIDs
// TODO: expiry
func (mgr *FCROfferMgr) selectOffersRange(maxOffers int, cidMin *cid.ContentID, cidMax *cid.ContentID) (res []cidoffer.CIDOffer, find bool) {
	sqlSelectOffer :=
		`select o.digest, o.provider_id, o.expiry, o.price, o.qos, o.signature
		from offer o, (select distinct digest from content 
				where content_id>=? and content_id<=?) c
		where o.digest=c.digest`

	return mgr.selectOffers(sqlSelectOffer, maxOffers, cidMin.ToString(), cidMax.ToString())
}

// selectOffersDigest retrieves a offer by a digest
// TODO: expiry
func (mgr *FCROfferMgr) selectOffersDigest(digest []byte) (*cidoffer.CIDOffer, bool) {
	sqlSelectOffer :=
		`select o.digest, o.provider_id, o.expiry, o.price, o.qos, o.signature
		from offer o where o.digest=?`

	offers, exist := mgr.selectOffers(sqlSelectOffer, 1, hex.EncodeToString(digest))
	if !exist {
		return nil, exist
	}
	return &offers[0], exist
}

// selectOffers retrieves offers by a select statement and its parameters
func (mgr *FCROfferMgr) selectOffers(sqlSelectOffer string, maxOffers int, parameters ...interface{}) (res []cidoffer.CIDOffer, find bool) {
	offerRows, err := mgr.db.Query(sqlSelectOffer, parameters...)
	if err != nil {
		return
	}
	var cnt int
	sqlSelectContent := `select content_id from content where digest=? order by content_no`
	for offerRows.Next() {
		var digestHex, providerHex string
		var expiry int64
		var price, qos uint64
		var signature string
		offerRows.Scan(&digestHex, &providerHex, &expiry, &price, &qos, &signature)
		providerID, _ := nodeid.NewNodeIDFromHexString(providerHex)

		cidRows, _ := mgr.db.Query(sqlSelectContent, digestHex)
		cids := []cid.ContentID{}
		for cidRows.Next() {
			var aCidStr string
			cidRows.Scan(&aCidStr)
			aCid, _ := cid.NewContentIDFromHexString(aCidStr)
			cids = append(cids, *aCid)
			find = true
		}

		offer, _ := cidoffer.NewCIDOffer(providerID, cids, price, expiry, qos)
		offer.SetSignature(signature)
		res = append(res, *offer)
		cnt++
		if cnt >= maxOffers && maxOffers > 0 {
			break
		}
	}
	return
}
