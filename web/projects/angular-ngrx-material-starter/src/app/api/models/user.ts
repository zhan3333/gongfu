export class User {
  id?: number;
  openid?: string;
  nickname?: string;
  sex?: Number;
  province?: string;
  city?: string;
  country?: string;
  headimgurl?: string;
  unionid?: string;
  phone?: string;
  role?: string;

  isAdmin() {
    return this.role === ROLE_ADMIN
  }

  isUser() {
    return this.role === ROLE_USER
  }

  isCoach() {
    return this.role === ROLE_COACH
  }
}

// 管理员
export const ROLE_ADMIN = 'admin'

// 用户
export const ROLE_USER = 'user'

// 教练
export const ROLE_COACH = 'coach'

// 教练信息
export class Coach {
  id?: number;
  userID?: number;
  level?: string; // 等级
  teachingSpace?: string; // 任教单位
  teachingAge?: string; // 任教年限
  teachingExperiences?: string[]; // 任教经历
}
