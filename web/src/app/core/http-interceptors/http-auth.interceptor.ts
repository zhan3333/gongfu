import { Injectable } from '@angular/core';
import {
  HttpRequest,
  HttpHandler,
  HttpEvent,
  HttpInterceptor, HttpHeaders
} from '@angular/common/http';
import { Observable } from 'rxjs';
import { AuthService } from '../auth/auth.service';
import { environment } from '../../../environments/environment';

@Injectable()
export class HttpAuthInterceptor implements HttpInterceptor {

  constructor(
    private authService: AuthService,
  ) {
  }

  intercept(request: HttpRequest<unknown>, next: HttpHandler): Observable<HttpEvent<unknown>> {
    if (request.url.startsWith('http') || request.url.indexOf('/i18n/') !== -1) {
      return next.handle(request)
    }
    const newReq = request.clone({
      headers: request.headers.set('Authorization', 'Bearer ' + this.authService.getAccessToken()),
      url: environment.api + request.url,
    })
    return next.handle(newReq);
  }
}
