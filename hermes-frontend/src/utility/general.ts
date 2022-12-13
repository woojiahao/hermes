export function all<T>(arr: T[], predicate: (el: T) => boolean) {
  return arr.filter(predicate).length === arr.length
}
