import { TeachingRecord } from './teaching-record';
import { StudyRecord } from './study-record';

export interface User {
  id: number;
  openid: string;
  nickname: string;
  sex: number;
  headimgurl: string;
  phone: string;
  roleNames: string[];
  uuid: string;
  teachingRecords: TeachingRecord[];
  studyRecords: StudyRecord[];
  coachStatus: string;
  coachRegisterDate: string;
}

export class UserClass {
  id = 0;
  openid = '';
  nickname = '';
  sex = 0;
  headimgurl = '';
  phone = '';
  roleNames: string[] = [];
  uuid = '';
  teachingRecords: TeachingRecord[] = [];
  studyRecords: StudyRecord[] = [];
  coachStatus: string = '';

  constructor(payload: Partial<UserClass>) {
    this.id = payload.id || 0;
    this.openid = payload.openid || '';
    this.nickname = payload.nickname || '';
    this.sex = payload.sex || 0;
    this.headimgurl = payload.headimgurl || '';
    this.phone = payload.phone || '';
    this.roleNames = payload.roleNames || [];
    this.uuid = payload.uuid || '';
    this.teachingRecords = payload.teachingRecords || [];
    this.studyRecords = payload.studyRecords || [];
  }

  hasRole(role: string): boolean {
    if (this.roleNames === undefined) {
      return false;
    }
    for (const roleName of this.roleNames) {
      if (roleName === role) {
        return true;
      }
    }
    return false;
  }

  hasAnyRole(roles: string[]): boolean {
    if (roles.length === 0) {
      return true;
    }
    if (this.roleNames === undefined) {
      return false;
    }
    for (const userRole of this.roleNames) {
      for (const checkRole of roles) {
        if (userRole === checkRole) {
          return true;
        }
      }
    }
    return false;
  }

  displayRoles(): string {
    if (this.roleNames === undefined) {
      return '';
    }
    return displayRoles(this.roleNames);
  }
}

export function displayRoles(roleNames: string[]): string {
  let ret = '';
  for (let i = 0; i < roleNames?.length; i++) {
    if (i === 0) {
      ret += displayRoleName(roleNames[i]);
    } else {
      ret += '|' + displayRoleName(roleNames[i]);
    }
  }
  return ret;
}

export function displayRoleName(name: string): string {
  switch (name) {
    case ROLE_ADMIN:
      return '管理员';
    case ROLE_COACH:
      return '教练';
    case ROLE_USER:
      return '用户';
    case ROLE_MEMBER:
      return '会员';
    default:
      return name;
  }
}

// 管理员
export const ROLE_ADMIN = 'admin';

// 用户
export const ROLE_USER = 'user';

// 教练
export const ROLE_COACH = 'coach';

export const ROLE_MEMBER = "member"

// 教练信息
export interface ICoach {
  id?: number;
  userID?: number;
  level?: string; // 等级
  teachingSpace?: string; // 任教单位
  teachingAge?: string; // 任教年限
  teachingExperiences?: string[]; // 任教经历
}

export class UsersPage {
  users: UserClass[] = [];
  page = 0;
  count = 0;
  limit = 0;

  constructor(payload: Partial<UsersPage>) {
    const users = [];
    for (const user of payload.users || []) {
      users.push(new UserClass(user));
    }
    this.users = users;
    this.page = payload.page || 0;
    this.count = payload.count || 0;
    this.limit = payload.limit || 0;
  }
}

export interface SimpleUser {
  id: number;
  uuid: number;
  headImgUrl: string;
  name: string;
}
