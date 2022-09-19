export interface CheckIn {
  id?: number,
  // 创建时间
  createdAt?: number,
  // 播放链接
  url?: string,
  fileName?: string,
  userName?: string,
  headImgUrl?: string,
  key?: string,
}

export interface CheckInExist {
  exists: boolean,
  checkIn?: CheckIn,
}
