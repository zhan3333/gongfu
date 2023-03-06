import { Injectable } from '@angular/core';
import { Coach, User, UsersPage } from '../models/user';
import { map } from 'rxjs/operators';
import { HttpClient } from '@angular/common/http';


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
    return this.http.get<Coach>(getCoach + '/' + userID)
  }

  // 获取角色名列表
  getRoleNames() {
    return this.http.get<string[]>('/admin/role-names')
  }
}


export interface GetUsersParams {
  page: number,
  limit: number,
  keyword: string,
  desc: boolean
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
