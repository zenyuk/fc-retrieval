import { Client } from './index'

describe('Client', () => {
  it('Hello', () => {
    const client = new Client()
    expect(client.hello()).toBe('Hello World')
  })
})
