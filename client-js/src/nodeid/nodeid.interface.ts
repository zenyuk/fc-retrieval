export class NodeID {
  id: string

  constructor(id: string) {
    this.id = id
  }

  public toString = (): string => this.id
}
