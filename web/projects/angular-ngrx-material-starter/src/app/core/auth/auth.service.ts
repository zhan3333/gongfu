import { Injectable } from '@angular/core';
import { LocalStorageService } from '../local-storage/local-storage.service';
import { Observable, ReplaySubject, Subject } from 'rxjs';
import { ActivatedRouteSnapshot, CanActivate, RouterStateSnapshot, UrlTree } from '@angular/router';

const KEY_IS_AUTHENTICATED = 'isAuthenticated'
const KEY_ACCESS_TOKEN = 'accessToken'
const KEY_AUTH = 'auth'

@Injectable({
  providedIn: 'root'
})
export class AuthService implements CanActivate {
  public isAuthenticated$ = new ReplaySubject<boolean>();
  public accessToken$ = new ReplaySubject<string>();

  private store

  constructor(store: LocalStorageService) {
    this.store = store
    console.log('auth init', this.isAuthenticated(), this.getAccessToken())
    this.isAuthenticated$.next(this.isAuthenticated())
    this.accessToken$.next(this.getAccessToken())
  }

  isAuthenticated(): boolean {
    return this.store.getItem(KEY_AUTH)[KEY_IS_AUTHENTICATED] || false
  }

  getAccessToken(): string {
    return this.store.getItem(KEY_AUTH)[KEY_ACCESS_TOKEN] || ''
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
    return this.isAuthenticated$;
  }
}
