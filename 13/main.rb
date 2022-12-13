def check_order(left, right)
  result = if left.is_a?(Integer) && right.is_a?(Integer)
    if left < right
      :correct
    elsif left == right
      :continue
    else
      :incorrect
    end
  elsif left.is_a?(Array) && right.is_a?(Array)
    correct = if left.size < right.size
      :correct
    elsif left.size > right.size
      :incorrect
    else
      :continue
    end

    left.each_with_index do |l, i|
      if i >= right.size
        correct = :incorrect
        break
      end

      sub_correct = check_order(l, right[i])
      if sub_correct != :continue
        correct = sub_correct
        break
      end
    end

    correct
  elsif left.is_a?(Integer) && right.is_a?(Array)
    check_order([left], right)
  elsif left.is_a?(Array) && right.is_a?(Integer)
    check_order(left, [right])
  end
  result
end

input = File.readlines("13/input.txt").map { eval(_1.chomp) }.compact.each_slice(2).to_a

result = 0

input.each_with_index do |(left, right), i|
  result += i + 1 if check_order(left, right) == :correct
end

puts("Part 1: #{result}")

input = input.flatten(1).concat([[[2]], [[6]]])
input.sort! { |a, b| check_order(a, b) == :correct ? -1 : 1 }
i = input.find_index { |i| i == [[2]]} + 1
j = input.find_index { |i| i == [[6]]} + 1
pp i, j
puts("Part 2: #{i * j}")
