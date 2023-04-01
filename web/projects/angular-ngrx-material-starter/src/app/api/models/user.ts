export class User {
  id = 0;
  openid = '';
  nickname = '';
  sex = 0;
  headimgurl = '';
  phone = '';
  roleNames: string[] = [];
  uuid = '';

  constructor(payload: Partial<User>) {
    this.id = payload.id || 0
    this.openid = payload.openid || ''
    this.nickname = payload.nickname || ''
    this.sex = payload.sex || 0
    this.headimgurl = payload.headimgurl || ''
    this.phone = payload.phone || ''
    this.roleNames = payload.roleNames || []
    this.uuid = payload.uuid || ''
  }

  hasRole(role: string): boolean {
    if (this.roleNames === undefined) {
      return false
    }
    for (const roleName of this.roleNames) {
      if (roleName === role) {
        return true
      }
    }
    return false
  }

  hasAnyRole(roles: string[]): boolean {
    if (roles.length === 0) {
      return true
    }
    if (this.roleNames === undefined) {
      return false
    }
    console.log(roles, this.roleNames)
    for (const userRole of this.roleNames) {
      for (const checkRole of roles) {
        if (userRole === checkRole) {
          return true
        }
      }
    }
    return false
  }

  displayRoles(): string {
    if (this.roleNames === undefined) {
      return ''
    }
    return displayRoles(this.roleNames)
  }
}

export function displayRoles(roleNames: string[]): string {
  let ret = ''
  for (let i = 0; i < roleNames?.length; i++) {
    if (i === 0) {
      ret += displayRoleName(roleNames[i])
    } else {
      ret += '|' + displayRoleName(roleNames[i])
    }
  }
  return ret
}

export function displayRoleName(name: string): string {
  switch (name) {
    case ROLE_ADMIN:
      return '管理员'
    case ROLE_COACH:
      return '教练'
    case ROLE_USER:
      return '会员'
    default:
      return name
  }
}

// 管理员
export const ROLE_ADMIN = 'admin'

// 用户
export const ROLE_USER = 'user'

// 教练
export const ROLE_COACH = 'coach'

// 教练信息
export interface ICoach {
  id?: number;
  userID?: number;
  level?: string; // 等级
  teachingSpace?: string; // 任教单位
  teachingAge?: string; // 任教年限
  teachingExperiences?: string[]; // 任教经历
}

// 显示教练等级
export function displayLevel(level: string | undefined) {
  if (level === undefined) {
    return '未知'
  }
  switch (level) {
    case '1-1':
      return '初级1'
  }
  return level;
}


export class UsersPage {
  users: User[] = [];
  page = 0;
  count = 0;
  limit = 0;

  constructor(payload: Partial<UsersPage>) {
    const users = []
    for (const user of payload.users || []) {
      users.push(new User(user))
    }
    this.users = users
    this.page = payload.page || 0
    this.count = payload.count || 0
    this.limit = payload.limit || 0
  }
}
