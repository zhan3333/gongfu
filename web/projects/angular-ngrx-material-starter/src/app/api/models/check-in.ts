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
  date?: string,
  // 日排名
  dayRank?: number,
}

export type CheckInList = CheckIn[]

export interface CheckInExist {
  exists: boolean,
  checkIn?: CheckIn,
}

export interface CheckInCount {
  id?: number;
  userName?: string;
  userID?: number;
  headImgUrl?: string;
  checkInCount?: number;
  checkInContinuous?: number;
}

export type CheckInCountList = CheckInCount[]
