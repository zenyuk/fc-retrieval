
import { NodeID } from './nodeid/nodeid.interface'
import { ContentID } from './cid/cid.interface'


export interface SubCIDOffer {
	providerID:  NodeID
	subCID:      ContentID
	merkleRoot:  string
	merkleProof: *fcrmerkletree.FCRMerkleProof
	price:       number
	expiry:      number
	qos:         number
	signature:   string
}
