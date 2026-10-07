const fs = require('fs');
const path = require('path');

const enPath = path.join(__dirname, 'i18n/locales/en.json');
const esPath = path.join(__dirname, 'i18n/locales/es.json');

const newEn = JSON.parse(process.argv[2]);
const newEs = JSON.parse(process.argv[3]);

function mergeDeep(target, source) {
  for (const key in source) {
    if (source[key] instanceof Object && key in target) {
      Object.assign(source[key], mergeDeep(target[key], source[key]));
    }
  }
  Object.assign(target || {}, source);
  return target;
}

const enData = fs.existsSync(enPath) ? JSON.parse(fs.readFileSync(enPath, 'utf8')) : {};
const esData = fs.existsSync(esPath) ? JSON.parse(fs.readFileSync(esPath, 'utf8')) : {};

mergeDeep(enData, newEn);
mergeDeep(esData, newEs);

fs.writeFileSync(enPath, JSON.stringify(enData, null, 2) + '\n');
fs.writeFileSync(esPath, JSON.stringify(esData, null, 2) + '\n');
