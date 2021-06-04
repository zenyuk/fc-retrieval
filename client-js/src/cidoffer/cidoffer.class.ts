import { ContentID } from '../cid/cid.interface'
import { NodeID } from '../nodeid/nodeid.interface'

export class CIDOffer {
  // GetProviderID returns the provider ID of this offer.
  getProviderID(): NodeID {
    return {} as NodeID
  }

  // GetCIDs returns the cids of this offer.
  getCIDs(): ContentID {
    return {} as ContentID
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
    return 'hello'
  }
}
