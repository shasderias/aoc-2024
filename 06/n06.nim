import std/strutils
import std/sequtils
import std/sets
import std/sugar

let example = readFile("example.txt")
let input = readFile("input.txt")

type
  Entity = enum
    Nil
    Block
    Guard
    Visited
    Cross

  Point = tuple[x: int, y: int]

proc `+`(a: Point, b: Point): Point =
  return (a.x + b.x, a.y + b.y)

proc `-`(a: Point, b: Point): Point =
  return (a.x - b.x, a.y - b.y)

proc rotr90(pt: Point): Point =
  return (-pt.y, pt.x)

const
  N = (0, -1)
  S = (0, 1)
  W = (-1, 0)
  E = (1, 0)

assert N.rotr90 == E
assert E.rotr90 == S
assert S.rotr90 == W
assert W.rotr90 == N

type Grid = ref object
  data: seq[Entity]
  stride: int
  height: int

func parseEntity(c: char): Entity =
  case c
  of '.':
    return Nil
  of '#':
    return Block
  of '^':
    return Guard
  else:
    raise newException(ValueError, "unsupported entity")

proc newGrid(input: string): Grid =
  let lines = splitLines(input)
  let stride = lines[0].len
  let data = cast[seq[char]](lines.join("")).map(parseEntity)
  return Grid(data: data, stride: stride, height: data.len div stride)

proc idxToPoint(g: Grid, idx: int): Point =
  return (idx - ((idx div g.stride) * g.stride), idx div g.stride)

proc posOf(g: Grid, ent: Entity): Point =
  let idx = find(g.data, ent)
  return g.idxToPoint(idx)

proc at(g: Grid, pt: Point): Entity =
  return g.data[pt.y * g.stride + pt.x]

proc inBounds(g: Grid, pt: Point): bool =
  return pt.x >= 0 and pt.y >= 0 and pt.x < g.stride and pt.y < g.height

proc set(g: var Grid, pt: Point, ent: Entity) =
  g.data[pt.y * g.stride + pt.x] = ent

proc count(g: Grid, ent: Entity): int =
  return count(g.data, ent)

proc loopTest(g: Grid, blockPos: Point, pos: Point, dir: Point): bool =
  var visited: HashSet[tuple[pos: Point, dir: Point]]

  var pos = pos
  var dir = dir

  while true:
    if visited.contains((pos, dir)):
      return true
    visited.incl((pos, dir))

    if not g.inBounds(pos + dir):
      return false

    while pos + dir == blockPos or g.at(pos + dir) == Block:
      dir = dir.rotr90()

    pos = pos + dir

  raise newException(ValueError, "unreachable code")

proc star1(input: string): Grid =
  var g = newGrid(input)
  var dir = N
  var pos = g.posOf(Guard)

  while true:
    g.set(pos, Visited)

    if not g.inBounds(pos + dir):
      break

    while g.at(pos + dir) == Block:
      dir = dir.rotr90()

    pos = pos + dir

  return g

proc star2(input: string): int =
  var refG = newGrid(input)
  var guardPos = refG.posOf(Guard)

  let g = star1(input)
  var path: seq[Point]

  for idx, ent in g.data.pairs:
    let pt = g.idxToPoint(idx)
    if ent == Visited and pt != guardPos:
      path.add(pt)

  var acc = 0
  for blockPos in path.items:
    if blockPos == guardPos:
      echo "-1"
      continue
    if refG.loopTest(blockPos, guardPos, N):
      acc += 1

  return acc

echo star1(example).count(Visited)
echo star1(input).count(Visited)
echo star2(example)
echo star2(input)
