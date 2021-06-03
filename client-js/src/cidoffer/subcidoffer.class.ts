import { ContentID } from '../cid/cid.interface'
import { KeyPair } from '../fcrcrypto/types'
import { FCRMerkleProof } from '../fcrMerkleTree/proof.class'
import { NodeID } from '../nodeid/nodeid.interface'

export class SubCIDOffer {
  // GetProviderID returns the provider ID of this offer.
  getProviderID(): NodeID {
    return {} as NodeID
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
  getMerkleProof(): FCRMerkleProof {
    return {} as FCRMerkleProof
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
  verify(pubKey: KeyPair) {}

  // VerifyMerkleProof is used to verify the sub cid is part of the merkle trie
  verifyMerkleProof() {}
}
