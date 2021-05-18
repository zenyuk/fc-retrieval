import { Client } from './index';

test('My Greeter', () => {
  const client = new Client();
  expect(client.hello()).toBe('Hello World');
});
