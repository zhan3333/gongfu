import { displayRoles, ICoach } from './user';
import { StudyRecord } from './study-record';
import { TeachingRecord } from './teaching-record';

export class Profile {
  id = 0;
  nickname = '';
  headimgurl = '';
  role?: string;
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
    console.log('paylod', payload);
    this.studyRecords = payload.studyRecords || [];
    this.teachingRecords = payload.teachingRecords || [];
  }

  displayRoles(): string {
    if (this.roleNames === undefined) {
      return '';
    }
    return displayRoles(this.roleNames);
  }
}
