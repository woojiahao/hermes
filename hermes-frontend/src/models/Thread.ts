// TODO: Display the OP's username
export default interface Thread {
  id: string
  is_published: boolean
  is_open: boolean
  title: string
  content: string
  tags: string[]
}
