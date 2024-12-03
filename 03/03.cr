example01 = File.read("example_01.txt")
example02 = File.read("example_02.txt")
input = File.read("input.txt")

def star1(memory)
  memory.scan(/mul\((?<a>\d+),(?<b>\d+)\)/).map() { |m|
    m["a"].to_i * m["b"].to_i
  }.sum
end

def star2(memory)
  acc, enabled = 0, true
  memory.scan(/mul\((?<a>\d+),(?<b>\d+)\)|do\(\)|don't\(\)/).each() { |m|
    case m.to_s
    when .starts_with?("don't")
      enabled = false
    when .starts_with?("do")
      enabled = true
    else
      acc += m["a"].to_i * m["b"].to_i if enabled
    end
  }
  acc
end

puts star1(example01)
puts star1(input)
puts star2(example02)
puts star2(input)
