import { NodeID } from '../nodeid/nodeid.interface'
import { ContentID } from '../cid/cid.interface'
import { FCRMerkleProof } from '../fcrMerkleTree/proof.class'

export interface SubCIDOffer {
  providerID: NodeID
  subCID: ContentID
  merkleRoot: string
  merkleProof: FCRMerkleProof
  price: number
  expiry: number
  qos: number
  signature: string
}
