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
