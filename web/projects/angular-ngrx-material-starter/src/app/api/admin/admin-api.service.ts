import { Injectable } from '@angular/core';
import { ICoach, User, UsersPage } from '../models/user';
import { map } from 'rxjs/operators';
import { HttpClient } from '@angular/common/http';
import { Course, CoursesPage } from '../models/course';
import { School } from '../models/school';


const getUsers = '/admin/users'
const getUser = '/admin/user'
const updateUser = '/admin/user'
const getCoach = '/admin/coach'

@Injectable({
  providedIn: 'root'
})
export class AdminApiService {

  constructor(
    private http: HttpClient,
  ) {
  }

  getUsers(data: GetUsersParams) {
    return this.http.get<UsersPage>(getUsers, {
      params: {
        page: data.page,
        limit: data.limit,
        keyword: data.keyword,
        desc: data.desc,
      }
    }).pipe(
      map(v => new UsersPage(v))
    )
  }

  // 获取用户信息
  getUser(id: number) {
    return this.http.get<User>(getUser + '/' + id).pipe(
      map(v => new User(v))
    )
  }

  updateUser(id: number, data: UpdateUserParams) {
    return this.http.put(updateUser + '/' + id, data)
  }

  getCoach(userID: number) {
    return this.http.get<ICoach>(getCoach + '/' + userID)
  }

  // 获取角色名列表
  getRoleNames() {
    return this.http.get<string[]>('/admin/role-names')
  }

  // 获取所有教练
  getCoaches() {
    return this.http.get<ISampleCoach[]>('/admin/coaches')
  }

  // 获取学校列表
  getSchools() {
    return this.http.get<School[]>('/admin/schools')
  }

  createCourse(data: CreateCourseInput) {
    return this.http.post('/admin/course', data)
  }

  // 获取课程列表
  getCoursesPage(input: GetCoursesPageInput) {
    return this.http.get('/admin/courses', {params: {...input}}).pipe(
      map(v => new CoursesPage(v))
    )
  }

  // 获取课程详情
  getCourse(id: number) {
    return this.http.get<Course>('/admin/courses/' + id)
  }

  updateCourse(id: number, course: any) {
    return this.http.put('/admin/courses/' + id, course)
  }

  deleteCourse(id: number) {
    return this.http.delete('/admin/courses/' + id)
  }
}

export interface CreateCourseInput {
  coachId: number,
  assistantCoachIds: number[],
  schoolId: number,
  startDate: string,
  startTime: string,
  managerId: number,
}

// 简易的教练信息
export interface ISampleCoach {
  id: number
  name: string
}

export interface GetUsersParams {
  page: number,
  limit: number,
  keyword: string,
  desc: boolean
}


export interface GetCoursesPageInput {
  page: number,
  limit: number,
  keyword: string,
  desc: boolean,
  coachId?: number,
}

export interface UpdateUserParams {
  level: string;
  teachingSpace: string;
  teachingAge: string;
  teachingExperiences: string[];
  phone: string;
  nickname: string;
  roleNames: string[];
}
