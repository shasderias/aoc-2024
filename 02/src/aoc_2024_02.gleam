import gleam/int
import gleam/io
import gleam/list
import gleam/result
import gleam/string
import simplifile

pub fn main() {
  let assert Ok(records) = simplifile.read(from: "example.txt")
  let _ = io.debug(run(records, is_safe_1))
  let assert Ok(records) = simplifile.read(from: "input.txt")
  let _ = io.debug(run(records, is_safe_1))
  let assert Ok(records) = simplifile.read(from: "example.txt")
  let _ = io.debug(run(records, is_safe_2))
  let assert Ok(records) = simplifile.read(from: "input.txt")
  let _ = io.debug(run(records, is_safe_2))
}

fn run(input: String, is_safe_fn: fn(List(Int)) -> Bool) {
  input
  |> string.split(on: "\r\n")
  |> list.map(with: fn(x) {
    string.split(x, on: " ")
    |> list.map(with: fn(x) { result.unwrap(int.parse(x), 0) })
    |> is_safe_fn
  })
  |> list.count(fn(x) { x == True })
}

fn is_asc(x: Int) {
  x > -4 && x < 0
}

fn is_dsc(x: Int) {
  x > 0 && x < 4
}

pub fn is_safe_1(report: List(Int)) -> Bool {
  let diff =
    list.window_by_2(report)
    |> list.map(fn(pair) { pair.0 - pair.1 })

  let all_asc = list.all(diff, is_asc)
  let all_dsc = list.all(diff, is_dsc)

  all_asc || all_dsc
}

pub fn is_safe_2(report: List(Int)) -> Bool {
  let assert Ok(report_without_first) = list.rest(report)

  let is_safe_with = fn(report: List(Int), test_fn: fn(Int) -> Bool) {
    report
    |> list.fold_until(#(0, False, False), fn(acc, v) {
      case acc.0 == 0, test_fn(acc.0 - v), acc.1 {
        True, _, _ -> list.Continue(#(v, False, False))
        _, True, _ -> list.Continue(#(v, acc.1, False))
        _, False, False -> list.Continue(#(acc.0, True, False))
        _, False, True -> list.Stop(#(0, False, True))
      }
    })
    |> fn(acc) { !acc.2 }
  }

  let asc_first_case = report_without_first |> is_safe_1
  let asc = report |> is_safe_with(is_asc)
  let dsc_first_case = report_without_first |> is_safe_1
  let dsc = report |> is_safe_with(is_dsc)

  asc || asc_first_case || dsc || dsc_first_case
}
