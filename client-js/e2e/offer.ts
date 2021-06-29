import { FilecoinRetrievalClient, Settings, ContentID, NodeID } from '../src/'

const getOffers = async() => {

  console.log("Start getOffers")

  // const privateKey = "d54193a9668ae59befa59498cdee16b78cdc8228d43814442a64588fd1648a29"
  // const url = "http://127.0.0.1:1234/rpc/v0"
  // const authToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.c1sdp8fhSbCsJfDKCC2xcqrngVzbIDsLeoggz0IHxWE"
    
  const settings: Settings = {
    // establishmentTTL: 1,
    // registerURL: "",
    // clientId: "",
    // logLevel: undefined,
    // logTarget: undefined,
    // logServiceName: undefined,
    // blockchainPrivateKey: undefined,
    // retrievalPrivateKey: undefined,
    // retrievalPrivateKeyVer: undefined,
    // walletPrivateKey: undefined,
    lotusAP: "http://127.0.0.1:1234/rpc/v0",
    // lotusAuthToken: undefined,
    // searchPrice = defaultSearchPrice: undefined,
    // offerPrice = defaultOfferPrice,: undefined,
    // topUpAmount = defaultTopUpAmount: undefined,
  } as Settings

  const mgr = new FilecoinRetrievalClient(settings)

  const cid: ContentID = {id: "" } as ContentID
  const gatewayID: NodeID = new NodeID("")

  const numDHT = 10
  const maxOffers = 10

  // mgr.AddActiveGateways([])
  mgr.findOffersDHTDiscoveryV2(cid, gatewayID, numDHT, maxOffers)
  mgr.findOffersStandardDiscoveryV2(cid, gatewayID, maxOffers)
  console.log("mgr:", mgr)
  
}

getOffers()


