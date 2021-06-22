export const encodeToString = (data: string): string => {
  return data
    .split('')
    .map(c => c.charCodeAt(0).toString(16).padStart(2, '0'))
    .join('');
};
export const decodeString = (data: string): string => {
  return data
    .split(/(\w\w)/g)
    .filter(p => !!p)
    .map(c => String.fromCharCode(parseInt(c, 16)))
    .join('');
};

/** decodeStringAsHexArray
 *
 * @param data raw data
 * @return Uint8Array returns the bytes represented by the hexadecimal string raw
 */
export const decodeStringAsHexArray = (data: string): Uint8Array => {
  return new Uint8Array(
    data
      .split(/(\w\w)/g)
      .filter(p => !!p)
      .map(c => parseInt(c, 16)),
  );
};

export const str2array = (str: string): Uint8Array => {
  const bufView = new Uint8Array(str.length);
  for (let i = 0, strLen = str.length; i < strLen; i++) {
    bufView[i] = str.charCodeAt(i);
  }
  return bufView;
};

export const array2str = (data: Uint8Array): string => {
  let out: string = '';
  data.forEach(value => (out += String.fromCharCode(value)));

  return out;
};

export const string2Bin = (str: string): number[] => {
  const result = [];
  for (let i = 0; i < str.length; i++) {
    result.push(str.charCodeAt(i));
  }
  return result;
};
export const bin2String = (array: number[]): string => {
  return String.fromCharCode.apply(String, array);
};
