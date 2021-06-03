import { Content } from './types'

export class FCRMerkleProof {
  path: string[] = []
  index: number[] = []

  verifyContent(content: Content, root: string): boolean {
    // TODO
    return true
  }

  marshalJSON(): string {
    // TODO
    return ''
  }

  unmarshalJSON(p: string) {
    // TODO
  }
}
