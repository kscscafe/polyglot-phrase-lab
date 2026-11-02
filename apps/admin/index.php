<?php
declare(strict_types=1);

$api = getenv('PHRASE_API_URL') ?: 'http://127.0.0.1:8080';
$message = '';
$phrases = [];
$context = stream_context_create(['http' => ['timeout' => 2]]);
$response = @file_get_contents($api . '/api/phrases', false, $context);
if ($response === false) {
    $message = 'APIに接続できません。Go APIを起動してください。';
} else {
    $decoded = json_decode($response, true);
    $phrases = is_array($decoded) ? $decoded : [];
}
function e(string $value): string { return htmlspecialchars($value, ENT_QUOTES | ENT_SUBSTITUTE, 'UTF-8'); }
?>
<!doctype html>
<html lang="ja"><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>Phrase Lab 教材管理</title>
<style>
body{margin:0;background:#f1f5f8;color:#13263a;font:16px system-ui,sans-serif}.wrap{max-width:900px;margin:48px auto;padding:0 20px}h1{font-size:28px}.note{color:#52687c}.card{background:white;border:1px solid #d9e3eb;border-radius:14px;padding:20px;margin:14px 0}.zh{font-size:22px;font-weight:700}.ja{color:#566b7d}.tag{display:inline-block;background:#e5f8fa;color:#087e8b;padding:4px 9px;border-radius:20px;font-size:12px}
</style>
<main class="wrap"><p class="tag">ADMIN · READ ONLY</p><h1>教材一覧</h1><p class="note">この画面は教材確認用です。編集APIはまだありません。</p>
<?php if ($message !== ''): ?><p class="card"><?= e($message) ?></p><?php endif; ?>
<?php foreach ($phrases as $phrase): ?>
<article class="card"><span class="tag"><?= e((string)($phrase['level'] ?? '')) ?></span><p class="zh"><?= e((string)($phrase['zh'] ?? '')) ?></p><p><?= e((string)($phrase['pinyin'] ?? '')) ?></p><p class="ja"><?= e((string)($phrase['ja'] ?? '')) ?></p></article>
<?php endforeach; ?>
</main></html>
