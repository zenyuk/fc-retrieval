export enum FCRMessageType {
  // Message originating from client
  ClientEstablishmentRequestType          = 100,
  ClientEstablishmentResponseType         = 101,
  ClientStandardDiscoverRequestType       = 102,
  ClientStandardDiscoverResponseType      = 103,
  ClientDHTDiscoverRequestType            = 104,
  ClientDHTDiscoverResponseType           = 105,
  ClientDHTOfferAckRequestType            = 106,
  ClientDHTOfferAckResponseType           = 107,
  ClientStandardDiscoverRequestV2Type     = 108,
  ClientStandardDiscoverResponseV2Type    = 109,
  ClientStandardDiscoverOfferRequestType  = 110,
  ClientStandardDiscoverOfferResponseType = 111,
  ClientDHTDiscoverRequestV2Type          = 112,
  ClientDHTDiscoverResponseV2Type         = 113,
  ClientDHTDiscoverOfferRequestType       = 114,
  ClientDHTDiscoverOfferResponseType      = 115,

  // Message originating from gateway
  GatewayListDHTOfferRequestType                          = 200,
  GatewayListDHTOfferResponseType                         = 201,
  GatewayListDHTOfferAckType                              = 202,
  GatewayDHTDiscoverRequestType                           = 203,
  GatewayDHTDiscoverResponseType                          = 204,
  GatewayPingRequestType                                  = 205,
  GatewayPingResponseType                                 = 206,
  GatewayNotifyProviderGroupCIDOfferSupportedRequestType  = 205,
  GatewayNotifyProviderGroupCIDOfferSupportedResponseType = 206,
  GatewayDHTDiscoverRequestV2Type                         = 207,
  GatewayDHTDiscoverResponseV2Type                        = 208,
  GatewayDHTDiscoverOfferRequestType                      = 209,
  GatewayDHTDiscoverOfferResponseType                     = 210,

  // Message originating from provider
  ProviderPublishGroupOfferRequestType  = 300,
  ProviderPublishGroupOfferResponseType = 301,
  ProviderPublishDHTOfferRequestType    = 302,
  ProviderPublishDHTOfferResponseType   = 303,

  // Message originating from gateway admin
  GatewayAdminInitialiseKeyRequestType                      = 400,
  GatewayAdminInitialiseKeyResponseType                     = 401,
  GatewayAdminGetReputationRequestType                      = 402,
  GatewayAdminGetReputationResponseType                     = 403,
  GatewayAdminSetReputationRequestType                      = 404,
  GatewayAdminSetReputationResponseType                     = 405,
  GatewayAdminForceRefreshRequestType                       = 406,
  GatewayAdminForceRefreshResponseType                      = 407,
  GatewayAdminUpdateGatewayGroupCIDOfferSupportRequestType  = 408,
  GatewayAdminUpdateGatewayGroupCIDOfferSupportResponseType = 409,
  GatewayAdminListDHTOfferRequestType                       = 410,
  GatewayAdminListDHTOfferResponseType                      = 411,
  GatewayAdminInitialiseKeyRequestV2Type                    = 412,

  // Message originating from provider admin
  ProviderAdminInitialiseKeyRequestType      = 500,
  ProviderAdminInitialiseKeyResponseType     = 501,
  ProviderAdminPublishGroupOfferRequestType  = 502,
  ProviderAdminPublishGroupOfferResponseType = 503,
  ProviderAdminPublishDHTOfferRequestType    = 504,
  ProviderAdminPublishDHTOfferResponseType   = 505,
  ProviderAdminGetPublishedOfferRequestType  = 506,
  ProviderAdminGetPublishedOfferResponseType = 507,
  ProviderAdminForceRefreshRequestType       = 508,
  ProviderAdminForceRefreshResponseType      = 509,
  ProviderAdminInitialiseKeyRequestV2Type    = 510,

  // Messages for basic protocol
  ProtocolChangeRequestType     = 900,
  ProtocolChangeResponseType    = 901,
  InvalidMessageResponseType    = 902,
  InsufficientFundsResponseType = 903,
}
