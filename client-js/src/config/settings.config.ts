import { BigNumber } from "bignumber.js";
import { NodeID } from '../nodeid/nodeid.interface'

// DefaultEstablishmentTTL is the default Time To Live used with Client - Gateway estalishment messages.
const defaultEstablishmentTTL = 100

// DefaultLogLevel is the default amount of logging to show.
const defaultLogLevel = "trace"

// DefaultLogTarget is the default output location of log output.
const defaultLogTarget = "STDOUT"

// DefaultLogServiceName is the default service name of logging.
const defaultLogServiceName = "client"

// defaultSearchPrice is the default search price.
const defaultSearchPrice = new BigNumber("1000000000000000")
// const defaultSearchPrice = 1_000_000_000_000_000

// defaultOfferPrice is the default offer price.
const defaultOfferPrice = new BigNumber("1000000000000000")
// const defaultOfferPrice = 1_000_000_000_000_000

// defaultTopUpAmount is the default top up amount.
const defaultTopUpAmount = new BigNumber("100000000000000000")
// const defaultTopUpAmount = 100_000_000_000_000_000

export class Settings {
  establishmentTTL: number
  registerURL: string
  clientId: NodeID
  logLevel: string
  logTarget: string
  logServiceName: string
  blockchainPrivateKey: any
  retrievalPrivateKey: any
  retrievalPrivateKeyVer: any
  walletPrivateKey: string
  lotusAP: string
  lotusAuthToken: string
  searchPrice: BigNumber
  offerPrice: BigNumber
  topUpAmount: BigNumber

  constructor({
    establishmentTTL = defaultEstablishmentTTL,
    registerURL,
    clientId,
    logLevel = defaultLogLevel,
    logTarget = defaultLogTarget,
    logServiceName = defaultLogServiceName,
    blockchainPrivateKey,
    retrievalPrivateKey,
    retrievalPrivateKeyVer,
    walletPrivateKey,
    lotusAP,
    lotusAuthToken,
    searchPrice = defaultSearchPrice,
    offerPrice = defaultOfferPrice,
    topUpAmount = defaultTopUpAmount,
  }: any) {
    this.establishmentTTL = establishmentTTL
    this.registerURL = registerURL
    this.clientId = clientId
    this.logLevel = logLevel
    this.logTarget = logTarget
    this.logServiceName = logServiceName
    this.blockchainPrivateKey = blockchainPrivateKey
    this.retrievalPrivateKey = retrievalPrivateKey
    this.retrievalPrivateKeyVer = retrievalPrivateKeyVer
    this.walletPrivateKey = walletPrivateKey
    this.lotusAP = lotusAP
    this.lotusAuthToken = lotusAuthToken
    this.searchPrice = searchPrice
    this.offerPrice = offerPrice
    this.topUpAmount = topUpAmount
  }
}

export const buildSettings = (registerURL: string) => {
  // TODO: Generate keys
  return new Settings({
    blockchainPrivateKey: undefined,
    client: {} as NodeID,
    lotusAP: "",
    lotusAuthToken: "",
    registerURL,
    retrievalPrivateKey: undefined,
    retrievalPrivateKeyVer: undefined,
    walletPrivateKey: "",
  })
}