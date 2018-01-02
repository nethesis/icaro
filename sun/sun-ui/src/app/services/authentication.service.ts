import { apConfig } from './../../global';
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import 'rxjs/add/operator/map';

@Injectable()
export class AuthenticationService {
  private apiLoginUrl = apConfig.API_ENDPOINT_URL;

  constructor(private http: HttpClient) {}

  /**
   * Make the Login Call in API, return the API response.
   * @param username
   * @param password
   */
  authentication(username: string, password: string) {
    return this.http
      .post<any>(this.apiLoginUrl, {
        username: username,
        password: password
      })
      .map(res => res);
  }
}
