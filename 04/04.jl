function countXmas(lines)
    acc = 0
    for line = lines
        acc += eachmatch(r"XMAS", line) |> collect |> length
        acc += eachmatch(r"SAMX", line) |> collect |> length
    end
    acc
    return acc
end

function toLines(mat)
    lines = []
    for line = eachrow(mat)
        push!(lines, String(line))
    end
    return lines
end

function rot45Lines(mat, stride)
    ht = (stride * 2) - 1
    lines = [Vector{UInt8}() for _ in 1:ht]
    row, cnt = 1, 0

    for n = 1:stride
        j = n
        for i = 1:n
            push!(lines[row], mat[i, j])
            cnt += 1
            if cnt == (row < stride ? row : stride - (row - stride))
                cnt = 0
                row += 1
            end
            j -= 1
        end
    end
    for n = 2:stride
        j = stride
        for i = n:stride
            push!(lines[row], mat[i, j])
            cnt += 1
            if cnt == (row < stride ? row : stride - (row - stride))
                cnt = 0
                row += 1
            end
            j -= 1
        end
    end

    strLines = []
    for line = lines
        push!(strLines, String(line))
    end
    return strLines
end

function star1(in)
    stride = split(String(deepcopy(in)), "\r\n") |> length
    mat = filter(c -> (c != 0x0d && c != 0x0a), in) |> x -> reshape(x, stride, stride) |> transpose
    a = mat |> toLines |> countXmas
    b = mat |> rotl90 |> toLines |> countXmas
    c = mat |> x -> rot45Lines(x, stride) |> countXmas
    d = mat |> rotl90 |> x -> rot45Lines(x, stride) |> countXmas

    println(a + b + c + d)
end

star1(read("example.txt"))
star1(read("input.txt"))

function star2(in)
    lines = split(String(in), "\r\n")
    stride = length(lines)

    acc = 0

    for x = 2:stride-1
        for y = 2:stride-1
            if lines[x][y] != 'A'
                continue
            end
            if (
                (lines[x-1][y-1] == 'M' && lines[x+1][y+1] == 'S') ||
                (lines[x-1][y-1] == 'S' && lines[x+1][y+1] == 'M')
               ) &&
               (
                (lines[x-1][y+1] == 'M' && lines[x+1][y-1] == 'S') ||
                (lines[x-1][y+1] == 'S' && lines[x+1][y-1] == 'M')
               )
                acc += 1
            end
        end
    end

    println(acc)
end

star2(read("example.txt"))
star2(read("input.txt"))