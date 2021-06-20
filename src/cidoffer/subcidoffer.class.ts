import { ContentID } from '../cid/cid.interface'
import { MerkleProof } from '../fcrMerkleTree/proof'
import { NodeID } from '../nodeid/nodeid.interface'

export class SubCIDOffer {
  providerID: NodeID
  subCID: ContentID
  merkleRoot: string
  merkleProof: MerkleProof
  price: number
  expiry: number
  qos: number
  signature: string

  constructor(
    providerID: NodeID,
    subCID: ContentID,
    merkleRoot: string,
    merkleProof: MerkleProof,
    price: number,
    expiry: number,
    qos: number,
    signature: string,
  ) {
    this.providerID = providerID
    this.subCID = subCID
    this.merkleRoot = merkleRoot
    this.merkleProof = merkleProof
    this.price = price
    this.expiry = expiry
    this.qos = qos
    this.signature = signature
  }

  // GetProviderID returns the provider ID of this offer.
  getProviderID(): string {
    return ''
  }

  // GetSubCID returns the sub cid of this offer.
  getSubCID(): ContentID {
    return {} as ContentID
  }

  // GetMerkleRoot returns the merkle root of this offer.
  getMerkleRoot(): string {
    return 'hello'
  }

  // GetMerkleProof returns the merkle proof of this offer.
  getMerkleProof(): MerkleProof {
    return {} as MerkleProof
  }

  // GetPrice returns the price of this offer.
  getPrice(): number {
    return 42
  }

  // GetExpiry returns the expiry of this offer.
  getExpiry(): number {
    return 42
  }

  // GetQoS returns the quality of service of this offer.
  getQoS(): number {
    return 42
  }

  // GetSignature returns the signature of this offer.
  getSignature(): string {
    return ''
  }

  // Verify is used to verify the offer with a given public key.
  verify(pubKey: string) {}

  // VerifyMerkleProof is used to verify the sub cid is part of the merkle trie
  verifyMerkleProof() {}
}
