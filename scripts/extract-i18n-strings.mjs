/**
 * One-off: extract the `strings` object from js/i18n.js and write web/src/landing/strings.ts
 * Run from repo root: node scripts/extract-i18n-strings.mjs
 */
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const root = path.resolve(__dirname, '..');
const code = fs.readFileSync(path.join(root, 'js/i18n.js'), 'utf8');

const startMarker = 'const strings = ';
const start = code.indexOf(startMarker);
if (start === -1) throw new Error('strings not found');
let pos = start + startMarker.length;
let depth = 0;
for (; pos < code.length; pos++) {
  const c = code[pos];
  if (c === '{') depth++;
  else if (c === '}') { depth--; if (depth === 0) break; }
}
const objStr = code.slice(start + startMarker.length, pos + 1);
const strings = new Function('return ' + objStr)();

const out = `/** Auto-generated from js/i18n.js — do not edit by hand. Run: node scripts/extract-i18n-strings.mjs */\nexport const strings: Record<string, Record<string, string>> = ${JSON.stringify(strings, null, 0)};\n`;
const outPath = path.join(root, 'web/src/landing/strings.ts');
fs.mkdirSync(path.dirname(outPath), { recursive: true });
fs.writeFileSync(outPath, out, 'utf8');
console.log('Written', outPath);
