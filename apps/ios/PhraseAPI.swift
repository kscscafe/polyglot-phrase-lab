import Foundation

struct Practice: Decodable {
    let phraseId: String
    let level: String
    let tokens: [String]
    let prompt: String
}

struct CheckResult: Decodable {
    let correct: Bool
    let zh: String
    let pinyin: String
    let ja: String
}

enum PhraseAPIError: Error { case invalidResponse, status(Int) }

final class PhraseAPI {
    private let baseURL: URL
    private let session: URLSession

    init(baseURL: URL = URL(string: "http://localhost:8080")!, session: URLSession = .shared) {
        self.baseURL = baseURL
        self.session = session
    }

    func newPractice(level: String) async throws -> Practice {
        var components = URLComponents(url: URL(string: "/api/practice/new", relativeTo: baseURL)!, resolvingAgainstBaseURL: false)!
        components.queryItems = [URLQueryItem(name: "level", value: level)]
        let (data, response) = try await session.data(from: components.url!)
        guard let http = response as? HTTPURLResponse else { throw PhraseAPIError.invalidResponse }
        guard (200..<300).contains(http.statusCode) else { throw PhraseAPIError.status(http.statusCode) }
        return try JSONDecoder().decode(Practice.self, from: data)
    }

    func check(phraseId: String, tokens: [String]) async throws -> CheckResult {
        var request = URLRequest(url: URL(string: "/api/practice/check", relativeTo: baseURL)!)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.httpBody = try JSONEncoder().encode(CheckRequest(phraseId: phraseId, tokens: tokens))
        let (data, response) = try await session.data(for: request)
        guard let http = response as? HTTPURLResponse else { throw PhraseAPIError.invalidResponse }
        guard (200..<300).contains(http.statusCode) else { throw PhraseAPIError.status(http.statusCode) }
        return try JSONDecoder().decode(CheckResult.self, from: data)
    }
}

private struct CheckRequest: Encodable {
    let phraseId: String
    let tokens: [String]
}
