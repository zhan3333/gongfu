import { displayRoles, ICoach, ROLE_ADMIN, ROLE_COACH, ROLE_MEMBER } from './user';
import { StudyRecord } from './study-record';
import { TeachingRecord } from './teaching-record';
import { MemberCourse } from '../admin/admin-api.service';

export class Profile {
  id = 0;
  nickname = '';
  headimgurl = '';
  roleNames: string[] = [];
  uuid = '';
  coach?: ICoach;
  studyRecords: StudyRecord[]; // 学习记录
  teachingRecords: TeachingRecord[]; // 授课记录
  memberCourses: MemberCourse[]; // 会员课程

  constructor(payload: Partial<Profile>) {
    this.id = payload.id || 0;
    this.nickname = payload.nickname || '';
    this.headimgurl = payload.headimgurl || '';
    this.roleNames = payload.roleNames || [];
    this.uuid = payload.uuid || '';
    this.coach = payload.coach;
    this.studyRecords = payload.studyRecords || [];
    this.teachingRecords = payload.teachingRecords || [];
    this.memberCourses = payload.memberCourses || [];
  }

  displayRoles(): string {
    if (this.roleNames === undefined) {
      return '';
    }
    return displayRoles(this.roleNames);
  }

  isCoach(): boolean {
    return -1 !== this.roleNames.indexOf(ROLE_COACH);
  }

  isAdmin(): boolean {
    return -1 !== this.roleNames.indexOf(ROLE_ADMIN);
  }

  // 是否是会员
  isMember(): boolean {
    return -1 !== this.roleNames.indexOf(ROLE_MEMBER);
  }
}
