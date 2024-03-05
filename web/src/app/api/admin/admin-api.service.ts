import { Injectable } from '@angular/core';
import { ICoach, User, UsersPage } from '../models/user';
import { map } from 'rxjs/operators';
import { HttpClient } from '@angular/common/http';
import { Course, CoursesPage } from '../models/course';
import { School } from '../models/school';

const getUsers = '/admin/users';
const getUser = '/admin/user';
const updateUser = '/admin/user';
const getCoach = '/admin/coach';

@Injectable({
  providedIn: 'root'
})
export class AdminApiService {
  constructor(private http: HttpClient) {
  }

  getUsers(data: GetUsersParams) {
    return this.http
      .get<UsersPage>(getUsers, {
        params: {
          page: data.page,
          limit: data.limit,
          keyword: data.keyword,
          desc: data.desc,
          roleIds: data.roleIds,
        }
      })
      .pipe(map((v) => new UsersPage(v)));
  }

  // 获取用户信息
  getUser(id: number) {
    return this.http
      .get<User>(getUser + '/' + id)
      .pipe(map((v) => new User(v)));
  }

  updateUser(id: number, data: UpdateUserParams) {
    return this.http.put(updateUser + '/' + id, data);
  }

  getCoach(userID: number) {
    return this.http.get<ICoach>(getCoach + '/' + userID);
  }

  // 获取角色名列表
  getRoleNames() {
    return this.http.get<string[]>('/admin/role-names');
  }

  // 获取所有教练
  getCoaches() {
    return this.http.get<ISampleCoach[]>('/admin/coaches');
  }

  // 获取学校列表
  getSchools() {
    return this.http.get<School[]>('/admin/schools');
  }

  createCourse(data: CreateCourseInput) {
    return this.http.post('/admin/course', data);
  }

  // 获取课程列表
  getCoursesPage(input: GetCoursesPageInput) {
    return this.http
      .get('/admin/courses', {params: {...input}})
      .pipe(map((v) => new CoursesPage(v)));
  }

  // 获取课程详情
  getCourse(id: number) {
    return this.http.get<Course>('/admin/courses/' + id);
  }

  updateCourse(id: number, course: any) {
    return this.http.put('/admin/courses/' + id, course);
  }

  deleteCourse(id: number) {
    return this.http.delete('/admin/courses/' + id);
  }

  editTeachingRecord(data: {
    id?: number;
    date: string;
    address: string;
    userId: number;
  }) {
    return this.http.post('/admin/teaching-record', data);
  }

  editStudyRecord(data: {
    id?: number;
    date: string;
    content: string;
    userId: number;
  }) {
    return this.http.post('/admin/study-record', data);
  }

  deleteTeachingRecord(id: number) {
    return this.http.delete('/admin/teaching-record/' + id);
  }

  deleteStudyRecord(id: number) {
    return this.http.delete('/admin/study-record/' + id);
  }

  getMemberCourses(userId: number) {
    return this.http.get<MemberCourse[]>('/admin/member-courses', {
      params: {
        userId
      }
    });
  }

  delMemberCourse(memberCourseId: number) {
    return this.http.delete(`/admin/member-courses/${memberCourseId}`);
  }

  createMemberCourse(data: CreateMemberCourseParams) {
    console.log('create', data)
    return this.http.post('/admin/member-courses', data);
  }

  updateMemberCourse(id: number, params: UpdateMemberCourseParams) {
    return this.http.put(`/admin/member-courses/${id}`, params);
  }
}

export interface CreateMemberCourseParams {
  name: string
  userId: number
  startTime: string
  endTime: string
  total: number
  remark: string
}

export interface MemberCourse {
  id: number
  name: string
  userId: number
  startTime: string
  endTime: string
  total: number
  remain: number
  remark: string
  status: string
}

export interface CreateCourseInput {
  coachId: number;
  assistantCoachIds: number[];
  schoolId: number;
  startDate: string;
  startTime: string;
  managerId: number;
}

export interface UpdateMemberCourseParams {
  name: string
  startTime: string
  endTime: string
  total: number
  remark: string
  status: string
}

// 简易的教练信息
export interface ISampleCoach {
  id: number;
  name: string;
}

export interface GetUsersParams {
  page: number;
  limit: number;
  keyword: string;
  desc: boolean;
  roleIds: number[];
}

export interface GetCoursesPageInput {
  page: number;
  limit: number;
  keyword: string;
  desc: boolean;
  coachId?: number;
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
