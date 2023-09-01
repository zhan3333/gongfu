import { displayRoles, ICoach, ROLE_COACH } from './user';
import { StudyRecord } from './study-record';
import { TeachingRecord } from './teaching-record';

export class Profile {
  id = 0;
  nickname = '';
  headimgurl = '';
  roleNames: string[] = [];
  uuid = '';
  coach?: ICoach;
  studyRecords: StudyRecord[]; // 学习记录
  teachingRecords: TeachingRecord[]; // 授课记录

  constructor(payload: Partial<Profile>) {
    this.id = payload.id || 0;
    this.nickname = payload.nickname || '';
    this.headimgurl = payload.headimgurl || '';
    this.roleNames = payload.roleNames || [];
    this.uuid = payload.uuid || '';
    this.coach = payload.coach;
    this.studyRecords = payload.studyRecords || [];
    this.teachingRecords = payload.teachingRecords || [];
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
}
