package reputation
// Copyright (C) 2020 ConsenSys Software Inc


const clientMaxReputation = 10000
const clientMinReputaiton = -10000
const clientInitialReputation = 10
const clientEstablishmentChallenge = -1
const clientOnChainDeposit = 1000

//Response with one or more CID Offers. Initial payment and final payment made.
const clientStdDiscOneCidOffer = 10

// Response with no CID Offers. Initial payment payment made.
const clientStdDiscNoCidOffers = 1

// Response with one or more CID Offers. Response message sent after one second prior to TTL expiry. Initial payment payment made.
const clientStdDiscLateCidOffers = 1

// Response with one or more CID Offers. Response message sent prior to one second prior to TTL expiry. Initial payment payment made but final payment not paid.
const clientStdDiscNonPayment = -100

// Response with one or more CID Offers from one or more Gateways. Initial payment and final payment made
const clientDhtDiscOneCidOffer = 10

// Response with no CID Offers. Initial payment made.
const clientDhtDiscNoCidOffers = 1

// Response with one or more CID Offers. Response message sent after one second prior to TTL expiry. Initial payment payment made.
const clientDhtDiscLateCidOffers = 1

// Response with one or more CID Offers. Response message sent prior to one second prior to TTL expiry. Initial payment payment made but final payments not paid.
const clientDhtDiscNonPayment = -300

// Micro-payment paid for content via Gateway. Note that there will be many micro-payments during content retrieval.
const clientMicroPayment = 1

// Invalid message received
const clientInvalidMessage = -10


