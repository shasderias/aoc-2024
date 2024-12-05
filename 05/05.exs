defmodule Common do
  def open(path) do
    {:ok, body} = File.read(path)

    rules =
      body
      |> String.split("\n\n")
      |> Enum.at(0)
      |> String.split("\n")
      |> Enum.map(fn x -> x |> String.split("|") |> List.to_tuple() end)

    updates =
      body
      |> String.split("\n\n")
      |> Enum.at(1)
      |> String.split("\n")
      |> Enum.map(fn x -> String.split(x, ",") end)

    {rules, updates}
  end

  def get_mid(l) do
    Enum.at(l, Integer.floor_div(length(l), 2)) |> Integer.parse() |> elem(0)
  end

  def swap(a, i1, i2) do
    e1 = Enum.at(a, i1)
    e2 = Enum.at(a, i2)

    a
    |> List.replace_at(i1, e2)
    |> List.replace_at(i2, e1)
  end
end

defmodule Star1 do
  def run({rules, updates}) do
    rule_map =
      rules
      |> Enum.reduce(%{}, fn r, acc ->
        r1 = elem(r, 0)
        Map.update(acc, r1, [r], fn existing -> existing ++ [r] end)
      end)

    updates
    |> Enum.map(fn update ->
      Enum.map(update, fn page_no -> Map.get(rule_map, page_no, []) end)
      |> List.flatten()
      |> Enum.reduce(0, fn rule, acc -> if test_rule(update, rule), do: acc, else: acc + 1 end)
      |> then(fn x -> if x == 0, do: Common.get_mid(update), else: 0 end)
    end)
    |> Enum.sum()
  end

  def test_rule(list, rule) do
    a_idx = Enum.find_index(list, fn x -> x == elem(rule, 0) end)
    b_idx = Enum.find_index(list, fn x -> x == elem(rule, 1) end)
    if a_idx == nil || b_idx == nil, do: true, else: a_idx < b_idx
  end
end

defmodule Star2 do
  def run({rules, updates}) do
    rule_map =
      rules
      |> Enum.reduce(%{}, fn r, acc ->
        r1 = elem(r, 0)
        Map.update(acc, r1, [r], fn existing -> existing ++ [r] end)
      end)

    updates
    |> Enum.filter(fn update ->
      update
      |> Enum.map(fn page_no -> Map.get(rule_map, page_no, []) end)
      |> List.flatten()
      |> Enum.reduce(0, fn rule, acc ->
        if Star1.test_rule(update, rule), do: acc, else: acc + 1
      end)
      |> then(fn x -> x > 0 end)
    end)
    |> Enum.map(fn update ->
      Enum.map(update, fn page_no -> Map.get(rule_map, page_no, []) end)
      |> List.flatten()
      |> Enum.filter(fn rule -> Enum.any?(update, fn x -> x == elem(rule, 1) end) end)
      |> then(fn relevant_rules ->
        relevant_rules
        |> Enum.map(fn x -> elem(x, 0) end)
        |> Enum.frequencies()
        |> Map.to_list()
        |> Enum.sort(fn a, b ->
          elem(a, 1) >= elem(b, 1)
        end)
        |> Enum.map(fn x -> elem(x, 0) end)
      end)
      |> then(fn x -> Enum.uniq(x ++ update) end)
      |> Common.get_mid()
    end)
    |> Enum.sum()
  end

  def test_rule(list, rule) do
    a_idx = Enum.find_index(list, fn x -> x == elem(rule, 0) end)
    b_idx = Enum.find_index(list, fn x -> x == elem(rule, 1) end)

    cond do
      a_idx == nil || b_idx == nil -> true
      a_idx < b_idx -> true
      true -> {a_idx, b_idx}
    end
  end
end

IO.puts(Star1.run(Common.open("example.txt")))
IO.puts(Star1.run(Common.open("input.txt")))
IO.inspect(Star2.run(Common.open("example.txt")))
IO.inspect(Star2.run(Common.open("input.txt")))
