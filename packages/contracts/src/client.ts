import { API_BASE, type CheckRequest, type Practice, type PracticeResult } from "./index.js";

async function readJSON<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE}${path}`, init);
  const body: unknown = await response.json();
  if (!response.ok) {
    const message = typeof body === "object" && body !== null && "error" in body
      ? String((body as { error: unknown }).error)
      : `HTTP ${response.status}`;
    throw new Error(message);
  }
  return body as T;
}

export function fetchPractice(level: string): Promise<Practice> {
  return readJSON(`/api/practice/new?level=${encodeURIComponent(level)}`);
}

export function submitAnswer(answer: CheckRequest): Promise<PracticeResult> {
  return readJSON("/api/practice/check", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(answer),
  });
}
