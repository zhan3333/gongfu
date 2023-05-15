// all support coach levels
export const Levels = [
  { key: '实习教练', value: '0-1' },
  { key: '初级一', value: '1-1' },
  { key: '初级二', value: '1-2' },
  { key: '初级三', value: '1-3' },

  { key: '中级一', value: '2-1' },
  { key: '中级二', value: '2-2' },
  { key: '中级三', value: '2-3' },

  { key: '高级一', value: '3-1' },
  { key: '高级二', value: '3-2' },
  { key: '高级三', value: '3-3' },

  { key: '特级', value: '4-1' }
];

// coach level display
export function displayLevel(level: string | undefined) {
  for (let lv of Levels) {
    if (lv.value === level) {
      return lv.key;
    }
  }
  return '未知';
}
