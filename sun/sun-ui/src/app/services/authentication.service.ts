import { apConfig } from './../../global';
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import 'rxjs/add/operator/map';
import { HttpHeaders } from '@angular/common/http/';

@Injectable()
export class AuthenticationService {
  private apiLoginUrl = apConfig.API_ENDPOINT_URL + '/login';

  constructor(private http: HttpClient) {}

  /**
   * Make the Login Call in API, return the API response.
   * @param username
   * @param password
   */
  authentication(username: string, password: string) {
    let myHeader = new HttpHeaders();
    myHeader.append('Content-Type', 'application/json');
    let obj = {
      username: username,
      password: password
    };
    let data = JSON.stringify(obj);
    return this.http
      .post<any>(this.apiLoginUrl, data, { headers: myHeader })
      .map(user => {
        // login successful if there's a jwt token in the response
        if (user && user.token) {
          // store user details and jwt token in local storage to keep user logged in between page refreshes
          localStorage.setItem('currentUser', JSON.stringify(user));
          localStorage.setItem('authenticationToken', user.token);
        }
        return user;
      });
  }

  logout() {
    localStorage.removeItem('currentUser');
  }
}
