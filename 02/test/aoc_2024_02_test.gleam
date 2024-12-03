import aoc_2024_02 as aoc
import gleeunit
import gleeunit/should

pub fn main() {
  gleeunit.main()
}

// gleeunit test functions end in `_test`
pub fn hello_world_test() {
  1
  |> should.equal(1)
}

pub fn is_safe_1_test() {
  aoc.is_safe_1([7, 6, 4, 2, 1])
  |> should.equal(True)

  aoc.is_safe_1([1, 2, 7, 8, 9])
  |> should.equal(False)
}

pub fn is_safe_2_test() {
  aoc.is_safe_2([7, 6, 4, 2, 1])
  |> should.equal(True)

  aoc.is_safe_2([1, 2, 7, 8, 9])
  |> should.equal(False)

  aoc.is_safe_2([9, 7, 6, 2, 1])
  |> should.equal(False)

  aoc.is_safe_2([1, 3, 2, 4, 5])
  |> should.equal(True)

  aoc.is_safe_2([8, 6, 4, 4, 1])
  |> should.equal(True)

  aoc.is_safe_2([1, 3, 6, 7, 9])
  |> should.equal(True)
}
