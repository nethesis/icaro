import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import 'rxjs/add/operator/map';

@Injectable()
export class AuthenticationService {
  constructor(private http: HttpClient) {}

  /**
   * Make the Login Call in API, return the API response.
   * @param username
   * @param password
   */
  authentication(username: string, password: string) {
    return this.http
      .post<any>('http://hstest.neth.eu:8080/api/login', {
        username: username,
        password: password
      })
      .map(res => res);
  }
}
