import { ICoach, displayRoles } from './user';

export class Profile {
  id = 0;
  nickname = '';
  headimgurl = '';
  role?: string;
  roleNames: string[] = [];
  uuid = '';
  coach?: ICoach;

  constructor(payload: Partial<Profile>) {
    this.id = payload.id || 0
    this.nickname = payload.nickname || ''
    this.headimgurl =payload.headimgurl || ''
    this.roleNames = payload.roleNames || []
    this.uuid = payload.uuid || ''
    this.coach = payload.coach
  }


  displayRoles(): string {
    if (this.roleNames === undefined) {
      return ''
    }
    return displayRoles(this.roleNames)
  }
}
