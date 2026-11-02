const API = "http://localhost:8080";
const level = document.querySelector("#level");
const choices = document.querySelector("#choices");
const answer = document.querySelector("#answer");
const result = document.querySelector("#result");
const prompt = document.querySelector("#prompt");
const checkButton = document.querySelector("#check");
let exercise = null;
let selected = [];

async function request(path, options) {
  const response = await fetch(`${API}${path}`, options);
  const body = await response.json();
  if (!response.ok) throw new Error(body.error ?? `HTTP ${response.status}`);
  return body;
}

function render() {
  answer.replaceChildren();
  if (selected.length === 0) {
    const placeholder = document.createElement("span");
    placeholder.className = "placeholder";
    placeholder.textContent = "ここに語句が並びます";
    answer.append(placeholder);
  } else {
    selected.forEach(({ token, index }) => {
      const button = document.createElement("button");
      button.className = "token";
      button.type = "button";
      button.textContent = token;
      button.title = "タップして語句を戻す";
      button.addEventListener("click", () => { selected = selected.filter((item) => item.index !== index); render(); });
      answer.append(button);
    });
  }
  choices.replaceChildren();
  exercise?.tokens.forEach((token, index) => {
    const button = document.createElement("button");
    button.className = "token";
    button.type = "button";
    button.textContent = token;
    button.disabled = selected.some((item) => item.index === index);
    button.addEventListener("click", () => { selected.push({ token, index }); render(); });
    choices.append(button);
  });
  checkButton.disabled = !exercise || selected.length !== exercise.tokens.length;
}

async function newExercise() {
  selected = [];
  exercise = null;
  result.textContent = "";
  checkButton.disabled = true;
  prompt.textContent = "問題を読み込んでいます…";
  try {
    exercise = await request(`/api/practice/new?level=${encodeURIComponent(level.value)}`);
    prompt.textContent = exercise.prompt;
    render();
  } catch (error) {
    prompt.textContent = "APIに接続できません。Go APIを先に起動してください。";
    result.textContent = error.message;
  }
}

document.querySelector("#clear").addEventListener("click", () => { selected = []; result.textContent = ""; render(); });
document.querySelector("#check").addEventListener("click", async () => {
  if (!exercise || selected.length !== exercise.tokens.length) return;
  try {
    const outcome = await request("/api/practice/check", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ phraseId: exercise.phraseId, tokens: selected.map((item) => item.token) }),
    });
    const status = document.createElement("strong");
    status.textContent = outcome.correct ? "正解です" : "もう一度確認してみましょう";
    if (!outcome.correct) status.className = "wrong";
    result.replaceChildren(status, document.createElement("br"), document.createTextNode(`${outcome.zh}　${outcome.pinyin}　${outcome.ja}`));
  } catch (error) { result.textContent = error.message; }
});
level.addEventListener("change", newExercise);
newExercise();
