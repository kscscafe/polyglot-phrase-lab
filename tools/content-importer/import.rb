#!/usr/bin/env ruby
# Convert UTF-8 CSV columns id,level,zh,pinyin,ja,tokens into the shared JSON format.
require "csv"
require "json"

input, output = ARGV
abort "usage: ruby import.rb input.csv output.json" unless input && output

required = %w[id level zh pinyin ja tokens]
rows = CSV.read(input, headers: true, encoding: "bom|utf-8")
missing = required - (rows.headers || [])
abort "missing CSV columns: #{missing.join(', ')}" unless missing.empty?

phrases = rows.map.with_index(2) do |row, line|
  record = required.to_h { |key| [key, row[key]&.strip] }
  blank = required.reject { |key| key == "pinyin" }.select { |key| record[key].to_s.empty? }
  abort "line #{line}: missing #{blank.join(', ')}" unless blank.empty?
  record["tokens"] = record["tokens"].split("|").map(&:strip)
  abort "line #{line}: tokens must not be empty" if record["tokens"].empty?
  record
end

File.write(output, JSON.pretty_generate(phrases) + "\n")
puts "Imported #{phrases.length} phrases to #{output}"
