const lat: Record<string, string> = {
  а: "a", б: "b", в: "v", г: "g", д: "d", е: "e", ё: "e", ж: "zh", з: "z", и: "i", й: "y",
  к: "k", л: "l", м: "m", н: "n", о: "o", п: "p", р: "r", с: "s", т: "t", у: "u", ф: "f",
  х: "h", ц: "ts", ч: "ch", ш: "sh", щ: "sch", ъ: "", ы: "y", ь: "", э: "e", ю: "yu", я: "ya",
};

export function transliterate(text: string) {
  return [...text]
    .map((char) => {
      const lower = lat[char.toLowerCase()];
      if (lower === undefined) return char;
      if (char === char.toLowerCase()) return lower;
      return lower.charAt(0).toUpperCase() + lower.slice(1);
    })
    .join("");
}

export function previewText(text: string, options: { translit: boolean; shortLink: boolean }) {
  let next = text;
  if (options.translit) next = transliterate(next);
  if (options.shortLink) next = next.replace(/https?:\/\/\S+/g, "http://localhost:8080/s/preview");
  return next;
}

export function smsParts(text: string) {
  const unicode = /[^\u0000-\u007f£$¥èéùìòÇØøÅåΔΦΓΛΩΠΨΣΘΞÆæßÉ¤¡ÄÖÑÜ§¿äöñüà^{}\\[~\]|€]/.test(text);
  const single = unicode ? 70 : 160;
  const multi = unicode ? 67 : 153;
  const length = [...text].length;
  const parts = length === 0 ? 0 : length <= single ? 1 : Math.ceil(length / multi);
  return { length, parts, limit: parts <= 1 ? single : multi };
}
