def solve(n) = File.read("6/input.txt").strip.split("").each_cons(n).find_index { _1.uniq.size == _1.size } + n
puts "Part 1: #{solve(4)}"
puts "Part 2: #{solve(14)}"
