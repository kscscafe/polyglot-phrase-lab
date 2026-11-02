export type Phrase = {
  id: string;
  level: string;
  zh: string;
  pinyin: string;
  ja: string;
  tokens: string[];
};

export type Practice = {
  phraseId: string;
  level: string;
  tokens: string[];
  prompt: string;
};

export type PracticeResult = {
  correct: boolean;
  zh: string;
  pinyin: string;
  ja: string;
};

export type CheckRequest = Pick<Practice, "phraseId" | "tokens">;

export const API_BASE = "http://localhost:8080";
