import { Client } from './index'
import { defaults } from './defaults'

describe('Client', () => {
  it('Hello', () => {
    const client = new Client(defaults)
    expect(client.findGateways()).toBe('Hello World')
  })
})
