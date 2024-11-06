import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

export interface Course {
  id: number;
  schoolName: string;
  startDate: string;
  startTime: string;
  managerName: string;
  coachName: string;
  assistantNames: string[];
  checkInAt: string;
  checkOutAt: string;
  summary: string;
  content: string;
}

@Injectable({
  providedIn: 'root'
})
export class CourseService {

  constructor(
    private http: HttpClient
  ) {
  }

  public getCourseList(userId: number) {
    return this.http.get<Course[]>('/admin/course/list', {
      params: {
        userId: userId,
      }
    });
  }
}
