export interface DefaultSettings {
  defaultEstablishmentTTL: number
  defaultLogLevel: string
  defaultLogTarget: string
  defaultLogServiceName: string
  defaultRegisterURL: string
}

export const defaults: DefaultSettings = {
  // DefaultEstablishmentTTL is the default Time To Live used with Client - Gateway estalishment messages.
  defaultEstablishmentTTL: 100,

  // DefaultLogLevel is the default amount of logging to show.
  defaultLogLevel: 'trace',

  // DefaultLogTarget is the default output location of log output.
  defaultLogTarget: 'STDOUT',

  // DefaultLogServiceName is the default service name of logging.
  defaultLogServiceName: 'client',

  // DefaultRegisterURL is the default location of the Register service.
  // register:9020 is the value that will work for the integration test system.
  defaultRegisterURL: 'http://localhost:9020',
}
