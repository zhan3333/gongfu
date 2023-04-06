import { School } from './school';
import { ISampleCoach } from '../admin/admin-api.service';

export interface Course {
  id: number,
  createdAt: number,
  school: School,
  startDate: string,
  startTime: string,
  manager: ISampleCoach,
  coach?: ISampleCoach
  assistantCoaches: ISampleCoach[]
  assistantCoachIds: number[]
  checkInByUser?: ISampleCoach
  checkOutByUser?: ISampleCoach
  checkInAt: number,
  checkOutAt: number,
  images?: string[],
  summary: string,
  content: string,
}

export class CoursesPage {
  items: Course[] = [];
  page = 0;
  count = 0;
  limit = 0;

  constructor(payload: Partial<CoursesPage>) {
    const items = []
    for (const item of payload.items || []) {
      items.push(item)
    }
    this.items = items
    this.page = payload.page || 0
    this.count = payload.count || 0
    this.limit = payload.limit || 0
  }
}
