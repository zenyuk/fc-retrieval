import { Settings } from './settings.interface'
import { defaults } from '../constants/defaults'

export class CreateSettings {
  settings: Settings

  constructor() {
    this.settings = defaults as Settings
  }

  setDefaultEstablishmentTTL(ttl: number) {
    this.settings.establishmentTTL = ttl
  }

  setDefaultLogLevel(logLevel: string) {
    // this.settings.defaultLogLevel = logLevel
    this.settings.logLevel = logLevel
  }

  setDefaultLogTarget(logTarget: string) {
    // this.settings.defaultLogTarget = logTarget
    this.settings.logTarget = logTarget
  }

  setDefaultLogServiceName(serviceName: string) {
    // this.settings.defaultLogServiceName = serviceName
    this.settings.logServiceName = serviceName
  }

  setDefaultRegisterURL(url: string) {
    // this.settings.defaultRegisterURL = url
    this.settings.registerURL = url
  }

  build(): Settings {
    return this.settings
  }
}
