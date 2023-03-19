import { ICoach } from './user';

export interface Course {
  id: number,
  createdAt: number,
  schoolStartAt: number,
  address: string,
  name: string,
  coach?: ICoach
  assistantCoaches: ICoach[]
  checkInBy?: ICoach
  checkOutBy?: ICoach
  checkInAt: number,
  checkOutAt: number,
  images: string[],
  summary: string,
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
