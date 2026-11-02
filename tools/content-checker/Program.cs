using System.Text.Json;

// Validate phrase JSON before import or publication.
if (args.Length != 1)
{
    Console.Error.WriteLine("usage: dotnet run -- path/to/phrases.json");
    return 2;
}

try
{
    await using var stream = File.OpenRead(args[0]);
    var phrases = await JsonSerializer.DeserializeAsync<List<Phrase>>(stream,
        new JsonSerializerOptions { PropertyNameCaseInsensitive = true });
    if (phrases is null || phrases.Count == 0) throw new InvalidDataException("dataset must be a non-empty array");
    var errors = new List<string>();
    var ids = new HashSet<string>(StringComparer.Ordinal);
    foreach (var (phrase, index) in phrases.Select((p, i) => (p, i)))
    {
        var row = index + 1;
        if (string.IsNullOrWhiteSpace(phrase.Id)) errors.Add($"row {row}: id is required");
        else if (!ids.Add(phrase.Id)) errors.Add($"row {row}: duplicate id '{phrase.Id}'");
        if (string.IsNullOrWhiteSpace(phrase.Level)) errors.Add($"row {row}: level is required");
        if (string.IsNullOrWhiteSpace(phrase.Zh) || string.IsNullOrWhiteSpace(phrase.Ja)) errors.Add($"row {row}: zh and ja are required");
        if (phrase.Tokens is null || phrase.Tokens.Length < 2) errors.Add($"row {row}: at least two tokens are required");
        else if (phrase.Zh is { } zh && !string.Concat(phrase.Tokens).TrimEnd('。', '！', '？').Equals(zh.TrimEnd('。', '！', '？'), StringComparison.Ordinal))
            errors.Add($"row {row}: tokens do not reconstruct zh (check tokenization/punctuation)");
    }
    foreach (var error in errors) Console.Error.WriteLine(error);
    if (errors.Count > 0) return 1;
    Console.WriteLine($"OK: {phrases.Count} phrase(s)");
    return 0;
}
catch (Exception ex)
{
    Console.Error.WriteLine($"invalid dataset: {ex.Message}");
    return 1;
}

internal sealed record Phrase(string? Id, string? Level, string? Zh, string? Pinyin, string? Ja, string[]? Tokens);
