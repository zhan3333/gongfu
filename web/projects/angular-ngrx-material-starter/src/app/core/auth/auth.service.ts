import { Injectable } from '@angular/core';
import { LocalStorageService } from '../local-storage/local-storage.service';
import { Observable, ReplaySubject } from 'rxjs';
import { ActivatedRouteSnapshot, Router, RouterStateSnapshot, UrlTree } from '@angular/router';
import { User } from '../../api/models/user';

const KEY_IS_AUTHENTICATED = 'isAuthenticated'
const KEY_ACCESS_TOKEN = 'accessToken'
const KEY_USER = 'user'
const KEY_AUTH = 'auth'

@Injectable({
  providedIn: 'root'
})
export class AuthService  {
  public isAuthenticated$ = new ReplaySubject<boolean>();
  public accessToken$ = new ReplaySubject<string>();
  public user$ = new ReplaySubject<User>();

  private store

  constructor(store: LocalStorageService, private router: Router) {
    this.store = store
    this.isAuthenticated$.next(this.isAuthenticated())
    this.accessToken$.next(this.getAccessToken())
    this.user$.next(this.getUser())
  }

  isAuthenticated(): boolean {
    return this.store.getItem(KEY_AUTH)[KEY_IS_AUTHENTICATED] || false
  }

  getAccessToken(): string {
    return this.store.getItem(KEY_AUTH)[KEY_ACCESS_TOKEN] || ''
  }

  getUser(): User {
    return new User(this.store.getItem(KEY_USER))
  }

  setUser(user: User) {
    this.store.setItem(KEY_USER, user)
    this.user$.next(user)
  }

  login(accessToken: string) {
    this.store.setItem(KEY_AUTH, {
      isAuthenticated: true,
      accessToken: accessToken
    })
    this.isAuthenticated$.next(true)
    this.accessToken$.next(accessToken)
  }

  logout() {
    this.store.removeItem(KEY_AUTH)
  }

  canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<boolean | UrlTree> | Promise<boolean | UrlTree> | boolean | UrlTree {
    if (this.isAuthenticated()) {
      return true
    } else {
      return this.router.parseUrl('/login?type=login')
    }
  }
}
